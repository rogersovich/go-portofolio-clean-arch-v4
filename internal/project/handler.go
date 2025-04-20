package project

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

	project := r.Group("/projects")
	{
		project.GET("", h.GetAll)
		project.GET("/:id", h.GetProjectByIdWithRelations)
		project.POST("/store", h.CreateProject)
		project.POST("/update", h.UpdateProject)
		project.POST("/update-statistic", h.UpdateProjectStatistic)
		project.POST("/delete", h.DeleteProject)
	}
}
