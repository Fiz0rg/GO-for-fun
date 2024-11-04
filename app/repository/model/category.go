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

type UserCategory struct {
	UserUUID          string `bson:"user_uuid" json:"user_uuid"`
	CategoryUUID      string `bson:"category_uuid" json:"category_uuid"`
	TotalIntervalTime int64  `bson:"total_interval_time"`
}
