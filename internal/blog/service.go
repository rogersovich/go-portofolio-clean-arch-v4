package blog

import (
	"fmt"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/author"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/topic"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

type Service interface {
	GetAllBlogs() ([]BlogResponse, error)
	CreateBlog(p CreateBlogRequest) (BlogResponse, error)
}

type service struct {
	authorService author.Service
	topicService  topic.Service
	blogRepo      Repository
}

func NewService(authorSvc author.Service, topicSvc topic.Service, r Repository) Service {
	return &service{
		authorService: authorSvc,
		topicService:  topicSvc,
		blogRepo:      r,
	}
}

func (s *service) GetAllBlogs() ([]BlogResponse, error) {
	datas, err := s.blogRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []BlogResponse
	for _, p := range datas {
		result = append(result, ToBlogResponse(p))
	}
	return result, nil
}

func (s *service) CreateBlog(p CreateBlogRequest) (BlogResponse, error) {
	//todo: Check Author Id
	//ex: author_id = 1
	_, err := s.authorService.GetAuthorById(p.AuthorID)

	if err != nil {
		err = fmt.Errorf("author_id %d not found", p.AuthorID)
		return BlogResponse{}, err
	}

	//todo: Check Topic Ids
	//ex: topic_ids = [1, 2, 3]
	topic_ids, _ := utils.ConvertStringSliceToIntSlice(p.TopicIds)
	_, err = s.topicService.CheckTopicIds(topic_ids)
	if err != nil {
		return BlogResponse{}, err
	}

	//todo: Create Blog
	return BlogResponse{}, nil
}
