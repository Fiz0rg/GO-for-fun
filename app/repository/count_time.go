package repository

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time_app/app/model"
	"time_app/db"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CountTimeRepo UpdateTimeAllCollectionRepository

type UpdateTimeAllCollectionRepositoryImpl struct {
	resource           *db.Resource
	intervalCollection *mongo.Collection
	timeAllCollection  *mongo.Collection
}

type UpdateTimeAllCollectionRepository interface {
	TimeCalculation() (int, error)
	getIntervalRecords(ctx context.Context) []model.UserCategory
	getTimeAllRecords(ctx context.Context) []model.TimeAll
}

func NewCountTimeRepository(resource *db.Resource) UpdateTimeAllCollectionRepository {
	intervalCollection := resource.DB.Collection("Interval")
	timeAllCollection := resource.DB.Collection("TimeAll")
	countTimeRepo := &UpdateTimeAllCollectionRepositoryImpl{
		resource:           resource,
		intervalCollection: intervalCollection,
		timeAllCollection:  timeAllCollection}
	return countTimeRepo
}

func updateTimeAllByTimeDays() {
	panic("NOT IMPLEMENTED")
}

func (repo *UpdateTimeAllCollectionRepositoryImpl) TimeCalculation() (int, error) {
	ctx, cancel := InitContext()
	defer cancel()
	var wg sync.WaitGroup

	userCategoryInterals := repo.getIntervalRecords(ctx)

	u := updateTimeAllByIntervals(ctx, &wg, repo.timeAllCollection, userCategoryInterals)
	if u != nil {
		log.Panicf("Fail operation by updateTimeAllByIntervals, %v", u)
	}

	d := deleteUnnecessaryIntervals(ctx, &wg, repo.intervalCollection, userCategoryInterals)
	if d != nil {
		log.Panicf("Fail operation by deleteUnnecessaryIntervals, %v", d)
	}

	wg.Wait()
	return http.StatusOK, nil

}

func deletePerfomBulkWrite(ctx context.Context, collection *mongo.Collection, deletions []mongo.WriteModel, waitGroup *sync.WaitGroup, errorDeleteChannel chan error) {
	defer waitGroup.Done()
	if len(deletions) > 0 {
		bulkOpt := options.BulkWrite().SetOrdered(false)
		_, err := collection.BulkWrite(ctx, deletions, bulkOpt)
		if err != nil {
			errorDeleteChannel <- fmt.Errorf("BulkWrite error: %v", err)
		} else {
			fmt.Printf("Bulk write completed for %d deleted\n", len(deletions))
		}
	}
}

func updatePerformBulkWrite(ctx context.Context, collection *mongo.Collection, updates []mongo.WriteModel, wg *sync.WaitGroup, errorChannel chan<- error) {
	defer wg.Done()
	if len(updates) > 0 {
		bulkOpts := options.BulkWrite().SetOrdered(false) // неупорядоченная обработка для монги
		_, err := collection.BulkWrite(ctx, updates, bulkOpts)
		if err != nil {
			errorChannel <- fmt.Errorf("BulkWrite error: %v", err)
		} else {
			fmt.Printf("Bulk write completed for %d updates\n", len(updates))
		}
	}
}

func (repo *UpdateTimeAllCollectionRepositoryImpl) getIntervalRecords(ctx context.Context) []model.UserCategory {
	pipeline := FormPipeline()
	cursor, err := repo.intervalCollection.Aggregate(ctx, pipeline)

	if err != nil {
		log.Printf("Get interval method error, %v", err)
	}

	var res []model.UserCategory
	err = cursor.All(ctx, &res)
	if err != nil {
		log.Printf("Error in decode result from interval collection, %v", err)
	}
	return res
}

