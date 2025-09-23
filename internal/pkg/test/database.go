package test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/config"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/db"
	"github.com/syahidfrd/go-boilerplate/internal/user"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

// Container wraps testcontainers with database and cache connections
type Container struct {
	PostgresContainer *postgres.PostgresContainer
	RedisContainer    *redisContainer.RedisContainer
	DB                *gorm.DB
	Redis             *redis.Client
}

// Terminate gracefully shuts down all containers
func (c *Container) Terminate() error {
	ctx := context.Background()

	if c.PostgresContainer != nil {
		if err := c.PostgresContainer.Terminate(ctx); err != nil {
			return err
		}
	}

	if c.RedisContainer != nil {
		if err := c.RedisContainer.Terminate(ctx); err != nil {
			return err
		}
	}

	return nil
}

// CleanupDatabase cleans all test data from database tables
func (c *Container) CleanupDatabase(t *testing.T) {
	t.Helper()

	if c.DB == nil {
		return
	}

	// Dynamically get all user tables from information_schema
	var tables []string
	err := c.DB.Raw(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_type = 'BASE TABLE'
		AND table_name NOT LIKE 'pg_%'
		AND table_name NOT LIKE 'sql_%'
	`).Scan(&tables).Error

	if err != nil {
		t.Logf("failed to get table list: %v", err)
		return
	}

	// Truncate all found tables
	for _, table := range tables {
		err = c.DB.Exec("TRUNCATE TABLE " + table + " CASCADE").Error
		if err != nil {
			t.Logf("failed to truncate table %s: %v", table, err)
		}
	}
}

// CleanupCache clears all Redis cache data
func (c *Container) CleanupCache(t *testing.T) {
	t.Helper()

	if c.Redis == nil {
		return
	}

	err := c.Redis.FlushAll().Err()
	if err != nil {
		t.Logf("failed to flush Redis cache: %v", err)
	}
}

// CleanupAll cleans both database and cache
func (c *Container) CleanupAll(t *testing.T) {
	c.CleanupDatabase(t)
	c.CleanupCache(t)
}

// SetupTestMain sets up shared containers for a test package
// Should be called from TestMain to initialize containers once per package
func SetupTestMain() (*Container, func() int) {
	ctx := context.Background()

	// Setup PostgreSQL
	postgresContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute)),
	)
	if err != nil {
		panic("failed to start postgres container: " + err.Error())
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		panic("failed to get postgres host: " + err.Error())
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		panic("failed to get postgres port: " + err.Error())
	}

	cfg := &config.Config{
		Database: config.Database{
			Host:        host,
			Port:        port.Int(),
			User:        "testuser",
			Password:    "testpass",
			Name:        "testdb",
			MaxIdleConn: 10,
			MaxOpenConn: 100,
		},
	}

	dbConn, err := db.NewPostgres(cfg)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	// Setup Redis
	redisContainer, err := redisContainer.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(wait.ForLog("Ready to accept connections")),
	)
	if err != nil {
		panic("failed to start redis container: " + err.Error())
	}

	redisHost, err := redisContainer.Host(ctx)
	if err != nil {
		panic("failed to get redis host: " + err.Error())
	}

	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		panic("failed to get redis port: " + err.Error())
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort.Port(),
	})

	// Test Redis connection
	err = redisClient.Ping().Err()
	if err != nil {
		panic("failed to connect to test redis: " + err.Error())
	}

	container := &Container{
		PostgresContainer: postgresContainer,
		RedisContainer:    redisContainer,
		DB:                dbConn,
		Redis:             redisClient,
	}

	// Return cleanup function for TestMain
	cleanup := func() int {
		if err := container.Terminate(); err != nil {
			println("failed to terminate containers:", err.Error())
			return 1
		}
		return 0
	}

	return container, cleanup
}

// RunStandardMigrations runs migrations for common models used across tests
func (c *Container) RunStandardMigrations(t *testing.T) {
	t.Helper()

	err := c.DB.AutoMigrate(&user.User{}, &user.Preference{})
	require.NoError(t, err)
}
