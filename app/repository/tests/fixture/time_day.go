package fixture

import (
	"context"
	"log"
	"math/rand"
	"time_app/app/repository/model"
	"time_app/db"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateOneTimeDay(
	resource *db.Resource,
	ctx context.Context,
	user *model.User,
	category *model.Category,
	item *model.TimeDay,
) model.TimeDay {
	if user == nil || category == nil {
		panic("You have to provide correct params")
	}
	if item == nil {
		item = defaultTimeDay(user, category)
	}

	c := resource.DB.Collection("TimeDay")
	stmt, err := c.InsertOne(ctx, item)
	if err != nil {
		log.Printf("Insert TimeDay Error, %v", err)
	}

	filter := bson.M{"_id": stmt.InsertedID}
	created_obj := c.FindOne(ctx, filter)

	var res model.TimeDay
	err = created_obj.Decode(&res)
	if err != nil {
		log.Printf("Decode TimeDay Error, %v", err)
	}
	return res
}

func defaultTimeDay(user *model.User, category *model.Category) *model.TimeDay {
	return &model.TimeDay{
		UUUID:        GenUUID(),
		UserUUID:     user.UUUID,
		CategoryUUID: category.UUUID,
		TimeDay:      GetTimeNow(),
		TimeTotal:    rand.Intn(1000),
	}
}

func CreateManyTimeDay(
	resource *db.Resource,
	ctx context.Context,
	user *model.User,
	categoryList *[]model.Category,
) []model.TimeDay {
	if user == nil {
		panic("You have to provide User")
	}

	itemList := make([]model.TimeDay, 0, len(*categoryList))
	for _, category := range *categoryList {
		obj := CreateOneTimeDay(resource, ctx, user, &category, nil)
		itemList = append(itemList, obj)
	}
	return itemList
}
