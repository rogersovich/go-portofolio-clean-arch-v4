package blog

import (
	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/author"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/topic"
	"gorm.io/gorm"
)

type handler struct {
	service Service
}

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	//* Create author repo & service
	authorRepo := author.NewRepository(db)
	authorService := author.NewService(authorRepo)

	//* Create topic repo & service
	topicRepo := topic.NewRepository(db)
	topicService := topic.NewService(topicRepo)

	blogRepo := NewRepository(db)
	blogService := NewService(authorService, topicService, blogRepo)
	h := handler{service: blogService}

	blog := r.Group("/blogs")
	{
		blog.GET("", h.GetAll)
		// blog.GET("/:id", h.GetBlogByIdWithRelations)
		blog.POST("/store", h.CreateBlog)
		// blog.POST("/update", h.UpdateBlog)
		// blog.POST("/update-statistic", h.UpdateBlogStatistic)
		// blog.POST("/delete", h.DeleteBlog)
	}
}
