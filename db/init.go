package db

import (
	"context"
	"fmt"
	"log"
	"net/url"
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Err get config for database, %v", err)
	}

	dbName := config.Database.MONGODB_DATABASE
	dbHOST := config.Database.MONGODB_HOST
	dbPORT := config.Database.MONGODB_PORT
	dbUSERNAME := url.QueryEscape(config.Database.MONGODB_USERNAME)
	dbPASSWORD := url.QueryEscape(config.Database.MONGODB_PASSWORD)

	authData := ""
	if allNotEmpty(dbPASSWORD, dbUSERNAME) {
		authData = fmt.Sprintf("%v:%v@", dbUSERNAME, dbPASSWORD)
	}
	uriForm := fmt.Sprintf("mongodb://%s%s:%s/authSource=%s&retryWrites=true&w=majority", authData, dbHOST, dbPORT, dbName)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uriForm))

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

func allNotEmpty(values ...string) bool {
	for _, v := range values {
		if v == "" {
			return false
		}
	}
	return true
}
