package helpers

import (
	"os"
	"runtime/debug"
)

func IsRaceEnabled() bool {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return false
	}

	for _, setting := range info.Settings {
		if setting.Key == "-race" && setting.Value == "true" {
			return true
		}
	}
	return false
}

func Env(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
