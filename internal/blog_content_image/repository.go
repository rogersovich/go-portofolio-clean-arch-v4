package blog_content_image

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]BlogContentImage, error)
	FindById(id int) (BlogContentImage, error)
	CreateBlogContentImage(p CreateBlogContentImageRequest) (BlogContentImage, error)
	UpdateBlogContentImage(p UpdateBlogContentImageDTO) (BlogContentImage, error)
	DeleteBlogContentImage(id int) (BlogContentImage, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]BlogContentImage, error) {
	var datas []BlogContentImage
	err := r.db.Find(&datas).Error
	return datas, err
}

func (r *repository) FindById(id int) (BlogContentImage, error) {
	var data BlogContentImage
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) CreateBlogContentImage(p CreateBlogContentImageRequest) (BlogContentImage, error) {
	data := BlogContentImage{
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		IsUsed:        p.IsUsed == "Y"}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateBlogContentImage(p UpdateBlogContentImageDTO) (BlogContentImage, error) {
	data := BlogContentImage{
		ID:            p.Id,
		BlogID:        p.BlogID,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		IsUsed:        p.IsUsed == "Y"}
	err := r.db.Updates(&data).Error
	return data, err
}

func (r *repository) DeleteBlogContentImage(id int) (BlogContentImage, error) {
	var data BlogContentImage

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return BlogContentImage{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return BlogContentImage{}, err
	}

	// Step 3: Return the data
	return data, nil
}
