package blog_topic

import (
	"slices"

	"gorm.io/gorm"
)

type Service interface {
	GetAllBlogTopics() ([]BlogTopicResponse, error)
	GetBlogTopicById(id int) (BlogTopicResponse, error)
	CreateBlogTopic(p CreateBlogTopicRequest) (BlogTopicResponse, error)
	BulkCreateBlogTopic(topic_ids []int, project_id int, tx *gorm.DB) error
	UpdateBlogTopic(p UpdateBlogTopicRequest) error
	DeleteBlogTopic(id int) (BlogTopic, error)
	BatchUpdateBlogTopic(topic_ids []int, blog_id int, tx *gorm.DB) error
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

func (s *service) GetBlogTopicById(id int) (BlogTopicResponse, error) {
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

func (s *service) BulkCreateBlogTopic(topic_ids []int, blog_id int, tx *gorm.DB) error {
	err := s.repo.BulkCreateBlogTopic(topic_ids, blog_id, tx)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateBlogTopic(p UpdateBlogTopicRequest) error {
	_, err := s.repo.FindById(p.ID)
	if err != nil {
		return err
	}

	err = s.repo.UpdateBlogTopic(p)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteBlogTopic(id int) (BlogTopic, error) {
	data, err := s.repo.DeleteBlogTopic(id)
	if err != nil {
		return BlogTopic{}, err
	}
	return data, nil
}

func (s *service) BulkDeleteHard(topic_ids []int, tx *gorm.DB) error {
	err := s.repo.BulkDeleteHard(topic_ids, tx)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) BatchUpdateBlogTopic(topic_ids []int, blog_id int, tx *gorm.DB) error {
	existingBlogTopics, err := s.repo.FindExistingBlogTopics(blog_id)
	if err != nil {
		return err
	}

	var existingtopic_ids []int
	var isNewTopic bool = true
	for _, item := range existingBlogTopics {
		if slices.Contains(topic_ids, item.TopicID) {
			isNewTopic = false
		}
		existingtopic_ids = append(existingtopic_ids, item.TopicID)
	}

	if !isNewTopic {
		return nil
	}

	err = s.repo.BulkDeleteHard(existingtopic_ids, tx)
	if err != nil {
		return err
	}

	err = s.repo.BulkCreateBlogTopic(topic_ids, blog_id, tx)
	if err != nil {
		return err
	}

	return nil
}
