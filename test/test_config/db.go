package test_config

import (
	"time_app/db"
)

var TestDB *db.Resource

func InitTestDB() error {
	var err error
	TestDB, err = TestInitResource()
	return err
}

func GetTestDB() *db.Resource {
	return TestDB
}
