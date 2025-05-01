package experience

import (
	"time"

	"gorm.io/gorm"
)

type Experience struct {
	ID                int        `json:"id" gorm:"primaryKey"`
	Position          string     `json:"position"`
	CompanyName       string     `json:"company_name"`
	WorkType          string     `json:"work_type"`
	Country           string     `json:"country"`
	City              *string    `json:"city"`
	SummaryHTML       string     `json:"summary_html"`
	FromDate          time.Time  `json:"from_date"`
	ToDate            *time.Time `json:"to_date"`
	CompImageUrl      string     `json:"comp_image_url"`
	CompImageFileName string     `json:"comp_image_file_name"`
	CompWebsiteUrl    string     `json:"comp_website_url"`
	IsCurrent         bool       `json:"is_current"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}
