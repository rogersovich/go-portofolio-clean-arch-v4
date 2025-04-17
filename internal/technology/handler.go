package technology

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

	project := r.Group("/technologies")
	{
		project.GET("", h.GetAll)
		project.GET("/:id", h.GetTechnologyById)
		project.POST("/store", h.CreateTechnology)
		project.POST("/update", h.UpdateTechnology)
		project.POST("/delete", h.DeleteTechnology)
	}
}
