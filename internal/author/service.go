package author

type Service interface {
	GetAllAuthors() ([]AuthorResponse, error)
	GetAuthorById(id string) (AuthorResponse, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllAuthors() ([]AuthorResponse, error) {
	authors, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []AuthorResponse
	for _, p := range authors {
		result = append(result, ToAuthorResponse(p))
	}
	return result, nil
}

func (s *service) GetAuthorById(id string) (AuthorResponse, error) {
	author, err := s.repo.FindById(id)
	if err != nil {
		return AuthorResponse{}, err
	}
	return ToAuthorResponse(author), nil
}
