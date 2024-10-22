package app

import (
	"time_app/app/api"
	"time_app/config"

	"time_app/db"
	_ "time_app/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartGin(config *config.Config) {

	r := gin.Default()
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	publicRoute := r.Group("/total_time")
	resource, err := db.InitResource(config)
	if err != nil {
		panic(err)
	}

	api.ApplyCountTimeAPI(publicRoute, resource)
	r.Routes()
	r.Run()
}
