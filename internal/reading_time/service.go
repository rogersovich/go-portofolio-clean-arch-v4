package reading_time

import "gorm.io/gorm"

type Service interface {
	GetAllReadingTimes(params GetAllReadingTimeParams) ([]ReadingTimeResponse, int, error)
	GetReadingTimeById(id int) (ReadingTimeResponse, error)
	CreateReadingTime(p CreateReadingTimeRequest, tx *gorm.DB) (ReadingTimeResponse, error)
	UpdateReadingTime(p UpdateReadingTimeRequest, tx *gorm.DB) error
	DeleteReadingTime(id int) (ReadingTime, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllReadingTimes(params GetAllReadingTimeParams) ([]ReadingTimeResponse, int, error) {
	datas, total, err := s.repo.FindAll(params)
	if err != nil {
		return nil, 0, err
	}

	var result []ReadingTimeResponse
	for _, p := range datas {
		result = append(result, ToReadingTimeResponse(p))
	}
	return result, total, nil
}

func (s *service) GetReadingTimeById(id int) (ReadingTimeResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return ReadingTimeResponse{}, err
	}
	return ToReadingTimeResponse(data), nil
}

func (s *service) CreateReadingTime(p CreateReadingTimeRequest, tx *gorm.DB) (ReadingTimeResponse, error) {
	data, err := s.repo.CreateReadingTime(p, tx)
	if err != nil {
		return ReadingTimeResponse{}, err
	}
	return ToReadingTimeResponse(data), nil
}

func (s *service) UpdateReadingTime(p UpdateReadingTimeRequest, tx *gorm.DB) error {
	_, err := s.repo.FindById(p.ID)
	if err != nil {
		return err
	}

	err = s.repo.UpdateReadingTime(p, tx)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteReadingTime(id int) (ReadingTime, error) {
	data, err := s.repo.DeleteReadingTime(id)
	if err != nil {
		return ReadingTime{}, err
	}
	return data, nil
}
