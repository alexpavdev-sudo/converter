package config

import (
	"encoding/base64"
	"fmt"
	"os"
)

type SessionKeys struct {
	AuthKey []byte
	EncKey  []byte
}

type Config struct {
	RedisAddr     string
	RedisPassword string
	SessionKeys   SessionKeys
}

func Load() (*Config, error) {
	// Декодируем ключи
	authKey, err := base64.StdEncoding.DecodeString(os.Getenv("SESSION_AUTH_KEY"))
	if err != nil {
		return nil, fmt.Errorf("invalid SESSION_AUTH_KEY base64: %w", err)
	}
	encKey, err := base64.StdEncoding.DecodeString(os.Getenv("SESSION_ENC_KEY"))
	if err != nil {
		return nil, fmt.Errorf("invalid SESSION_ENC_KEY base64: %w", err)
	}

	// Проверяем длину
	if len(authKey) != 64 {
		return nil, fmt.Errorf("SESSION_AUTH_KEY decoded length %d, expected 64", len(authKey))
	}
	if len(encKey) != 32 {
		return nil, fmt.Errorf("SESSION_ENC_KEY decoded length %d, expected 32", len(encKey))
	}

	return &Config{
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		SessionKeys: SessionKeys{
			AuthKey: authKey,
			EncKey:  encKey,
		},
	}, nil
}
