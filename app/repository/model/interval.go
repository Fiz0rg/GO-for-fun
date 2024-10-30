package model

type Interval struct {
	UUID         string `bson:"uuid"`
	UserUUID     string `bson:"user_uuid"`
	CategoryUUID string `bson:"category_uuid"`
	StartedAt    int64  `bson:"started_at"`
	EndAt        *int64 `bson:"end_at,omitempty"`
}

// type IntervalPart struct {
// 	UUID               string `bson:"uuid"`
// 	UserUUID           string `bson:"user_uuid"`
// 	TimeTotal          int    `bson:"time_total"`
// 	GreatThanEqualDate int
// 	LessThanEqualDate  int
// }

type UserCategory struct {
	UserUUID          string `bson:"user_uuid" json:"user_uuid"`
	CategoryUUID      string `bson:"category_uuid" json:"category_uuid"`
	TotalIntervalTime int64  `bson:"total_interval_time"`
}
