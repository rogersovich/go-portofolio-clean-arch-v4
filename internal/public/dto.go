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

type AboutPublicResponse struct {
	ID              int    `json:"about_id"`
	Title           string `json:"title"`
	DescriptionHTML string `json:"description_html"`
	AvatarUrl       string `json:"avatar_url"`
	AvatarFileName  string `json:"avatar_file_name"`
}

type TechnologyProfilePublicResponse struct {
	ID           int    `json:"tech_id"`
	Name         string `json:"name"`
	LogoUrl      string `json:"logo_url"`
	LogoFileName string `json:"logo_file_name"`
}

type ProfilePublicResponse struct {
	About        AboutPublicResponse               `json:"about"`
	Technologies []TechnologyProfilePublicResponse `json:"technologies"`
}
