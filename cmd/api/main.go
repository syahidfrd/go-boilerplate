package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/syahidfrd/go-boilerplate/config"
	httpDelivery "github.com/syahidfrd/go-boilerplate/delivery/http"
	"github.com/syahidfrd/go-boilerplate/infrastructure/datastore"
	appMiddleware "github.com/syahidfrd/go-boilerplate/middleware"
	pgsqlRepository "github.com/syahidfrd/go-boilerplate/repository/pgsql"
	redisRepository "github.com/syahidfrd/go-boilerplate/repository/redis"
	"github.com/syahidfrd/go-boilerplate/usecase"
	"github.com/syahidfrd/go-boilerplate/utils/logger"
)

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
	authorUC := usecase.NewAuthorUsecase(authorRepo, redisRepo)

	// Setup middleware manager
	middManager := appMiddleware.NewMiddlewareManager(appLogger)

	// Setup route engine & middleware
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middManager.GenerateCID())
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		middManager.InboundLog(c, reqBody, resBody)
	}))

	// Setup handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "i am alive")
	})

	httpDelivery.NewAuthorHandler(e, middManager, authorUC)

	e.Logger.Fatal(e.Start(":" + configApp.ServerPORT))
}
