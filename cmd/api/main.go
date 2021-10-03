package main

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/syahidfrd/go-boilerplate/config"
	httpDelivery "github.com/syahidfrd/go-boilerplate/delivery/http"
	appMiddleware "github.com/syahidfrd/go-boilerplate/delivery/http/middleware"
	"github.com/syahidfrd/go-boilerplate/domain"
	"github.com/syahidfrd/go-boilerplate/infrastructure/datastore"
	"github.com/syahidfrd/go-boilerplate/repository"
	"github.com/syahidfrd/go-boilerplate/usecase"
)

func main() {

	var (
		configApp        *config.Config          = config.LoadConfig()
		dbInstance       *sql.DB                 = datastore.NewDatabase(configApp.DatabaseURL)
		authorRepository domain.AuthorRepository = repository.NewAuthorRepository(dbInstance)
		authorUsecase    domain.AuthorUsecase    = usecase.NewAuthorUsecase(authorRepository)
	)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(appMiddleware.GenerateCorrelationID())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "i am alive")
	})

	api := e.Group("/api/v1")
	httpDelivery.NewAuthorHandler(api, authorUsecase)

	e.Logger.Fatal(e.Start(":" + configApp.ServerPORT))
}
