package fixture

import (
	"context"
	"time_app/app/repository/model"
	"time_app/db"

	"go.mongodb.org/mongo-driver/mongo"
)

func GenTimeAllData(resource *db.Resource, ctx context.Context) ([]model.TimeAll, mongo.Collection) {
	timeAllCollection := resource.DB.Collection("TimeAll")
	itemList := []model.TimeAll{
		{
			UUID:         "uuid-3",
			UserUUID:     "user-1",
			CategoryUUID: "category-001",
			TimeTotal:    100,
		},
		{
			UUID:         "uuid-4",
			UserUUID:     "user-1",
			CategoryUUID: "category-002",
			TimeTotal:    300,
		},
	}
	for _, item := range itemList {
		timeAllCollection.InsertOne(ctx, item)
	}
	return itemList, *timeAllCollection
}
