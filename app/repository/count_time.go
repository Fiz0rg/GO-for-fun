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
	UpdateTimeAll() (int, error)
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

func (repo *UpdateTimeAllCollectionRepositoryImpl) UpdateTimeAll() (int, error) {
	ctx, cancel := InitContext()
	defer cancel()

	userCategoryInterals := repo.getIntervalRecords(ctx)

	// Размер пакета, в котором будет 100 записей (можно изменить кол-во)
	batchSize := 100
	var wg sync.WaitGroup

	// Канал горутины, которая будет обрабатывать отправленные в неё пакеты batch
	updatesChannel := make(chan []mongo.WriteModel, 10)
	// Канал для ошибок
	errorChannel := make(chan error, 10)

	// Устанавливаем воркеры, которые будут работать с bulk (осторожно с ядрами)
	numWorkers := 3
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			performBulkWrite(ctx, repo.timeAllCollection, <-updatesChannel, &wg, errorChannel)
		}()
	}

	// Запихиваем операции с коллекцией в пакеты и отправляем обрабатываться
	var updates []mongo.WriteModel
	for i, item := range userCategoryInterals {

		filter := bson.M{
			"user_uuid":     item.UserUUID,
			"category_uuid": item.CategoryUUID,
		}
		update := bson.M{
			"$inc": bson.M{"time_total": item.TotalIntervalTime},
		}

		updateModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		updates = append(updates, updateModel)

		// Отправка пакета (batch) на выполнение при достижении размера
		if len(updates) == batchSize {
			updatesChannel <- updates
			updates = nil
		}

		// Просто чтобы видна была работа крч
		if i%10000 == 0 {
			fmt.Printf("Processed %d results\n", i)
		}
	}

	// отправляем то, что не отправилось
	if len(updates) > 0 {
		updatesChannel <- updates
	}
	// закрываем канал только когда всё закрыто
	close(updatesChannel)
	// ждём завершения всех воркеров
	wg.Wait()

	return http.StatusOK, nil

}

func performBulkWrite(ctx context.Context, collection *mongo.Collection, updates []mongo.WriteModel, wg *sync.WaitGroup, errorChannel chan<- error) {
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
	pipeline := formPipeline()
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

func formPipeline() mongo.Pipeline {
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
