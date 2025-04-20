package blog

import (
	"time"

	"gorm.io/gorm"
)

type Blog struct {
	ID              int    `json:"id" gorm:"primaryKey"`
	StatisticID     int    `json:"statistic_id"`
	ReadingTimeID   int    `json:"reading_time_id"`
	AuthorID        int    `json:"author_id"`
	Title           string `json:"title"`
	DescriptionHTML string `json:"description_html"`
	BannerUrl       string `json:"banner_url"`
	BannerFileName  string `json:"banner_file_name"`
	Summary         string `json:"summary"`
	Status          string `json:"status"`
	PublishedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
