package main

import (
	"net/http"
	"time"

	_ "github.com/syahidfrd/go-boilerplate/docs"
	"github.com/syahidfrd/go-boilerplate/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/syahidfrd/go-boilerplate/config"
	httpDelivery "github.com/syahidfrd/go-boilerplate/delivery/http"
	appMiddleware "github.com/syahidfrd/go-boilerplate/delivery/middleware"
	"github.com/syahidfrd/go-boilerplate/infrastructure/datastore"
	pgsqlRepository "github.com/syahidfrd/go-boilerplate/repository/pgsql"
	redisRepository "github.com/syahidfrd/go-boilerplate/repository/redis"
	"github.com/syahidfrd/go-boilerplate/usecase"
	"github.com/syahidfrd/go-boilerplate/utils/logger"
)

// @title Go Boilerplate
// @version 1.0.4
// @termsOfService http://swagger.io/terms/
func main() {
	// Load config
	configApp := config.LoadConfig()

	// Setup logger
	appLogger := logger.NewApiLogger(configApp)
	appLogger.InitLogger()

	// Setup infra
	dbInstance, err := datastore.NewDatabase(configApp.DatabaseURL)
	utils.PanicIfNeeded(err)

	cacheInstance, err := datastore.NewCache(configApp.CacheURL)
	utils.PanicIfNeeded(err)

	// Setup repository
	redisRepo := redisRepository.NewRedisRepository(cacheInstance)
	authorRepo := pgsqlRepository.NewPgsqlAuthorRepository(dbInstance)

	// Setup usecase
	ctxTimeout := time.Duration(configApp.ContextTimeout) * time.Second
	authorUC := usecase.NewAuthorUsecase(authorRepo, redisRepo, ctxTimeout)

	// Setup app middleware
	appMiddleware := appMiddleware.NewMiddleware(appLogger)

	// Setup route engine & middleware
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(appMiddleware.RequestID())
	e.Use(appMiddleware.Logger())
	e.Use(middleware.Recover())

	// Setup handler
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "i am alive")
	})

	httpDelivery.NewAuthorHandler(e, appMiddleware, authorUC)

	e.Logger.Fatal(e.Start(":" + configApp.ServerPORT))
}
