package image

import (
	"fmt"
	"github.com/chai2010/webp"
	"image"
	"image/draw"
	"image/jpeg"
	"io/fs"
	"os"
)

type Webp2JPGHandler struct {
}

func NewWebp2JPGHandler() *Webp2JPGHandler {
	return &Webp2JPGHandler{}
}

func (h *Webp2JPGHandler) Convert(inputPath, outputPath string, perm fs.FileMode) (int64, error) {
	var size int64 = 0
	quality := 80

	fileIn, err := os.Open(inputPath)
	if err != nil {
		return size, err
	}
	defer fileIn.Close()

	// Декодируем WebP изображение
	img, err := webp.Decode(fileIn)
	if err != nil {
		return size, err
	}

	outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return size, err
	}
	defer outputFile.Close()

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)

	options := &jpeg.Options{
		Quality: quality,
	}

	err = jpeg.Encode(outputFile, rgba, options)
	if err != nil {
		return size, err
	}

	fileInfo, err := outputFile.Stat()
	if err != nil {
		fmt.Printf("Warning: could not get file info: %v\n", err)
	} else {
		size = fileInfo.Size()
	}

	return size, err
}
