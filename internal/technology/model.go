package technology

import (
	"time"

	"gorm.io/gorm"
)

type Technology struct {
	ID              int    `json:"id" gorm:"primaryKey"`
	Name            string `json:"name"`
	DescriptionHTML string `json:"description_html"`
	LogoUrl         string `json:"logo_url"`
	LogoFileName    string `json:"logo_file_name"`
	IsMajor         bool   `json:"is_major"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
