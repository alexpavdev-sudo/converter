package main

import (
	"context"
	"converter/app"
	"converter/controllers"
	"converter/helpers"
	"converter/services/cleanup"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"log"
	"time"
)

func init() {
	if helpers.IsRaceEnabled() {
		fmt.Println("[WARN] Детектор гонок ВКЛЮЧЕН")
	}
}

func main() {
	app.Init(false)
	defer app.App().DeInit()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cleanup.Start(ctx)

	r := gin.Default()
	initMiddleware(r)
	initHandlers(r)

	if err := r.Run(); err != nil {
		log.Fatal("Server error:", err)
	}
}

func initHandlers(r *gin.Engine) {
	filesGroup := r.Group("/files")
	{
		filesGroup.POST("/upload", controllers.Upload)
		filesGroup.GET("/", controllers.GetFiles)
		filesGroup.GET("/:id", controllers.GetFile)
		filesGroup.GET("/error/:id", controllers.GetFileError)
		filesGroup.DELETE("/:id", controllers.DeleteFile)
		filesGroup.GET("/download/:id", controllers.DownloadFile)
	}
	formatsGroup := r.Group("/formats")
	{
		formatsGroup.GET("/", controllers.GetFormats)
	}
	r.GET("/user/profile", controllers.NewUserHandler().GetProfile)
	r.GET("/ping", controllers.Ping)
}

func initMiddleware(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(sessions.Sessions(app.App().Config.SessionConfig.SessionName, *app.App().SessionStore))

	r.Use(func(c *gin.Context) {
		patchedReq := csrf.PlaintextHTTPRequest(c.Request)
		c.Request = patchedReq
		c.Next()
	})

	r.Use(adapter.Wrap(csrf.Protect(
		app.App().Config.BaseConfig.CsrfKey,
		csrf.Secure(false),
		csrf.Path("/"),
	)))
}
