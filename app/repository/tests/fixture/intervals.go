package fixture

import (
	"context"
	"log"
	"math/rand"
	"time_app/app/repository/model"
	"time_app/db"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateOneInterval(
	resource *db.Resource,
	ctx context.Context,
	user *model.User,
	category *model.Category,
	i *model.Interval,
) model.Interval {

	if user == nil || category == nil {
		panic("You have to provide correct params")
	}

	if i == nil {
		i = defaultInterval(user, category)
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

func defaultInterval(user *model.User, category *model.Category) *model.Interval {
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
		UUID:         GenUUID(),
		UserUUID:     user.UUUID,
		CategoryUUID: category.UUUID,
		StartedAt:    timeNow - randmomInt,
		EndAt:        endAtTime,
	}
}

func CreateManyIntervals(
	resource *db.Resource,
	ctx context.Context,
	user *model.User,
	categoryList *[]model.Category,
	amount *int,
) []model.Interval {
	if amount == nil {
		a := 4
		amount = &a
	}

	if user == nil {
		panic("You have to provide User")
	}

	itemList := make([]model.Interval, 0, (len(*categoryList) * *amount))

	for _, category := range *categoryList {
		for i := 0; i < *amount; i++ {
			obj := CreateOneInterval(resource, ctx, user, &category, nil)
			itemList = append(itemList, obj)
		}
	}
	return itemList
}
