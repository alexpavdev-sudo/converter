package converter

import (
	"converter/app"
	"converter/components/converters/factory"
	"converter/config"
	"converter/entities"
	"converter/helpers"
	"converter/repositories"
	"converter/services/notify"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"log"
	"os"
)

const PermFile = 0600

type Converter struct {
	fileId uint
	repo   repositories.FileRepositoryInterface
	db     *gorm.DB
}

func NewConverter(fileId uint) *Converter {
	return &Converter{fileId: fileId, repo: app.App().FileRepo, db: app.App().DB}
}

func (c *Converter) Run() error {
	log.Printf("received a message: %d", c.fileId)
	var err error
	if err = os.MkdirAll(config.ConvertedDir, 0700); err != nil {
		return fmt.Errorf("failed to create directory")
	}

	file, err := c.repo.GetFileById(c.fileId)
	if err != nil {
		return fmt.Errorf("error db: %s", err.Error())
	}
	if file.Status != entities.StatusQueued {
		return fmt.Errorf("error status: %s", file.Status.String())
	}

	err = c.repo.SetStatus(file.ID, entities.StatusProcessing)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic: %v", r)
			_ = c.repo.SetStatusError(file.ID, "panic")
			c.notifyError(&file, "panic")
		} else if err != nil {
			_ = c.repo.SetStatusError(file.ID, err.Error())
			c.notifyError(&file, err.Error())
		}
	}()
	if err != nil {
		return fmt.Errorf("error set status")
	}

	processedPath, err := c.generateUniqueProcessedPath()
	if err != nil {
		return fmt.Errorf("error generate unique processed path: %s", err.Error())
	}
	err = c.repo.SetProcessedPath(file.ID, processedPath)
	if err != nil {
		return fmt.Errorf("error: failed to update processed_path")
	}

	outputPath := entities.ProcessedPathFull(sql.NullString{String: processedPath, Valid: true}, file.Format)
	converter, err := factory.Factory{}.Create(file, outputPath)
	if err != nil {
		return fmt.Errorf("error factory converter: %s", err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			_ = converter.Rollback()
		} else if err != nil {
			_ = converter.Rollback()
		}
	}()
	size, err := converter.Convert(file, PermFile)
	if err != nil {
		return fmt.Errorf("error convert: %s", err.Error())
	}

	err = c.repo.SetStatusProcessed(file.ID, size)
	if err != nil {
		return fmt.Errorf("error: failed to update status processed")
	}
	exist, _ := c.repo.ExistFile(file.ID)
	if !exist {
		err = fmt.Errorf("error: the file %d has been deleted", file.ID)
		return err
	}

	c.notify(&file)
	log.Printf("done: %d", c.fileId)

	return nil
}

func (c *Converter) notify(file *entities.File) {
	err := notify.NotifyService{}.NotifyConvertFileSuccess(file.ID, file.OriginalName)
	if err != nil {
		log.Printf("error save notification: %s", err.Error())
	}
}

func (c *Converter) notifyError(file *entities.File, error string) {
	err := notify.NotifyService{}.NotifyConvertFileError(file.ID, file.OriginalName, error)
	if err != nil {
		log.Printf("error save notification: %s", err.Error())
	}
}

func (c *Converter) generateUniqueProcessedPath() (string, error) {
	const maxAttempts = 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		randStr, err := helpers.GenerateRandomStoredName(128)
		if err != nil {
			return "", err
		}

		var exists bool
		err = c.db.Raw(
			"SELECT EXISTS(SELECT 1 FROM files WHERE processed_path = ?)",
			randStr,
		).Scan(&exists).Error
		if err != nil {
			return "", fmt.Errorf("database error: %w", err)
		}

		if !exists {
			return randStr, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique processed_path after %d attempts", maxAttempts)
}
