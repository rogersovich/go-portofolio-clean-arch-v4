package project_content_image

import (
	"context"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

type Service interface {
	GetAllProjectContentImages() ([]ProjectContentImageResponse, error)
	GetProjectContentImageById(id int) (ProjectContentImageResponse, error)
	CreateProjectContentImage(p CreateProjectContentImageRequest) (ProjectContentImageResponse, error)
	UpdateProjectContentImage(p UpdateProjectContentImageDTO, oldPath string, newFilePath string) (ProjectContentImageUpdateResponse, error)
	DeleteProjectContentImage(id int) (ProjectContentImageResponse, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllProjectContentImages() ([]ProjectContentImageResponse, error) {
	datas, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []ProjectContentImageResponse
	for _, p := range datas {
		result = append(result, ToProjectContentImageResponse(p))
	}
	return result, nil
}

func (s *service) GetProjectContentImageById(id int) (ProjectContentImageResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return ProjectContentImageResponse{}, err
	}
	return ToProjectContentImageResponse(data), nil
}

func (s *service) CreateProjectContentImage(p CreateProjectContentImageRequest) (ProjectContentImageResponse, error) {
	data, err := s.repo.CreateProjectContentImage(p)
	if err != nil {
		return ProjectContentImageResponse{}, err
	}
	return ToProjectContentImageResponse(data), nil
}

func (s *service) UpdateProjectContentImage(p UpdateProjectContentImageDTO, oldPath string, newFilePath string) (ProjectContentImageUpdateResponse, error) {
	data, err := s.repo.UpdateProjectContentImage(p)
	if err != nil {
		return ProjectContentImageUpdateResponse{}, err
	}

	// 3. Optional: Delete old file from MinIO
	if oldPath != newFilePath {
		err = utils.DeleteFromMinio(context.Background(), oldPath) // ignore error or handle if needed
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}

	return ToProjectContentImageUpdateResponse(data), nil
}

func (s *service) DeleteProjectContentImage(id int) (ProjectContentImageResponse, error) {
	data, err := s.repo.DeleteProjectContentImage(id)
	if err != nil {
		return ProjectContentImageResponse{}, err
	}
	return ToProjectContentImageResponse(data), nil
}
