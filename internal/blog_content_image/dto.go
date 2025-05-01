package blog_content_image

type CreateBlogContentImageRequest struct {
	ImageUrl      string `json:"image_url" validate:"required"`
	ImageFileName string `json:"image_file_name" validate:"required"`
}

type UpdateBlogContentImageRequest struct {
	Id     int  `json:"id" validate:"required"`
	BlogID *int `json:"blog_id" validate:"required"`
}

type UpdateBlogContentImageDTO struct {
	Id            int
	BlogID        *int
	ImageUrl      string
	ImageFileName string
}

type BlogContentImageResponse struct {
	ID            int    `json:"id"`
	BlogID        *int   `json:"blog_id"`
	ImageUrl      string `json:"image_url"`
	ImageFileName string `json:"image_file_name"`
	CreatedAt     string `json:"created_at"`
}

type BlogContentImageUpdateResponse struct {
	ID            int    `json:"id"`
	BlogID        *int   `json:"blog_id"`
	ImageUrl      string `json:"image_url"`
	ImageFileName string `json:"image_file_name"`
}

type BlogContentImageDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

type BlogContentImageExistingResponse struct {
	ID       int    `json:"id"`
	BlogID   *int   `json:"blog_id"`
	ImageUrl string `json:"image_url"`
}

type BlogContentImageBulkUpdateDTO struct {
	ImageUrls []string `json:"image_urls"`
	BlogID    int      `json:"blog_id"`
}

func ToBlogContentImageResponse(p BlogContentImage) BlogContentImageResponse {
	return BlogContentImageResponse{
		ID:            p.ID,
		BlogID:        p.BlogID,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		CreatedAt:     p.CreatedAt.Format("2006-01-02"),
	}
}

func ToBlogContentImageUpdateResponse(p BlogContentImage) BlogContentImageUpdateResponse {
	return BlogContentImageUpdateResponse{
		ID:            p.ID,
		BlogID:        p.BlogID,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
	}
}
