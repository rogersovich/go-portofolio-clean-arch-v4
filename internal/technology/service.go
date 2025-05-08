package technology

import (
	"context"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

type Service interface {
	GetAllTechnologies(params GetAllTechnologyParams) ([]TechnologyResponse, int, error)
	GetTechnologyById(id int) (TechnologyResponse, error)
	CreateTechnology(p CreateTechnologyRequest) (TechnologyResponse, error)
	UpdateTechnology(p UpdateTechnologyRequest) error
	DeleteTechnology(id int) (Technology, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllTechnologies(params GetAllTechnologyParams) ([]TechnologyResponse, int, error) {
	datas, total, err := s.repo.FindAll(params)
	if err != nil {
		return nil, 0, err
	}

	var result []TechnologyResponse
	for _, p := range datas {
		result = append(result, ToTechnologyResponse(p))
	}
	return result, total, nil
}

func (s *service) GetTechnologyById(id int) (TechnologyResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return TechnologyResponse{}, err
	}
	return ToTechnologyResponse(data), nil
}

func (s *service) CreateTechnology(p CreateTechnologyRequest) (TechnologyResponse, error) {
	logoRes, err := utils.HandlUploadFile(p.LogoFile, "technology")
	if err != nil {
		return TechnologyResponse{}, err
	}

	payload := CreateTechnologyDTO{
		Name:            p.Name,
		DescriptionHTML: p.DescriptionHTML,
		LogoUrl:         logoRes.FileURL,
		LogoFileName:    logoRes.FileName,
		IsMajor:         p.IsMajor == "Y",
		Link:            p.Link,
	}

	data, err := s.repo.CreateTechnology(payload)
	if err != nil {
		_ = utils.DeleteFromMinio(context.Background(), logoRes.FileName)
		return TechnologyResponse{}, err
	}
	return ToTechnologyResponse(data), nil
}

func (s *service) UpdateTechnology(p UpdateTechnologyRequest) error {
	//todo: Get Technology
	technology, err := s.repo.FindById(p.ID)
	if err != nil {
		return err
	}

	//todo: set oldFileName
	oldFileName := ""
	if technology.LogoFileName != "" {
		oldFileName = technology.LogoFileName
	}

	var newFileURL string
	var newFileName string

	//todo: Upload File
	if p.LogoFile != nil {
		logoRes, err := utils.HandlUploadFile(p.LogoFile, "technology")
		if err != nil {
			return err
		}

		newFileURL = logoRes.FileURL
		newFileName = logoRes.FileName
	} else {
		newFileURL = technology.LogoUrl
		newFileName = technology.LogoFileName
	}

	payload := UpdateTechnologyDTO{
		ID:              p.ID,
		Name:            p.Name,
		DescriptionHTML: p.DescriptionHTML,
		LogoUrl:         newFileURL,
		LogoFileName:    newFileName,
		IsMajor:         p.IsMajor == "Y",
		Link:            p.Link,
	}

	err = s.repo.UpdateTechnology(payload)
	if err != nil {
		_ = utils.DeleteFromMinio(context.Background(), newFileName)
		return err
	}

	//todo: Delete Old Image
	if oldFileName != newFileName {
		_ = utils.DeleteFromMinio(context.Background(), oldFileName)
	}

	return nil
}

func (s *service) DeleteTechnology(id int) (Technology, error) {
	data, err := s.repo.DeleteTechnology(id)
	if err != nil {
		return Technology{}, err
	}
	return data, nil
}