func (repo *UpdateTimeAllCollectionRepositoryImpl) getTimeAllRecords(ctx context.Context) []model.TimeAll {
	timeAllList := []model.TimeAll{}
	cursor, err := repo.timeAllCollection.Find(ctx, bson.M{})

	if err != nil {
		log.Printf("Didnt get records from TimeAll collection, %v", err)
		return timeAllList
	}

	for cursor.Next(ctx) {
		var timeAll model.TimeAll
		err = cursor.Decode(&timeAll)
		if err != nil {
			logrus.Print(err)
		}
		timeAllList = append(timeAllList, timeAll)
	}
	cursor.Close(ctx)
	return timeAllList
}

func deleteUnnecessaryIntervals(
	ctx context.Context,
	wg *sync.WaitGroup,
	interval_collection *mongo.Collection,
	userCategoryIntervals []model.UserCategory,
) error {
	deleteChannel := make(chan []mongo.WriteModel, 5)
	errorDeleteChannel := make(chan error, 10)

	numDeleteWorkers := 2
	for i := 0; i < numDeleteWorkers; i++ {
		wg.Add(1)
		go func() {
			deletePerfomBulkWrite(ctx, interval_collection, <-deleteChannel, wg, errorDeleteChannel)
		}()
	}

	batchSize := 100

	var deleteIntervals []mongo.WriteModel
	for _, item := range userCategoryIntervals {
		for _, intervalUUID := range item.UUIDList {
			intervalfilter := bson.M{"uuid": intervalUUID}
			stmt := mongo.NewDeleteOneModel().SetFilter(intervalfilter)
			deleteIntervals = append(deleteIntervals, stmt)
		}
		if len(deleteIntervals) == batchSize {
			deleteChannel <- deleteIntervals
			deleteIntervals = nil
		}
	}
	if len(deleteIntervals) > 0 {
		deleteChannel <- deleteIntervals
	}

	close(deleteChannel)
	return nil
}

func updateTimeAllByIntervals(
	ctx context.Context,
	wg *sync.WaitGroup,
	timeAllCollection *mongo.Collection,
	userCategoryInterals []model.UserCategory,
) error {
	batchSize := 50

	updatesChannel := make(chan []mongo.WriteModel, 5)
	errorUpdateChannel := make(chan error, 10)

	numUpdateWorkers := 2
	for i := 0; i < numUpdateWorkers; i++ {
		wg.Add(1)
		go func() {
			updatePerformBulkWrite(ctx, timeAllCollection, <-updatesChannel, wg, errorUpdateChannel)
		}()
	}
	var updateTimeAll []mongo.WriteModel
	for _, item := range userCategoryInterals {

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
	return nil
}

func FormPipeline() mongo.Pipeline {
	pipeline := mongo.Pipeline{

		// Stage 1: Фильтруем записи (убираем записи без end_at и последние в группе)
		{
			{Key: "$match", Value: bson.D{
				{Key: "end_at", Value: bson.D{{Key: "$ne", Value: nil}}},
				{Key: "rowNum", Value: bson.D{{Key: "$ne", Value: 1}}},
			}},
		},
		// Stage 2: Группируем по user_uuid и category_uuid
		{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "user_uuid", Value: "$user_uuid"},
					{Key: "category_uuid", Value: "$category_uuid"},
				}},
				{Key: "total_interval_time", Value: bson.D{
					{Key: "$sum", Value: bson.D{
						{Key: "$subtract", Value: bson.A{"$end_at", "$started_at"}},
					}},
				}},
				{Key: "uuid_list", Value: bson.D{
					{Key: "$push", Value: "$uuid"},
				}},
			}},
		},
		// Stage 3: Проверяем, что в группе больше одной записи
		{
			{Key: "$match", Value: bson.D{
				{Key: "$expr", Value: bson.D{
					{Key: "$gt", Value: bson.A{
						bson.D{{Key: "$size", Value: "$uuid_list"}},
						1,
					}},
				}},
			}},
		},
		// Stage 4: Формируем финальный вывод
		{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "user_uuid", Value: "$_id.user_uuid"},
				{Key: "category_uuid", Value: "$_id.category_uuid"},
				{Key: "total_interval_time", Value: 1},
				{Key: "uuid_list", Value: 1},
			}},
		},
	}
	return pipeline
}
