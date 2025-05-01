package project_content_image

type CreateProjectContentImageRequest struct {
	ImageUrl      string `json:"image_url" validate:"required"`
	ImageFileName string `json:"image_file_name" validate:"required"`
}

type UpdateProjectContentImageRequest struct {
	Id        int  `json:"id" validate:"required"`
	ProjectID *int `json:"project_id" validate:"required"`
}

type UpdateProjectContentImageDTO struct {
	Id            int
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
