package uploader

import (
	"context"
	"converter/app"
	"converter/components/cache"
	"converter/dto/inner"
	"converter/entities"
	"converter/helpers"
	"converter/repositories"
	"converter/services/user"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type StreamFileUploader struct {
	reader      *multipart.Reader
	db          *gorm.DB
	maxFileSize int64
	maxSize     int64
	savedFiles  []entities.File
	userService *user.UserService
}

func (u *StreamFileUploader) CountSavedFiles() int {
	return len(u.savedFiles)
}

func NewStreamFileUploader(reader *multipart.Reader, maxFileSize int64, maxSize int64, session sessions.Session) *StreamFileUploader {
	return &StreamFileUploader{
		reader:      reader,
		db:          app.App().DB,
		maxFileSize: maxFileSize,
		maxSize:     maxSize,
		userService: user.NewUserService(session),
	}
}

func (u *StreamFileUploader) Upload() error {
	var formats []string

	for {
		part, err := u.nextPart(u.reader)
		if err != nil {
			return err
		}
		if part == nil {
			break
		}

		err = func() error {
			defer part.Close()

			switch part.FormName() {
			case "formats":
				formatBytes, err := io.ReadAll(part)
				if err != nil {
					return fmt.Errorf("failed to read format value: %w", err)
				}
				formats = append(formats, string(formatBytes))
			case "images":

				if part.FileName() != "" {
					fileIdx := len(u.savedFiles)
					if fileIdx >= len(formats) {
						return fmt.Errorf("missing format for file %s", part.FileName())
					}
					fileRecord, err := u.saveFilePart(part, formats[fileIdx])
					if err != nil {
						return fmt.Errorf("failed to save file %s: %w", part.FileName(), err)
					}
					u.savedFiles = append(u.savedFiles, *fileRecord)
				}
			}

			return nil
		}()

		if err != nil {
			return err
		}
	}

	if len(u.savedFiles) == 0 {
		return fmt.Errorf("no files uploaded")
	}

	return nil
}

func (u *StreamFileUploader) nextPart(reader *multipart.Reader) (*multipart.Part, error) {
	part, err := reader.NextPart()
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return part, nil
}

func (u *StreamFileUploader) saveFilePart(part *multipart.Part, format string) (*entities.File, error) {
	storedName, err := u.generateUniqueStoredName()
	if err != nil {
		return nil, err
	}
	ext := strings.TrimPrefix(filepath.Ext(part.FileName()), ".")
	var personalDir string
	if u.userService.IsAuthenticated() {
		personalDir, err = u.userService.UserPersonalDir()
	} else {
		personalDir, err = u.userService.GuestPersonalDir(false)
	}
	if err != nil {
		return nil, err
	}
	personalDirFull, err := u.userService.GuestPersonalDir(true)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(personalDirFull, 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory")
	}
	finalPathFull := filepath.Join(personalDirFull, storedName+"."+ext)
	fileRes, err := os.OpenFile(finalPathFull, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	needRemoveOnExit := true
	defer func() {
		_ = fileRes.Close()
		if needRemoveOnExit {
			_ = os.Remove(finalPathFull)
		}
	}()

	written, err := io.Copy(fileRes, io.LimitReader(part, u.maxFileSize))
	if err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	tx := app.App().StartTransaction()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	fileRecord := &entities.File{
		StoredName:    storedName,
		Extension:     strings.ToLower(ext),
		OriginalName:  part.FileName(),
		Path:          filepath.Join(personalDir, storedName),
		ProcessedPath: sql.NullString{String: "", Valid: false},
		Format:        strings.ToLower(format),
		Size:          written,
		Status:        entities.StatusQueued,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := tx.Create(fileRecord).Error; err != nil {
		return nil, fmt.Errorf("error save file: %s", err.Error())
	}

	if u.userService.IsAuthenticated() {
		//todo
	} else {
		guestID, err := u.userService.GuestId()
		if err != nil {
		}
		guestFile := &entities.GuestFile{
			FileID:  fileRecord.ID,
			GuestID: guestID,
		}
		if err := tx.Create(guestFile).Error; err != nil {
			return nil, fmt.Errorf("error save guest_file")
		}

		if err := tx.Commit().Error; err != nil {
			return nil, fmt.Errorf("commit failed: %w", err)
		}
		clearCache(guestID)
		sendConversion(fileRecord.ID)
	}
	needRemoveOnExit = false
	return fileRecord, nil
}

func clearCache(guestId uint) {
	cache, err := cache.CachedFactory{}.Create()
	tag := repositories.CachedFileRepository{}.TagGuest(guestId)
	if err == nil {
		_ = cache.DeleteByTag([]string{tag})
	}
}

func exit(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func sendConversion(fileId uint) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	exit(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	exit(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durability
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	exit(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(&inner.MessageDto{FileID: fileId})
	if err != nil {
		log.Printf("error: %s", err)
	}
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	exit(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}

func (u *StreamFileUploader) generateUniqueStoredName() (string, error) {
	for {
		randStr, err := helpers.GenerateRandomStoredName(128)
		if err != nil {
			return "", err
		}

		var exists bool
		err = u.db.Raw(
			"SELECT EXISTS(SELECT 1 FROM files WHERE stored_name = ?)",
			randStr,
		).Scan(&exists).Error
		if err != nil {
			return "", fmt.Errorf("database error: %w", err)
		}

		if !exists {
			return randStr, nil
		}
	}
}
