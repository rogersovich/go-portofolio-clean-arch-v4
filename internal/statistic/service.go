package statistic

type Service interface {
	GetAllStatistics() ([]StatisticResponse, error)
	GetStatisticById(id string) (StatisticResponse, error)
	CreateStatistic(p CreateStatisticRequest) (StatisticResponse, error)
	UpdateStatistic(p UpdateStatisticRequest) (StatisticUpdateResponse, error)
	DeleteStatistic(id int) (Statistic, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllStatistics() ([]StatisticResponse, error) {
	technologies, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []StatisticResponse
	for _, p := range technologies {
		result = append(result, ToStatisticResponse(p))
	}
	return result, nil
}

func (s *service) GetStatisticById(id string) (StatisticResponse, error) {
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

func (s *service) UpdateStatistic(p UpdateStatisticRequest) (StatisticUpdateResponse, error) {
	data, err := s.repo.UpdateStatistic(p)
	if err != nil {
		return StatisticUpdateResponse{}, err
	}

	return ToStatisticUpdateResponse(data), nil
}

func (s *service) DeleteStatistic(id int) (Statistic, error) {
	data, err := s.repo.DeleteStatistic(id)
	if err != nil {
		return Statistic{}, err
	}
	return data, nil
}
