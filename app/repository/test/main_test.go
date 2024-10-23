package test

import (
	"log"
	"os"
	"testing"
	"time_app/test/test_config"
)

func TestMain(m *testing.M) {
	err := test_config.InitTestDB()
	if err != nil {
		log.Fatalf("Cannot connect to database, %v", err)
	}
	code := m.Run()

	os.Exit(code)

}
