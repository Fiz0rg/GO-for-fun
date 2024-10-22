package config

import "os"

type DatabaseConfig struct {
	MONGODB_HOST     string
	MONGODB_PORT     string
	MONGODB_DATABASE string
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		MONGODB_HOST:     os.Getenv("MONGODB_HOST"),
		MONGODB_PORT:     os.Getenv("MONGODB_PORT"),
		MONGODB_DATABASE: os.Getenv("MONGODB_DATABASE"),
	}
}
