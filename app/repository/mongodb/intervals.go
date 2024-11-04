package mongodb

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"time_app/app/repository"
	"time_app/app/repository/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *UpdateTimeRepositoryImpl) getIntervalRecords(ctx context.Context) []model.Interval {
	pipeline := FormPipeline()
	cursor, err := repo.intervalCollection.Aggregate(ctx, pipeline)

	if err != nil {
		log.Printf("Get interval method error, %v", err)
	}

	var res []model.Interval
	err = cursor.All(ctx, &res)
	if err != nil {
		log.Printf("Error in decode result from interval collection, %v", err)
	}
	return res
}

func deleteIntervalsPerfomBulkWrite(ctx context.Context, collection *mongo.Collection, deletions []mongo.WriteModel, waitGroup *sync.WaitGroup, errorDeleteChannel chan error) {
	defer waitGroup.Done()
	if len(deletions) > 0 {
		bulkOpt := options.BulkWrite().SetOrdered(false)
		_, err := collection.BulkWrite(ctx, deletions, bulkOpt)
		if err != nil {
			errorDeleteChannel <- fmt.Errorf("BulkWrite error: %v", err)
		} else {
			fmt.Printf("Completed %d deleted intervals\n", len(deletions))
		}
	}
}

func deleteUnnecessaryIntervals(
	ctx context.Context,
	wg *sync.WaitGroup,
	interval_collection *mongo.Collection,
	intervals []model.Interval,
) error {
	deleteChannel := make(chan []mongo.WriteModel, 5)
	errorDeleteChannel := make(chan error, 10)
	numDeleteWorkers := 2
	for i := 0; i < numDeleteWorkers; i++ {
		wg.Add(1)
		go func() {
			deleteIntervalsPerfomBulkWrite(ctx, interval_collection, <-deleteChannel, wg, errorDeleteChannel)
		}()
	}

	batchSize := 100

	var deleteIntervals []mongo.WriteModel
	for _, i := range intervals {
		intervalfilter := bson.M{"uuid": i.UUID}
		stmt := mongo.NewDeleteOneModel().SetFilter(intervalfilter)
		deleteIntervals = append(deleteIntervals, stmt)

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

func SplitIntervals(intervals []model.Interval) []model.Interval {
	var subIntervals []model.Interval

	for _, i := range intervals {
		start := time.Unix(i.StartedAt, 0)
		end := time.Unix(*i.EndAt, 0)

		for current := start; current.Before(end); {
			endOfDay := time.Date(current.Year(), current.Month(), current.Day(), 23, 59, 59, 0, current.Location())
			if endOfDay.After(end) {
				endOfDay = end
			}
			i := model.Interval{
				UUID:         i.UUID,
				UserUUID:     i.UserUUID,
				CategoryUUID: i.CategoryUUID,
				StartedAt:    current.Unix(),
				EndAt:        repository.IntPtr(endOfDay.Unix()),
			}
			subIntervals = append(subIntervals, i)
			current = endOfDay.Add(time.Second)
		}
	}
	return subIntervals
}

func FormPipeline() mongo.Pipeline {
	pipeline := mongo.Pipeline{
		// Stage 1: EndAt != None
		{
			{Key: "$match", Value: bson.D{
				{Key: "end_at", Value: bson.D{{Key: "$ne", Value: nil}}},
			}},
		},
		// Stage 2: Sort by user_uuid, category_uuid, and started_at
		{
			{Key: "$sort", Value: bson.D{
				{Key: "user_uuid", Value: 1},
				{Key: "category_uuid", Value: 1},
				{Key: "started_at", Value: 1},
			}},
		},
		// Stage 3: Group and collect records per user_uuid and category_uuid
		{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "user_uuid", Value: "$user_uuid"},
					{Key: "category_uuid", Value: "$category_uuid"},
				}},
				{Key: "records", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
			}},
		},
		// Stage 4: Filter groups with more than one record
		{
			{Key: "$match", Value: bson.D{
				{Key: "records.1", Value: bson.D{{Key: "$exists", Value: true}}},
			}},
		},
		// Stage 5: Slice to exclude the last record
		{
			{Key: "$project", Value: bson.D{
				{Key: "records", Value: bson.D{
					{Key: "$slice", Value: bson.A{
						"$records",
						bson.D{{Key: "$subtract", Value: bson.A{
							bson.D{{Key: "$size", Value: "$records"}},
							1,
						}}},
					}},
				}},
			}},
		},
		// Stage 6: Unwind records
		{
			{Key: "$unwind", Value: "$records"},
		},
		// Stage 7: Return records
		{
			{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$records"}}},
		},
	}
	return pipeline
}
