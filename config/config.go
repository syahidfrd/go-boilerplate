package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config ...
type Config struct {
	ServerPORT  string
	DatabaseURL string
	CacheURL    string
}

// LoadConfig will load config from environment variable
func LoadConfig() (config *Config) {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	serverPORT := os.Getenv("SERVER_PORT")
	databaseURL := os.Getenv("DATABASE_URL")
	cacheURL := os.Getenv("CACHE_URL")

	return &Config{
		ServerPORT:  serverPORT,
		DatabaseURL: databaseURL,
		CacheURL:    cacheURL,
	}
}
