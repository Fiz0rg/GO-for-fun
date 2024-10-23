package test_config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"time_app/db"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	Database DatabaseConfig
}

func LoadTestConfig() (*config, error) {

	if err := godotenv.Load("../../../.env"); err != nil {
		log.Panicf("ERROR GETTING ENV, %v", err)
	}
	return &config{
		Database: loadDatabaseConfig(),
	}, nil
}

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

func TestInitResource() (*db.Resource, error) { // хз норм ли так (папка не тестов)
	config, err_conf := LoadTestConfig()
	if err_conf != nil {
		log.Printf("Something wrong with config, %v", err_conf)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	dbName := config.Database.MONGODB_DATABASE
	dbPort := config.Database.MONGODB_PORT
	dbHost := config.Database.MONGODB_HOST

	URI := fmt.Sprintf("mongodb://%s:%s/authSource=timeappdb&retryWrites=true&w=majority", dbHost, dbPort)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB!")

	defer cancel()

	return &db.Resource{DB: client.Database(dbName)}, nil
}
