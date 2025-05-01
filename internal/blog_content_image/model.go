package blog_content_image

import (
	"time"

	"gorm.io/gorm"
)

type BlogContentImage struct {
	ID            int    `json:"id" gorm:"primaryKey"`
	BlogID        *int   `json:"blog_id"`
	ImageUrl      string `json:"image_url"`
	ImageFileName string `json:"image_file_name"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
