package test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
	"time_app/app/repository"
	"time_app/app/repository/model"
	"time_app/app/repository/mongodb"
	"time_app/app/repository/tests/fixture"
	"time_app/db"
	test_config "time_app/tests"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestTimeAllByIntervals(t *testing.T) {
	ctx, cancel := repository.InitContext(1 * time.Second)
	defer cancel()
	db := test_config.TestDB
	if db == nil {
		t.Fatal("DATABASE IS NIL")
	}
	arrangeIntervalList, intervalCollection := fixture.GenIntervalsData(db, ctx)
	arrangeTimeAllList, timeAllCollection := fixture.GenTimeAllData(db, ctx)

	countTimeRepo := mongodb.NewCountTimeRepository(db)

	err := countTimeRepo.TimeCalculation()
	if err != nil {
		t.Fatalf("Something wrong in repo, %v", err)
	}
	arrangeResult := getArrangeResult(arrangeIntervalList, arrangeTimeAllList)
	resultRepo := getTimeAllRecords(ctx, db, t)

	intervalList := getIntervalsRecords(ctx, db, t)

	assert.Equal(t, arrangeResult, resultRepo)
	assert.Empty(t, intervalList)

	t.Cleanup(func() {
		time.Sleep(300 * time.Microsecond)
		test_config.CleanupTestData(db)
	})
}

func getArrangeResult(intervals []model.Interval, timeAll []model.TimeAll) []model.TimeAll {
	intervalTimeSubtraction := subtractionIntervalsTime(intervals)

	for i := range timeAll {
		for _, interval := range intervalTimeSubtraction {
			if timeAll[i].UserUUID == interval.UserUUID && timeAll[i].CategoryUUID == interval.CategoryUUID {
				timeAll[i].TimeTotal += interval.TimeTotal
			}
		}
	}
	return timeAll
}

func getIntervalsRecords(ctx context.Context, resource *db.Resource, t *testing.T) []model.Interval {
	pipeline := mongodb.FormPipeline()
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

func subtractionIntervalsTime(intervals []model.Interval) []model.TimeAll {
	timeMap := make(map[string]*model.TimeAll)

	for _, interval := range intervals {
		if interval.EndAt == nil {
			continue
		}

		key := fmt.Sprintf("%s-%s", interval.UserUUID, interval.CategoryUUID)
		if _, exists := timeMap[key]; !exists {
			timeMap[key] = &model.TimeAll{
				UserUUID:     interval.UserUUID,
				CategoryUUID: interval.CategoryUUID,
				TimeTotal:    0,
			}
		}
		timeMap[key].TimeTotal += *interval.EndAt - interval.StartedAt
	}

	var timeTotals []model.TimeAll
	for _, timeAll := range timeMap {
		timeTotals = append(timeTotals, *timeAll)
	}
	log.Printf("RESULT , %v", timeTotals)
	return timeTotals
}
