package repositories

import (
	"converter/components/cache"
	"converter/entities"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type CachedFileRepository struct {
	FileDbRepository
	cache cache.Cache
}

func NewCachedFileRepository(db *gorm.DB) (*CachedFileRepository, error) {
	cache, err := cache.CachedFactory{}.Create()
	if err != nil {
		return nil, err
	}
	return &CachedFileRepository{
		FileDbRepository: *NewFileRepository(db),
		cache:            cache,
	}, nil
}

func (r *CachedFileRepository) CloseRepo() error {
	return r.cache.Close()
}

func (r *CachedFileRepository) GetFiles(guestId uint) ([]entities.File, error) {
	key := r.Key(fmt.Sprintf("files:%d", guestId))
	tags := []string{r.TagGuest(guestId), r.TagAll()}
	return cacheOrFetch(r, key, tags, func() ([]entities.File, error) {
		return r.FileDbRepository.GetFiles(guestId)
	}, nil)
}

func (r *CachedFileRepository) GetCountFiles(guestId uint) (int64, error) {
	key := r.Key(fmt.Sprintf("countFiles:%d", guestId))
	tags := []string{r.TagGuest(guestId), r.TagAll()}
	return cacheOrFetch(r, key, tags, func() (int64, error) {
		return r.FileDbRepository.GetCountFiles(guestId)
	}, nil)
}

func (r *CachedFileRepository) GetFile(guestId uint, fileId uint) (entities.File, error) {
	key := r.Key(fmt.Sprintf("file:guest:%d:file:%d", guestId, fileId))
	tags := []string{r.TagGuest(guestId), r.TagAll()}
	isValid := func(f entities.File) bool { return f.ID != 0 }
	return cacheOrFetch(r, key, tags, func() (entities.File, error) {
		return r.FileDbRepository.GetFile(guestId, fileId)
	}, isValid)
}

func (r *CachedFileRepository) GetFileById(fileId uint) (entities.File, error) {
	key := r.Key(fmt.Sprintf("fileById:%d", fileId))
	tags := []string{r.TagAll()}
	isValid := func(f entities.File) bool { return f.ID != 0 }
	return cacheOrFetch(r, key, tags, func() (entities.File, error) {
		return r.FileDbRepository.GetFileById(fileId)
	}, isValid)
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

func cacheOrFetch[T any](
	r *CachedFileRepository,
	key string,
	tags []string,
	fetch func() (T, error),
	isValid func(T) bool, // необязательная проверка валидности результата
) (T, error) {
	var result T
	err := r.cache.Get(key, &result)
	if err == nil {
		if isValid == nil || isValid(result) {
			return result, nil
		}
	}

	result, err = fetch()
	if err != nil {
		return result, err
	}

	if isValid != nil && !isValid(result) {
		return result, nil
	}

	if setErr := r.cache.Set(key, result, tags); setErr != nil {
		log.Printf("failed to set cache: %v", setErr)
	}
	return result, nil
}
