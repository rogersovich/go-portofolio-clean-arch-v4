package project_technology

import (
	"time"

	"gorm.io/gorm"
)

type ProjectTechnology struct {
	ID           int `json:"id" gorm:"primaryKey"`
	ProjectID    int `json:"project_id"`
	TechnologyID int `json:"technology_id"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
