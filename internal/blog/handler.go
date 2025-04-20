package blog

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

	blog := r.Group("/blogs")
	{
		blog.GET("", h.GetAll)
		// blog.GET("/:id", h.GetBlogByIdWithRelations)
		blog.POST("/store", h.CreateBlog)
		// blog.POST("/update", h.UpdateBlog)
		// blog.POST("/update-statistic", h.UpdateBlogStatistic)
		// blog.POST("/delete", h.DeleteBlog)
	}
}
