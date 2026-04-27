package converters

import (
	"converter/components/converters/image"
	"converter/entities"
	"errors"
	"io/fs"
)

type ConverterInterface interface {
	Convert(inputPath, outputPath string, perm fs.FileMode) (int64, error)
}

type Factory struct{}

func (f Factory) Create(file entities.File) (ConverterInterface, error) {
	switch file.Extension {
	case "jpg", "jpeg":
		if file.Format == "webp" {
			return image.NewJPG2WebpHandler(), nil
		}
	}

	return nil, errors.New("not found converter")
}
