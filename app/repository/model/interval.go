package model

type Interval struct {
	UUID         string `bson:"uuid"`
	UserUUID     string `bson:"user_uuid"`
	CategoryUUID string `bson:"category_uuid"`
	StartedAt    int64  `bson:"started_at"`
	EndAt        *int64 `bson:"end_at,omitempty"`
}

type UserCategory struct {
	UserUUID          string `bson:"user_uuid" json:"user_uuid"`
	CategoryUUID      string `bson:"category_uuid" json:"category_uuid"`
	TotalIntervalTime int64  `bson:"total_interval_time"`
}
