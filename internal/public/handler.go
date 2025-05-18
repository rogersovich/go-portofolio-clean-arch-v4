package public

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	service Service
}

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	//* Create author repo & service
	// authorRepo := author.NewRepository(db)
	// authorService := author.NewService(authorRepo)

	publicRepo := NewRepository(db)
	service := NewService(publicRepo)
	h := handler{service: service}

	r.GET("/profile", h.GetProfile)
	r.GET("/blogs", h.GetPublicBlogs)
	r.GET("/blogs/:slug", h.GetPublicBlogBySlug)
	r.GET("/testimonials", h.GetPublicTestimonials)
	r.GET("/topics", h.GetPublicTopics)
	r.GET("/projects", h.GetPublicProjects)
	r.GET("/projects/:slug", h.GetPublicProjectBySlug)
	r.GET("/technologies", h.GetPublicTechnologies)
	r.GET("/authors", h.GetPublicAuthors)
	r.GET("/experiences", h.GetPublicExperiences)
}
