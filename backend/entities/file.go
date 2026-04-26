package entities

import (
	"converter/config"
	"path/filepath"
	"time"
)

type File struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	StoredName    string     `gorm:"column:stored_name;type:varchar(64);not null;uniqueIndex" json:"-"`
	Extension     string     `gorm:"column:extension;type:varchar(20);not null" json:"extension"`
	OriginalName  string     `gorm:"column:original_name;type:varchar(255);not null" json:"original_name"`
	Path          string     `gorm:"column:path;type:text;not null" json:"-"`
	Format        string     `gorm:"column:format;type:varchar(50);not null" json:"format"`
	Size          int64      `gorm:"column:size;type:bigint;not null" json:"size"`
	Status        FileStatus `gorm:"column:status;type:tinyint;not null;default:0;index" json:"status"`
	ProcessedPath string     `gorm:"column:processed_path;type:text" json:"-"`
	CreatedAt     time.Time  `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;type:timestamp;autoUpdateTime" json:"updated_at"`
}

func (File) TableName() string {
	return "files"
}

func (f *File) PathFull() string {
	if f.Path == "" {
		return ""
	}
	return filepath.Join(config.UploadDir, f.Path)
}

type FileStatus uint8

const (
	StatusQueued     FileStatus = iota // 0 - В очереди
	StatusProcessing                   // 1 - В обработке
	StatusProcessed                    // 2 - Обработан
	StatusError                        // 3 - Ошибка
)

func (s FileStatus) String() string {
	switch s {
	case StatusQueued:
		return "queued"
	case StatusProcessing:
		return "processing"
	case StatusProcessed:
		return "processed"
	case StatusError:
		return "error"
	default:
		return "unknown"
	}
}
