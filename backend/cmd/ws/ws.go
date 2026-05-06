package main

import (
	"converter/app"
	"converter/services/user"
	"converter/ws"
	"github.com/gin-contrib/cors"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"log"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	app.Init(false)
	defer app.App().DeInit()

	r := gin.Default()
	initMiddleware(r)

	hub := ws.NewHub()
	go hub.Run()
	go ws.StartNotificationsWatcher(app.App().DB, hub)

	r.GET("/ws", func(c *gin.Context) {
		userService := user.NewSessionUserService(sessions.Default(c))
		guestID, err := userService.GuestID()
		if err != nil {
			app.Fail(c, 500, "1", err.Error())
			return
		}

		ws.ServeWs(hub, guestID, c.Writer, c.Request)
	})

	if err := r.Run(); err != nil {
		log.Fatal("Websocket server error:", err)
	}
}

func initMiddleware(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(sessions.Sessions(app.App().Config.SessionConfig.SessionName, *app.App().SessionStore))

	r.Use(adapter.Wrap(csrf.Protect(
		app.App().Config.BaseConfig.CsrfKey,
		csrf.Secure(true),
		csrf.Path("/"),
	)))
}
