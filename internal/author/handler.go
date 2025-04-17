package author

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

	author := r.Group("/authors")
	{
		author.GET("", h.GetAll)
		author.GET("/:id", h.GetAuthorById)
		author.POST("/store", h.CreateAuthor)
		author.POST("/update", h.UpdateAuthor)
		author.POST("/delete", h.DeleteAuthor)
	}
}
