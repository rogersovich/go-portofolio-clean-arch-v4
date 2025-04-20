package project_content_temp_image

import (
	"time"

	"gorm.io/gorm"
)

type ProjectContentTempImages struct {
	ID            int    `json:"id" gorm:"primaryKey"`
	ImageUrl      string `json:"image_url"`
	ImageFileName string `json:"image_file_name"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
