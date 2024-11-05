package test

import (
	"context"
	"testing"
	"time"
	"time_app/app/repository"
	"time_app/app/repository/model"
	"time_app/app/repository/tests/fixture"
	"time_app/db"
	test_config "time_app/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntervalPipeline(t *testing.T) {
	ctx, cancel := repository.InitContext(5 * time.Second)
	defer cancel()
	db := test_config.TestDB
	require.NotNil(t, db, "Database connection should not be nil")

	user := fixture.CreateUser(db, ctx)

	intervalEndANil(db, ctx, user)
	intervalEndAtNotNill(db, ctx, user)
	arrangeInterval := genArrangeIntervalList(db, ctx, user)

	pipelineIntervals := getIntervalsRecords(ctx, db, t)

	assert.Equal(t, pipelineIntervals, arrangeInterval, "Intervals arrays not equal")
	t.Cleanup(func() {
		time.Sleep(10 * time.Microsecond)
		test_config.CleanupTestData(db)
	})
}

func genArrangeIntervalList(db *db.Resource, ctx context.Context, user model.User) []model.Interval {
	var categoryAmount = 2
	var intervalAmount = 4

	categoryList := fixture.CreateManyCategories(db, ctx, &user, &categoryAmount)
	rawIntervalList := fixture.CreateManyIntervals(db, ctx, &user, &categoryList, &intervalAmount)

	filteredIntervals := filterIntervalsEqualPipeline(rawIntervalList)
	return filteredIntervals
}

func intervalEndAtNotNill(db *db.Resource, ctx context.Context, user model.User) {
	category := fixture.CreateOneCategory(db, ctx, &user, nil)
	endAt := int64(fixture.GetTimeNow() + 100)
	obj := model.Interval{
		UUID:         fixture.GenUUID(),
		UserUUID:     user.UUUID,
		CategoryUUID: category.UUUID,
		StartedAt:    fixture.GetTimeNow() - 100,
		EndAt:        &endAt,
	}

	fixture.CreateOneInterval(db, ctx, &user, &category, &obj)
}

func intervalEndANil(db *db.Resource, ctx context.Context, user model.User) {
	category := fixture.CreateOneCategory(db, ctx, &user, nil)
	obj := model.Interval{
		UUID:         fixture.GenUUID(),
		UserUUID:     user.UUUID,
		CategoryUUID: category.UUUID,
		StartedAt:    fixture.GetTimeNow() - 100,
		EndAt:        nil,
	}

	fixture.CreateOneInterval(db, ctx, &user, &category, &obj)
}
