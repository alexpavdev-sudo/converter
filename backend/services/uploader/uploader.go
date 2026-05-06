package uploader

import (
	"converter/app"
	"converter/components/cache"
	"converter/components/queue_conversion"
	"converter/config"
	"converter/entities"
	"converter/helpers"
	"converter/repositories"
	"converter/services/user"
	"database/sql"
	"fmt"
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
	userService user.UserService
	queue       queue_conversion.ConverterQueue
}

func (u *StreamFileUploader) CountSavedFiles() int {
	return len(u.savedFiles)
}

func NewStreamFileUploader(reader *multipart.Reader, maxFileSize int64, maxSize int64, userService user.UserService, queue queue_conversion.ConverterQueue) *StreamFileUploader {
	return &StreamFileUploader{
		reader:      reader,
		db:          app.App().DB,
		maxFileSize: maxFileSize,
		maxSize:     maxSize,
		userService: userService,
		queue:       queue,
	}
}

func (u *StreamFileUploader) Upload() error {
	err := u.checkAccess()
	if err != nil {
		return err
	}

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

func (u *StreamFileUploader) checkAccess() error {
	if u.userService.IsAuthenticated() {
		//todo
	} else {
		guestId, err := u.userService.InitGuestID()
		if err != nil {
			return err
		}

		var count int64
		err = u.db.Model(&entities.File{}).
			Joins("INNER JOIN guest_files ON guest_files.file_id = files.id").
			Where("guest_files.guest_id = ?", guestId).
			Where("files.status IN (?)", []entities.FileStatus{entities.StatusQueued, entities.StatusProcessing}).
			Count(&count).Error
		if err != nil {
			return err
		}
		if count >= config.ProcessingFilesMaxCount {
			return fmt.Errorf("исчерпан лимит одновременной обработки файлов: максимально до %d", config.ProcessingFilesMaxCount)
		}
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
		guestID, err := u.userService.InitGuestID()
		if err != nil {
			return nil, err
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
		u.sendConversion(fileRecord.ID)
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

func (u *StreamFileUploader) sendConversion(fileId uint) {
	err := u.queue.Push(fileId)
	if err != nil {
		log.Println(err)
	}
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
