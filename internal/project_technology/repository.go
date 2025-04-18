package project_technology

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]ProjectTechnology, error)
	FindById(id string) (ProjectTechnology, error)
	CreateProjectTechnology(p CreateProjectTechnologyRequest) (ProjectTechnology, error)
	UpdateProjectTechnology(p UpdateProjectTechnologyRequest) (ProjectTechnology, error)
	DeleteProjectTechnology(id int) (ProjectTechnology, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]ProjectTechnology, error) {
	var datas []ProjectTechnology
	err := r.db.Find(&datas).Error
	return datas, err
}

func (r *repository) FindById(id string) (ProjectTechnology, error) {
	var data ProjectTechnology
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) CreateProjectTechnology(p CreateProjectTechnologyRequest) (ProjectTechnology, error) {
	data := ProjectTechnology{
		ProjectID:    p.ProjectID,
		TechnologyID: p.TechnologyID}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateProjectTechnology(p UpdateProjectTechnologyRequest) (ProjectTechnology, error) {
	data := ProjectTechnology{
		ID:           p.Id,
		ProjectID:    p.ProjectID,
		TechnologyID: p.TechnologyID}
	err := r.db.Updates(&data).Error
	return data, err
}

func (r *repository) DeleteProjectTechnology(id int) (ProjectTechnology, error) {
	var data ProjectTechnology

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return ProjectTechnology{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return ProjectTechnology{}, err
	}

	// Step 3: Return the data
	return data, nil
}
