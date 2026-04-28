package delete

import (
	"converter/app"
	"converter/entities"
	"converter/helpers"
	"gorm.io/gorm"
	"os"
)

type DeleteService struct {
	db *gorm.DB
}

func NewDeleteService() *DeleteService {
	return &DeleteService{
		db: app.App().DB,
	}
}

func DeleteFile(file entities.File) error {
	tx := app.App().StartTransaction()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if err := tx.Where("id = ?", file.ID).Delete(&file).Error; err != nil {
		return err
	}

	if helpers.ExistsFile(file.PathFull()) {
		if err := os.Remove(file.PathFull()); err != nil {
			return err
		}
	}
	if helpers.ExistsFile(file.ProcessedPathFull()) {
		if err := os.Remove(file.ProcessedPathFull()); err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	//todo cache

	return nil
}
