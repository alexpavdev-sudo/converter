package config

import (
	"encoding/base64"
	"fmt"
	"os"
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

type BaseConfig struct {
	CsrfKey []byte
	DbUrl   string
}

func GetBaseConfig(isConsole bool) (*BaseConfig, error) {
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
