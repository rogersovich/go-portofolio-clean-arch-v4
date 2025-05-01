package testimonial

type CreateTestimonialRequest struct {
	Name      string  `json:"name" binding:"required"`
	Via       *string `json:"via"`
	Role      *string `json:"role"`
	WorkingAt *string `json:"working_at"`
}

type UpdateTestimonialRequest struct {
	ID        int     `json:"id" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Via       *string `json:"via"`
	Role      *string `json:"role"`
	WorkingAt *string `json:"working_at"`
}

type TestimonialResponse struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Via       *string `json:"via"`
	Role      *string `json:"role"`
	WorkingAt *string `json:"working_at"`
	CreatedAt string  `json:"created_at"`
}

type TestimonialDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

func ToTestimonialResponse(p Testimonial) TestimonialResponse {
	return TestimonialResponse{
		ID:        p.ID,
		Name:      p.Name,
		Via:       p.Via,
		Role:      p.Role,
		WorkingAt: p.WorkingAt,
		CreatedAt: p.CreatedAt.Format("2006-01-02"),
	}
}
