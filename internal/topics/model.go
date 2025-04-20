package topics

import (
	"time"

	"gorm.io/gorm"
)

type Topic struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
