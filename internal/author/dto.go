package author

type CreateAuthorRequest struct {
	Name           string `json:"name" validate:"required"`
	AvatarUrl      string `json:"avatar_url" validate:"required"`
	AvatarFileName string `json:"avatar_file_name" validate:"required"`
}

type UpdateAuthorRequest struct {
	Id   int    `json:"id" validate:"required,numeric,number"`
	Name string `json:"name" validate:"required"`
}

type UpdateAuthorDTO struct {
	Id             uint   `json:"id" validate:"required"`
	Name           string `json:"name" validate:"required"`
	AvatarUrl      string `json:"avatar_url"`
	AvatarFileName string `json:"avatar_file_name"`
}

type AuthorResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	AvatarUrl      string `json:"avatar_url"`
	AvatarFileName string `json:"avatar_file_name"`
	CreatedAt      string `json:"created_at"`
}

type AuthorUpdateResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	AvatarUrl      string `json:"avatar_url"`
	AvatarFileName string `json:"avatar_file_name"`
}

type AuthorDeleteRequest struct {
	ID int `json:"id" binding:"required"`
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

func ToAuthorUpdateResponse(p Author) AuthorUpdateResponse {
	return AuthorUpdateResponse{
		ID:             p.ID,
		Name:           p.Name,
		AvatarUrl:      p.AvatarUrl,
		AvatarFileName: p.AvatarFileName,
	}
}
