package blog

import (
	"mime/multipart"
	"time"
)

type CreateBlogRequest struct {
	TopicIds        []string
	AuthorID        int    `validate:"required"`
	Title           string `validate:"required"`
	DescriptionHTML string `validate:"required"`
	BannerFile      *multipart.FileHeader
	Summary         string `validate:"required"`
	IsPublished     string `validate:"required,oneof=Y N"`
}

type CreateBlogDTO struct {
	TopicIds        []string
	AuthorID        int
	StatisticID     int
	ReadingTimeID   int
	Title           string
	DescriptionHTML string
	BannerUrl       string
	BannerFileName  string
	Summary         string
	Status          string
	PublishedAt     *time.Time
}

type BlogResponse struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	DescriptionHTML string  `json:"description_html"`
	BannerUrl       string  `json:"banner_url"`
	BannerFileName  string  `json:"banner_file_name"`
	Summary         string  `json:"summary"`
	Status          string  `json:"status"`
	PublishedAt     *string `json:"published_at"`
	CreatedAt       string  `json:"created_at"`
}

func ToBlogResponse(p Blog) BlogResponse {
	var publishedAtPointer *string
	if p.PublishedAt != nil {
		formattedPublishedAt := p.PublishedAt.Format("2006-01-02")
		publishedAtPointer = &formattedPublishedAt
	}
	return BlogResponse{
		ID:              p.ID,
		Title:           p.Title,
		DescriptionHTML: p.DescriptionHTML,
		BannerUrl:       p.BannerUrl,
		BannerFileName:  p.BannerFileName,
		Summary:         p.Summary,
		Status:          p.Status,
		PublishedAt:     publishedAtPointer,
		CreatedAt:       p.CreatedAt.Format("2006-01-02"),
	}
}
