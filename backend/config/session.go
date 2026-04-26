package config

import (
	"encoding/base64"
	"fmt"
	"os"
)

const (
	sessionName string = "sess_id"
)

type SessionConfig struct {
	AuthKey     []byte
	EncKey      []byte
	SessionName string
}

func GetSessionConfig() (*SessionConfig, error) {
	authKey, err := base64.StdEncoding.DecodeString(os.Getenv("SESSION_AUTH_KEY"))
	if err != nil {
		return nil, fmt.Errorf("invalid SESSION_AUTH_KEY base64: %w", err)
	}
	encKey, err := base64.StdEncoding.DecodeString(os.Getenv("SESSION_ENC_KEY"))
	if err != nil {
		return nil, fmt.Errorf("invalid SESSION_ENC_KEY base64: %w", err)
	}

	if len(authKey) != 64 {
		return nil, fmt.Errorf("SESSION_AUTH_KEY decoded length %d, expected 64", len(authKey))
	}
	if len(encKey) != 32 {
		return nil, fmt.Errorf("SESSION_ENC_KEY decoded length %d, expected 32", len(encKey))
	}

	return &SessionConfig{
		AuthKey:     authKey,
		EncKey:      encKey,
		SessionName: sessionName,
	}, nil
}
