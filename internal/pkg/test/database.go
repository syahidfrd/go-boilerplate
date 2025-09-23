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

// SetupPostgresContainer creates a PostgreSQL testcontainer with migrations
func SetupPostgresContainer(t *testing.T) *Container {
	t.Helper()

	ctx := context.Background()

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute)),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

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
	require.NoError(t, err)

	return &Container{
		PostgresContainer: postgresContainer,
		DB:                dbConn,
	}
}

// MigrateAll runs migrations for all common models
func (c *Container) MigrateAll(t *testing.T) {
	t.Helper()

	err := db.AutoMigrate(c.DB, &user.User{}, &user.Preference{})
	require.NoError(t, err)
}

// MigrateModels runs migrations for specific models
func (c *Container) MigrateModels(t *testing.T, models ...any) {
	t.Helper()

	err := db.AutoMigrate(c.DB, models...)
	require.NoError(t, err)
}

// CleanupTables truncates all tables for test isolation
func (c *Container) CleanupTables(t *testing.T, tables ...string) {
	t.Helper()

	for _, table := range tables {
		err := c.DB.Exec("TRUNCATE TABLE " + table + " CASCADE").Error
		require.NoError(t, err)
	}
}

// SetupRedisContainer creates a Redis testcontainer
func SetupRedisContainer(t *testing.T) *Container {
	t.Helper()

	ctx := context.Background()

	redisContainer, err := redisContainer.RunContainer(ctx,
		testcontainers.WithImage("redis:7-alpine"),
		testcontainers.WithWaitStrategy(wait.ForLog("Ready to accept connections")),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate redis container: %s", err)
		}
	})

	host, err := redisContainer.Host(ctx)
	require.NoError(t, err)

	port, err := redisContainer.MappedPort(ctx, "6379")
	require.NoError(t, err)

	redisClient := redis.NewClient(&redis.Options{
		Addr: host + ":" + port.Port(),
	})

	// Test connection
	err = redisClient.Ping().Err()
	require.NoError(t, err)

	return &Container{
		RedisContainer: redisContainer,
		Redis:          redisClient,
	}
}

// SetupFullContainer creates both PostgreSQL and Redis testcontainers
func SetupFullContainer(t *testing.T) *Container {
	t.Helper()

	ctx := context.Background()

	// Setup PostgreSQL
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute)),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate postgres container: %s", err)
		}
	})

	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

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
	require.NoError(t, err)

	// Setup Redis
	redisContainer, err := redisContainer.RunContainer(ctx,
		testcontainers.WithImage("redis:7-alpine"),
		testcontainers.WithWaitStrategy(wait.ForLog("Ready to accept connections")),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate redis container: %s", err)
		}
	})

	redisHost, err := redisContainer.Host(ctx)
	require.NoError(t, err)

	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	require.NoError(t, err)

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort.Port(),
	})

	// Test Redis connection
	err = redisClient.Ping().Err()
	require.NoError(t, err)

	return &Container{
		PostgresContainer: postgresContainer,
		RedisContainer:    redisContainer,
		DB:                dbConn,
		Redis:             redisClient,
	}
}