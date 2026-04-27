package converters

import (
	"converter/components/converters/image"
	"converter/entities"
	"errors"
)

type ConverterInterface interface {
	Convert(inputPath, outputPath string) error
}

type Factory struct{}

func (f Factory) Create(file entities.File) (ConverterInterface, error) {
	if file.Extension == "jpg" && file.Format == "webp" {
		return image.NewJPG2WebpHandler(), nil
	}

	return nil, errors.New("not found converter")
}
