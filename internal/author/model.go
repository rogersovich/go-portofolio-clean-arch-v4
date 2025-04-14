package author

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	Name           string `json:"name"`
	AvatarUrl      string `json:"avatar_url"`
	AvatarFileName string `json:"avatar_file_name"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
