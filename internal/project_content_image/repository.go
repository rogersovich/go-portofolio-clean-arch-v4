package project_content_image

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]ProjectContentImage, error)
	FindById(id int) (ProjectContentImage, error)
	CreateProjectContentImage(p CreateProjectContentImageRequest) (ProjectContentImage, error)
	UpdateProjectContentImage(p UpdateProjectContentImageDTO) (ProjectContentImage, error)
	DeleteProjectContentImage(id int) (ProjectContentImage, error)
	CountUnusedProjectImages(ids []string) (total int, err error)
	BatchUpdateProjectImages(projectImages []string, project_id int, tx *gorm.DB) error
	FindImageExist(image_urls []string, project_id int) ([]ProjectImagesFindResponse, error)
	FindImageNotExist(image_urls []string, project_id int) ([]ProjectImagesFindResponse, error)
	BatchUpdateImagesById(ids []int, project_id int, tx *gorm.DB) error
	BulkDeleteHardByImageUrls(image_urls []string, tx *gorm.DB) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]ProjectContentImage, error) {
	var datas []ProjectContentImage
	err := r.db.Find(&datas).Error
	return datas, err
}

func (r *repository) FindById(id int) (ProjectContentImage, error) {
	var data ProjectContentImage
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) CreateProjectContentImage(p CreateProjectContentImageRequest) (ProjectContentImage, error) {
	data := ProjectContentImage{
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		IsUsed:        p.IsUsed == "Y"}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateProjectContentImage(p UpdateProjectContentImageDTO) (ProjectContentImage, error) {
	data := ProjectContentImage{
		ID:            p.Id,
		ProjectID:     p.ProjectID,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		IsUsed:        p.IsUsed == "Y"}
	err := r.db.Updates(&data).Error
	return data, err
}

func (r *repository) DeleteProjectContentImage(id int) (ProjectContentImage, error) {
	var data ProjectContentImage

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return ProjectContentImage{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return ProjectContentImage{}, err
	}

	// Step 3: Return the data
	return data, nil
}

func (r *repository) CountUnusedProjectImages(ids []string) (total int, err error) {
	err = r.db.Raw(`
		SELECT COUNT(*) FROM project_content_images 
		WHERE image_url IN ? AND
		project_id IS NULL AND
		deleted_at IS NULL
	`, ids).Scan(&total).Error
	return total, err
}

func (r *repository) BatchUpdateProjectImages(projectImages []string, project_id int, tx *gorm.DB) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	err := db.Model(&ProjectContentImage{}).
		Where("image_url IN ?", projectImages).
		Updates(map[string]interface{}{
			"project_id": project_id,
		}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) BatchUpdateImagesById(ids []int, project_id int, tx *gorm.DB) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	err := db.Table("project_content_images").
		Where("id IN ?", ids).
		Update("project_id", project_id).Error
	return err
}

func (r *repository) FindImageExist(image_urls []string, project_id int) ([]ProjectImagesFindResponse, error) {
	var data []ProjectImagesFindResponse
	err := r.db.Table("project_content_images").
		Where("image_url IN ? AND (project_id = ? OR project_id IS NULL)", image_urls, project_id).
		Select("id, project_id, image_url").
		Find(&data).Error
	return data, err
}

func (r *repository) FindImageNotExist(image_urls []string, project_id int) ([]ProjectImagesFindResponse, error) {
	var data []ProjectImagesFindResponse
	err := r.db.Table("project_content_images").
		Where("project_id = ? AND image_url NOT IN ?", project_id, image_urls).
		Select("id, project_id, image_url").
		Find(&data).Error
	return data, err
}

func (r *repository) BulkDeleteHardByImageUrls(image_urls []string, tx *gorm.DB) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	// Create a raw SQL query to delete records with IDs in the slice
	query := "DELETE FROM project_content_images WHERE image_url IN ?"

	// Execute the raw query
	if err := db.Exec(query, image_urls).Error; err != nil {
		return err
	}

	return nil
}
