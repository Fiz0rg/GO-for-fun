package model

type TimeDay struct {
	UUUID        string `bson:"uuid"`
	UserUUID     string `bson:"user_uuid"`
	CategoryUUID string `bson:"category_uuid"`
	TimeDay      int64  `bson:"time_day"`
	TimeTotal    int    `bson:"time_total"`
}
