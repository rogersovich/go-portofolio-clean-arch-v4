package author

type CreateAuthorRequest struct {
	Name       string      `json:"name" validate:"required"`
	AvatarFile interface{} `json:"avatar_file"`
}

type AuthorResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
	CreatedAt string `json:"created_at"`
}

func ToAuthorResponse(p Author) AuthorResponse {
	return AuthorResponse{
		ID:        p.ID,
		Name:      p.Name,
		AvatarUrl: p.AvatarUrl,
		CreatedAt: p.CreatedAt.Format("2006-01-02"),
	}
}
