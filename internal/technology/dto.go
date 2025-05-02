package technology

import (
	"mime/multipart"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

type CreateTechnologyRequest struct {
	Name            string                `json:"name" validate:"required"`
	DescriptionHTML string                `json:"description_html" validate:"required"`
	LogoFile        *multipart.FileHeader `json:"logo_file"`
	IsMajor         string                `json:"is_major" validate:"required,oneof=Y N"`
	Link            *string               `json:"link"`
}

type UpdateTechnologyRequest struct {
	ID              int                   `json:"id" validate:"required,numeric,number"`
	Name            string                `json:"name" validate:"required"`
	DescriptionHTML string                `json:"description_html" validate:"required"`
	LogoFile        *multipart.FileHeader `json:"logo_file"`
	IsMajor         string                `json:"is_major" validate:"required,oneof=Y N"`
	Link            *string               `json:"link"`
}

type CreateTechnologyDTO struct {
	Name            string
	DescriptionHTML string
	LogoUrl         string
	LogoFileName    string
	IsMajor         bool
	Link            *string
}

type UpdateTechnologyDTO struct {
	ID              int
	Name            string
	DescriptionHTML string
	LogoUrl         string
	LogoFileName    string
	IsMajor         bool
	Link            *string
}

type TechnologyResponse struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	DescriptionHTML string  `json:"description_html"`
	LogoUrl         string  `json:"logo_url"`
	LogoFileName    string  `json:"logo_file_name"`
	IsMajor         string  `json:"is_major"`
	CreatedAt       string  `json:"created_at"`
	Link            *string `json:"link"`
}

type TechnologyUpdateResponse struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	DescriptionHTML string  `json:"description_html"`
	LogoUrl         string  `json:"logo_url"`
	LogoFileName    string  `json:"logo_file_name"`
	IsMajor         string  `json:"is_major"`
	Link            *string `json:"link"`
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
		Link:            p.Link,
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
		Link:            p.Link,
		IsMajor:         utils.BoolToYN(p.IsMajor),
	}
}
