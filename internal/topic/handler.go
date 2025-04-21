package topic

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

	topic := r.Group("/topics")
	{
		topic.GET("", h.GetAll)
		topic.GET("/:id", h.GetTopicById)
		topic.POST("/store", h.CreateTopic)
		topic.POST("/update", h.UpdateTopic)
		topic.POST("/delete", h.DeleteTopic)
		topic.POST("/check-has-ids", h.CheckTopicIds)
	}
}
