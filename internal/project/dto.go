package project

import (
	"mime/multipart"
	"time"
)

type CreateProjectRequest struct {
	Title         string `validate:"required"`
	Description   string `validate:"required"`
	ImageFile     *multipart.FileHeader
	RepositoryUrl *string
	Summary       string `validate:"required"`
	IsPublished   string `validate:"required,oneof=Y N"`
	TechnologyIds []string
	ContentImages []string
}

type UpdateProjectRequest struct {
	Id              int                   `validate:"required"`
	Title           string                `validate:"required"`
	Description     string                `validate:"required"`
	ImageFile       *multipart.FileHeader `validate:"omitempty"`
	RepositoryUrl   *string
	Summary         string                       `validate:"required"`
	IsPublished     string                       `validate:"required,oneof=Y N"`
	TechnologyIds   []ProjectTechUpdatePayload   `json:"technology_ids" validate:"required,dive"`
	ContentImageIds []ProjectImagesUpdatePayload `json:"content_image_ids" validate:"required,dive"`
}

type CreateProjectDTO struct {
	Id                   int
	TechnologyIds        []string
	ProjectContentImages []string
	Title                string
	Description          string
	RepositoryUrl        *string
	ImageUrl             string
	ImageFileName        string
	Summary              string
	Status               string
	PublishedAt          *time.Time
}

type UpdateProjectDTO struct {
	Id              int
	Title           string
	Description     string
	RepositoryUrl   *string
	ImageUrl        string
	ImageFileName   string
	Summary         string
	Status          string
	PublishedAt     *time.Time
	TechnologyIds   []ProjectTechUpdatePayload
	ContentImageIds []ProjectImagesUpdatePayload
}

type ProjectResponse struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	ImageUrl      string  `json:"image_url"`
	ImageFileName string  `json:"image_file_name"`
	RepositoryUrl *string `json:"repository_url"`
	Summary       string  `json:"summary"`
	Status        string  `json:"status"`
	PublishedAt   *string `json:"published_at"`
	CreatedAt     string  `json:"created_at"`
}

type RawProjectRelationResponse struct {
	ID                  int        `json:"id"`
	StatisticID         int        `json:"statistic_id"`
	StatisticViews      int        `json:"statistic_views"`
	StatisticLikes      int        `json:"statistic_likes"`
	StatisticType       string     `json:"statistic_type"`
	ProjectTechnologyID int        `json:"project_technology_id"`
	TechnologyID        int        `json:"technology_id"`
	TechnologyName      string     `json:"technology_name"`
	ProjectImgID        int        `json:"project_img_id"`
	ProjectImgFileName  string     `json:"project_img_file_name"`
	ProjectImgUrl       string     `json:"project_img_url"`
	Title               string     `json:"title"`
	Description         string     `json:"description"`
	ImageUrl            string     `json:"image_url"`
	ImageFileName       string     `json:"image_file_name"`
	RepositoryUrl       *string    `json:"repository_url"`
	Summary             string     `json:"summary"`
	Status              string     `json:"status"`
	PublishedAt         *time.Time `json:"published_at"`
	CreatedAt           time.Time  `json:"created_at"`
}

type ProjectRelationResponse struct {
	ID            int                       `json:"id"`
	StatisticID   int                       `json:"statistic_id"`
	Statistic     ProjectStatisticDTO       `json:"statistic"`
	Technologies  []ProjectTechnologiesDTO  `json:"technologies"`
	ContentImages []ProjectContentImagesDTO `json:"images"`
	Title         string                    `json:"title"`
	Description   string                    `json:"description"`
	ImageUrl      string                    `json:"image_url"`
	ImageFileName string                    `json:"image_file_name"`
	RepositoryUrl *string                   `json:"repository_url"`
	Summary       string                    `json:"summary"`
	Status        string                    `json:"status"`
	PublishedAt   *string                   `json:"published_at"`
	CreatedAt     string                    `json:"created_at"`
}

type ProjectTechnologiesDTO struct {
	ProjectTechID int    `json:"project_tech_id"`
	TechID        int    `json:"tech_id"`
	TechName      string `json:"tech_name"`
}

type ProjectContentImagesDTO struct {
	ProjectImageID int    `json:"project_image_id"`
	ImageUrl       string `json:"image_url"`
	ImageFileName  string `json:"image_file_name"`
}

type ProjectStatisticDTO struct {
	ID        int    `json:"id"`
	Likes     int    `json:"likes"`
	Views     int    `json:"views"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

type ProjectUpdateResponse struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	ImageUrl      string  `json:"image_url"`
	ImageFileName string  `json:"image_file_name"`
	RepositoryUrl *string `json:"repository_url"`
	Summary       string  `json:"summary"`
	Status        string  `json:"status"`
	PublishedAt   *string `json:"published_at"`
}

type ProjectTechUpdatePayload struct {
	ID     int `json:"id" validate:"required,gt=0"`
	TechID int `json:"tech_id" validate:"required,gt=0"`
}

type ProjectImagesUpdatePayload struct {
	ID            int    `json:"id" validate:"required,gt=0"`
	ImageUrl      string `json:"image_url" validate:"required"`
	ImageFileName string `json:"image_file_name" validate:"required"`
}

type ProjectDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

func ToProjectResponse(p Project) ProjectResponse {
	var publishedAtPointer *string
	if p.PublishedAt != nil {
		formattedPublishedAt := p.PublishedAt.Format("2006-01-02")
		publishedAtPointer = &formattedPublishedAt
	}
	return ProjectResponse{
		ID:            p.ID,
		Title:         p.Title,
		Description:   p.Description,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		RepositoryUrl: p.RepositoryUrl,
		Summary:       p.Summary,
		Status:        p.Status,
		PublishedAt:   publishedAtPointer,
		CreatedAt:     p.CreatedAt.Format("2006-01-02"),
	}
}

func ToProjectUpdateResponse(p Project) ProjectUpdateResponse {
	var publishedAtPointer *string
	if p.PublishedAt != nil {
		formattedPublishedAt := p.PublishedAt.Format("2006-01-02")
		publishedAtPointer = &formattedPublishedAt
	}
	return ProjectUpdateResponse{
		ID:            p.ID,
		Title:         p.Title,
		Description:   p.Description,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		RepositoryUrl: p.RepositoryUrl,
		Summary:       p.Summary,
		Status:        p.Status,
		PublishedAt:   publishedAtPointer,
	}
}
