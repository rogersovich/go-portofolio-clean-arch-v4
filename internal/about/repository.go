package about

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]About, error)
	FindById(id int) (About, error)
	CreateAbout(p CreateAboutDTO) (About, error)
	UpdateAbout(p UpdateAboutDTO) error
	DeleteAbout(id int) error
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
	updateMap := map[string]interface{}{
		"id":               p.ID,
		"title":            p.Title,
		"description_html": p.DescriptionHTML,
		"avatar_url":       p.AvatarUrl,
		"avatar_file_name": p.AvatarFileName,
		"is_used":          p.IsUsed,
	}
	err := r.db.Table("abouts").Where("id = ?", p.ID).Updates(updateMap).Error
	return err
}

func (r *repository) DeleteAbout(id int) error {

	// Hard Delete
	if err := r.db.Unscoped().Where("id = ?", id).Delete(&About{}).Error; err != nil {
		return err
	}

	// Return the data
	return nil
}
