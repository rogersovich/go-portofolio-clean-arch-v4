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

	project := r.Group("/abouts")
	{
		project.GET("", h.GetAll)
		project.GET("/:id", h.GetAboutById)
		project.POST("/store", h.CreateAbout)
		project.POST("/update", h.UpdateAbout)
		project.POST("/delete", h.DeleteAbout)
	}
}
