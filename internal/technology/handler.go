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

	technology := r.Group("/technologies")
	{
		technology.GET("", h.GetAll)
		technology.GET("/:id", h.GetTechnologyById)
		technology.POST("/store", h.CreateTechnology)
		technology.POST("/update", h.UpdateTechnology)
		technology.POST("/delete", h.DeleteTechnology)
	}
}
