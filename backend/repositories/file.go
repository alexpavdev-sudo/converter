package repositories

import (
	"converter/entities"
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
