package project_content_temp_image

import "github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"

type CreateProjectContentTempImgRequest struct {
	ImageUrl      string `json:"image_url" binding:"required"`
	ImageFileName string `json:"image_file_name" binding:"required"`
	IsUsed        string `json:"is_used" binding:"required,oneof=Y N"`
}

type UpdateProjectContentTempImgRequest struct {
	Id            int    `json:"id" binding:"required"`
	ImageUrl      string `json:"image_url" binding:"required"`
	ImageFileName string `json:"image_file_name" binding:"required"`
	IsUsed        string `json:"is_used" binding:"required,oneof=Y N"`
}

type ProjectContentTempImgResponse struct {
	ID            int    `json:"id"`
	ImageUrl      string `json:"image_url"`
	ImageFileName string `json:"image_file_name"`
	IsUsed        string `json:"is_used"`
	CreatedAt     string `json:"created_at"`
}

type ProjectContentTempImgUpdateResponse struct {
	ID            int    `json:"id"`
	ImageUrl      string `json:"image_url"`
	ImageFileName string `json:"image_file_name"`
	IsUsed        string `json:"is_used"`
}

type ProjectContentTempImgDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

func ToProjectContentTempImgResponse(p ProjectContentTempImages) ProjectContentTempImgResponse {
	return ProjectContentTempImgResponse{
		ID:            p.ID,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		IsUsed:        utils.BoolToYN(p.IsUsed),
		CreatedAt:     p.CreatedAt.Format("2006-01-02"),
	}
}

func ToProjectContentTempImgUpdateResponse(p ProjectContentTempImages) ProjectContentTempImgUpdateResponse {
	return ProjectContentTempImgUpdateResponse{
		ID:            p.ID,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		IsUsed:        utils.BoolToYN(p.IsUsed),
	}
}
