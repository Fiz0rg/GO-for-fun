package model

type TimeAll struct {
	UUID         string `bson:"uuid"`
	UserUUID     string `bson:"user_uuid"`
	CategoryUUID string `bson:"category_uuid"`
	TimeTotal    int    `bson:"time_total"`
}
