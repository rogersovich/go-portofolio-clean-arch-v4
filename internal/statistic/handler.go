package statistic

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	service Service
}

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	h := handler{service: service}

	statistic := r.Group("/statistics")
	{
		statistic.GET("", h.GetAll)
		statistic.GET("/:id", h.GetStatisticById)
		statistic.POST("/store", h.CreateStatistic)
		statistic.POST("/update", h.UpdateStatistic)
		statistic.POST("/delete", h.DeleteStatistic)
	}
}
