package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/syahidfrd/go-boilerplate/internal/auth"
	"github.com/syahidfrd/go-boilerplate/internal/health"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/cache"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/config"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/db"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/jwt"
	"github.com/syahidfrd/go-boilerplate/internal/todo"
	"github.com/syahidfrd/go-boilerplate/internal/user"
)

// Server represents the HTTP server with its router
type Server struct {
	router *http.ServeMux
}

// NewServer creates and configures a new HTTP server with all dependencies and routes
func NewServer() *Server {
	// Load configuration
	cfg := config.LoadEnv()

	// Initialize database
	dbConn, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// Auto migrate models
	if err := db.AutoMigrate(dbConn, &user.User{}, &user.Preference{}, &todo.Todo{}); err != nil {
		log.Fatal().Err(err).Msg("failed to auto migrate database")
	}

	// Initialize Redis cache
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.CacheURL,
	})
	redisCache := cache.NewRedis(redisClient)

	// Initialize services
	userStore := user.NewStore(dbConn)
	userService := user.NewService(userStore)
	jwtService := jwt.NewService(cfg.AppSecret)
	authService := auth.NewService(userService, jwtService)

	todoStore := todo.NewStore(dbConn)
	todoService := todo.NewService(todoStore, redisCache)

	healthStore := health.NewStore(dbConn, redisClient)
	healthService := health.NewService(healthStore)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	todoHandler := todo.NewHandler(todoService)
	healthHandler := health.NewHandler(healthService)

	// Initialize middleware
	jwtMiddleware := auth.NewJWTMiddleware(jwtService)

	// Configure HTTP routes
	r := http.NewServeMux()

	// Public routes
	r.Handle("GET /", http.HandlerFunc(rootHandler))
	r.Handle("GET /health", http.HandlerFunc(healthHandler.Health))

	// Auth routes
	r.Handle("POST /api/auth/signup", http.HandlerFunc(authHandler.SignUp))
	r.Handle("POST /api/auth/signin", http.HandlerFunc(authHandler.SignIn))

	// Todo routes (protected)
	r.Handle("POST /api/todos", jwtMiddleware.Authenticate(http.HandlerFunc(todoHandler.Create)))
	r.Handle("GET /api/todos", jwtMiddleware.Authenticate(http.HandlerFunc(todoHandler.GetByUserID)))
	r.Handle("GET /api/todos/{id}", jwtMiddleware.Authenticate(http.HandlerFunc(todoHandler.GetByID)))
	r.Handle("PUT /api/todos/{id}", jwtMiddleware.Authenticate(http.HandlerFunc(todoHandler.Update)))
	r.Handle("PATCH /api/todos/{id}/toggle", jwtMiddleware.Authenticate(http.HandlerFunc(todoHandler.ToggleComplete)))
	r.Handle("DELETE /api/todos/{id}", jwtMiddleware.Authenticate(http.HandlerFunc(todoHandler.Delete)))

	return &Server{
		router: r,
	}
}

// rootHandler handles the root endpoint "/"
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// Run starts the HTTP server on the specified port with graceful shutdown handling
func (s *Server) Run(port int) {
	addr := fmt.Sprintf(":%d", port)

	// Apply middleware chain for logging, recovery, CORS, etc.
	h := chainMiddleware(
		s.router,
		recoverMiddleware,
		loggerMiddleware(func(w http.ResponseWriter, r *http.Request) bool { return r.URL.Path == "/health" }),
		realIPMiddleware,
		requestIDMiddleware,
		corsMiddleware,
	)

	// Configure HTTP server with timeouts
	httpSrv := http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	// Setup graceful shutdown channels
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Graceful shutdown goroutine
	go func() {
		<-quit
		log.Info().Msg("server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpSrv.SetKeepAlivesEnabled(false)
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Msg("could not gracefully shutdown the server")
		}
		close(done)
	}()

	// Start the server
	log.Info().Msgf("server serving on port %d", port)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msgf("could not listen on %s", addr)
	}

	// Wait for shutdown to complete
	<-done
	log.Info().Msg("server stopped")
}
