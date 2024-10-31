package app

import (
	"log"
	"time_app/app/repository/mongodb"
	"time_app/config"
	"time_app/db"
)

func RepoInit() mongodb.UpdateTimeAllCollectionRepository {
	config, err := config.LoadConfig()
	if err != nil {
		log.Printf("Something wrong with config, %v", err)
	}

	resource, err := db.InitResource(config)
	if err != nil {
		panic(err)
	}

	r := mongodb.NewCountTimeRepository(resource)

	return r
}
