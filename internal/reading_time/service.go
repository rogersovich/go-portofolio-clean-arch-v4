package reading_time

type Service interface {
	GetAllReadingTimes() ([]ReadingTimeResponse, error)
	GetReadingTimeById(id string) (ReadingTimeResponse, error)
	CreateReadingTime(p CreateReadingTimeRequest) (ReadingTimeResponse, error)
	UpdateReadingTime(p UpdateReadingTimeRequest) (ReadingTimeUpdateResponse, error)
	DeleteReadingTime(id int) (ReadingTime, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllReadingTimes() ([]ReadingTimeResponse, error) {
	datas, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []ReadingTimeResponse
	for _, p := range datas {
		result = append(result, ToReadingTimeResponse(p))
	}
	return result, nil
}

func (s *service) GetReadingTimeById(id string) (ReadingTimeResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return ReadingTimeResponse{}, err
	}
	return ToReadingTimeResponse(data), nil
}

func (s *service) CreateReadingTime(p CreateReadingTimeRequest) (ReadingTimeResponse, error) {
	data, err := s.repo.CreateReadingTime(p)
	if err != nil {
		return ReadingTimeResponse{}, err
	}
	return ToReadingTimeResponse(data), nil
}

func (s *service) UpdateReadingTime(p UpdateReadingTimeRequest) (ReadingTimeUpdateResponse, error) {
	data, err := s.repo.UpdateReadingTime(p)
	if err != nil {
		return ReadingTimeUpdateResponse{}, err
	}

	return ToReadingTimeUpdateResponse(data), nil
}

func (s *service) DeleteReadingTime(id int) (ReadingTime, error) {
	data, err := s.repo.DeleteReadingTime(id)
	if err != nil {
		return ReadingTime{}, err
	}
	return data, nil
}
