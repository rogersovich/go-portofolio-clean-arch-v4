package technology

import (
	"context"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

type Service interface {
	GetAllTechnologies() ([]TechnologyResponse, error)
	GetTechnologyById(id string) (TechnologyResponse, error)
	CreateTechnology(p CreateTechnologyRequest) (TechnologyResponse, error)
	UpdateTechnology(p UpdateTechnologyDTO, oldPath string, newFilePath string) (TechnologyUpdateResponse, error)
	DeleteTechnology(id int) (Technology, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllTechnologies() ([]TechnologyResponse, error) {
	technologies, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []TechnologyResponse
	for _, p := range technologies {
		result = append(result, ToTechnologyResponse(p))
	}
	return result, nil
}

func (s *service) GetTechnologyById(id string) (TechnologyResponse, error) {
	technology, err := s.repo.FindById(id)
	if err != nil {
		return TechnologyResponse{}, err
	}
	return ToTechnologyResponse(technology), nil
}

func (s *service) CreateTechnology(p CreateTechnologyRequest) (TechnologyResponse, error) {
	technology, err := s.repo.CreateTechnology(p)
	if err != nil {
		return TechnologyResponse{}, err
	}
	return ToTechnologyResponse(technology), nil
}

func (s *service) UpdateTechnology(p UpdateTechnologyDTO, oldPath string, newFilePath string) (TechnologyUpdateResponse, error) {
	technology, err := s.repo.UpdateTechnology(p)
	if err != nil {
		return TechnologyUpdateResponse{}, err
	}

	// 3. Optional: Delete old file from MinIO
	if oldPath != newFilePath {
		err = utils.DeleteFromMinio(context.Background(), oldPath) // ignore error or handle if needed
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}

	return ToTechnologyUpdateResponse(technology), nil
}

func (s *service) DeleteTechnology(id int) (Technology, error) {
	technology, err := s.repo.DeleteTechnology(id)
	if err != nil {
		return Technology{}, err
	}
	return technology, nil
}
