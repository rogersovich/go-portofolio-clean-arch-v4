package experience

import (
	"mime/multipart"
	"time"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

type CreateExperienceRequest struct {
	Position       string                `json:"position" validate:"required"`
	CompanyName    string                `json:"company_name" validate:"required"`
	WorkType       string                `json:"work_type" validate:"required,oneof=Office Remote Hybrid"`
	Country        string                `json:"country" validate:"required"`
	City           *string               `json:"city" validate:"required"`
	SummaryHTML    string                `json:"summary_html" validate:"required"`
	FromDate       string                `json:"from_date" validate:"required"`
	ToDate         *string               `json:"to_date"`
	CompImageFile  *multipart.FileHeader `json:"comp_image_file"`
	CompWebsiteUrl string                `json:"comp_website_url" validate:"required"`
	IsCurrent      string                `json:"is_current" validate:"required,oneof=Y N"`
}

type UpdateExperienceRequest struct {
	ID             int                   `json:"id" validate:"required,numeric,number"`
	Position       string                `json:"position" validate:"required"`
	CompanyName    string                `json:"company_name" validate:"required"`
	WorkType       string                `json:"work_type" validate:"required,oneof=Office Remote Hybrid"`
	Country        string                `json:"country" validate:"required"`
	City           *string               `json:"city" validate:"required"`
	SummaryHTML    string                `json:"summary_html" validate:"required"`
	FromDate       string                `json:"from_date" validate:"required"`
	ToDate         *string               `json:"to_date"`
	CompImageFile  *multipart.FileHeader `json:"comp_image_file"`
	CompWebsiteUrl string                `json:"comp_website_url" validate:"required"`
	IsCurrent      string                `json:"is_current" validate:"required,oneof=Y N"`
}

type CreateExperienceDTO struct {
	Position          string
	CompanyName       string
	WorkType          string
	Country           string
	City              *string
	SummaryHTML       string
	FromDate          time.Time
	ToDate            *time.Time
	CompImageUrl      string
	CompImageFileName string
	CompWebsiteUrl    string
	IsCurrent         bool
}

type UpdateExperienceDTO struct {
	ID                int
	Position          string
	CompanyName       string
	WorkType          string
	Country           string
	City              *string
	SummaryHTML       string
	FromDate          time.Time
	ToDate            *time.Time
	CompImageUrl      string
	CompImageFileName string
	CompWebsiteUrl    string
	IsCurrent         bool
}

type ExperienceResponse struct {
	ID                int     `json:"id"`
	Position          string  `json:"position"`
	CompanyName       string  `json:"company_name"`
	WorkType          string  `json:"work_type"`
	Country           string  `json:"country"`
	City              *string `json:"city"`
	SummaryHTML       string  `json:"summary_html"`
	FromDate          string  `json:"from_date"`
	ToDate            *string `json:"to_date"`
	CompImageUrl      string  `json:"comp_image_url"`
	CompImageFileName string  `json:"comp_image_file_name"`
	CompWebsiteUrl    string  `json:"comp_website_url"`
	IsCurrent         string  `json:"is_current"`
	CreatedAt         string  `json:"created_at"`
}

type ExperienceDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

func ToExperienceResponse(p Experience) ExperienceResponse {
	var toDateAtPointer *string
	if p.ToDate != nil {
		formattedToDate := p.ToDate.Format("2006-01-02")
		toDateAtPointer = &formattedToDate
	}

	return ExperienceResponse{
		ID:                p.ID,
		Position:          p.Position,
		CompanyName:       p.CompanyName,
		WorkType:          p.WorkType,
		Country:           p.Country,
		City:              p.City,
		SummaryHTML:       p.SummaryHTML,
		FromDate:          p.FromDate.Format("2006-01-02"),
		ToDate:            toDateAtPointer,
		CompWebsiteUrl:    p.CompWebsiteUrl,
		CompImageUrl:      p.CompImageUrl,
		CompImageFileName: p.CompImageFileName,
		IsCurrent:         utils.BoolToYN(p.IsCurrent),
		CreatedAt:         p.CreatedAt.Format("2006-01-02"),
	}
}
