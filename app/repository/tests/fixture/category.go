package fixture

import (
	"context"
	"log"
	"time_app/app/repository/model"
	"time_app/db"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateOneCategory(
	resource *db.Resource,
	ctx context.Context,
	user *model.User,
	category *model.Category,
) model.Category {
	if user == nil {
		panic("User not provided")
	}

	if category == nil {
		pos := int(1)
		obj := model.Category{
			UUUID:    GenUUID(),
			UserUUID: user.NAME,
			NAME:     genName(),
			ICON:     genString(),
			COLOR:    genString(),
			ACTIVE:   true,
			POSITION: &pos,
		}
		category = &obj
	}

	c := resource.DB.Collection("Category")
	stmt, err := c.InsertOne(ctx, category)
	if err != nil {
		log.Printf("Insert Category Error, %v", err)
	}

	filter := bson.M{"_id": stmt.InsertedID}
	created_obj := c.FindOne(ctx, filter)

	var res model.Category
	err = created_obj.Decode(&res)
	if err != nil {
		log.Printf("Decode User Error, %v", err)
	}
	return res
}

func CreateManyCategories(
	resource *db.Resource,
	ctx context.Context,
	user *model.User,
	amount *int,
) []model.Category {
	if user == nil {
		panic("User not provided")
	}
	if amount == nil {
		a := 4
		amount = &a
	}
	itemList := make([]model.Category, 0, *amount)

	for i := 0; i < *amount; i++ {
		obj := CreateOneCategory(resource, ctx, user, nil)
		itemList = append(itemList, obj)
	}

	return itemList

}
