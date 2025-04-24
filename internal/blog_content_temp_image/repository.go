package blog_content_temp_image

import (
	"context"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]BlogContentTempImages, error)
	FindById(id string) (BlogContentTempImages, error)
	CreateBlogContentTempImg(p CreateBlogContentTempImgRequest) (BlogContentTempImages, error)
	UpdateBlogContentTempImg(p UpdateBlogContentTempImgRequest) (BlogContentTempImages, error)
	DeleteBlogContentTempImg(id int) (BlogContentTempImages, error)
	CountTempImages(tempImages []CountTempImagesDTO) (total int, err error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]BlogContentTempImages, error) {
	var datas []BlogContentTempImages
	err := r.db.Find(&datas).Error
	return datas, err
}

func (r *repository) FindById(id string) (BlogContentTempImages, error) {
	var data BlogContentTempImages
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) CreateBlogContentTempImg(p CreateBlogContentTempImgRequest) (BlogContentTempImages, error) {
	data := BlogContentTempImages{
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateBlogContentTempImg(p UpdateBlogContentTempImgRequest) (BlogContentTempImages, error) {
	data := BlogContentTempImages{
		ID:            p.Id,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName}
	err := r.db.Updates(&data).Error
	return data, err
}

func (r *repository) DeleteBlogContentTempImg(id int) (BlogContentTempImages, error) {
	var data BlogContentTempImages

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return BlogContentTempImages{}, err // return if not found or any error
	}

	oldFileName := data.ImageFileName

	// Step 2: Permanently Delete
	if err := r.db.Unscoped().Delete(&data).Error; err != nil {
		return BlogContentTempImages{}, err
	}

	err := utils.DeleteFromMinio(context.Background(), oldFileName)
	if err != nil {
		utils.Logger.Error(err.Error())
	}

	// Step 3: Return the data
	return data, nil
}

func (r *repository) CountTempImages(tempImages []CountTempImagesDTO) (total int, err error) {
	// Prepare slices for image URLs and IDs
	var imageUrls []string
	var ids []int

	// Extract image URLs and IDs from tempImages
	for _, tempImage := range tempImages {
		imageUrls = append(imageUrls, tempImage.ImageUrl)
		ids = append(ids, tempImage.ID)
	}

	err = r.db.Raw(`
		SELECT COUNT(*) FROM blog_content_temp_images 
		WHERE 
			image_url IN ? AND
			id IN ? AND
			deleted_at IS NULL
	`, imageUrls, ids).Scan(&total).Error

	return total, err
}
