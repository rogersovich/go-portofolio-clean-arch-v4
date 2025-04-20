package statistic

import (
	"time"

	"gorm.io/gorm"
)

type Statistic struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Likes     *int   `json:"likes"`
	Views     *int   `json:"views"`
	Type      string `json:"type"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
