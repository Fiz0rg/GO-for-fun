package main

import (
	"log"
	"time_app/app"
	"time_app/config"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Printf("Something wrong with config, %v", err)
	}

	app.StartGin(config)
}
