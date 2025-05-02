package about

import (
	"time"

	"gorm.io/gorm"
)

type About struct {
	ID              int    `json:"id" gorm:"primaryKey"`
	Title           string `json:"title"`
	DescriptionHTML string `json:"description_html"`
	AvatarUrl       string `json:"avatar_url"`
	AvatarFileName  string `json:"avatar_file_name"`
	IsUsed          bool   `json:"is_used"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
