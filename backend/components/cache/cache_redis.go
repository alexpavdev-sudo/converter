package cache

import (
	"context"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	redisStore "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	marsh       *marshaler.Marshaler
	cache       *cache.Cache[any]
	ttl         time.Duration
	redisClient *redis.Client
}

func NewRedisCache(redisAddr string, redisPassword string, ttl time.Duration) (*RedisCache, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       1,
	})

	redisStore := redisStore.NewRedis(redisClient, store.WithExpiration(ttl))
	cacheManager := cache.New[any](redisStore)
	marsh := marshaler.New(cacheManager)

	return &RedisCache{
		marsh:       marsh,
		cache:       cacheManager,
		ttl:         ttl,
		redisClient: redisClient,
	}, nil
}

func (c *RedisCache) Set(key, object any, tags []string) error {
	ctx := context.Background()
	if err := c.marsh.Set(ctx, key, object, store.WithTags(tags)); err != nil {
		return err
	}
	return nil
}

func (c *RedisCache) Get(key, object any) error {
	ctx := context.Background()
	_, err := c.marsh.Get(ctx, key, object)
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisCache) Delete(key any) error {
	ctx := context.Background()
	return c.marsh.Delete(ctx, key)
}

func (c *RedisCache) DeleteByTag(tags []string) error {
	ctx := context.Background()
	return c.marsh.Invalidate(ctx, store.WithInvalidateTags(tags))
}

func (c *RedisCache) Close() error {
	err := c.redisClient.Close()
	if err != nil {
		return err
	}
	return nil
}
