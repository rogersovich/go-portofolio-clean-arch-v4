package about

type CreateAboutRequest struct {
	Title           string `json:"title" validate:"required"`
	DescriptionHTML string `json:"description_html" validate:"required"`
	AvatarUrl       string `json:"avatar_url" validate:"required"`
	AvatarFileName  string `json:"avatar_file_name" validate:"required"`
}

type UpdateAboutRequest struct {
	Id              int    `json:"id" validate:"required,numeric,number"`
	Title           string `json:"title" validate:"required"`
	DescriptionHTML string `json:"description_html" validate:"required"`
}

type UpdateAboutDTO struct {
	Id              uint   `json:"id" validate:"required"`
	Title           string `json:"title" validate:"required"`
	DescriptionHTML string `json:"description_html" validate:"required"`
	AvatarUrl       string `json:"avatar_url"`
	AvatarFileName  string `json:"avatar_file_name"`
}

type AboutResponse struct {
	ID              uint   `json:"id"`
	Title           string `json:"title"`
	DescriptionHTML string `json:"description_html" validate:"required"`
	AvatarUrl       string `json:"avatar_url"`
	AvatarFileName  string `json:"avatar_file_name"`
	CreatedAt       string `json:"created_at"`
}

type AboutUpdateResponse struct {
	ID              uint   `json:"id"`
	Title           string `json:"title"`
	DescriptionHTML string `json:"description_html" validate:"required"`
	AvatarUrl       string `json:"avatar_url"`
	AvatarFileName  string `json:"avatar_file_name"`
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

func ToAboutUpdateResponse(p About) AboutUpdateResponse {
	return AboutUpdateResponse{
		ID:              p.ID,
		Title:           p.Title,
		DescriptionHTML: p.DescriptionHTML,
		AvatarUrl:       p.AvatarUrl,
		AvatarFileName:  p.AvatarFileName,
	}
}
