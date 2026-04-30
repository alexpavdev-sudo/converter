package repositories

import (
	"converter/entities"
)

type FileRepositoryInterface interface {
	GetFiles(guestId uint) ([]entities.File, error)
	GetCountFiles(guestId uint) (int64, error)
	GetFile(guestId uint, fileId uint) (entities.File, error)
	GetFileById(fileId uint) (entities.File, error)

	SetProcessedPath(fileID uint, processedPath string) error

	SetStatus(fileId uint, status entities.FileStatus) error
	SetStatusProcessed(fileID uint, size int64) error
	SetStatusError(fileID uint, msgErr string) error

	ExistFile(fileID uint) (bool, error)

	CloseRepo() error
}
