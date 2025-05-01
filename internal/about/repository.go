package about

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]About, error)
	FindById(id int) (About, error)
	CreateAbout(p CreateAboutDTO) (About, error)
	UpdateAbout(p UpdateAboutDTO) error
	DeleteAbout(id int) (About, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]About, error) {
	var abouts []About
	err := r.db.Find(&abouts).Error
	return abouts, err
}

func (r *repository) FindById(id int) (About, error) {
	var about About
	err := r.db.Where("id = ?", id).First(&about).Error
	return about, err
}

func (r *repository) CreateAbout(p CreateAboutDTO) (About, error) {
	about := About{
		Title:           p.Title,
		DescriptionHTML: p.DescriptionHTML,
		AvatarUrl:       p.AvatarUrl,
		AvatarFileName:  p.AvatarFileName}
	err := r.db.Create(&about).Error
	return about, err
}

func (r *repository) UpdateAbout(p UpdateAboutDTO) error {
	about := About{
		ID:              p.ID,
		Title:           p.Title,
		DescriptionHTML: p.DescriptionHTML,
		AvatarUrl:       p.AvatarUrl,
		AvatarFileName:  p.AvatarFileName}
	err := r.db.Updates(&about).Error
	return err
}

func (r *repository) DeleteAbout(id int) (About, error) {
	var about About

	// Step 1: Find by ID
	if err := r.db.First(&about, id).Error; err != nil {
		return About{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&about).Error; err != nil {
		return About{}, err
	}

	// Step 3: Return the data
	return about, nil
}
