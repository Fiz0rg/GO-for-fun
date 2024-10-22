package db

import (
	"context"
	"fmt"
	"log"
	"time"
	"time_app/config"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Resource struct {
	DB *mongo.Database
}

func InitResource(config *config.Config) (*Resource, error) {
	dbName := config.Database.MONGODB_DATABASE
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Err get config for database, %v", err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017/authSource=timeappdb&retryWrites=true&w=majority"))

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

	return &Resource{DB: client.Database(dbName)}, nil
}
