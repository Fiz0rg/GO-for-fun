package test

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

func cleanupTestData(collectionList []mongo.Collection) {
	for _, coll := range collectionList {
		coll.Database().Drop(context.TODO())
	}
}

func printConsoleTestPass(funcName string) {
	log.Println()
	log.Printf("%v PASSED", funcName)
	log.Println()
}
