package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// Config represents the main application configuration structure
// It contains all environment variables and nested configuration objects
type Config struct {
	AppSecret string `env:"APP_SECRET"`
	CacheURL  string `env:"CACHE_URL"`
	Database  Database
}

// Database represents the database connection configuration
// It contains all database-related environment variables
type Database struct {
	Host        string `env:"DATABASE_HOST"`
	Port        int    `env:"DATABASE_PORT"`
	User        string `env:"DATABASE_USER"`
	Password    string `env:"DATABASE_PASSWORD"`
	Name        string `env:"DATABASE_NAME"`
	MaxIdleConn int    `env:"DATABASE_MAX_IDLE_CONN"`
	MaxOpenConn int    `env:"DATABASE_MAX_OPEN_CONN"`
}

// DataSourceName returns a PostgreSQL connection string formatted with the database configuration.
func (d Database) DataSourceName() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Name)
}

// LoadEnv loads environment variables and returns a configured Config instance.
// It uses the env package to parse environment variables into the Config struct.
func LoadEnv() *Config {
	// Load .env file if it exists (for development)
	_ = godotenv.Load() // Ignore error if .env doesn't exist

	var c Config
	if err := env.Parse(&c); err != nil {
		log.Fatal().Msgf("failed to load env: %s", err.Error())
	}

	return &c
}
