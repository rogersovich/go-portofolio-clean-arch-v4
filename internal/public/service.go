package public

import "github.com/rogersovich/go-portofolio-clean-arch-v4/internal/author"

type Service interface {
	GetAllPublicAuthors(params AuthorPublicParams) ([]AuthorPublicResponse, error)
}

type service struct {
	authorService author.Service
	repo          Repository
}

func NewService(authorSvc author.Service, r Repository) Service {
	return &service{
		authorService: authorSvc,
		repo:          r,
	}
}

func (s *service) GetAllPublicAuthors(params AuthorPublicParams) ([]AuthorPublicResponse, error) {
	data, err := s.repo.FindAllAuthors(params)
	if err != nil {
		return nil, err
	}

	return data, nil
}
