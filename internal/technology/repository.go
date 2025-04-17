package technology

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Technology, error)
	FindById(id string) (Technology, error)
	CreateTechnology(p CreateTechnologyRequest) (Technology, error)
	UpdateTechnology(p UpdateTechnologyDTO) (Technology, error)
	DeleteTechnology(id int) (Technology, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Technology, error) {
	var technologies []Technology
	err := r.db.Find(&technologies).Error
	return technologies, err
}

func (r *repository) FindById(id string) (Technology, error) {
	var technology Technology
	err := r.db.Where("id = ?", id).First(&technology).Error
	return technology, err
}

func (r *repository) CreateTechnology(p CreateTechnologyRequest) (Technology, error) {
	technology := Technology{
		Name:            p.Name,
		DescriptionHTML: p.DescriptionHTML,
		LogoUrl:         p.LogoUrl,
		LogoFileName:    p.LogoFileName,
		IsMajor:         p.IsMajor == "Y"}
	err := r.db.Create(&technology).Error
	return technology, err
}

func (r *repository) UpdateTechnology(p UpdateTechnologyDTO) (Technology, error) {
	technology := Technology{
		ID:              p.Id,
		Name:            p.Name,
		DescriptionHTML: p.DescriptionHTML,
		LogoUrl:         p.LogoUrl,
		LogoFileName:    p.LogoFileName,
		IsMajor:         p.IsMajor == "Y"}
	err := r.db.Updates(&technology).Error
	return technology, err
}

func (r *repository) DeleteTechnology(id int) (Technology, error) {
	var technology Technology

	// Step 1: Find by ID
	if err := r.db.First(&technology, id).Error; err != nil {
		return Technology{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&technology).Error; err != nil {
		return Technology{}, err
	}

	// Step 3: Return the data
	return technology, nil
}
