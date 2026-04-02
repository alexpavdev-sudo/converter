package session

import (
	"converter/internal/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"net/http"
)

// NewRedisStore создаёт хранилище сессий на основе Redis.
func NewRedisStore(cfg *config.Config) (sessions.Store, error) {
	store, err := redis.NewStore(
		10,                      // размер пула соединений
		"tcp",                   // сеть
		cfg.RedisAddr,           // адрес Redis
		"",                      // опционально: ключ для кластеризации (не используется)
		cfg.RedisPassword,       // пароль
		cfg.SessionKeys.AuthKey, // ключ аутентификации
		cfg.SessionKeys.EncKey,  // ключ шифрования
	)
	if err != nil {
		return nil, err
	}

	// Общие опции для всех сессий
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return store, nil
}
