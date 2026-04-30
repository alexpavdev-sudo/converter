package formater

import (
	"fmt"
	"strings"
)

type Category string

const (
	CategoryDocument Category = "document"
	CategoryImage    Category = "image"
	CategoryVideo    Category = "video"
)

type Format struct {
	Ext      string   `json:"ext"`
	Category Category `json:"category"`
	Name     string   `json:"name"`
	MimeType string   `json:"mime_type"`
}

var formats = map[string]Format{
	//"pdf":  {Ext: "pdf", Category: CategoryDocument, Name: "PDF Document", MimeType: "application/pdf"},
	//"doc":  {Ext: "doc", Category: CategoryDocument, Name: "Word Document", MimeType: "application/msword"},
	//"docx": {Ext: "docx", Category: CategoryDocument, Name: "Word Document", MimeType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
	//"txt":  {Ext: "txt", Category: CategoryDocument, Name: "Text File", MimeType: "text/plain"},
	"jpg": {Ext: "jpg", Category: CategoryImage, Name: "JPEG Image", MimeType: "image/jpeg"},
	"png": {Ext: "png", Category: CategoryImage, Name: "PNG Image", MimeType: "image/png"},
	//"gif":  {Ext: "gif", Category: CategoryImage, Name: "GIF Image", MimeType: "image/gif"},
	"webp": {Ext: "webp", Category: CategoryImage, Name: "WebP Image", MimeType: "image/webp"},
	"mp4":  {Ext: "mp4", Category: CategoryVideo, Name: "MP4 Video", MimeType: "video/mp4"},
	"avi":  {Ext: "avi", Category: CategoryVideo, Name: "AVI Video", MimeType: "video/x-msvideo"},
	//"mov":  {Ext: "mov", Category: CategoryVideo, Name: "QuickTime Video", MimeType: "video/quicktime"},
}

var conversionMap = map[string][]string{
	"pdf":  {"doc", "docx", "txt"},
	"doc":  {"pdf", "docx", "txt"},
	"docx": {"pdf", "doc", "txt"},
	"txt":  {"pdf", "doc", "docx"},

	"jpg":  {"png", "webp", "gif"},
	"jpeg": {"png", "webp", "gif"},
	"png":  {"jpg", "webp", "gif"},
	"gif":  {"jpg", "png", "webp"},
	"webp": {"jpg", "png", "gif"},

	"mp4": {"avi", "mov"},
	"avi": {"mp4", "mov"},
	"mov": {"mp4", "avi"},
}

type FormatService struct{}

func NewFormatService() *FormatService {
	return &FormatService{}
}

func (s *FormatService) GetFormat(ext string) (Format, bool) {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	f, ok := formats[ext]
	return f, ok
}

func (s *FormatService) GetFormats() map[string]Format {
	return formats
}

func (s *FormatService) IsValidFormat(ext string) bool {
	_, ok := formats[strings.ToLower(strings.TrimPrefix(ext, "."))]
	return ok
}

func (s *FormatService) GetPossibleConversions(fromExt string) ([]Format, error) {
	fromExt = strings.ToLower(strings.TrimPrefix(fromExt, "."))

	_, ok := formats[fromExt]
	if !ok {
		return nil, fmt.Errorf("format %s not supported", fromExt)
	}

	targets, ok := conversionMap[fromExt]
	if !ok {
		return []Format{}, nil
	}

	result := make([]Format, 0, len(targets))
	for _, targetExt := range targets {
		if targetFormat, ok := formats[targetExt]; ok {
			result = append(result, targetFormat)
		}
	}

	return result, nil
}

func (s *FormatService) CanConvert(fromExt, toExt string) bool {
	fromExt = strings.ToLower(fromExt)
	toExt = strings.ToLower(toExt)

	targets, ok := conversionMap[fromExt]
	if !ok {
		return fromExt == toExt
	}

	for _, target := range targets {
		if target == toExt {
			return true
		}
	}
	return false
}

func (s *FormatService) GetFormatsByCategory(category Category) []Format {
	result := make([]Format, 0)
	for _, f := range formats {
		if f.Category == category {
			result = append(result, f)
		}
	}
	return result
}
