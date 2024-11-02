package fixture

import (
	"context"
	"log"
	"time_app/app/repository/model"
	"time_app/db"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateUser(
	resource *db.Resource,
	ctx context.Context,
) model.User {
	c := resource.DB.Collection("User")
	obj := model.User{
		UUUID: genUUID(),
		NAME:  genUUID(),
		EMAIL: "some@.com",
	}
	stmt, err := c.InsertOne(ctx, obj)
	if err != nil {
		log.Printf("Insert User ERROR, %v", err)
	}
	filter := bson.M{"_id": stmt.InsertedID}
	created_obj := c.FindOne(ctx, filter)

	var res model.User
	err = created_obj.Decode(&res)
	if err != nil {
		log.Printf("Decode User Error, %v", err)
	}
	return res
}
