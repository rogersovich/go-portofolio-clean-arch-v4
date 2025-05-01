package project

import (
	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/project_content_image"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/project_technology"
	"gorm.io/gorm"
)

type handler struct {
	service Service
}

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {

	//* Create project_technology repo & service
	projectTechRepo := project_technology.NewRepository(db)
	projectTechService := project_technology.NewService(projectTechRepo)

	//* Create project_content_image repo & service
	projectImagesRepo := project_content_image.NewRepository(db)
	projectImagesService := project_content_image.NewService(projectImagesRepo)

	projectRepo := NewRepository(db)
	projectService := NewService(projectTechService, projectImagesService, projectRepo)

	h := handler{service: projectService}

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
