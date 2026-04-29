package config

import (
	"os"
)

type RedisConfig struct {
	RedisAddr     string
	RedisPassword string
}

func getRedisConfig() (*RedisConfig, error) {
	return &RedisConfig{
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}, nil
}
