package model

type Category struct {
	UUUID    string `bson:"uuid"`
	UserUUID string `bson:"user_uuid"`
	NAME     string `bson:"NAME"`
	ICON     string `bson:"ICON"`
	COLOR    string `bson:"COLOR"`
	POSITION *int   `bson:"POSITION"`
	ACTIVE   bool   `bson:"ACTIVE"`
}
