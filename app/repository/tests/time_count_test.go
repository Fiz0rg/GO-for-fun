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
)

func TestTimeCount(t *testing.T) {
	ctx, cancel := repository.InitContext(5 * time.Second)
	defer cancel()
	db := test_config.TestDB
	if db == nil {
		t.Fatal("DATABASE IS NIL")
	}

	arrangeTimeDay, arrangeTimeAll := arrangeData(db, ctx)

	countTimeRepo := mongodb.NewCountTimeRepository(db)
	err := countTimeRepo.TimeCalculation()
	if err != nil {
		t.Fatalf("Something wrong in repo, %v", err)
	}
	resultTimeAll := getTimeAllRecords(ctx, db, t)
	resultTimeDay := getTimeDayRecords(ctx, db, t)

	intervalListAfterCount := getIntervalsRecords(ctx, db, t)
	assert.Equal(t, arrangeTimeAll, resultTimeAll)
	assert.Equal(t, arrangeTimeDay, resultTimeDay)
	assert.Empty(t, intervalListAfterCount)

	t.Cleanup(func() {
		time.Sleep(3000 * time.Microsecond)
		test_config.CleanupTestData(db)
	})
}

func arrangeData(db *db.Resource, ctx context.Context) ([]model.TimeDay, []model.TimeAll) {
	var categoryAmount int = 3
	var intervalAmount int = 10

	user := fixture.CreateUser(db, ctx)
	categoryList := fixture.CreateManyCategories(db, ctx, &user, &categoryAmount)

	arrangeIntervslList := fixture.CreateManyIntervals(db, ctx, &user, &categoryList, &intervalAmount)
	arrangeTimeAllList := fixture.CreateManyTimeAll(db, ctx, &user, &categoryList)
	arrangeTimeDayList := fixture.CreateManyTimeDay(db, ctx, &user, &categoryList)

	arrangeTimeAll, arrangeTimeDay := getArrangeResult(arrangeIntervslList, arrangeTimeAllList, arrangeTimeDayList)

	return arrangeTimeDay, arrangeTimeAll
}

func getTimeDayRecords(ctx context.Context, db *db.Resource, t *testing.T) []model.TimeDay {
	c := db.DB.Collection("TimeDay")
	stmt, err := c.Find(ctx, bson.M{})
	if err != nil {
		t.Fatalf("Find TimeDay Error, %v", err)
	}
	var res []model.TimeDay
	err = stmt.All(ctx, &res)
	if err != nil {
		log.Printf("Decode TimeDay Error, %v", err)
	}
	return res
}

func getArrangeResult(intervals []model.Interval, timeAll []model.TimeAll, timeDays []model.TimeDay) ([]model.TimeAll, []model.TimeDay) {
	filteredInterval := filterExpectedUUIDs(intervals)

	for i := range timeAll {
		for _, interval := range filteredInterval {
			if timeAll[i].UserUUID == interval.UserUUID && timeAll[i].CategoryUUID == interval.CategoryUUID {
				timeAll[i].TimeTotal += *interval.EndAt - interval.StartedAt
			}
		}
	}

	splittedIntervalsByDays := mongodb.SplitIntervals(filteredInterval)
	for t := range timeDays {
		for _, i := range splittedIntervalsByDays {

			sameUser := timeDays[t].UserUUID == i.UserUUID
			sameCategory := timeDays[t].CategoryUUID == i.CategoryUUID
			withinTimeRange := i.StartedAt <= timeDays[t].TimeDay && timeDays[t].TimeDay <= int64(*i.EndAt)

			if sameUser && sameCategory && withinTimeRange {
				timeDays[t].TimeTotal += int(*i.EndAt) - int(i.StartedAt)
			}
		}
	}

	return timeAll, timeDays
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
