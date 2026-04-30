package converters

import (
	"converter/entities"
	"converter/helpers"
	"io/fs"
	"log"
	"os"
)

type ConverterInterface interface {
	Convert(file entities.File, perm fs.FileMode) (int64, error)
	Rollback() error
	GetOutputPath() string
}

type BaseConverter struct{}

func (c BaseConverter) Rollback(outputPath string) error {
	log.Println("rollback convert")
	if helpers.ExistsFile(outputPath) {
		if err := os.Remove(outputPath); err != nil {
			return err
		}
	}
	return nil
}
