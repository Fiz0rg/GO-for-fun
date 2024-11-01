package fixture

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time_app/app/repository/model"
	"time_app/db"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateOneInterval(
	resource *db.Resource,
	ctx context.Context,
	i *model.Interval,
) model.Interval {
	if i == nil {
		i = defaultInterval()
	}
	c := resource.DB.Collection("Interval")
	insert, err := c.InsertOne(ctx, i)
	if err != nil {
		log.Panicf("Cant create interval for test, %v", err)
	}
	filter := bson.M{"_id": insert.InsertedID}
	stmt := c.FindOne(ctx, filter)

	var res model.Interval
	err = stmt.Decode(&res)
	if err != nil {
		log.Printf("Interval DECODE error, %v", err)
	}
	return res
}

func defaultInterval() *model.Interval {
	randmomInt := int64(rand.Intn(1000))
	timeNow := GetTimeNow()

	endAtTime := func() *int64 {
		if rand.Intn(2) == 0 {
			value := timeNow + randmomInt
			return &value
		}
		return nil
	}()

	return &model.Interval{
		UUID:         genUUID(),
		UserUUID:     "user-1",
		CategoryUUID: "category-1",
		StartedAt:    timeNow - randmomInt,
		EndAt:        endAtTime,
	}
}

func CreateManyIntervals(
	resource *db.Resource,
	ctx context.Context,
	amount *int,
	itemList *[]model.Interval,
) []model.Interval {
	if amount == nil {
		a := 4
		amount = &a
	}
	if itemList == nil {
		emtpyList := make([]model.Interval, 0, *amount)
		itemList = &emtpyList
		for i := 0; i < *amount; i++ {
			item := defaultInterval()
			*itemList = append(*itemList, *item)
		}
	}

	docs := make([]interface{}, len(*itemList))
	for i, v := range *itemList {
		docs[i] = v
	}

	c := resource.DB.Collection("Interval")
	_, err := c.InsertMany(ctx, docs)
	if err != nil {
		fmt.Printf("InsertMany intervals ERROR, %v", err)
	}

	stmt, err := c.Find(ctx, bson.D{})
	if err != nil {
		fmt.Printf("Find (ALL) intervals ERROR, %v", err)
	}

	var res []model.Interval
	err = stmt.All(ctx, &res)
	if err != nil {
		fmt.Printf("Decode ERROR intervals, %v", err)
	}

	return res
}

// func GenIntervalsData(resource *db.Resource, ctx context.Context) ([]model.Interval, mongo.Collection) {
// 	endAt := int64(500)
// 	intervalCollections := resource.DB.Collection("Interval")
// 	intervalsInPipline := []model.Interval{
// 		{
// 			UUID:         "uuid-1",
// 			UserUUID:     "user-1",
// 			CategoryUUID: "category-001",
// 			StartedAt:    400,
// 			EndAt:        &endAt,
// 		},
// 	}

// 	intervalsOutPipline := []model.Interval{
// 		{
// 			UUID:         "uuid-2",
// 			UserUUID:     "user-1",
// 			CategoryUUID: "category-001",
// 			StartedAt:    400,
// 			EndAt:        &endAt,
// 		},
// 		{
// 			UUID:         "uuid-3",
// 			UserUUID:     "user-1",
// 			CategoryUUID: "category-002",
// 			StartedAt:    400,
// 			EndAt:        &endAt,
// 		},
// 		{
// 			UUID:         "uuid-4",
// 			UserUUID:     "user-1",
// 			CategoryUUID: "category-003",
// 			StartedAt:    400,
// 			EndAt:        nil,
// 		},
// 	}
// 	allIntervals := append(intervalsInPipline, intervalsOutPipline...)

// 	// Перед отправкой в insertMany без этого переобразования нельзя
// 	docs := make([]interface{}, len(allIntervals))
// 	for i, v := range allIntervals {
// 		docs[i] = v
// 	}

// 	_, err := intervalCollections.InsertMany(ctx, docs)
// 	if err != nil {
// 		log.Fatalf("Intervals not inserted, %v", err)
// 	}

// 	return intervalsInPipline, *intervalCollections
// }
