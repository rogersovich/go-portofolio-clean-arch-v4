package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/author"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"data":    nil,
			"message": "route not found",
		})
	})

	api := r.Group("/api")
	author.RegisterRoutes(api, db)
	return r
}
