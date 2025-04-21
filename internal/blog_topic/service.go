package blog_topic

type Service interface {
	GetAllBlogTopics() ([]BlogTopicResponse, error)
	GetBlogTopicById(id string) (BlogTopicResponse, error)
	CreateBlogTopic(p CreateBlogTopicRequest) (BlogTopicResponse, error)
	UpdateBlogTopic(p UpdateBlogTopicRequest) (BlogTopicUpdateResponse, error)
	DeleteBlogTopic(id int) (BlogTopic, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllBlogTopics() ([]BlogTopicResponse, error) {
	datas, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []BlogTopicResponse
	for _, p := range datas {
		result = append(result, ToBlogTopicResponse(p))
	}
	return result, nil
}

func (s *service) GetBlogTopicById(id string) (BlogTopicResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return BlogTopicResponse{}, err
	}
	return ToBlogTopicResponse(data), nil
}

func (s *service) CreateBlogTopic(p CreateBlogTopicRequest) (BlogTopicResponse, error) {
	data, err := s.repo.CreateBlogTopic(p)
	if err != nil {
		return BlogTopicResponse{}, err
	}
	return ToBlogTopicResponse(data), nil
}

func (s *service) UpdateBlogTopic(p UpdateBlogTopicRequest) (BlogTopicUpdateResponse, error) {
	data, err := s.repo.UpdateBlogTopic(p)
	if err != nil {
		return BlogTopicUpdateResponse{}, err
	}

	return ToBlogTopicUpdateResponse(data), nil
}

func (s *service) DeleteBlogTopic(id int) (BlogTopic, error) {
	data, err := s.repo.DeleteBlogTopic(id)
	if err != nil {
		return BlogTopic{}, err
	}
	return data, nil
}
