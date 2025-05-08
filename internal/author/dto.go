package author

import "mime/multipart"

type CreateAuthorRequest struct {
	Name       string                `json:"name" validate:"required"`
	AvatarFile *multipart.FileHeader `json:"avatar_file"`
}

type UpdateAuthorRequest struct {
	ID         int                   `json:"id" validate:"required,numeric,number"`
	Name       string                `json:"name" validate:"required"`
	AvatarFile *multipart.FileHeader `json:"avatar_file"`
}

type CreateAuthorDTO struct {
	Name           string
	AvatarUrl      string
	AvatarFileName string
}

type UpdateAuthorDTO struct {
	ID             int
	Name           string
	AvatarUrl      string
	AvatarFileName string
}

type AuthorResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	AvatarUrl      string `json:"avatar_url"`
	AvatarFileName string `json:"avatar_file_name"`
	CreatedAt      string `json:"created_at"`
}

type AuthorDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

type GetAllAuthorParams struct {
	Limit     int `binding:"required"`
	Page      int `binding:"required"`
	Sort      string
	Order     string
	Name      string
	CreatedAt []string
}

func ToAuthorResponse(p Author) AuthorResponse {
	return AuthorResponse{
		ID:             p.ID,
		Name:           p.Name,
		AvatarUrl:      p.AvatarUrl,
		AvatarFileName: p.AvatarFileName,
		CreatedAt:      p.CreatedAt.Format("2006-01-02"),
	}
}
