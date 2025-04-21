package blog_topic

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

	blog_topic := r.Group("/blog-topics")
	{
		blog_topic.GET("", h.GetAll)
		blog_topic.GET("/:id", h.GetBlogTopicById)
		blog_topic.POST("/store", h.CreateBlogTopic)
		blog_topic.POST("/update", h.UpdateBlogTopic)
		blog_topic.POST("/delete", h.DeleteBlogTopic)
	}
}
