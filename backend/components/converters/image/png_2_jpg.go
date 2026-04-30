package image

import (
	"converter/components/converters"
	"converter/entities"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/fs"
	"os"
)

type PNG2JPGHandler struct {
	converters.BaseConverter
	outputPath string
}

func NewPNG2JPGHandler(outputPath string) *PNG2JPGHandler {
	return &PNG2JPGHandler{outputPath: outputPath}
}

func (h *PNG2JPGHandler) GetOutputPath() string {
	return h.outputPath
}

func (h *PNG2JPGHandler) Rollback() error {
	return converters.BaseConverter{}.Rollback(h.GetOutputPath())
}

func (h *PNG2JPGHandler) Convert(file entities.File, perm fs.FileMode) (int64, error) {
	var size int64 = 0
	quality := 80
	inputPath := file.PathFull()

	fileIn, err := os.Open(inputPath)
	if err != nil {
		return size, err
	}
	defer fileIn.Close()

	img, err := png.Decode(fileIn)
	if err != nil {
		return size, err
	}

	outputFile, err := os.OpenFile(h.GetOutputPath(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
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
