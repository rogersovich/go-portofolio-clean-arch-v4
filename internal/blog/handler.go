package blog

import (
	"github.com/gin-gonic/gin"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/author"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/blog_content_image"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/blog_content_temp_image"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/blog_topic"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/reading_time"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/statistic"
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

	//* Create statistic repo & service
	statisticRepo := statistic.NewRepository(db)
	statisticService := statistic.NewService(statisticRepo)

	//* Create readingTime repo & service
	readingTimeRepo := reading_time.NewRepository(db)
	readingTimeService := reading_time.NewService(readingTimeRepo)

	//* Create blogTopic repo & service
	blogTopicRepo := blog_topic.NewRepository(db)
	blogTopicService := blog_topic.NewService(blogTopicRepo)

	//* Create blogContentImage repo & service
	blogContentImageRepo := blog_content_image.NewRepository(db)
	blogContentImageService := blog_content_image.NewService(blogContentImageRepo)

	//* Create blogContentTempImage repo & service
	blogContentTempImageRepo := blog_content_temp_image.NewRepository(db)
	blogContentTempImageService := blog_content_temp_image.NewService(blogContentTempImageRepo)

	blogRepo := NewRepository(db)
	blogService := NewService(
		authorService,
		topicService,
		statisticService,
		readingTimeService,
		blogTopicService,
		blogContentImageService,
		blogContentTempImageService,
		blogRepo, db)
	h := handler{service: blogService}

	blog := r.Group("/blogs")
	{
		blog.GET("", h.GetAll)
		blog.GET("/:id", h.GetBlogByIdWithRelations)
		blog.POST("/store", h.CreateBlog)
		blog.POST("/update", h.UpdateBlog)
		// blog.POST("/delete", h.DeleteBlog)
	}
}
