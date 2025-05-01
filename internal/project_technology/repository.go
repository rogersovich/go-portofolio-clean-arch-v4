package project_technology

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]ProjectTechnology, error)
	FindById(id int) (ProjectTechnology, error)
	CreateProjectTechnology(p CreateProjectTechnologyRequest) (ProjectTechnology, error)
	UpdateProjectTechnology(p UpdateProjectTechnologyRequest) error
	DeleteProjectTechnology(id int) (ProjectTechnology, error)
	CountTechnologiesByIDs(ids []int) (total int, err error)
	BulkCreateTechnologies(tech_ids []ProjectTechnology, tx *gorm.DB) error
	FindExistingProjectTechnologies(project_id int) ([]ProjectTechnologyExistingResponse, error)
	BulkDeleteHard(tech_ids []int, tx *gorm.DB) error
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

func (r *repository) FindById(id int) (ProjectTechnology, error) {
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

func (r *repository) UpdateProjectTechnology(p UpdateProjectTechnologyRequest) error {
	data := ProjectTechnology{
		ID:           p.ID,
		ProjectID:    p.ProjectID,
		TechnologyID: p.TechnologyID}
	err := r.db.Updates(&data).Error
	return err
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

func (r *repository) CountTechnologiesByIDs(ids []int) (total int, err error) {
	err = r.db.Raw(`
		SELECT COUNT(*) FROM technologies 
		WHERE id IN ? AND
		deleted_at IS NULL
	`, ids).Scan(&total).Error
	return total, err
}

func (r *repository) BulkCreateTechnologies(tech_ids []ProjectTechnology, tx *gorm.DB) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	if err := db.Create(&tech_ids).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) FindExistingProjectTechnologies(project_id int) ([]ProjectTechnologyExistingResponse, error) {
	var data []ProjectTechnologyExistingResponse
	err := r.db.Table("project_technologies").
		Where("project_id = ?", project_id).
		Select("id, project_id, technology_id").
		Find(&data).Error
	return data, err
}

func (r *repository) BulkDeleteHard(tech_ids []int, tx *gorm.DB) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	// Create a raw SQL query to delete records with IDs in the slice
	query := "DELETE FROM project_technologies WHERE technology_id IN ?"

	// Execute the raw query
	if err := db.Exec(query, tech_ids).Error; err != nil {
		return err
	}

	return nil
}
