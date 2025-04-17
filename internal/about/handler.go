package about

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

	about := r.Group("/abouts")
	{
		about.GET("", h.GetAll)
		about.GET("/:id", h.GetAboutById)
		about.POST("/store", h.CreateAbout)
		about.POST("/update", h.UpdateAbout)
		about.POST("/delete", h.DeleteAbout)
	}
}
