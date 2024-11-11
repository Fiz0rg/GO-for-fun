package mongodb

import (
	"log"
	"sync"
	"time"
	"time_app/app/repository"
	"time_app/db"

	"go.mongodb.org/mongo-driver/mongo"
)

type CountTimeRepositoryImpl struct {
	resource           *db.Resource
	intervalCollection *mongo.Collection
	timeAllCollection  *mongo.Collection
	timeDayCollection  *mongo.Collection
}

type CountTimeRepository interface {
	TimeCalculation() error
}

func NewCountTimeRepository(resource *db.Resource) CountTimeRepository {
	timeDayCollection := resource.DB.Collection("TimeDay")
	intervalCollection := resource.DB.Collection("Interval")
	timeAllCollection := resource.DB.Collection("TimeAll")
	countTimeRepo := &CountTimeRepositoryImpl{
		resource:           resource,
		intervalCollection: intervalCollection,
		timeAllCollection:  timeAllCollection,
		timeDayCollection:  timeDayCollection,
	}
	return countTimeRepo
}

func (repo *CountTimeRepositoryImpl) TimeCalculation() error {
	ctx, cancel := repository.InitContext(1 * time.Second)
	defer cancel()

	intervals := repo.getIntervalRecords(ctx)

	var wgUpdate sync.WaitGroup
	updateTimeAllCollection(ctx, &wgUpdate, repo.timeAllCollection, intervals)
	updateTimeDayCollection(ctx, &wgUpdate, repo.timeDayCollection, intervals)
	wgUpdate.Wait()

	var wgDelete sync.WaitGroup
	d := deleteUnnecessaryIntervals(ctx, &wgDelete, repo.intervalCollection, intervals)
	if d != nil {
		log.Panicf("Fail operation by deleteUnnecessaryIntervals, %v", d)
	}
	wgDelete.Wait()

	return nil
}
