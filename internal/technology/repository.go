package technology

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Technology, error)
	FindById(id int) (Technology, error)
	CreateTechnology(p CreateTechnologyDTO) (Technology, error)
	UpdateTechnology(p UpdateTechnologyDTO) error
	DeleteTechnology(id int) (Technology, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Technology, error) {
	var datas []Technology
	err := r.db.Find(&datas).Error
	return datas, err
}

func (r *repository) FindById(id int) (Technology, error) {
	var data Technology
	err := r.db.Where("id = ?", id).First(&data).Error
	if err == gorm.ErrRecordNotFound {
		return Technology{}, gorm.ErrRecordNotFound
	}
	return data, err
}

func (r *repository) CreateTechnology(p CreateTechnologyDTO) (Technology, error) {
	data := Technology{
		Name:            p.Name,
		DescriptionHTML: p.DescriptionHTML,
		LogoUrl:         p.LogoUrl,
		LogoFileName:    p.LogoFileName,
		IsMajor:         p.IsMajor,
	}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateTechnology(p UpdateTechnologyDTO) error {

	updateMap := map[string]interface{}{
		"id":               p.ID,
		"name":             p.Name,
		"description_html": p.DescriptionHTML,
		"logo_url":         p.LogoUrl,
		"logo_file_name":   p.LogoFileName,
		"is_major":         p.IsMajor,
	}
	err := r.db.Table("technologies").Where("id = ?", p.ID).Updates(updateMap).Error
	if err != nil {
		return err
	}

	return err
}

func (r *repository) DeleteTechnology(id int) (Technology, error) {
	var data Technology

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return Technology{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return Technology{}, err
	}

	// Step 3: Return the data
	return data, nil
}
