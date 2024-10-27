package config

import "os"

type DatabaseConfig struct {
	ENV string

	MONGODB_USERNAME string
	MONGODB_PASSWORD string

	MONGODB_HOST     string
	MONGODB_PORT     string
	MONGODB_DATABASE string
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		ENV:              os.Getenv("ENV"),
		MONGODB_USERNAME: os.Getenv("MONGODB_USERNAME"),
		MONGODB_PASSWORD: os.Getenv("MONGODB_PASSWORD"),
		MONGODB_HOST:     os.Getenv("MONGODB_HOST"),
		MONGODB_PORT:     os.Getenv("MONGODB_PORT"),
		MONGODB_DATABASE: os.Getenv("MONGODB_DATABASE"),
	}
}
