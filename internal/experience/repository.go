package experience

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Experience, error)
	FindById(id int) (Experience, error)
	CreateExperience(p CreateExperienceDTO) (Experience, error)
	UpdateExperience(p UpdateExperienceDTO) error
	DeleteExperience(id int) (Experience, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Experience, error) {
	var datas []Experience
	err := r.db.Find(&datas).Error
	return datas, err
}

func (r *repository) FindById(id int) (Experience, error) {
	var data Experience
	err := r.db.Where("id = ?", id).First(&data).Error
	if err == gorm.ErrRecordNotFound {
		return Experience{}, gorm.ErrRecordNotFound
	}
	return data, err
}

func (r *repository) CreateExperience(p CreateExperienceDTO) (Experience, error) {
	data := Experience{
		Position:          p.Position,
		CompanyName:       p.CompanyName,
		WorkType:          p.WorkType,
		Country:           p.Country,
		City:              p.City,
		SummaryHTML:       p.SummaryHTML,
		FromDate:          p.FromDate,
		ToDate:            p.ToDate,
		CompImageUrl:      p.CompImageUrl,
		CompImageFileName: p.CompImageFileName,
		CompWebsiteUrl:    p.CompWebsiteUrl,
		IsCurrent:         p.IsCurrent}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateExperience(p UpdateExperienceDTO) error {
	// Create a map with only the fields that are non-zero
	updateMap := map[string]interface{}{
		"position":             p.Position,
		"company_name":         p.CompanyName,
		"work_type":            p.WorkType,
		"country":              p.Country,
		"city":                 p.City,
		"summary_html":         p.SummaryHTML,
		"from_date":            p.FromDate,
		"to_date":              p.ToDate,
		"comp_image_url":       p.CompImageUrl,
		"comp_image_file_name": p.CompImageFileName,
		"comp_website_url":     p.CompWebsiteUrl,
		"is_current":           p.IsCurrent,
	}

	err := r.db.Table("experiences").Where("id = ?", p.ID).Updates(updateMap).Error
	if err != nil {
		return err
	}

	return err
}

func (r *repository) DeleteExperience(id int) (Experience, error) {
	var data Experience

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return Experience{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return Experience{}, err
	}

	// Step 3: Return the data
	return data, nil
}
