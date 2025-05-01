package experience

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

	experience := r.Group("/experiences")
	{
		experience.GET("", h.GetAll)
		experience.GET("/:id", h.GetExperienceById)
		experience.POST("/store", h.CreateExperience)
		experience.POST("/update", h.UpdateExperience)
		experience.POST("/delete", h.DeleteExperience)
	}
}
