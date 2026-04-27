package helpers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"path/filepath"
)

func GenerateRandomStoredName(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("ошибка генерации имени файла")
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func GetFileNameWithoutExt(filePath string) string {
	filename := filepath.Base(filePath)
	ext := filepath.Ext(filename)
	return filename[:len(filename)-len(ext)]
}
