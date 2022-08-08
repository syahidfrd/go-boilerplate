package main

import (
	"net/http"
	"time"

	_ "github.com/syahidfrd/go-boilerplate/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/syahidfrd/go-boilerplate/config"
	httpDelivery "github.com/syahidfrd/go-boilerplate/delivery/http"
	"github.com/syahidfrd/go-boilerplate/infrastructure/datastore"
	appMiddleware "github.com/syahidfrd/go-boilerplate/middleware"
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
	dbInstance := datastore.NewDatabase(configApp.DatabaseURL)
	cacheInstance := datastore.NewCache(configApp.CacheURL)

	// Setup repository
	redisRepo := redisRepository.NewRedisRepository(cacheInstance)
	authorRepo := pgsqlRepository.NewPgsqlAuthorRepository(dbInstance)

	// Setup usecase
	ctxTimeout := time.Duration(configApp.ContextTimeout) * time.Second
	authorUC := usecase.NewAuthorUsecase(authorRepo, redisRepo, ctxTimeout)

	// Setup middleware manager
	middManager := appMiddleware.NewMiddlewareManager(appLogger)

	// Setup route engine & middleware
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middManager.GenerateCID())
	e.Use(middManager.InboundLog)
	e.Use(middleware.Recover())

	// Setup handler
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "i am alive")
	})

	httpDelivery.NewAuthorHandler(e, middManager, authorUC)

	e.Logger.Fatal(e.Start(":" + configApp.ServerPORT))
}
