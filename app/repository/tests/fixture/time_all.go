package fixture

import (
	"context"
	"log"
	"math/rand"
	"time_app/app/repository/model"
	"time_app/db"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateOneTimeAll(
	resource *db.Resource,
	ctx context.Context,
	category *model.Category,
	user *model.User,
	obj *model.TimeAll,
) model.TimeAll {

	if user == nil || category == nil {
		panic("You have to provide correct params")
	}

	c := resource.DB.Collection("TimeAll")
	if obj == nil {
		obj = defaultTimeAll(user, category)

	}
	stmt, err := c.InsertOne(ctx, obj)
	if err != nil {
		log.Printf("Insert TimeAll ERROR, %v", err)
	}

	fltr := bson.M{"_id": stmt.InsertedID}
	created_obj := c.FindOne(ctx, fltr)

	var res model.TimeAll
	err = created_obj.Decode(&res)
	if err != nil {
		log.Printf("Decode TimeAll ERROR, %v", err)
	}
	return res
}

func defaultTimeAll(user *model.User, category *model.Category) *model.TimeAll {
	t := model.TimeAll{
		UUID:         GenUUID(),
		UserUUID:     user.UUUID,
		CategoryUUID: category.UUUID,
		TimeTotal:    int64(rand.Intn(1000)),
	}
	return &t
}

func CreateManyTimeAll(
	resource *db.Resource,
	ctx context.Context,
	user *model.User,
	categoryList *[]model.Category,
) []model.TimeAll {
	if user == nil {
		panic("You have to provide User")
	}
	itemList := make([]model.TimeAll, 0, len(*categoryList))
	for _, category := range *categoryList {
		obj := CreateOneTimeAll(resource, ctx, &category, user, nil)
		itemList = append(itemList, obj)
	}
	return itemList
}
