package repositories

import (
	"converter/components/cache"
	"converter/entities"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type CachedFileRepository struct {
	repo  *FileRepository
	cache cache.Cache
}

func NewCachedFileRepository(db *gorm.DB) (*CachedFileRepository, error) {
	cache, err := cache.CachedFactory{}.Create()
	if err != nil {
		return nil, err
	}
	return &CachedFileRepository{
		repo:  NewFileRepository(db),
		cache: cache,
	}, nil
}

func (r *CachedFileRepository) CloseRepo() error {
	return r.cache.Close()
}

func (r *CachedFileRepository) SetProcessedPath(fileID uint, processedPath string) error {
	return r.repo.SetProcessedPath(fileID, processedPath)
}

func (r *CachedFileRepository) SetStatusProcessed(fileID uint, size int64) error {
	return r.repo.SetStatusProcessed(fileID, size)
}

func (r *CachedFileRepository) SetStatusError(fileID uint, msgErr string) error {
	return r.repo.SetStatusError(fileID, msgErr)
}

func (r *CachedFileRepository) SetStatus(fileId uint, status entities.FileStatus) error {
	return r.repo.SetStatus(fileId, status)
}

func (r *CachedFileRepository) GetFiles(guestId uint) ([]entities.File, error) {
	key := r.Key(fmt.Sprintf("files:%d", guestId))

	var files []entities.File
	err := r.cache.Get(key, &files)
	if err == nil {
		return files, nil
	}

	files, err = r.repo.GetFiles(guestId)
	if err != nil {
		return nil, err
	}

	if err := r.cache.Set(key, files, []string{
		r.TagGuest(guestId),
		r.TagAll(),
	}); err != nil {
		log.Printf("failed to set cache: %v", err)
	}

	return files, nil
}

func (r *CachedFileRepository) GetCountFiles(guestId uint) (int64, error) {
	key := r.Key(fmt.Sprintf("countFiles:%d", guestId))

	var count int64
	err := r.cache.Get(key, &count)
	if err == nil {
		return count, nil
	}

	count, err = r.repo.GetCountFiles(guestId)
	if err != nil {
		return 0, err
	}

	if err := r.cache.Set(key, count, []string{
		r.TagGuest(guestId),
		r.TagAll(),
	}); err != nil {
		log.Printf("failed to set cache: %v", err)
	}

	return count, nil
}

func (r *CachedFileRepository) GetFile(guestId uint, fileId uint) (entities.File, error) {
	key := r.Key(fmt.Sprintf("file:guest:%d:file:%d", guestId, fileId))

	var file entities.File
	err := r.cache.Get(key, &file)
	if err == nil && file.ID != 0 {
		return file, nil
	}

	file, err = r.repo.GetFile(guestId, fileId)
	if err != nil {
		return file, err
	}

	if file.ID != 0 {
		if err := r.cache.Set(key, file, []string{
			r.TagGuest(guestId),
			r.TagAll(),
		}); err != nil {
			log.Printf("failed to set cache: %v", err)
		}
	}

	return file, nil
}

func (r *CachedFileRepository) GetFileById(fileId uint) (entities.File, error) {
	key := r.Key(fmt.Sprintf("fileById:%d", fileId))

	var file entities.File
	err := r.cache.Get(key, &file)
	if err == nil && file.ID != 0 {
		return file, nil
	}

	file, err = r.repo.GetFileById(fileId)
	if err != nil {
		return file, err
	}

	if file.ID != 0 {
		if err := r.cache.Set(key, file, []string{
			r.TagAll(),
		}); err != nil {
			log.Printf("failed to set cache: %v", err)
		}
	}

	return file, nil
}

func (r CachedFileRepository) Key(k string) string {
	return fmt.Sprintf("file:repo:%s", k)
}

func (r CachedFileRepository) TagAll() string {
	return r.Key("all")
}

func (r CachedFileRepository) TagGuest(guestId uint) string {
	return r.Key(fmt.Sprintf("guest:%d", guestId))
}
