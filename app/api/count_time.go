package api

import (
	"log"
	"net/http"
	"time_app/app/repository"
	"time_app/db"

	"github.com/gin-gonic/gin"
)

func ApplyCountTimeAPI(app *gin.RouterGroup, resource *db.Resource) {
	countTimeRepo := repository.NewCountTimeRepository(resource)
	app.GET("/count", updateTimeAllCollection(countTimeRepo))
}

// @Summary Get total count time
// @Description Получение всех интервалов
// @Tags Count total time
// @Produce json
// @Success 200 {list} []model.Interval
// @Failure 500 {object} error
// @Router /total_time/count [get]
func updateTimeAllCollection(intervalRepo repository.UpdateTimeAllCollectionRepository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		code, err := intervalRepo.UpdateTimeAll()

		if err != nil {
			log.Printf("Error update time all repo: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		responce := map[string]interface{}{
			"code": code,
			"err":  err,
		}

		ctx.JSON(code, responce)
	}

}
