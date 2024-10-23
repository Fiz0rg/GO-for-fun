package test_config

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func CleanupTestData(collectionList []mongo.Collection) {
	for _, coll := range collectionList {
		coll.Database().Drop(context.TODO())
	}
}
