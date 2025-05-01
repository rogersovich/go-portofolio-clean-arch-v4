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
	CountExistingProjectImages(projectImages []ProjectImagesExistingPayload) (total int, err error)
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

func (r *repository) CountExistingProjectImages(projectImages []ProjectImagesExistingPayload) (total int, err error) {
	var image_urls []string
	for _, v := range projectImages {
		image_urls = append(image_urls, v.ImageUrl)
	}
	err = r.db.Raw(`
		SELECT COUNT(*) FROM project_content_temp_images 
		WHERE image_url IN ? AND
		deleted_at IS NULL
	`, image_urls).Scan(&total).Error
	return total, err
}
