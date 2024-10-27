package app

import (
	"log"
	"time_app/app/repository"
	"time_app/config"
	"time_app/db"
)

func RepoInit() repository.UpdateTimeAllCollectionRepository {
	config, err := config.LoadConfig()
	if err != nil {
		log.Printf("Something wrong with config, %v", err)
	}

	resource, err := db.InitResource(config)
	if err != nil {
		panic(err)
	}

	r := repository.NewCountTimeRepository(resource)

	return r
}
