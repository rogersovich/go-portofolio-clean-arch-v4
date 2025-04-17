package project_content_image

import (
	"time"

	"gorm.io/gorm"
)

type ProjectContentImage struct {
	ID            int    `json:"id" gorm:"primaryKey"`
	ProjectID     *int   `json:"project_id"`
	ImageUrl      string `json:"image_url"`
	ImageFileName string `json:"image_file_name"`
	IsUsed        bool   `json:"is_major"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
