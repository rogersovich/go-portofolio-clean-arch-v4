package project_technology

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

	project_technologies := r.Group("/project-technologies")
	{
		project_technologies.GET("", h.GetAll)
		project_technologies.GET("/:id", h.GetProjectTechnologyById)
		project_technologies.POST("/store", h.CreateProjectTechnology)
		project_technologies.POST("/update", h.UpdateProjectTechnology)
		project_technologies.POST("/delete", h.DeleteProjectTechnology)
	}
}
