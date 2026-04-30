package factory

import (
	"converter/components/converters"
	"converter/components/converters/image"
	"converter/components/converters/video"
	"converter/components/formater"
	"converter/entities"
	"errors"
)

type Factory struct{}

func (f Factory) Create(file entities.File, outputPath string) (converters.ConverterInterface, error) {
	formatService := formater.NewFormatService()

	if formatService.CanConvert(file.Extension, file.Format) {
		switch file.Extension {
		case "jpg", "jpeg":
			if file.Format == "webp" {
				return image.NewJPG2WebpHandler(outputPath), nil
			}
		case "png":
			if file.Format == "jpg" {
				return image.NewPNG2JPGHandler(outputPath), nil
			}
		case "webp":
			if file.Format == "jpg" {
				return image.NewWebp2JPGHandler(outputPath), nil
			}
		case "mp4":
			if file.Format == "avi" {
				return video.NewMP4ToAVIHandler(outputPath), nil
			}
		}
		return nil, errors.New("not found converter")
	}

	return nil, errors.New("cannot be converted")
}
