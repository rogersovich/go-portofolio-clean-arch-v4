package topics

type Service interface {
	GetAllTopics() ([]TopicResponse, error)
	GetTopicById(id string) (TopicResponse, error)
	CreateTopic(p CreateTopicRequest) (TopicResponse, error)
	UpdateTopic(p UpdateTopicRequest) (TopicUpdateResponse, error)
	DeleteTopic(id int) (Topic, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllTopics() ([]TopicResponse, error) {
	datas, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []TopicResponse
	for _, p := range datas {
		result = append(result, ToTopicResponse(p))
	}
	return result, nil
}

func (s *service) GetTopicById(id string) (TopicResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return TopicResponse{}, err
	}
	return ToTopicResponse(data), nil
}

func (s *service) CreateTopic(p CreateTopicRequest) (TopicResponse, error) {
	data, err := s.repo.CreateTopic(p)
	if err != nil {
		return TopicResponse{}, err
	}
	return ToTopicResponse(data), nil
}

func (s *service) UpdateTopic(p UpdateTopicRequest) (TopicUpdateResponse, error) {
	data, err := s.repo.UpdateTopic(p)
	if err != nil {
		return TopicUpdateResponse{}, err
	}

	return ToTopicUpdateResponse(data), nil
}

func (s *service) DeleteTopic(id int) (Topic, error) {
	data, err := s.repo.DeleteTopic(id)
	if err != nil {
		return Topic{}, err
	}
	return data, nil
}
