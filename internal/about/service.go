package about

import (
	"context"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

type Service interface {
	GetAllAbouts() ([]AboutResponse, error)
	GetAboutById(id int) (AboutResponse, error)
	CreateAbout(p CreateAboutRequest) (AboutResponse, error)
	UpdateAbout(p UpdateAboutRequest) error
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

func (s *service) GetAboutById(id int) (AboutResponse, error) {
	about, err := s.repo.FindById(id)
	if err != nil {
		return AboutResponse{}, err
	}
	return ToAboutResponse(about), nil
}

func (s *service) CreateAbout(p CreateAboutRequest) (AboutResponse, error) {
	avatarRes, err := utils.HandlUploadFile(p.AvatarFile, "about")
	if err != nil {
		return AboutResponse{}, err
	}

	payload := CreateAboutDTO{
		Title:           p.Title,
		DescriptionHTML: p.DescriptionHTML,
		AvatarUrl:       avatarRes.FileURL,
		AvatarFileName:  avatarRes.FileName,
	}

	about, err := s.repo.CreateAbout(payload)
	if err != nil {
		_ = utils.DeleteFromMinio(context.Background(), avatarRes.FileName)
		return AboutResponse{}, err
	}
	return ToAboutResponse(about), nil
}

func (s *service) UpdateAbout(p UpdateAboutRequest) error {
	//todo: Get About
	about, err := s.repo.FindById(p.ID)
	if err != nil {
		return err
	}

	//todo: set oldFileName
	oldFileName := ""
	if about.AvatarFileName != "" {
		oldFileName = about.AvatarFileName
	}

	var newFileURL string
	var newFileName string

	//todo: Upload File
	if p.AvatarFile != nil {
		logoRes, err := utils.HandlUploadFile(p.AvatarFile, "about")
		if err != nil {
			return err
		}

		newFileURL = logoRes.FileURL
		newFileName = logoRes.FileName
	} else {
		newFileURL = about.AvatarUrl
		newFileName = about.AvatarFileName
	}

	payload := UpdateAboutDTO{
		ID:              p.ID,
		Title:           p.Title,
		DescriptionHTML: p.DescriptionHTML,
		AvatarUrl:       newFileURL,
		AvatarFileName:  newFileName,
	}

	err = s.repo.UpdateAbout(payload)
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

func (s *service) DeleteAbout(id int) (About, error) {
	about, err := s.repo.DeleteAbout(id)
	if err != nil {
		return About{}, err
	}
	return about, nil
}
