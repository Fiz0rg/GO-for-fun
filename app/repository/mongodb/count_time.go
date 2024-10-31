package mongodb

import (
	"log"
	"sync"
	"time"
	"time_app/app/repository"
	"time_app/db"

	"go.mongodb.org/mongo-driver/mongo"
)

var CountTimeRepo UpdateTimeAllCollectionRepository

type UpdateTimeRepositoryImpl struct {
	resource           *db.Resource
	intervalCollection *mongo.Collection
	timeAllCollection  *mongo.Collection
	timeDayCollection  *mongo.Collection
}

type UpdateTimeAllCollectionRepository interface {
	TimeCalculation() error
}

func NewCountTimeRepository(resource *db.Resource) UpdateTimeAllCollectionRepository {
	timeDayCollection := resource.DB.Collection("TimeDay")
	intervalCollection := resource.DB.Collection("Interval")
	timeAllCollection := resource.DB.Collection("TimeAll")
	countTimeRepo := &UpdateTimeRepositoryImpl{
		resource:           resource,
		intervalCollection: intervalCollection,
		timeAllCollection:  timeAllCollection,
		timeDayCollection:  timeDayCollection,
	}
	return countTimeRepo
}

func (repo *UpdateTimeRepositoryImpl) TimeCalculation() error {
	ctx, cancel := repository.InitContext(1 * time.Second)
	defer cancel()
	var wg sync.WaitGroup

	intervals := repo.getIntervalRecords(ctx)
	updateTimeAllCollection(ctx, &wg, repo.timeAllCollection, intervals)
	updateTimeDayCollection(ctx, &wg, repo.timeDayCollection, intervals)

	d := deleteUnnecessaryIntervals(ctx, &wg, repo.intervalCollection, intervals)
	if d != nil {
		log.Panicf("Fail operation by deleteUnnecessaryIntervals, %v", d)
	}

	wg.Wait()
	return nil
}
