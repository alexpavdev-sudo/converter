package delete

import (
	"converter/app"
	"converter/components/cache"
	"converter/entities"
	"converter/helpers"
	"converter/repositories"
	"os"
)

type DeleteService struct {
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
	clearCache()

	return nil
}

func clearCache() {
	cache, err := cache.CachedFactory{}.Create()
	tag := repositories.CachedFileRepository{}.TagAll()
	if err == nil {
		_ = cache.DeleteByTag([]string{tag})
	}
}
