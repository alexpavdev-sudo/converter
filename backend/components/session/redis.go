package session

import (
	"converter/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"log"
	"net/http"
)

func NewRedisStore(redisCfg *config.RedisConfig, sessionCfg *config.SessionConfig) *sessions.Store {
	store, err := redis.NewStoreWithDB(
		10,
		"tcp",
		redisCfg.RedisAddr,
		"",
		redisCfg.RedisPassword,
		"0",
		sessionCfg.AuthKey,
		sessionCfg.EncKey,
	)
	if err != nil {
		log.Fatal("Redis store error:", err)
	}

	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   config.SessionDuration,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	storeSession, _ := store.(sessions.Store)
	return &storeSession
}
