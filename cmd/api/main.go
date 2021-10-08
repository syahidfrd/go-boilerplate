package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/syahidfrd/go-boilerplate/config"
	httpDelivery "github.com/syahidfrd/go-boilerplate/delivery/http"
	httpDeliveryMiddleware "github.com/syahidfrd/go-boilerplate/delivery/http/middleware"
	"github.com/syahidfrd/go-boilerplate/infrastructure/datastore"
	pgsqlRepository "github.com/syahidfrd/go-boilerplate/repository/pgsql"
	"github.com/syahidfrd/go-boilerplate/repository/redis"
	"github.com/syahidfrd/go-boilerplate/usecase"
)

func main() {

	// Load config
	configApp := config.LoadConfig()

	// Setup infra
	dbInstance := datastore.NewDatabase(configApp.DatabaseURL)
	cacheInstance := datastore.NewCache(configApp.CacheURL)

	// Setup repository
	redisRepository := redis.NewRedisRepository(cacheInstance)
	authorRepository := pgsqlRepository.NewPgsqlAuthorRepository(dbInstance)

	// Setup usecase
	authorUsecase := usecase.NewAuthorUsecase(authorRepository, redisRepository)

	// Setup route engine & middleware
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(httpDeliveryMiddleware.GenerateCorrelationID())

	// Setup handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "i am alive")
	})

	api := e.Group("/api/v1")
	httpDelivery.NewAuthorHandler(api, authorUsecase)

	e.Logger.Fatal(e.Start(":" + configApp.ServerPORT))
}
