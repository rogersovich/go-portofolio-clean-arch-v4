package blog_content_image

import (
	"context"
	"fmt"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
	"gorm.io/gorm"
)

type Service interface {
	GetAllBlogContentImages() ([]BlogContentImageResponse, error)
	GetBlogContentImageById(id int) (BlogContentImageResponse, error)
	CreateBlogContentImage(p CreateBlogContentImageRequest) (BlogContentImageResponse, error)
	UpdateBlogContentImage(p UpdateBlogContentImageDTO, oldPath string, newFilePath string) (BlogContentImageUpdateResponse, error)
	DeleteBlogContentImage(id int) (BlogContentImageResponse, error)
	CheckHasBlogImages(image_urls []string) error
	MarkImagesUsedByBlog(image_urls []string, blog_id int, tx *gorm.DB) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllBlogContentImages() ([]BlogContentImageResponse, error) {
	datas, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []BlogContentImageResponse
	for _, p := range datas {
		result = append(result, ToBlogContentImageResponse(p))
	}
	return result, nil
}

func (s *service) GetBlogContentImageById(id int) (BlogContentImageResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return BlogContentImageResponse{}, err
	}
	return ToBlogContentImageResponse(data), nil
}

func (s *service) CreateBlogContentImage(p CreateBlogContentImageRequest) (BlogContentImageResponse, error) {
	data, err := s.repo.CreateBlogContentImage(p)
	if err != nil {
		return BlogContentImageResponse{}, err
	}
	return ToBlogContentImageResponse(data), nil
}

func (s *service) UpdateBlogContentImage(p UpdateBlogContentImageDTO, oldPath string, newFilePath string) (BlogContentImageUpdateResponse, error) {
	data, err := s.repo.UpdateBlogContentImage(p)
	if err != nil {
		return BlogContentImageUpdateResponse{}, err
	}

	// 3. Optional: Delete old file from MinIO
	if oldPath != newFilePath {
		err = utils.DeleteFromMinio(context.Background(), oldPath) // ignore error or handle if needed
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}

	return ToBlogContentImageUpdateResponse(data), nil
}

func (s *service) DeleteBlogContentImage(id int) (BlogContentImageResponse, error) {
	data, err := s.repo.DeleteBlogContentImage(id)
	if err != nil {
		return BlogContentImageResponse{}, err
	}
	return ToBlogContentImageResponse(data), nil
}

func (s *service) CheckHasBlogImages(image_urls []string) error {
	total, err := s.repo.CheckHasBlogImages(image_urls)
	if err != nil {
		return err
	}

	if total != len(image_urls) {
		err := fmt.Errorf("some project_content_images not found in database")
		return err
	}
	return nil
}

func (s *service) MarkImagesUsedByBlog(image_urls []string, blog_id int, tx *gorm.DB) error {
	payload := BlogContentImageBulkUpdateDTO{
		ImageUrls: image_urls,
		BlogID:    blog_id,
	}
	err := s.repo.MarkImagesUsedByBlog(payload, tx)
	if err != nil {
		return err
	}

	return nil
}
