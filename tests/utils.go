package test_config

import (
	"context"
	"time_app/db"

	"go.mongodb.org/mongo-driver/mongo"
)

func CleanupTestData(resource *db.Resource) {
	intervalCollection := resource.DB.Collection("Interval")
	timeDayCollection := resource.DB.Collection("TimeDay")
	timeAllCollection := resource.DB.Collection("TimeAll")
	userCollection := resource.DB.Collection("User")
	categoryCollection := resource.DB.Collection("Category")

	collectionList := []mongo.Collection{*intervalCollection, *timeAllCollection, *timeDayCollection,
		*userCollection, *categoryCollection}
	for _, coll := range collectionList {
		coll.Database().Drop(context.TODO())
	}
}
