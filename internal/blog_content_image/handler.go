package blog_content_image

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

	blog_content_image := r.Group("/blog-content-images")
	{
		blog_content_image.GET("", h.GetAll)
		blog_content_image.GET("/:id", h.GetBlogContentImageById)
		blog_content_image.POST("/store", h.CreateBlogContentImage)
		blog_content_image.POST("/update", h.UpdateBlogContentImage)
		blog_content_image.POST("/delete", h.DeleteBlogContentImage)
	}
}
