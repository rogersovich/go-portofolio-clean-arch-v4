package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/about"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/auth"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/author"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/blog"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/blog_content_image"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/blog_topic"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/experience"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/project"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/project_content_image"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/project_technology"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/public"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/reading_time"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/statistic"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/technology"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/testimonial"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/topic"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/user"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(utils.RecoveryWithLogger())
	r.Use(utils.LoggerMiddleware())

	// Configure CORS options
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"http://43.134.162.211:3000",
		"https://www.dimasroger.com",
		"https://dimasroger.com",
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true

	// Apply CORS middleware
	r.Use(cors.New(corsConfig))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"data":    nil,
			"message": "route not found",
		})
	})

	api := r.Group("/api")
	{
		auth.RegisterRoutes(api, db)

		// Apply JWT middleware to other routes
		api.Use(utils.JWTMiddleware()) // Protect all subsequent routes

		user.RegisterRoutes(api, db)
		author.RegisterRoutes(api, db)
		about.RegisterRoutes(api, db)
		technology.RegisterRoutes(api, db)
		statistic.RegisterRoutes(api, db)
		project_content_image.RegisterRoutes(api, db)
		project_technology.RegisterRoutes(api, db)
		project.RegisterRoutes(api, db)
		topic.RegisterRoutes(api, db)
		reading_time.RegisterRoutes(api, db)
		blog.RegisterRoutes(api, db)
		blog_topic.RegisterRoutes(api, db)
		blog_content_image.RegisterRoutes(api, db)
		experience.RegisterRoutes(api, db)
		testimonial.RegisterRoutes(api, db)
	}

	// Define the public API group
	apiPublic := r.Group("/api-public")
	{
		public.RegisterRoutes(apiPublic, db)
	}

	return r
}
