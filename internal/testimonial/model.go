package testimonial

import (
	"time"

	"gorm.io/gorm"
)

type Testimonial struct {
	ID        int     `json:"id" gorm:"primaryKey"`
	Name      string  `json:"name"`
	Via       *string `json:"via"`
	Role      *string `json:"role"`
	WorkingAt *string `json:"working_at"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
