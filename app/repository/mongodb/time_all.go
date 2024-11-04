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

func updateTimeAllPerformBulkWrite(ctx context.Context, collection *mongo.Collection, updates []mongo.WriteModel, wg *sync.WaitGroup, errorChannel chan<- error) {
	defer wg.Done()
	if len(updates) > 0 {
		bulkOpts := options.BulkWrite().SetOrdered(false) // неупорядоченная обработка для монги
		_, err := collection.BulkWrite(ctx, updates, bulkOpts)
		if err != nil {
			errorChannel <- fmt.Errorf("bulk write error (timeall): %v", err)
		} else {
			fmt.Printf("TimeAll %d updates by intervals \n", len(updates))
		}
	}
}

func updateTimeAllCollection(
	ctx context.Context,
	wg *sync.WaitGroup,
	timeAllCollection *mongo.Collection,
	interals []model.Interval,
) {
	userCategoryMap := make(map[string]model.UserCategory)

	for _, i := range interals {
		key := i.UserUUID + i.CategoryUUID
		if item, exists := userCategoryMap[key]; exists {
			item.TotalIntervalTime += (*i.EndAt - i.StartedAt)
			userCategoryMap[key] = item
		} else {
			newItem := model.UserCategory{
				UserUUID:          i.UserUUID,
				CategoryUUID:      i.CategoryUUID,
				TotalIntervalTime: *i.EndAt - i.StartedAt,
			}
			userCategoryMap[key] = newItem
		}
	}
	batchSize := 50
	updatesChannel := make(chan []mongo.WriteModel, 5)
	errorUpdateChannel := make(chan error, 10)

	numUpdateWorkers := 2
	for i := 0; i < numUpdateWorkers; i++ {
		wg.Add(1)
		go func() {
			updateTimeAllPerformBulkWrite(ctx, timeAllCollection, <-updatesChannel, wg, errorUpdateChannel)
		}()
	}
	var updateTimeAll []mongo.WriteModel
	for _, item := range userCategoryMap {
		timeAllFilter := bson.M{
			"user_uuid":     item.UserUUID,
			"category_uuid": item.CategoryUUID,
		}
		update := bson.M{
			"$inc": bson.M{"time_total": item.TotalIntervalTime},
		}

		updateRequest := mongo.NewUpdateOneModel().SetFilter(timeAllFilter).SetUpdate(update).SetUpsert(true)
		updateTimeAll = append(updateTimeAll, updateRequest)

		if len(updateTimeAll) == batchSize {
			updatesChannel <- updateTimeAll
			updateTimeAll = nil
		}
	}

	if len(updateTimeAll) > 0 {
		updatesChannel <- updateTimeAll
	}

	close(updatesChannel)
}
