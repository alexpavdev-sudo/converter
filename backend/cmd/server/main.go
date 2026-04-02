package main

import (
	"converter/internal/config"
	"converter/internal/handlers"
	"converter/internal/session"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 1. Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config error:", err)
	}

	// 2. Создаём хранилище сессий
	store, err := session.NewRedisStore(cfg)
	if err != nil {
		log.Fatal("Redis store error:", err)
	}

	// 3. Настраиваем Gin
	r := gin.Default()
	r.Use(sessions.Sessions("mysession", store))

	// 4. Регистрируем маршруты
	r.GET("/ping", handlers.Ping)

	// 5. Запуск
	if err := r.Run(); err != nil {
		log.Fatal("Server error:", err)
	}
}
