package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPORT  string
	DatabaseURL string
}

func LoadConfig() (config *Config) {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	serverPORT := os.Getenv("SERVER_PORT")
	databaseURL := os.Getenv("DATABASE_URL")

	return &Config{
		ServerPORT:  serverPORT,
		DatabaseURL: databaseURL,
	}
}
