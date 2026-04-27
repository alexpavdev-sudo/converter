package dto

import (
	"converter/components/mapper"
	"converter/entities"
	"time"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

type File struct {
	ID            uint                `json:"id"`
	Extension     string              `json:"extension"`
	OriginalName  string              `json:"original_name"`
	Format        string              `json:"format"`
	Size          int64               `json:"size"`
	SizeProcessed int64               `json:"size_processed"`
	Status        entities.FileStatus `json:"status"`
	StatusLabel   string              `json:"status_label"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}

var fileToDTO = mapper.New(func(e entities.File) File {
	return File{
		ID:            e.ID,
		Extension:     e.Extension,
		OriginalName:  e.OriginalName,
		Format:        e.Format,
		Size:          e.Size,
		SizeProcessed: e.SizeProcessed,
		Status:        e.Status,
		StatusLabel:   e.Status.String(),
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
})

func FileToDTO() *mapper.Mapper[entities.File, File] {
	return fileToDTO
}
