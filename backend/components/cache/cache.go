package cache

import (
	"converter/config"
)

type Cache interface {
	Set(key, object any, tags []string) error
	Get(key, object any) error
	Delete(key any) error
	DeleteByTag(tags []string) error
	Close() error
}

type CachedFactory struct {
}

func (f CachedFactory) Create() (Cache, error) {
	cfg := config.GetConfig()
	return NewRedisCache(cfg.RedisConfig.RedisAddr, cfg.RedisConfig.RedisPassword, config.CacheTTL)
}
