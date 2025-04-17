package project_content_image

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

	project_content_image := r.Group("/project-content-images")
	{
		project_content_image.GET("", h.GetAll)
		project_content_image.GET("/:id", h.GetProjectContentImageById)
		project_content_image.POST("/store", h.CreateProjectContentImage)
		project_content_image.POST("/update", h.UpdateProjectContentImage)
		project_content_image.POST("/delete", h.DeleteProjectContentImage)
	}
}
