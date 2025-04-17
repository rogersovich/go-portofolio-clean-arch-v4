package technology

import "github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"

type CreateTechnologyRequest struct {
	Name            string `json:"name" validate:"required"`
	DescriptionHTML string `json:"description_html" validate:"required"`
	LogoUrl         string `json:"logo_url" validate:"required"`
	LogoFileName    string `json:"logo_file_name" validate:"required"`
	IsMajor         string `json:"is_major" validate:"required,oneof=Y N"`
}

type UpdateTechnologyRequest struct {
	Id              int    `json:"id" validate:"required,numeric,number"`
	Name            string `json:"name" validate:"required"`
	DescriptionHTML string `json:"description_html" validate:"required"`
	IsMajor         string `json:"is_major" validate:"required,oneof=Y N"`
}

type UpdateTechnologyDTO struct {
	Id              uint
	Name            string
	DescriptionHTML string
	LogoUrl         string
	LogoFileName    string
	IsMajor         string
}

type TechnologyResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	DescriptionHTML string `json:"description_html"`
	LogoUrl         string `json:"logo_url"`
	LogoFileName    string `json:"logo_file_name"`
	IsMajor         string `json:"is_major"`
	CreatedAt       string `json:"created_at"`
}

type TechnologyUpdateResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	DescriptionHTML string `json:"description_html"`
	LogoUrl         string `json:"logo_url"`
	LogoFileName    string `json:"logo_file_name"`
	IsMajor         string `json:"is_major"`
}

type TechnologyDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

func ToTechnologyResponse(p Technology) TechnologyResponse {
	return TechnologyResponse{
		ID:              p.ID,
		Name:            p.Name,
		DescriptionHTML: p.DescriptionHTML,
		LogoUrl:         p.LogoUrl,
		LogoFileName:    p.LogoFileName,
		IsMajor:         utils.BoolToYN(p.IsMajor),
		CreatedAt:       p.CreatedAt.Format("2006-01-02"),
	}
}

func ToTechnologyUpdateResponse(p Technology) TechnologyUpdateResponse {
	return TechnologyUpdateResponse{
		ID:              p.ID,
		Name:            p.Name,
		DescriptionHTML: p.DescriptionHTML,
		LogoUrl:         p.LogoUrl,
		LogoFileName:    p.LogoFileName,
		IsMajor:         utils.BoolToYN(p.IsMajor),
	}
}
