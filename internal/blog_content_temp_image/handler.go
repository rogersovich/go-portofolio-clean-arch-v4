package blog_content_temp_image

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

	blog_content_temp_image := r.Group("/blog-content-temp-images")
	{
		blog_content_temp_image.GET("", h.GetAll)
		blog_content_temp_image.GET("/:id", h.GetBlogContentTempImgById)
		blog_content_temp_image.POST("/store", h.CreateBlogContentTempImg)
		blog_content_temp_image.POST("/update", h.UpdateBlogContentTempImg)
		blog_content_temp_image.POST("/delete", h.DeleteBlogContentTempImg)
	}
}
