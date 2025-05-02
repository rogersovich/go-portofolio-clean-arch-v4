package public

type AuthorPublicResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	AvatarUrl      string `json:"avatar_url"`
	AvatarFileName string `json:"avatar_file_name"`
	CreatedAt      string `json:"created_at"`
}

type AuthorPublicParams struct {
	Page  int
	Limit int
	Sort  string
	Order string
	Name  string
}
