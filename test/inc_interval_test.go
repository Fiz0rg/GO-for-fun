package test

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
	"time_app/app/model"
	"time_app/app/repository"
	"time_app/db"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestUpdateTimeAllByIntervals(t *testing.T) {
	db, err := testInitResource()
	if err != nil {
		t.Fatalf("Failed to connect MONGO, %v", err)
	}

	ctx, cancel := repository.InitContext()
	defer cancel()

	arrangeIntervalList, intervalCollection := genIntervalsInMongo(db, ctx)
	arrangeTimeAllList, timeAllCollection := genTimeAllsInMongo(db, ctx)

	countTimeRepo := repository.NewCountTimeRepository(db)

	status_code, err := countTimeRepo.UpdateTimeAll()
	if err != nil {
		t.Fatalf("Something wrong in repo, %v", err)
	}
	if status_code != 200 {
		t.Fatalf("Status code != 200, %v", err)
	}

	arrangeResult := getArrangeResult(arrangeIntervalList, arrangeTimeAllList)
	resultRepo := getTimeAllRecords(ctx, db, t)

	assert.Equal(t, arrangeResult, resultRepo)

	t.Cleanup(func() {
		collectionList := []mongo.Collection{intervalCollection, timeAllCollection}
		time.Sleep(300 * time.Microsecond)
		cleanupTestData(collectionList)

		pc, _, _, _ := runtime.Caller(0) // Получаем PC текущей функции
		funcName := runtime.FuncForPC(pc).Name()
		printConsoleTestPass(funcName)
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
	collection := resource.DB.Collection("Interval")
	records, err := collection.Find(ctx, bson.M{})
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
