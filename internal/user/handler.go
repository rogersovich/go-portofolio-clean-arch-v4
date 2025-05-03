package user

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

	user := r.Group("/users")
	{
		user.GET("", h.GetAll)
		user.GET("/:id", h.GetUserById)
		user.POST("/update", h.UpdateUser)
		user.POST("/delete", h.DeleteUser)
	}
}
