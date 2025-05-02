package author

import (
	"context"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

type Service interface {
	GetAllAuthors() ([]AuthorResponse, error)
	GetAuthorById(id int) (AuthorResponse, error)
	CreateAuthor(p CreateAuthorRequest) (AuthorResponse, error)
	UpdateAuthor(p UpdateAuthorRequest) error
	DeleteAuthor(id int) error
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
	avatarRes, err := utils.HandlUploadFile(p.AvatarFile, "author")
	if err != nil {
		return AuthorResponse{}, err
	}

	payload := CreateAuthorDTO{
		Name:           p.Name,
		AvatarUrl:      avatarRes.FileURL,
		AvatarFileName: avatarRes.FileName,
	}

	author, err := s.repo.CreateAuthor(payload)
	if err != nil {
		_ = utils.DeleteFromMinio(context.Background(), avatarRes.FileName)
		return AuthorResponse{}, err
	}
	return ToAuthorResponse(author), nil
}

func (s *service) UpdateAuthor(p UpdateAuthorRequest) error {
	//todo: Get Author
	author, err := s.repo.FindById(p.ID)
	if err != nil {
		return err
	}

	//todo: set oldFileName
	oldFileName := ""
	if author.AvatarFileName != "" {
		oldFileName = author.AvatarFileName
	}

	var newFileURL string
	var newFileName string

	//todo: Upload File
	if p.AvatarFile != nil {
		logoRes, err := utils.HandlUploadFile(p.AvatarFile, "author")
		if err != nil {
			return err
		}

		newFileURL = logoRes.FileURL
		newFileName = logoRes.FileName
	} else {
		newFileURL = author.AvatarUrl
		newFileName = author.AvatarFileName
	}

	payload := UpdateAuthorDTO{
		ID:             p.ID,
		Name:           p.Name,
		AvatarUrl:      newFileURL,
		AvatarFileName: newFileName,
	}

	err = s.repo.UpdateAuthor(payload)
	if err != nil {
		_ = utils.DeleteFromMinio(context.Background(), newFileName)
		return err
	}

	//todo: Delete Old Image
	if oldFileName != newFileName {
		_ = utils.DeleteFromMinio(context.Background(), oldFileName)
	}

	return nil
}

func (s *service) DeleteAuthor(id int) error {
	author, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	err = s.repo.DeleteAuthor(id)
	if err != nil {
		return err
	}

	//todo: Delete Old Image
	_ = utils.DeleteFromMinio(context.Background(), author.AvatarFileName)

	return nil
}
