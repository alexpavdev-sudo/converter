package app

import (
	"converter/components/session"
	"converter/config"
	"converter/dto/web"
	"converter/repositories"
	"database/sql"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	instance *Application
	once     sync.Once
)

type Application struct {
	isConsole    bool
	Config       *config.Config
	SessionStore *sessions.Store
	DB           *gorm.DB
	FileRepo     repositories.FileRepository
}

func App() *Application {
	if instance == nil {
		panic("App not initialized")
	}
	return instance
}

func Init(isConsole bool) {
	once.Do(func() {
		config.Init(isConsole)
		cfg := config.GetConfig()
		db := openDB(cfg.BaseConfig.DbUrl)

		var sessionStore *sessions.Store
		if !isConsole {
			sessionStore = session.NewRedisStore(cfg.RedisConfig, cfg.SessionConfig)
		}

		instance = &Application{
			Config:       cfg,
			SessionStore: sessionStore,
			DB:           db,
			FileRepo:     getFileRepo(db),
			isConsole:    isConsole,
		}
	})
}

func getCacheFileRepo(db *gorm.DB) repositories.FileRepository {
	repo, err := repositories.NewCachedFileRepository(db)
	if err != nil {
		log.Fatal("Config error:", err)
	}
	return repo
}

func getFileRepo(db *gorm.DB) repositories.FileRepository {
	return repositories.NewFileRepository(db)
}

func openDB(DbUrl string) *gorm.DB {
	var err error
	DB, err := gorm.Open(postgres.Open(DbUrl), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				ParameterizedQueries:      true,
				Colorful:                  true,
			},
		),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected successfully")
	return DB
}

func (app *Application) closeDB() error {
	sqlDB, err := app.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (app *Application) DeInit() {
	app.FileRepo.CloseRepo()
	app.closeDB()
}

func (app *Application) StartTransaction() *gorm.DB {
	return app.DB.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, web.ResponseDto{
		Success: true,
		Data:    data,
	})
}

func Fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, web.ResponseDto{
		Success: false,
		Error:   &web.ErrorInfoDto{Code: code, Message: message},
	})
}
