package about

import "mime/multipart"

type CreateAboutRequest struct {
	Title           string                `json:"title" validate:"required"`
	DescriptionHTML string                `json:"description_html" validate:"required"`
	AvatarFile      *multipart.FileHeader `json:"avatar_file"`
}

type UpdateAboutRequest struct {
	ID              int                   `json:"id" validate:"required,numeric,number"`
	Title           string                `json:"title" validate:"required"`
	DescriptionHTML string                `json:"description_html" validate:"required"`
	AvatarFile      *multipart.FileHeader `json:"avatar_file"`
}

type CreateAboutDTO struct {
	Title           string
	DescriptionHTML string
	AvatarUrl       string
	AvatarFileName  string
}

type UpdateAboutDTO struct {
	ID              int
	Title           string
	DescriptionHTML string
	AvatarUrl       string
	AvatarFileName  string
}

type AboutResponse struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	DescriptionHTML string `json:"description_html" validate:"required"`
	AvatarUrl       string `json:"avatar_url"`
	AvatarFileName  string `json:"avatar_file_name"`
	CreatedAt       string `json:"created_at"`
}

type AboutDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

func ToAboutResponse(p About) AboutResponse {
	return AboutResponse{
		ID:              p.ID,
		Title:           p.Title,
		DescriptionHTML: p.DescriptionHTML,
		AvatarUrl:       p.AvatarUrl,
		AvatarFileName:  p.AvatarFileName,
		CreatedAt:       p.CreatedAt.Format("2006-01-02"),
	}
}
