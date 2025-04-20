package project_content_temp_image

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

	project_content_temp_image := r.Group("/project-content-temp-images")
	{
		project_content_temp_image.GET("", h.GetAll)
		project_content_temp_image.GET("/:id", h.GetProjectContentTempImgById)
		project_content_temp_image.POST("/store", h.CreateProjectContentTempImg)
		project_content_temp_image.POST("/update", h.UpdateProjectContentTempImg)
		project_content_temp_image.POST("/delete", h.DeleteProjectContentTempImg)
	}
}
