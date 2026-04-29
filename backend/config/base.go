package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const UploadDir = "/home/appuser/files"
const ConvertedDir = "/home/appuser/files_converted"
const PrefixGuestDir = "guest_"
const SessionDuration = 86400

// const SessionDuration = 5

const CleanupInterval = 300

// const CleanupInterval = 5

const MaxSize = 300 * 1024 * 1024
const MaxSizeFile = 250 * 1024 * 1024
const CacheTTL = 10 * time.Minute

var (
	instance *Config
	once     sync.Once
)

type BaseConfig struct {
	CsrfKey []byte
	DbUrl   string
}

type Config struct {
	BaseConfig    *BaseConfig
	SessionConfig *SessionConfig
	RedisConfig   *RedisConfig
}

func GetConfig() *Config {
	if instance == nil {
		panic("Config not initialized")
	}
	return instance
}

func Init(isConsole bool) {
	once.Do(func() {
		baseCfg, err := getBaseConfig(isConsole)
		if err != nil {
			log.Fatal("Config error:", err)
		}

		var sessionCfg *SessionConfig
		if !isConsole {
			sessionCfg, err = getSessionConfig()
			if err != nil {
				log.Fatal("Config error:", err)
			}
		}
		redisCfg, err := getRedisConfig()
		if err != nil {
			log.Fatal("Config error:", err)
		}
		instance = &Config{
			BaseConfig:    baseCfg,
			SessionConfig: sessionCfg,
			RedisConfig:   redisCfg,
		}
	})
}

func getBaseConfig(isConsole bool) (*BaseConfig, error) {
	var csrfKey []byte
	if !isConsole {
		var err error
		csrfKey, err = base64.StdEncoding.DecodeString(os.Getenv("CSRF_KEY"))
		if err != nil {
			return nil, fmt.Errorf("invalid CSRF_KEY base64: %w", err)
		}

		if len(csrfKey) != 32 {
			return nil, fmt.Errorf("CSRF_KEY decoded length %d, expected 32", len(csrfKey))
		}
	}

	return &BaseConfig{
		CsrfKey: csrfKey,
		DbUrl:   os.Getenv("DB_URL"),
	}, nil
}
