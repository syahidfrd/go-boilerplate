package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/syahidfrd/go-boilerplate/utils/crypto"
	"github.com/syahidfrd/go-boilerplate/utils/jwt"

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
)

// @title Go Boilerplate
// @version 1.0.4
// @termsOfService http://swagger.io/terms/
// @securityDefinitions.apikey JwtToken
// @in header
// @name Authorization
func main() {
	// Load config
	configApp := config.LoadConfig()

	// Setup infra
	dbInstance, err := datastore.NewDatabase(configApp.DatabaseURL)
	utils.PanicIfNeeded(err)

	cacheInstance, err := datastore.NewCache(configApp.CacheURL)
	utils.PanicIfNeeded(err)

	// Setup repository
	redisRepo := redisRepository.NewRedisRepository(cacheInstance)
	todoRepo := pgsqlRepository.NewPgsqlTodoRepository(dbInstance)
	userRepo := pgsqlRepository.NewPgsqlUserRepository(dbInstance)

	// Setup Service
	cryptoSvc := crypto.NewCryptoService()
	jwtSvc := jwt.NewJWTService(configApp.JWTSecretKey)

	// Setup usecase
	ctxTimeout := time.Duration(configApp.ContextTimeout) * time.Second
	todoUC := usecase.NewTodoUsecase(todoRepo, redisRepo, ctxTimeout)
	authUC := usecase.NewAuthUsecase(userRepo, cryptoSvc, jwtSvc, ctxTimeout)

	// Setup app middleware
	appMiddleware := appMiddleware.NewMiddleware(jwtSvc)

	// Setup route engine & middleware
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(appMiddleware.Logger(nil))

	// Setup handler
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "i am alive")
	})

	httpDelivery.NewTodoHandler(e, appMiddleware, todoUC)
	httpDelivery.NewAuthHandler(e, appMiddleware, authUC)

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configApp.ContextTimeout)*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
