package main

import (
	"fmt"
	"time_app/app"

	"github.com/robfig/cron/v3"
)

func main() {

	a := app.RepoInit()

	c := cron.New()
	c.AddFunc("* * * * *", func() {
		err := a.TimeCalculation()
		if err != nil {
			fmt.Printf("Time Calculation Error, %v", err)
		} else {
			fmt.Printf("Task was completed!")
		}
	})

	c.Start()

	select {}

}
