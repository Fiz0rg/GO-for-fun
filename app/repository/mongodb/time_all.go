package mongodb

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time_app/app/repository/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func updateTimeAllPerformBulkWrite(ctx context.Context, collection *mongo.Collection, updates []mongo.WriteModel, wg *sync.WaitGroup, errorChannel chan<- error) {
	defer wg.Done()
	if len(updates) > 0 {
		log.Printf("UPDATES IS GETE, %v", len(updates))
		bulkOpts := options.BulkWrite().SetOrdered(false) // неупорядоченная обработка для монги
		_, err := collection.BulkWrite(ctx, updates, bulkOpts)
		if err != nil {
			errorChannel <- fmt.Errorf("BulkWrite error: %v", err)
		} else {
			fmt.Printf("Bulk write completed for %d updates\n", len(updates))
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
		log.Printf("INTERVAL , %v", i)
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
			log.Printf("newITEM , %v", newItem)
			userCategoryMap[key] = newItem
		}
	}
	log.Printf("USERCAT, %v", userCategoryMap)
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
