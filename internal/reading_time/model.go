package reading_time

import (
	"time"

	"gorm.io/gorm"
)

type ReadingTime struct {
	ID               int    `json:"id" gorm:"primaryKey"`
	Minutes          int    `json:"minutes"`
	TextLength       int    `json:"text_length"`
	EstimatedSeconds int    `json:"estimated_seconds"`
	WordCount        int    `json:"word_count"`
	Type             string `json:"type"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
