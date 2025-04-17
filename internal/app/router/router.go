package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/about"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/author"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/project_content_image"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/statistic"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/technology"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(utils.RecoveryWithLogger())
	r.Use(utils.LoggerMiddleware())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"data":    nil,
			"message": "route not found",
		})
	})

	api := r.Group("/api")
	author.RegisterRoutes(api, db)
	about.RegisterRoutes(api, db)
	technology.RegisterRoutes(api, db)
	statistic.RegisterRoutes(api, db)
	project_content_image.RegisterRoutes(api, db)

	return r
}
