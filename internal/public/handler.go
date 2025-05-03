package public

import (
	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/author"
	"gorm.io/gorm"
)

type handler struct {
	service Service
}

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	//* Create author repo & service
	authorRepo := author.NewRepository(db)
	authorService := author.NewService(authorRepo)

	publicRepo := NewRepository(db)
	service := NewService(authorService, publicRepo)
	h := handler{service: service}

	r.GET("/profile", h.GetProfile)
	r.GET("/blogs", h.GetPublicBlogs)
	r.GET("/blogs/:slug", h.GetPublicBlogBySlug)
}
