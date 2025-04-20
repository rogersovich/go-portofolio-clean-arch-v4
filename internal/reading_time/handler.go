package reading_time

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

	reading_time := r.Group("/reading-times")
	{
		reading_time.GET("", h.GetAll)
		reading_time.GET("/:id", h.GetReadingTimeById)
		reading_time.POST("/store", h.CreateReadingTime)
		reading_time.POST("/update", h.UpdateReadingTime)
		reading_time.POST("/delete", h.DeleteReadingTime)
	}
}
