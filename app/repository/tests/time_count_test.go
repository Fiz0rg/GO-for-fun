package test

import (
	"context"
	"testing"
	"time"
	"time_app/app/repository"
	"time_app/app/repository/model"
	"time_app/app/repository/mongodb"
	"time_app/app/repository/tests/fixture"
	"time_app/db"
	test_config "time_app/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestTimeCount(t *testing.T) {
	ctx, cancel := repository.InitContext(5 * time.Second)
	defer cancel()
	db := test_config.TestDB
	require.NotNil(t, db, "Database connection should not be nil")

	expectedTimeDay, expectedTimeAll := setupTestData(db, ctx)

	countTimeRepo := mongodb.NewCountTimeRepository(db)
	err := countTimeRepo.TimeCalculation()
	assert.NoError(t, err, "TimeCalculation method should be without errors")

	resultTimeAll := getTimeAllRecords(ctx, db, t)
	resultTimeDay := getTimeDayRecords(ctx, db, t)

	intervalListAfterCount := getIntervalsRecords(ctx, db, t)

	assert.Equal(t, expectedTimeAll, resultTimeAll, "TimeAll missmatch")
	assert.Equal(t, expectedTimeDay, resultTimeDay, "TimeDay missmatch")
	assert.Empty(t, intervalListAfterCount, "Intervals list after repo should me empty")

	t.Cleanup(func() {
		time.Sleep(10 * time.Microsecond)
		test_config.CleanupTestData(db)
	})
}

func setupTestData(db *db.Resource, ctx context.Context) ([]model.TimeDay, []model.TimeAll) {
	var categoryAmount int = 4
	var intervalAmount int = 5

	user := fixture.CreateUser(db, ctx)
	categoryList := fixture.CreateManyCategories(db, ctx, &user, &categoryAmount)

	arrangeIntervslList := fixture.CreateManyIntervals(db, ctx, &user, &categoryList, &intervalAmount)
	arrangeTimeAllList := fixture.CreateManyTimeAll(db, ctx, &user, &categoryList)
	arrangeTimeDayList := fixture.CreateManyTimeDay(db, ctx, &user, &categoryList)

	arrangeTimeAll, arrangeTimeDay := getArrangeResult(arrangeIntervslList, arrangeTimeAllList, arrangeTimeDayList)

	return arrangeTimeDay, arrangeTimeAll
}

func isSameCategoryAndUser(interval model.Interval, timeDay model.TimeDay) bool {
	return timeDay.UserUUID == interval.UserUUID && timeDay.CategoryUUID == interval.CategoryUUID
}

func isWithinDayRange(interval model.Interval, timeDay model.TimeDay) bool {
	return interval.StartedAt <= timeDay.TimeDay && timeDay.TimeDay <= int64(*interval.EndAt)
}

func getArrangeResult(intervals []model.Interval, timeAll []model.TimeAll, timeDays []model.TimeDay) ([]model.TimeAll, []model.TimeDay) {
	filteredInterval := filterIntervalsEqualPipeline(intervals)

	for i := range timeAll {
		for _, interval := range filteredInterval {
			if timeAll[i].UserUUID == interval.UserUUID && timeAll[i].CategoryUUID == interval.CategoryUUID {
				timeAll[i].TimeTotal += *interval.EndAt - interval.StartedAt
			}
		}
	}

	splittedIntervalsByDays := mongodb.SplitIntervals(filteredInterval)
	for t := range timeDays {
		for _, interval := range splittedIntervalsByDays {
			if isSameCategoryAndUser(interval, timeDays[t]) && isWithinDayRange(interval, timeDays[t]) {
				timeDays[t].TimeTotal += int(*interval.EndAt) - int(interval.StartedAt)
			}
		}
	}

	return timeAll, timeDays
}

func getIntervalsRecords(ctx context.Context, resource *db.Resource, t *testing.T) []model.Interval {
	pipeline := mongodb.FormPipeline()
	collection := resource.DB.Collection("Interval")
	records, err := collection.Aggregate(ctx, pipeline)
	assert.NoError(t, err, "Failed to aggregate intervals")

	var decodeItems []model.Interval
	err = records.All(ctx, &decodeItems)
	assert.NoError(t, err, "Failed to decode intervals")

	return decodeItems
}

func getTimeAllRecords(ctx context.Context, resource *db.Resource, t *testing.T) []model.TimeAll {
	collection := resource.DB.Collection("TimeAll")
	records, err := collection.Find(ctx, bson.M{})
	assert.NoError(t, err, "Failed to find TimeAll")

	var res []model.TimeAll
	err = records.All(ctx, &res)
	assert.NoError(t, err, "Failed to decode TimeAll")

	return res
}

func getTimeDayRecords(ctx context.Context, db *db.Resource, t *testing.T) []model.TimeDay {
	c := db.DB.Collection("TimeDay")
	stmt, err := c.Find(ctx, bson.M{})
	assert.NoError(t, err, "Failed to find TimeDay")

	var res []model.TimeDay
	err = stmt.All(ctx, &res)
	assert.NoError(t, err, "Failed to decode TimeDay")
	return res
}
