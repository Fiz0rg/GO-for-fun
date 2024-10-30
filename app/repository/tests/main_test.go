package test

import (
	"log"
	"os"
	"testing"
	test_config "time_app/tests"
)

func TestMain(m *testing.M) {
	err := test_config.InitTestDB()
	if err != nil {
		log.Fatalf("Cannot connect to database, %v", err)
	}
	code := m.Run()

	os.Exit(code)

}
