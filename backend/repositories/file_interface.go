package repositories

import (
	"converter/entities"
)

type FileRepositoryInterface interface {
	GetFiles(guestId uint) ([]entities.File, error)
	GetCountFiles(guestId uint) (int64, error)
	GetFile(guestId uint, fileId uint) (entities.File, error)
	CloseRepo() error
}
