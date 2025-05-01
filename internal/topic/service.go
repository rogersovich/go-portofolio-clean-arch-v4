package topic

import "fmt"

type Service interface {
	GetAllTopics() ([]TopicResponse, error)
	GetTopicById(id int) (TopicResponse, error)
	CreateTopic(p CreateTopicRequest) (TopicResponse, error)
	UpdateTopic(p UpdateTopicRequest) error
	DeleteTopic(id int) (Topic, error)
	CheckTopicIds(ids []int) ([]TopicHasCheckResponse, error)
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

func (s *service) GetTopicById(id int) (TopicResponse, error) {
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

func (s *service) UpdateTopic(p UpdateTopicRequest) error {
	_, err := s.repo.FindById(p.ID)
	if err != nil {
		return err
	}

	err = s.repo.UpdateTopic(p)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteTopic(id int) (Topic, error) {
	data, err := s.repo.DeleteTopic(id)
	if err != nil {
		return Topic{}, err
	}
	return data, nil
}

func (s *service) CheckTopicIds(ids []int) ([]TopicHasCheckResponse, error) {
	data, err := s.repo.CheckTopicIds(ids)
	if err != nil {
		return nil, err
	}

	if len(data) != len(ids) {
		err := fmt.Errorf("some topic_ids not found in database")
		return nil, err
	}

	var res []TopicHasCheckResponse
	for _, p := range data {
		res = append(res, TopicHasCheckResponse{
			ID:   p.ID,
			Name: p.Name,
		})
	}

	return res, nil
}
