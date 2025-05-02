package blog

import (
	"mime/multipart"
	"time"
)

type CreateBlogRequest struct {
	TopicIds        []int
	ContentImages   []string
	AuthorID        int    `validate:"required"`
	Title           string `validate:"required"`
	DescriptionHTML string `validate:"required"`
	BannerFile      *multipart.FileHeader
	Summary         string `validate:"required"`
	IsPublished     string `validate:"required,oneof=Y N"`
}

type CreateBlogDTO struct {
	TopicIds        []int
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
	StatisticID     int     `json:"statistic_id"`
	ReadingTimeID   int     `json:"reading_time_id"`
	AuthorID        int     `json:"author_id"`
	Title           string  `json:"title"`
	DescriptionHTML string  `json:"description_html"`
	BannerUrl       string  `json:"banner_url"`
	BannerFileName  string  `json:"banner_file_name"`
	Summary         string  `json:"summary"`
	Status          string  `json:"status"`
	PublishedAt     *string `json:"published_at"`
	CreatedAt       string  `json:"created_at"`
}

type RawBlogRelationResponse struct {
	ID                          int        `json:"id"`
	Title                       string     `json:"title"`
	DescriptionHTML             string     `json:"description_html"`
	BannerUrl                   string     `json:"banner_url"`
	BannerFileName              string     `json:"banner_file_name"`
	Summary                     string     `json:"summary"`
	Status                      string     `json:"status"`
	PublishedAt                 *time.Time `json:"published_at"`
	CreatedAt                   time.Time  `json:"created_at"`
	AuthorID                    int        `json:"author_id"`
	AuthorName                  string     `json:"author_name"`
	ReadingTimeID               int        `json:"reading_time_id"`
	ReadingTimeMinutes          int        `json:"reading_time_minutes"`
	ReadingTimeTextLength       int        `json:"reading_time_text_length"`
	ReadingTimeEstimatedSeconds float64    `json:"reading_time_estimated_seconds"`
	ReadingTimeWordCount        int        `json:"reading_time_word_count"`
	ReadingTimeType             string     `json:"reading_time_type"`
	StatisticID                 int        `json:"statistic_id"`
	StatisticLikes              int        `json:"statistic_likes"`
	StatisticViews              int        `json:"statistic_views"`
	StatisticType               string     `json:"statistic_type"`
	TopicID                     int        `json:"topic_id"`
	TopicName                   string     `json:"topic_name"`
	BlogContentImageID          int        `json:"blog_content_image_id"`
	BlogContentImageUrl         string     `json:"blog_content_image_url"`
	BlogContentImageFileName    string     `json:"blog_content_image_file_name"`
}

type BlogAuthorDTO struct {
	AuthorID   int    `json:"id"`
	AuthorName string `json:"name"`
}

type BlogReadingTimeDTO struct {
	ReadingTimeID               int     `json:"id"`
	ReadingTimeMinutes          int     `json:"minutes"`
	ReadingTimeTextLength       int     `json:"text_length"`
	ReadingTimeEstimatedSeconds float64 `json:"estimated_seconds"`
	ReadingTimeWordCount        int     `json:"word_count"`
	ReadingTimeType             string  `json:"type"`
}

type BlogStatisticDTO struct {
	StatisticID    int    `json:"id"`
	StatisticLikes int    `json:"likes"`
	StatisticViews int    `json:"views"`
	StatisticType  string `json:"type"`
}

type BlogTopicDTO struct {
	TopicID   int    `json:"id"`
	TopicName string `json:"name"`
}

type BlogContentImageDTO struct {
	BlogContentImageID       int    `json:"id"`
	BlogContentImageUrl      string `json:"image_url"`
	BlogContentImageFileName string `json:"image_file_name"`
}

type UpdateBlogRequest struct {
	ID              int                  `validate:"required"`
	TopicIds        []UpdateBlogTopicDTO `validate:"required,dive"`
	ContentImages   []string             `validate:"required,dive"`
	AuthorID        int                  `validate:"required"`
	Title           string               `validate:"required"`
	DescriptionHTML string               `validate:"required"`
	BannerFile      *multipart.FileHeader
	Summary         string `validate:"required"`
	IsPublished     string `validate:"required,oneof=Y N"`
}

type UpdateBlogDTO struct {
	ID              int
	TopicIds        []UpdateBlogTopicDTO
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

type UpdateBlogTopicDTO struct {
	TopicID int `json:"topic_id"`
}

type BlogRelationResponse struct {
	ID              int                   `json:"id"`
	Title           string                `json:"title"`
	DescriptionHTML string                `json:"description_html"`
	BannerUrl       string                `json:"banner_url"`
	BannerFileName  string                `json:"banner_file_name"`
	Summary         string                `json:"summary"`
	Status          string                `json:"status"`
	PublishedAt     *string               `json:"published_at"`
	CreatedAt       string                `json:"created_at"`
	Author          *BlogAuthorDTO        `json:"author"`
	ReadingTime     *BlogReadingTimeDTO   `json:"reading_time"`
	Statistic       *BlogStatisticDTO     `json:"statistic"`
	Topics          []BlogTopicDTO        `json:"topics"`
	ContentImages   []BlogContentImageDTO `json:"content_image"`
}

type BlogUpdateResponse struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	DescriptionHTML string  `json:"description_html"`
	BannerUrl       string  `json:"banner_url"`
	BannerFileName  string  `json:"banner_file_name"`
	Summary         string  `json:"summary"`
	Status          string  `json:"status"`
	PublishedAt     *string `json:"published_at"`
	StatisticID     int     `json:"statistic_id"`
	ReadingTimeID   int     `json:"reading_time_id"`
	AuthorID        int     `json:"author_id"`
}

type BlogDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

func ToBlogResponse(p Blog) BlogResponse {
	var publishedAtPointer *string
	if p.PublishedAt != nil {
		formattedPublishedAt := p.PublishedAt.Format("2006-01-02")
		publishedAtPointer = &formattedPublishedAt
	}
	return BlogResponse{
		ID:              p.ID,
		StatisticID:     p.StatisticID,
		ReadingTimeID:   p.ReadingTimeID,
		AuthorID:        p.AuthorID,
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

func ToBlogUpdateResponse(p Blog) BlogUpdateResponse {
	var publishedAtPointer *string
	if p.PublishedAt != nil {
		formattedPublishedAt := p.PublishedAt.Format("2006-01-02")
		publishedAtPointer = &formattedPublishedAt
	}
	return BlogUpdateResponse{
		ID:              p.ID,
		Title:           p.Title,
		DescriptionHTML: p.DescriptionHTML,
		BannerUrl:       p.BannerUrl,
		BannerFileName:  p.BannerFileName,
		Summary:         p.Summary,
		Status:          p.Status,
		PublishedAt:     publishedAtPointer,
		StatisticID:     p.StatisticID,
		ReadingTimeID:   p.ReadingTimeID,
		AuthorID:        p.AuthorID,
	}
}
