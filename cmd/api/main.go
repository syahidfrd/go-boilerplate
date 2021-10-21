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
)

func main() {

	// Load config
	configApp := config.LoadConfig()

	// Setup infra
	dbInstance := datastore.NewDatabase(configApp.DatabaseURL)
	cacheInstance := datastore.NewCache(configApp.CacheURL)

	// Setup repository
	redisRepository := redisRepository.NewRedisRepository(cacheInstance)
	authorRepository := pgsqlRepository.NewPgsqlAuthorRepository(dbInstance)

	// Setup usecase
	authorUsecase := usecase.NewAuthorUsecase(authorRepository, redisRepository)

	// Setup middleware manager
	middManager := appMiddleware.NewMiddlewareManager()

	// Setup route engine & middleware
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middManager.GenerateCorrelationID())

	// Setup handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "i am alive")
	})

	api := e.Group("/api/v1")
	httpDelivery.NewAuthorHandler(api, authorUsecase)

	e.Logger.Fatal(e.Start(":" + configApp.ServerPORT))
}
