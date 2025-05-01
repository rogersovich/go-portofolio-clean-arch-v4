package blog_content_image

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]BlogContentImage, error)
	FindById(id int) (BlogContentImage, error)
	CreateBlogContentImage(p CreateBlogContentImageDTO) (BlogContentImage, error)
	UpdateBlogContentImage(p UpdateBlogContentImageDTO) error
	DeleteBlogContentImage(id int) (BlogContentImage, error)
	CountUnlinkedImages(image_urls []string) (total int, err error)
	MarkImagesUsedByBlog(p BlogContentImageBulkUpdateDTO, tx *gorm.DB) error
	CountImagesLinkedToBlog(image_urls []string, blog_id int) (total int, err error)
	FindImageExist(image_urls []string, blog_id int) ([]BlogContentImageExistingResponse, error)
	FindImageNotExist(image_urls []string, blog_id int) ([]BlogContentImageExistingResponse, error)
	BatchUpdateImagesById(ids []int, blog_id int, tx *gorm.DB) error
	BulkDeleteHardByImageUrls(image_urls []string, tx *gorm.DB) error
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

func (r *repository) FindImageExist(image_urls []string, blog_id int) ([]BlogContentImageExistingResponse, error) {
	var data []BlogContentImageExistingResponse
	err := r.db.Table("blog_content_images").
		Where("image_url IN ? AND (blog_id = ? OR blog_id IS NULL)", image_urls, blog_id).
		Select("id, blog_id, image_url").
		Find(&data).Error
	return data, err
}

func (r *repository) FindImageNotExist(image_urls []string, blog_id int) ([]BlogContentImageExistingResponse, error) {
	var data []BlogContentImageExistingResponse
	err := r.db.Table("blog_content_images").
		Where("blog_id = ? AND image_url NOT IN ?", blog_id, image_urls).
		Select("id, blog_id, image_url").
		Find(&data).Error
	return data, err
}

func (r *repository) CreateBlogContentImage(p CreateBlogContentImageDTO) (BlogContentImage, error) {
	data := BlogContentImage{
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateBlogContentImage(p UpdateBlogContentImageDTO) error {
	data := BlogContentImage{
		ID:            p.ID,
		BlogID:        p.BlogID,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName}
	err := r.db.Updates(&data).Error
	return err
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

func (r *repository) CountUnlinkedImages(image_urls []string) (total int, err error) {
	err = r.db.Raw(`
		SELECT COUNT(*) FROM blog_content_images 
		WHERE image_url IN ? AND
		blog_id IS NULL AND
		deleted_at IS NULL
	`, image_urls).Scan(&total).Error
	return total, err
}

func (r *repository) MarkImagesUsedByBlog(p BlogContentImageBulkUpdateDTO, tx *gorm.DB) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	// UPDATE to TABLE PROJECT CONTENT IMAGES
	err := db.Model(&BlogContentImage{}).
		Where("image_url IN ?", p.ImageUrls).
		Updates(map[string]interface{}{
			"blog_id": p.BlogID,
		}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CountImagesLinkedToBlog(image_urls []string, blog_id int) (total int, err error) {
	err = r.db.Raw(`
		SELECT COUNT(*) FROM blog_content_images 
		WHERE image_url IN ? AND
		blog_id = ? AND
		deleted_at IS NULL
	`, image_urls, blog_id).Scan(&total).Error
	return total, err
}

func (r *repository) BatchUpdateImagesById(ids []int, blog_id int, tx *gorm.DB) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	err := db.Table("blog_content_images").
		Where("id IN ?", ids).
		Update("blog_id", blog_id).Error
	return err
}

func (r *repository) BulkDeleteHardByImageUrls(image_urls []string, tx *gorm.DB) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	// Create a raw SQL query to delete records with IDs in the slice
	query := "DELETE FROM blog_content_images WHERE image_url IN ?"

	// Execute the raw query
	if err := db.Exec(query, image_urls).Error; err != nil {
		return err
	}

	return nil
}
