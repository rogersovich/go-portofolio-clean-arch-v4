package blog

import (
	"fmt"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/author"
)

type Service interface {
	GetAllBlogs() ([]BlogResponse, error)
	CreateBlog(p CreateBlogRequest) (BlogResponse, error)
}

type service struct {
	authorService author.Service
	blogRepo      Repository
}

func NewService(authorSvc author.Service, r Repository) Service {
	return &service{
		authorService: authorSvc,
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
	_, err := s.authorService.GetAuthorById(p.AuthorID)

	if err != nil {
		err = fmt.Errorf("author_id %d not found", p.AuthorID)
		return BlogResponse{}, err // DB error
	}

	//todo: Check Topic Ids

	//todo: Create Blog
	return BlogResponse{}, nil
}
