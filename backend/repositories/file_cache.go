package repositories

import (
	"context"
	"converter/entities"
	"fmt"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	redisStore "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"time"
)

type CachedFileRepository struct {
	repo        *FileRepository
	marsh       *marshaler.Marshaler
	cache       *cache.Cache[any]
	ttl         time.Duration
	redisClient *redis.Client
}

func NewCachedFileRepository(db *gorm.DB, redisAddr string, redisPassword string, ttl time.Duration) (*CachedFileRepository, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       1,
	})

	redisStore := redisStore.NewRedis(redisClient, store.WithExpiration(ttl))
	cacheManager := cache.New[any](redisStore)
	marsh := marshaler.New(cacheManager)

	return &CachedFileRepository{
		repo:        NewFileRepository(db),
		marsh:       marsh,
		cache:       cacheManager,
		ttl:         ttl,
		redisClient: redisClient,
	}, nil
}

func (r *CachedFileRepository) CloseRepo() error {
	err := r.redisClient.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *CachedFileRepository) GetFiles(guestId uint) ([]entities.File, error) {
	key := r.key(fmt.Sprintf("files:%d", guestId))
	ctx := context.Background()

	var files []entities.File
	_, err := r.marsh.Get(ctx, key, &files)
	if err == nil {
		return files, nil
	}

	files, err = r.repo.GetFiles(guestId)
	if err != nil {
		return nil, err
	}

	if err := r.marsh.Set(ctx, key, files, store.WithTags([]string{
		r.tagGuest(guestId),
		r.tagAll(),
	})); err != nil {
		log.Printf("failed to set cache: %v", err)
	}

	return files, nil
}

func (r *CachedFileRepository) GetCountFiles(guestId uint) (int64, error) {
	key := r.key(fmt.Sprintf("countFiles:%d", guestId))
	ctx := context.Background()

	var count int64
	_, err := r.marsh.Get(ctx, key, &count)
	if err == nil {
		return count, nil
	}

	count, err = r.repo.GetCountFiles(guestId)
	if err != nil {
		return 0, err
	}

	if err := r.marsh.Set(ctx, key, count, store.WithTags([]string{
		r.tagGuest(guestId),
		r.tagAll(),
	})); err != nil {
		log.Printf("failed to set cache: %v", err)
	}

	return count, nil
}

func (r *CachedFileRepository) GetFile(guestId uint, fileId uint) (entities.File, error) {
	key := r.key(fmt.Sprintf("file:guest:%d:file:%d", guestId, fileId))
	ctx := context.Background()

	var file entities.File
	_, err := r.marsh.Get(ctx, key, &file)
	if err == nil && file.ID != 0 {
		return file, nil
	}

	file, err = r.repo.GetFile(guestId, fileId)
	if err != nil {
		return file, err
	}

	if file.ID != 0 {
		if err := r.marsh.Set(ctx, key, file, store.WithTags([]string{
			r.tagGuest(guestId),
			r.tagAll(),
		})); err != nil {
			log.Printf("failed to set cache: %v", err)
		}
	}

	return file, nil
}

func (r *CachedFileRepository) key(k string) string {
	return fmt.Sprintf("file:repo:%s", k)
}

func (r *CachedFileRepository) tagAll() string {
	return r.key("all")
}

func (r *CachedFileRepository) tagGuest(guestId uint) string {
	return r.key(fmt.Sprintf("guest:%d", guestId))
}

func (r *CachedFileRepository) tagFile(fileId uint) string {
	return r.key(fmt.Sprintf("file:%d", fileId))
}

func (r *CachedFileRepository) InvalidateGuest(guestId uint) error {
	ctx := context.Background()

	return r.marsh.Invalidate(ctx, store.WithInvalidateTags([]string{
		r.tagGuest(guestId),
	}))
}

func (r *CachedFileRepository) InvalidateFile(fileId uint) error {
	ctx := context.Background()

	return r.marsh.Invalidate(ctx, store.WithInvalidateTags([]string{
		r.tagFile(fileId),
	}))
}

func (r *CachedFileRepository) InvalidateAll() error {
	ctx := context.Background()
	return r.marsh.Invalidate(ctx, store.WithInvalidateTags([]string{
		r.tagAll(),
	}))
}
