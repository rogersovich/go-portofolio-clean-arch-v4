package about

import (
	"context"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

type Service interface {
	GetAllAbouts() ([]AboutResponse, error)
	GetAboutById(id string) (AboutResponse, error)
	CreateAbout(p CreateAboutRequest) (AboutResponse, error)
	UpdateAbout(p UpdateAboutDTO, oldPath string, newFilePath string) (AboutUpdateResponse, error)
	DeleteAbout(id int) (About, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllAbouts() ([]AboutResponse, error) {
	abouts, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []AboutResponse
	for _, p := range abouts {
		result = append(result, ToAboutResponse(p))
	}
	return result, nil
}

func (s *service) GetAboutById(id string) (AboutResponse, error) {
	about, err := s.repo.FindById(id)
	if err != nil {
		return AboutResponse{}, err
	}
	return ToAboutResponse(about), nil
}

func (s *service) CreateAbout(p CreateAboutRequest) (AboutResponse, error) {
	about, err := s.repo.CreateAbout(p)
	if err != nil {
		return AboutResponse{}, err
	}
	return ToAboutResponse(about), nil
}

func (s *service) UpdateAbout(p UpdateAboutDTO, oldPath string, newFilePath string) (AboutUpdateResponse, error) {
	about, err := s.repo.UpdateAbout(p)
	if err != nil {
		return AboutUpdateResponse{}, err
	}

	// 3. Optional: Delete old file from MinIO
	if oldPath != newFilePath {
		err = utils.DeleteFromMinio(context.Background(), oldPath) // ignore error or handle if needed
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}

	return ToAboutUpdateResponse(about), nil
}

func (s *service) DeleteAbout(id int) (About, error) {
	about, err := s.repo.DeleteAbout(id)
	if err != nil {
		return About{}, err
	}
	return about, nil
}
