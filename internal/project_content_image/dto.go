package project_content_image

import "mime/multipart"

type CreateProjectContentImageRequest struct {
	ImageFile *multipart.FileHeader `json:"image_file" validate:"required"`
}

type UpdateProjectContentImageRequest struct {
	ID        int                   `json:"id" validate:"required"`
	ProjectID *int                  `json:"project_id" validate:"required"`
	ImageFile *multipart.FileHeader `json:"image_file" validate:"required"`
}

type CreateProjectContentImageDTO struct {
	ProjectID     *int
	ImageUrl      string
	ImageFileName string
}

type UpdateProjectContentImageDTO struct {
	ID            int
	ProjectID     *int
	ImageUrl      string
	ImageFileName string
}

type ProjectContentImageResponse struct {
	ID            int    `json:"id"`
	ProjectID     *int   `json:"project_id"`
	ImageUrl      string `json:"image_url"`
	ImageFileName string `json:"image_file_name"`
	CreatedAt     string `json:"created_at"`
}

type ProjectContentImageUpdateResponse struct {
	ID            int    `json:"id"`
	ProjectID     *int   `json:"project_id"`
	ImageUrl      string `json:"image_url"`
	ImageFileName string `json:"image_file_name"`
}

type ProjectContentImageDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

type ProjectImagesExistingPayload struct {
	ID            int    `json:"id"`
	ImageUrl      string `json:"image_url"`
	ImageFileName string `json:"image_file_name"`
}

type ProjectImagesFindResponse struct {
	ID        int    `json:"id"`
	ProjectID *int   `json:"project_id"`
	ImageUrl  string `json:"image_url"`
}

type ProjectImagesBulkUpdateDTO struct {
	ProjectID int      `json:"project_id"`
	ImageUrls []string `json:"image_urls"`
}

func ToProjectContentImageResponse(p ProjectContentImage) ProjectContentImageResponse {
	return ProjectContentImageResponse{
		ID:            p.ID,
		ProjectID:     p.ProjectID,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		CreatedAt:     p.CreatedAt.Format("2006-01-02"),
	}
}

func ToProjectContentImageUpdateResponse(p ProjectContentImage) ProjectContentImageUpdateResponse {
	return ProjectContentImageUpdateResponse{
		ID:            p.ID,
		ProjectID:     p.ProjectID,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
	}
}
