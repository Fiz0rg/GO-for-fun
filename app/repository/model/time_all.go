package model

type TimeAll struct {
	UUID         string `bson:"uuid"`
	UserUUID     string `bson:"user_uuid"`
	CategoryUUID string `bson:"category_uuid"`
	TimeTotal    int64  `bson:"time_total"`
}
