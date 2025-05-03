package testimonial

import "github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"

type CreateTestimonialRequest struct {
	Name      string  `json:"name" binding:"required"`
	Via       *string `json:"via"`
	Role      *string `json:"role"`
	WorkingAt *string `json:"working_at"`
}

type CreateTestimonialDTO struct {
	Name      string  `json:"name"`
	Via       *string `json:"via"`
	Role      *string `json:"role"`
	WorkingAt *string `json:"working_at"`
	IsUsed    bool    `json:"is_used"`
}

type UpdateTestimonialRequest struct {
	ID        int     `json:"id" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Via       *string `json:"via"`
	Role      *string `json:"role"`
	WorkingAt *string `json:"working_at"`
	IsUsed    string  `json:"is_used" binding:"required,oneof=Y N"`
}

type UpdateTestimonialDTO struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Via       *string `json:"via"`
	Role      *string `json:"role"`
	WorkingAt *string `json:"working_at"`
	IsUsed    bool    `json:"is_used"`
}

type TestimonialResponse struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Via       *string `json:"via"`
	Role      *string `json:"role"`
	WorkingAt *string `json:"working_at"`
	IsUsed    string  `json:"is_used"`
	CreatedAt string  `json:"created_at"`
}

type TestimonialDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

type TestimonialChangeStatusRequest struct {
	ID     int    `json:"id" binding:"required"`
	IsUsed string `json:"is_used" binding:"required,oneof=Y N"`
}

type TestimonialChangeMultiStatusRequest struct {
	IDs    []int  `json:"ids" binding:"required"`
	IsUsed string `json:"is_used" binding:"required,oneof=Y N"`
}

func ToTestimonialResponse(p Testimonial) TestimonialResponse {
	return TestimonialResponse{
		ID:        p.ID,
		Name:      p.Name,
		Via:       p.Via,
		Role:      p.Role,
		WorkingAt: p.WorkingAt,
		IsUsed:    utils.BoolToYN(p.IsUsed),
		CreatedAt: p.CreatedAt.Format("2006-01-02"),
	}
}
