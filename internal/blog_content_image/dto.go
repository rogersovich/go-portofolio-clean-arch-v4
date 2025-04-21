package blog_content_image

import "github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"

type CreateBlogContentImageRequest struct {
	ImageUrl      string `json:"image_url" validate:"required"`
	ImageFileName string `json:"image_file_name" validate:"required"`
	IsUsed        string `json:"is_used"`
}

type UpdateBlogContentImageRequest struct {
	Id     int    `json:"id" validate:"required"`
	BlogID *int   `json:"blog_id" validate:"required"`
	IsUsed string `json:"is_used"`
}

type UpdateBlogContentImageDTO struct {
	Id            int
	BlogID        *int
	ImageUrl      string
	ImageFileName string
	IsUsed        string
}

type BlogContentImageResponse struct {
	ID            int    `json:"id"`
	BlogID        *int   `json:"blog_id"`
	ImageUrl      string `json:"image_url"`
	ImageFileName string `json:"image_file_name"`
	IsUsed        string `json:"is_used"`
	CreatedAt     string `json:"created_at"`
}

type BlogContentImageUpdateResponse struct {
	ID            int    `json:"id"`
	BlogID        *int   `json:"blog_id"`
	ImageUrl      string `json:"image_url"`
	ImageFileName string `json:"image_file_name"`
	IsUsed        string `json:"is_used"`
}

type BlogContentImageDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

func ToBlogContentImageResponse(p BlogContentImage) BlogContentImageResponse {
	return BlogContentImageResponse{
		ID:            p.ID,
		BlogID:        p.BlogID,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		IsUsed:        utils.BoolToYN(p.IsUsed),
		CreatedAt:     p.CreatedAt.Format("2006-01-02"),
	}
}

func ToBlogContentImageUpdateResponse(p BlogContentImage) BlogContentImageUpdateResponse {
	return BlogContentImageUpdateResponse{
		ID:            p.ID,
		BlogID:        p.BlogID,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		IsUsed:        utils.BoolToYN(p.IsUsed),
	}
}
