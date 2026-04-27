package image

import (
	"github.com/chai2010/webp"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

type JPG2WebpHandler struct {
}

func NewJPG2WebpHandler() *JPG2WebpHandler {
	return &JPG2WebpHandler{}
}

func (h *JPG2WebpHandler) Convert(inputPath, outputPath string) error {
	quality := 80
	fileIn, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer fileIn.Close()

	img, err := jpeg.Decode(fileIn)
	if err != nil {
		return err
	}

	outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Конвертируем в RGBA, если нужно (WebP предпочитает RGBA)
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)

	options := &webp.Options{
		Lossless: false,
		Quality:  float32(quality),
	}

	return webp.Encode(outputFile, rgba, options)
}
