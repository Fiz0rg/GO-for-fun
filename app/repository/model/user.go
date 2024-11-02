package model

type User struct {
	UUUID string `bson:"UUID"`
	NAME  string `bson:"NAME"`
	EMAIL string `bson:"EMAIL"`
}
