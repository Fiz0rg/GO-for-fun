package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	app      AppConfig
	Database DatabaseConfig
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	return &Config{
		app:      loadAppConfig(),
		Database: loadDatabaseConfig(),
	}, nil
}
