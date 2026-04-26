package delete

import (
	"converter/app"
	"converter/entities"
	"github.com/gin-contrib/sessions"
	"gorm.io/gorm"
	"os"
)

type DeleteService struct {
	db      *gorm.DB
	session sessions.Session
}

func NewDeleteService(session sessions.Session) *DeleteService {
	return &DeleteService{
		db:      app.App().DB,
		session: session,
	}
}

func (s *DeleteService) DeleteFile(guestId uint, fileId uint) error {
	tx := app.App().StartTransaction()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	var file entities.File
	err := tx.Model(&entities.File{}).
		Select("files.*").
		Joins("INNER JOIN guest_files ON guest_files.file_id = files.id").
		Where("guest_files.guest_id = ? AND guest_files.file_id = ?", guestId, fileId).
		Take(&file).Error
	if err != nil {
		return err
	}
	if err := tx.Where("id = ?", file.ID).Delete(&file).Error; err != nil {
		return err
	}
	if err := os.Remove(file.PathFull()); err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	app.ClearCache(guestId)
	return nil
}
