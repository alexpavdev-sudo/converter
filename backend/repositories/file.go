package repositories

import (
	"converter/entities"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{
		db: db,
	}
}

func (r *FileRepository) CloseRepo() error {
	return nil
}

func (r *FileRepository) UpdateError(fileID uint, details string) error {
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	err := tx.Model(&entities.File{}).
		Where("id = ?", fileID).
		Updates(map[string]interface{}{"Status": entities.StatusError}).Error
	if err != nil {
		return err
	}

	errorModel := &entities.Error{
		FileID:  fileID,
		Details: details,
	}
	if err := tx.Create(errorModel).Error; err != nil {
		return fmt.Errorf("error save error entity")
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func (r *FileRepository) UpdateProcessed(fileID uint, processedPath string, size int64) error {
	return r.db.Model(&entities.File{}).
		Where("id = ?", fileID).
		Updates(map[string]interface{}{"Status": entities.StatusProcessed, "processed_path": processedPath, "size_processed": size}).Error
}

func (r *FileRepository) SetStatus(fileId uint, status entities.FileStatus) error {
	result := r.db.Model(&entities.File{}).
		Where("id = ?", fileId).
		Updates(map[string]interface{}{"Status": status})

	return result.Error
}

func (r *FileRepository) GetFiles(guestId uint) ([]entities.File, error) {
	var files []entities.File
	err := r.db.Model(&entities.File{}).
		Select("files.*").
		Joins("INNER JOIN guest_files ON guest_files.file_id = files.id").
		Where("guest_files.guest_id = ?", guestId).
		Order("files.created_at DESC").
		Scan(&files).Error

	if err != nil {
		return nil, err
	}
	if files == nil {
		files = []entities.File{}
	}
	return files, nil
}

func (r *FileRepository) GetCountFiles(guestId uint) (int64, error) {
	var count int64
	err := r.db.Model(&entities.File{}).
		Joins("INNER JOIN guest_files ON guest_files.file_id = files.id").
		Where("guest_files.guest_id = ?", guestId).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *FileRepository) GetFile(guestId uint, fileId uint) (entities.File, error) {
	var file entities.File
	err := r.db.Model(&entities.File{}).
		Select("files.*").
		Joins("INNER JOIN guest_files ON guest_files.file_id = files.id").
		Where("guest_files.guest_id = ? AND guest_files.file_id = ?", guestId, fileId).
		Order("files.created_at DESC").
		Scan(&file).Error

	if err != nil {
		return file, err
	}
	if file.ID == 0 {
		return file, gorm.ErrRecordNotFound
	}

	return file, nil
}

func (r *FileRepository) GetFileById(fileId uint) (entities.File, error) {
	var file entities.File
	err := r.db.Model(&entities.File{}).
		Select("*").
		Where("id = ?", fileId).
		First(&file).Error

	if err != nil {
		return file, err
	}
	if file.ID == 0 {
		return file, gorm.ErrRecordNotFound
	}

	return file, nil
}
