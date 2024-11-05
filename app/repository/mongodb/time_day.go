package mongodb

import (
	"context"
	"fmt"
	"sync"
	"time_app/app/repository/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func updateTimeDayCollection(
	ctx context.Context,
	waitGroup *sync.WaitGroup,
	collection *mongo.Collection,
	intervals []model.Interval,
) {
	intervals = SplitIntervals(intervals)

	batchSize := 30
	updatesChanes := make(chan []mongo.WriteModel, 10)
	errorChanel := make(chan error, 10)
	defer close(errorChanel)

	numUpdateWorkers := 2
	for i := 0; i < numUpdateWorkers; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for updateBatch := range updatesChanes {
				updateTimeDayPerformBulkWrite(ctx, collection, updateBatch, errorChanel)
			}
		}()
	}
	var updateTimeDay []mongo.WriteModel
	for _, item := range intervals {
		timeDayfilter := bson.M{
			"user_uuid":     item.UserUUID,
			"category_uuid": item.CategoryUUID,
			"time_day": bson.M{
				"$gte": item.StartedAt,
				"$lte": item.EndAt,
			},
		}
		update := bson.M{
			"$inc": bson.M{"time_total": *item.EndAt - item.StartedAt},
		}
		updateRequest := mongo.NewUpdateOneModel().SetFilter(timeDayfilter).SetUpdate(update).SetUpsert(true)
		updateTimeDay = append(updateTimeDay, updateRequest)

		if len(updateTimeDay) == batchSize {
			updatesChanes <- updateTimeDay
			updateTimeDay = nil
		}
	}
	if len(updateTimeDay) > 0 {
		updatesChanes <- updateTimeDay
	}
	close(updatesChanes)
}

func updateTimeDayPerformBulkWrite(ctx context.Context, collection *mongo.Collection, updates []mongo.WriteModel, errorChanel chan<- error) {
	if len(updates) > 0 {
		bulkOpts := options.BulkWrite().SetOrdered(false)
		_, err := collection.BulkWrite(ctx, updates, bulkOpts)
		if err != nil {
			errorChanel <- fmt.Errorf("bulk write err, %v", err)
		} else {
			fmt.Printf("Updated %d TimeDay by intervals (Bulk Write)\n", len(updates))
		}
	}
}
