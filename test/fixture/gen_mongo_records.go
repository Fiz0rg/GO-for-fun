package fixture

import (
	"context"
	"log"
	"time_app/app/model"
	"time_app/db"

	"go.mongodb.org/mongo-driver/mongo"
)

func GenIntervalsInMongo(resource *db.Resource, ctx context.Context) ([]model.Interval, mongo.Collection) {
	endAt := int64(500)
	intervalCollections := resource.DB.Collection("Interval")
	intervalsInPipline := []model.Interval{
		{
			UUID:         "uuid-1",
			UserUUID:     "user-1",
			CategoryUUID: "category-001",
			StartedAt:    400,
			EndAt:        &endAt,
		},
		{
			UUID:         "uuid-2",
			UserUUID:     "user-1",
			CategoryUUID: "category-001",
			StartedAt:    400,
			EndAt:        &endAt,
		},
	}

	intervalsOutPipline := []model.Interval{
		{
			UUID:         "uuid-3",
			UserUUID:     "user-1",
			CategoryUUID: "category-002",
			StartedAt:    400,
			EndAt:        &endAt,
		},
		{
			UUID:         "uuid-3",
			UserUUID:     "user-1",
			CategoryUUID: "category-002",
			StartedAt:    400,
			EndAt:        nil,
		},
	}
	allIntervals := append(intervalsInPipline, intervalsOutPipline...)

	// Перед отправкой в insertMany без этого переобразования нельзя
	docs := make([]interface{}, len(allIntervals))
	for i, v := range allIntervals {
		docs[i] = v
	}

	_, err := intervalCollections.InsertMany(ctx, docs)
	if err != nil {
		log.Fatalf("Intervals not inserted, %v", err)
	}

	return intervalsInPipline, *intervalCollections
}

func GenTimeAllsInMongo(resource *db.Resource, ctx context.Context) ([]model.TimeAll, mongo.Collection) {
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
