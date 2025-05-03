package testimonial

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

	testimonial := r.Group("/testimonials")
	{
		testimonial.GET("", h.GetAll)
		testimonial.GET("/:id", h.GetTestimonialById)
		testimonial.POST("/store", h.CreateTestimonial)
		testimonial.POST("/update", h.UpdateTestimonial)
		testimonial.POST("/delete", h.DeleteTestimonial)
		testimonial.POST("/change-status", h.ChangeStatusTestimonial)
		testimonial.POST("/bulk-change-status", h.ChangeMultiStatusTestimonial)
	}
}
