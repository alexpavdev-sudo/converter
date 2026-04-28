package repositories

import (
	"converter/entities"
)

type FileRepositoryInterface interface {
	GetFiles(guestId uint) ([]entities.File, error)
	GetCountFiles(guestId uint) (int64, error)
	GetFile(guestId uint, fileId uint) (entities.File, error)
	GetFileById(fileId uint) (entities.File, error)
	SetStatus(fileId uint, status entities.FileStatus) error
	UpdateProcessed(fileID uint, processedPath string, size int64) error
	UpdateError(fileID uint, msgErr string) error
	CloseRepo() error
}
