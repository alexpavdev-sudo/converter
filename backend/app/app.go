package app

import (
	"converter/components/session"
	"converter/config"
	"converter/dto"
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

type Config struct {
	BaseConfig    *config.BaseConfig
	SessionConfig *config.SessionConfig
	RedisConfig   *config.RedisConfig
}

type Application struct {
	isConsole    bool
	Config       Config
	SessionStore *sessions.Store
	DB           *gorm.DB
	FileRepo     repositories.FileRepositoryInterface
}

func App() *Application {
	if instance == nil {
		panic("App not initialized")
	}
	return instance
}

func Init(isConsole bool) {
	once.Do(func() {
		cfg := getConfig(isConsole)
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

func ClearCache(guestID uint) {
	//todo
	if repo, ok := App().FileRepo.(*repositories.CachedFileRepository); ok {
		err := repo.InvalidateGuest(guestID)
		if err != nil {
			log.Printf("failed to invalidate cache guest: %v", err)
		}
	}
}

func getCacheFileRepo(db *gorm.DB, cfg *config.RedisConfig) repositories.FileRepositoryInterface {
	repo, err := repositories.NewCachedFileRepository(db, cfg.RedisAddr, cfg.RedisAddr, 10*time.Minute)
	if err != nil {
		log.Fatal("Config error:", err)
	}
	return repo
}

func getFileRepo(db *gorm.DB) repositories.FileRepositoryInterface {
	return repositories.NewFileRepository(db)
}

func getConfig(isConsole bool) Config {
	baseCfg, err := config.GetBaseConfig(isConsole)
	if err != nil {
		log.Fatal("Config error:", err)
	}

	var sessionCfg *config.SessionConfig
	if !isConsole {
		sessionCfg, err = config.GetSessionConfig()
		if err != nil {
			log.Fatal("Config error:", err)
		}
	}
	redisCfg, err := config.GetRedisConfig()
	if err != nil {
		log.Fatal("Config error:", err)
	}
	return Config{
		BaseConfig:    baseCfg,
		SessionConfig: sessionCfg,
		RedisConfig:   redisCfg,
	}
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
	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    data,
	})
}

func Fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, dto.Response{
		Success: false,
		Error:   &dto.ErrorInfo{Code: code, Message: message},
	})
}
