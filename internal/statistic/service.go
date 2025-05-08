package statistic

import "gorm.io/gorm"

type Service interface {
	GetAllStatistics(params GetAllStatisticParams) ([]StatisticResponse, int, error)
	GetStatisticById(id int) (StatisticResponse, error)
	CreateStatistic(p CreateStatisticRequest) (StatisticResponse, error)
	CreateStatisticWithTx(p CreateStatisticRequest, tx *gorm.DB) (StatisticResponse, error)
	UpdateStatistic(p UpdateStatisticRequest) error
	DeleteStatistic(id int) (Statistic, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllStatistics(params GetAllStatisticParams) ([]StatisticResponse, int, error) {
	datas, total, err := s.repo.FindAll(params)
	if err != nil {
		return nil, 0, err
	}

	var result []StatisticResponse
	for _, p := range datas {
		result = append(result, ToStatisticResponse(p))
	}
	return result, total, nil
}

func (s *service) GetStatisticById(id int) (StatisticResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return StatisticResponse{}, err
	}
	return ToStatisticResponse(data), nil
}

func (s *service) CreateStatistic(p CreateStatisticRequest) (StatisticResponse, error) {
	data, err := s.repo.CreateStatistic(p)
	if err != nil {
		return StatisticResponse{}, err
	}
	return ToStatisticResponse(data), nil
}

func (s *service) CreateStatisticWithTx(p CreateStatisticRequest, tx *gorm.DB) (StatisticResponse, error) {
	data, err := s.repo.CreateStatisticWithTx(p, tx)
	if err != nil {
		return StatisticResponse{}, err
	}
	return ToStatisticResponse(data), nil
}

func (s *service) UpdateStatistic(p UpdateStatisticRequest) error {
	_, err := s.repo.FindById(p.ID)
	if err != nil {
		return err
	}

	err = s.repo.UpdateStatistic(p)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteStatistic(id int) (Statistic, error) {
	data, err := s.repo.DeleteStatistic(id)
	if err != nil {
		return Statistic{}, err
	}
	return data, nil
}
