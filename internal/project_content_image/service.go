package project_content_image

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
	"gorm.io/gorm"
)

type Service interface {
	GetAllProjectContentImages() ([]ProjectContentImageResponse, error)
	GetProjectContentImageById(id int) (ProjectContentImageResponse, error)
	CreateProjectContentImage(p CreateProjectContentImageRequest) (ProjectContentImageResponse, error)
	UpdateProjectContentImage(p UpdateProjectContentImageRequest) error
	DeleteProjectContentImage(id int) (ProjectContentImageResponse, error)
	CountUnusedProjectImages(ids []string) error
	BatchUpdateProjectImages(projectImages []string, project_id int, tx *gorm.DB) error
	SyncProjectImages(image_urls []string, project_id int, tx *gorm.DB) ([]ProjectImagesFindResponse, error)
	BulkDeleteHardByImageUrls(image_urls []string, tx *gorm.DB) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllProjectContentImages() ([]ProjectContentImageResponse, error) {
	datas, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []ProjectContentImageResponse
	for _, p := range datas {
		result = append(result, ToProjectContentImageResponse(p))
	}
	return result, nil
}

func (s *service) GetProjectContentImageById(id int) (ProjectContentImageResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return ProjectContentImageResponse{}, err
	}
	return ToProjectContentImageResponse(data), nil
}

func (s *service) CreateProjectContentImage(p CreateProjectContentImageRequest) (ProjectContentImageResponse, error) {
	imageRes, err := utils.HandlUploadFile(p.ImageFile, "project")
	if err != nil {
		return ProjectContentImageResponse{}, err
	}

	payload := CreateProjectContentImageDTO{
		ProjectID:     nil,
		ImageUrl:      imageRes.FileURL,
		ImageFileName: imageRes.FileName,
	}

	data, err := s.repo.CreateProjectContentImage(payload)
	if err != nil {
		_ = utils.DeleteFromMinio(context.Background(), imageRes.FileName)
		return ProjectContentImageResponse{}, err
	}
	return ToProjectContentImageResponse(data), nil
}

func (s *service) UpdateProjectContentImage(p UpdateProjectContentImageRequest) error {
	//todo: Get Data
	projectImage, err := s.repo.FindById(p.ID)
	if err != nil {
		return err
	}

	//todo: set oldFileName
	oldFileName := ""
	if projectImage.ImageFileName != "" {
		oldFileName = projectImage.ImageFileName
	}

	var newFileURL string
	var newFileName string

	//todo: Upload File
	if p.ImageFile != nil {
		imageRes, err := utils.HandlUploadFile(p.ImageFile, "project")
		if err != nil {
			return err
		}

		newFileURL = imageRes.FileURL
		newFileName = imageRes.FileName
	} else {
		newFileURL = projectImage.ImageUrl
		newFileName = projectImage.ImageFileName
	}

	payload := UpdateProjectContentImageDTO{
		ID:            p.ID,
		ProjectID:     p.ProjectID,
		ImageUrl:      newFileURL,
		ImageFileName: newFileName,
	}

	err = s.repo.UpdateProjectContentImage(payload)
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

func (s *service) DeleteProjectContentImage(id int) (ProjectContentImageResponse, error) {
	data, err := s.repo.DeleteProjectContentImage(id)
	if err != nil {
		return ProjectContentImageResponse{}, err
	}

	_ = utils.DeleteFromMinio(context.Background(), data.ImageFileName)

	return ToProjectContentImageResponse(data), nil
}

func (s *service) CountUnusedProjectImages(ids []string) error {
	total, err := s.repo.CountUnusedProjectImages(ids)
	if err != nil {
		return err
	}

	if total != len(ids) {
		err := fmt.Errorf("some project_images not found in database")
		return err
	}
	return nil
}

func (s *service) BatchUpdateProjectImages(projectImages []string, project_id int, tx *gorm.DB) error {
	return s.repo.BatchUpdateProjectImages(projectImages, project_id, tx)
}

func (s *service) SyncProjectImages(
	image_urls []string,
	project_id int,
	tx *gorm.DB) (
	imageNotExist []ProjectImagesFindResponse,
	err error,
) {
	// 1. Cek apakah ada image baru di konten yang belum ada di database
	imageExist, err := s.repo.FindImageExist(image_urls, project_id)
	if err != nil {
		return imageNotExist, err
	}

	// Buat map untuk cepat cek
	imageURLMap := make(map[string]*ProjectImagesFindResponse)
	for _, img := range imageExist {
		imageURLMap[img.ImageUrl] = &img
	}

	var imageIDsToUpdate []int
	for _, url := range image_urls {
		img, found := imageURLMap[url]
		if found && img.ProjectID == nil {
			imageIDsToUpdate = append(imageIDsToUpdate, img.ID)
		}
	}

	if len(imageIDsToUpdate) > 0 {
		if err := s.repo.BatchUpdateImagesById(imageIDsToUpdate, project_id, tx); err != nil {
			return imageNotExist, err
		}
	}

	// 4. Hapus image lama yang tidak ada di konten lagi
	imageNotExist, err = s.repo.FindImageNotExist(image_urls, project_id)
	if err != nil {
		return imageNotExist, err
	}

	return imageNotExist, nil
}

func (s *service) BulkDeleteHardByImageUrls(image_urls []string, tx *gorm.DB) error {
	err := s.repo.BulkDeleteHardByImageUrls(image_urls, tx)
	if err != nil {
		return err
	}

	bucketName := os.Getenv("MINIO_BUCKET")
	images_key, _ := utils.MinioParseURLToImageKey(image_urls, bucketName)
	batchSize := 3

	err = utils.DeleteBulkImagesInBatches(bucketName, images_key, batchSize)
	if err != nil {
		log.Fatalf("Failed to delete images: %v", err)
	}

	return nil
}
