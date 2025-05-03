package project

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID            int     `json:"id" gorm:"primaryKey"`
	StatisticID   int     `json:"statistic_id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	ImageUrl      string  `json:"image_url"`
	ImageFileName string  `json:"image_file_name"`
	RepositoryUrl *string `json:"repository_url"`
	Summary       string  `json:"summary"`
	Status        string  `json:"status"`
	Slug          string  `json:"slug"`
	PublishedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
