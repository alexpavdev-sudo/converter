package converter

import (
	"converter/app"
	"converter/components/converters"
	"converter/config"
	"converter/entities"
	"converter/helpers"
	"converter/repositories"
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

	//todo проверть статус
	file, err := c.repo.GetFileById(c.fileId)
	if err != nil {
		return fmt.Errorf("error db: %s", err.Error())
	}

	err = c.repo.SetStatus(file.ID, entities.StatusProcessing)
	defer func() {
		if err != nil {
			c.repo.SetStatusError(file.ID, err.Error())
		}
	}()
	if err != nil {
		return fmt.Errorf("error set status")
	}

	converter, err := converters.Factory{}.Create(file)
	if err != nil {
		return fmt.Errorf("error factory converter: %s", err.Error())
	}

	processedPath, err := c.generateUniqueProcessedPath()
	if err != nil {
		return fmt.Errorf("error generate unique processed path: %s", err.Error())
	}
	err = c.repo.SetProcessedPath(file.ID, processedPath)
	if err != nil {
		return fmt.Errorf("failed to update processed_path")
	}

	size, err := converter.Convert(file.PathFull(), entities.ProcessedPathFull(sql.NullString{String: processedPath, Valid: true}, file.Format), PermFile)
	if err != nil {
		return fmt.Errorf("error convert: %s", err.Error())
	}

	err = c.repo.SetStatusProcessed(file.ID, size)
	if err != nil {
		return fmt.Errorf("failed to update status processed")
	}

	err = c.repo.SetStatus(file.ID, entities.StatusProcessed)
	if err != nil {
		return fmt.Errorf("error set status")
	}

	log.Printf("done: %d", c.fileId)
	return nil
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
