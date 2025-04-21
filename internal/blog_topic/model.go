package blog_topic

import (
	"time"

	"gorm.io/gorm"
)

type BlogTopic struct {
	ID        int `json:"id" gorm:"primaryKey"`
	BlogID    int `json:"blog_id"`
	TopicID   int `json:"topic_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
