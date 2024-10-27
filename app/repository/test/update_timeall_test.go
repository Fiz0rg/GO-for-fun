package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"time_app/app/model"
	"time_app/app/repository"
	"time_app/db"
	"time_app/test/fixture"
	"time_app/test/test_config"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestUpdateTimeAllByIntervals(t *testing.T) {
	ctx, cancel := repository.InitContext(1 * time.Second)
	defer cancel()
	db := test_config.TestDB
	if db == nil {
		t.Fatal("DATABASE IS NIL")
	}
	arrangeIntervalList, intervalCollection := fixture.GenIntervalsInMongo(db, ctx)
	arrangeTimeAllList, timeAllCollection := fixture.GenTimeAllsInMongo(db, ctx)

	countTimeRepo := repository.NewCountTimeRepository(db)

	status_code, err := countTimeRepo.TimeCalculation()
	if err != nil {
		t.Fatalf("Something wrong in repo, %v", err)
	}
	if status_code != 200 {
		t.Fatalf("Status code != 200, %v", err)
	}

	arrangeResult := getArrangeResult(arrangeIntervalList, arrangeTimeAllList)
	resultRepo := getTimeAllRecords(ctx, db, t)

	intervalList := getIntervalsRecords(ctx, db, t)

	assert.Equal(t, arrangeResult, resultRepo)
	assert.Empty(t, intervalList)

	t.Cleanup(func() {
		collectionList := []mongo.Collection{intervalCollection, timeAllCollection}
		time.Sleep(300 * time.Microsecond)
		test_config.CleanupTestData(collectionList)
	})
}

func getArrangeResult(intervals []model.Interval, timeAll []model.TimeAll) []model.TimeAll {
	intervalTimeSubtraction := subtractionIntervalsTime(intervals)

	for i := range timeAll {
		for _, interval := range intervalTimeSubtraction {
			if timeAll[i].UserUUID == interval.UserUUID && timeAll[i].CategoryUUID == interval.CategoryUUID {
				timeAll[i].TimeTotal += int(interval.TimeTotal)
			}
		}
	}
	return timeAll
}

func getIntervalsRecords(ctx context.Context, resource *db.Resource, t *testing.T) []model.Interval {
	pipeline := repository.FormPipeline()
	collection := resource.DB.Collection("Interval")
	records, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		t.Fatalf("Err getting interval records, %v", err)
	}

	var res []model.Interval
	err = records.All(ctx, &res)
	if err != nil {
		fmt.Printf("Decode interval Error, %v", err)
	}

	return res
}

func getTimeAllRecords(ctx context.Context, resource *db.Resource, t *testing.T) []model.TimeAll {
	collection := resource.DB.Collection("TimeAll")
	records, err := collection.Find(ctx, bson.M{})
	if err != nil {
		t.Fatalf("Err getting TimeAll records, %v", err)
	}

	var res []model.TimeAll
	err = records.All(ctx, &res)
	if err != nil {
		fmt.Printf("Decode TimeAll Error, %v", err)
	}

	return res
}

type TimeAll struct {
	UserUUID     string
	CategoryUUID string
	TimeTotal    int64
}

func subtractionIntervalsTime(intervals []model.Interval) []TimeAll {
	timeMap := make(map[string]*TimeAll)

	for _, interval := range intervals {
		if interval.EndAt == nil {
			continue
		}

		key := fmt.Sprintf("%s-%s", interval.UserUUID, interval.CategoryUUID)
		if _, exists := timeMap[key]; !exists {
			timeMap[key] = &TimeAll{
				UserUUID:     interval.UserUUID,
				CategoryUUID: interval.CategoryUUID,
				TimeTotal:    0,
			}
		}
		timeMap[key].TimeTotal += *interval.EndAt - interval.StartedAt
	}

	var timeTotals []TimeAll
	for _, timeAll := range timeMap {
		timeTotals = append(timeTotals, *timeAll)
	}
	return timeTotals
}
