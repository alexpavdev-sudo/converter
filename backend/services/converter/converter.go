package converter

import (
	"converter/app"
	"converter/components/converters"
	"converter/config"
	"converter/entities"
	"converter/helpers"
	"database/sql"
	"fmt"
	"log"
	"os"
)

const PermFile = 0600

type Converter struct {
	fileId uint
}

func NewConverter(fileId uint) *Converter {
	return &Converter{fileId: fileId}
}

func (c *Converter) Run() {
	log.Printf("received a message: %d", c.fileId)

	if err := os.MkdirAll(config.ConvertedDir, 0700); err != nil {
		log.Printf("failed to create directory")
		return
	}

	file, err := app.App().FileRepo.GetFileById(c.fileId)
	if err != nil {
		log.Printf("error db: %s", err.Error())
		return
	}
	result := app.App().DB.Model(&entities.File{}).
		Where("id = ?", file.ID).
		Updates(map[string]interface{}{"Status": entities.StatusProcessing})

	converter, err := converters.Factory{}.Create(file)
	if err != nil {
		log.Printf("error factory converter: %s", err.Error())
		return
	}
	processedPath, err := c.generateUniqueProcessedPath()
	if err != nil {
		log.Printf("error generate unique processed path: %s", err.Error())
		return
	}
	size, err := converter.Convert(file.PathFull(), entities.ProcessedPathFull(sql.NullString{String: processedPath, Valid: true}, file.Format), PermFile)
	if err != nil {
		log.Printf("error convert: %s", err.Error())
		return
	}

	result = app.App().DB.Model(&entities.File{}).
		Where("id = ?", file.ID).
		Updates(map[string]interface{}{"processed_path": processedPath, "Status": entities.StatusProcessed, "size_processed": size})

	if result.Error != nil {
		log.Printf("failed to update processed_path")
		return
	}

	log.Printf("Done")
}

func (c *Converter) generateUniqueProcessedPath() (string, error) {
	const maxAttempts = 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		randStr, err := helpers.GenerateRandomStoredName(128)
		if err != nil {
			return "", err
		}

		var exists bool
		err = app.App().DB.Raw(
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
