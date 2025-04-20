package author

import (
	"context"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

type Service interface {
	GetAllAuthors() ([]AuthorResponse, error)
	GetAuthorById(id int) (AuthorResponse, error)
	CreateAuthor(p CreateAuthorRequest) (AuthorResponse, error)
	UpdateAuthor(p UpdateAuthorDTO, oldPath string, newFilePath string) (AuthorUpdateResponse, error)
	DeleteAuthor(id int) (Author, error)
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

func (s *service) GetAuthorById(id int) (AuthorResponse, error) {
	author, err := s.repo.FindById(id)
	if err != nil {
		return AuthorResponse{}, err
	}
	return ToAuthorResponse(author), nil
}

func (s *service) CreateAuthor(p CreateAuthorRequest) (AuthorResponse, error) {
	author, err := s.repo.CreateAuthor(p)
	if err != nil {
		return AuthorResponse{}, err
	}
	return ToAuthorResponse(author), nil
}

func (s *service) UpdateAuthor(p UpdateAuthorDTO, oldPath string, newFilePath string) (AuthorUpdateResponse, error) {
	author, err := s.repo.UpdateAuthor(p)
	if err != nil {
		return AuthorUpdateResponse{}, err
	}

	// 3. Optional: Delete old file from MinIO
	if oldPath != newFilePath {
		err = utils.DeleteFromMinio(context.Background(), oldPath) // ignore error or handle if needed
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}

	return ToAuthorUpdateResponse(author), nil
}

func (s *service) DeleteAuthor(id int) (Author, error) {
	author, err := s.repo.DeleteAuthor(id)
	if err != nil {
		return Author{}, err
	}
	return author, nil
}
