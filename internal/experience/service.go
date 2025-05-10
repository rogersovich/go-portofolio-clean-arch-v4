package experience

import (
	"context"
	"time"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

type Service interface {
	GetAllExperiences(params GetAllExperienceParams) ([]ExperienceResponse, int, error)
	GetExperienceById(id int) (ExperienceResponse, error)
	CreateExperience(p CreateExperienceRequest) (ExperienceResponse, error)
	UpdateExperience(p UpdateExperienceRequest) error
	DeleteExperience(id int) (Experience, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllExperiences(params GetAllExperienceParams) ([]ExperienceResponse, int, error) {
	datas, total, err := s.repo.FindAll(params)
	if err != nil {
		return nil, 0, err
	}

	var result []ExperienceResponse
	for _, p := range datas {
		result = append(result, ToExperienceResponse(p))
	}
	return result, total, nil
}

func (s *service) GetExperienceById(id int) (ExperienceResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return ExperienceResponse{}, err
	}
	return ToExperienceResponse(data), nil
}

func (s *service) CreateExperience(p CreateExperienceRequest) (ExperienceResponse, error) {
	imageFile, err := utils.HandlUploadFile(p.CompImageFile, "experience")
	if err != nil {
		return ExperienceResponse{}, err
	}

	fromDate, _ := utils.ParseStringToTime(p.FromDate, "2006-01-02")
	isCurrent := p.IsCurrent == "Y"
	var toDateFiltered *time.Time
	if isCurrent {
		toDateFiltered = nil
	} else {
		toDate, _ := utils.ParseStringPtrToTimePtr(p.ToDate, "2006-01-02")
		toDateFiltered = toDate
	}
	payload := CreateExperienceDTO{
		Position:          p.Position,
		CompanyName:       p.CompanyName,
		WorkType:          p.WorkType,
		Country:           p.Country,
		City:              p.City,
		SummaryHTML:       p.SummaryHTML,
		FromDate:          fromDate,
		ToDate:            toDateFiltered,
		CompImageUrl:      imageFile.FileURL,
		CompImageFileName: imageFile.FileName,
		CompWebsiteUrl:    p.CompWebsiteUrl,
		IsCurrent:         p.IsCurrent == "Y",
	}

	data, err := s.repo.CreateExperience(payload)
	if err != nil {
		_ = utils.DeleteFromMinio(context.Background(), imageFile.FileName)
		return ExperienceResponse{}, err
	}
	return ToExperienceResponse(data), nil
}

func (s *service) UpdateExperience(p UpdateExperienceRequest) error {
	experience, err := s.repo.FindById(p.ID)
	if err != nil {
		return err
	}

	//todo: set oldFileName
	oldFileName := ""
	if experience.CompImageFileName != "" {
		oldFileName = experience.CompImageFileName
	}

	//todo: Upload File
	var newFileURL string
	var newFileName string

	if p.CompImageFile != nil {
		imageRes, err := utils.HandlUploadFile(p.CompImageFile, "project")
		if err != nil {
			return err
		}

		newFileURL = imageRes.FileURL
		newFileName = imageRes.FileName
	} else {
		newFileURL = experience.CompImageUrl // keep existing if not updated
		newFileName = experience.CompImageFileName
	}

	//todo: Prepare Payload
	fromDate, _ := utils.ParseStringToTime(p.FromDate, "2006-01-02")
	isCurrent := p.IsCurrent == "Y"
	var toDateFiltered *time.Time
	if isCurrent {
		toDateFiltered = nil
	} else {
		toDate, _ := utils.ParseStringPtrToTimePtr(p.ToDate, "2006-01-02")
		toDateFiltered = toDate
	}
	payload := UpdateExperienceDTO{
		ID:                p.ID,
		Position:          p.Position,
		CompanyName:       p.CompanyName,
		WorkType:          p.WorkType,
		Country:           p.Country,
		City:              p.City,
		SummaryHTML:       p.SummaryHTML,
		FromDate:          fromDate,
		ToDate:            toDateFiltered,
		CompImageUrl:      newFileURL,
		CompImageFileName: newFileName,
		CompWebsiteUrl:    p.CompWebsiteUrl,
		IsCurrent:         isCurrent,
	}

	//todo: Update Experience
	err = s.repo.UpdateExperience(payload)
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

func (s *service) DeleteExperience(id int) (Experience, error) {
	data, err := s.repo.DeleteExperience(id)
	if err != nil {
		return Experience{}, err
	}
	return data, nil
}
