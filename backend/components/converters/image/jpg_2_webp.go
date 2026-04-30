package image

import (
	"converter/components/converters"
	"converter/entities"
	"fmt"
	"github.com/chai2010/webp"
	"image"
	"image/draw"
	"image/jpeg"
	"io/fs"
	"os"
)

type JPG2WebpHandler struct {
	outputPath string
}

func NewJPG2WebpHandler(outputPath string) *JPG2WebpHandler {
	return &JPG2WebpHandler{outputPath: outputPath}
}

func (h *JPG2WebpHandler) GetOutputPath() string {
	return h.outputPath
}

func (h *JPG2WebpHandler) Rollback() error {
	return converters.BaseConverter{}.Rollback(h.GetOutputPath())
}

func (h *JPG2WebpHandler) Convert(file entities.File, perm fs.FileMode) (int64, error) {
	var size int64 = 0
	quality := 80
	inputPath := file.PathFull()
	fileIn, err := os.Open(inputPath)
	if err != nil {
		return size, err
	}
	defer fileIn.Close()

	img, err := jpeg.Decode(fileIn)
	if err != nil {
		return size, err
	}

	outputFile, err := os.OpenFile(h.GetOutputPath(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return size, err
	}
	defer outputFile.Close()

	// Конвертируем в RGBA, если нужно (WebP предпочитает RGBA)
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)

	options := &webp.Options{
		Lossless: false,
		Quality:  float32(quality),
	}

	err = webp.Encode(outputFile, rgba, options)

	fileInfo, err := outputFile.Stat()
	if err != nil {
		fmt.Printf("Warning: could not get file info: %v\n", err)
	} else {
		size = fileInfo.Size()
	}

	return size, err
}
