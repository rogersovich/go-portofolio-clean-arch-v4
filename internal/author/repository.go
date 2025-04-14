package author

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Author, error)
	FindById(id string) (Author, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Author, error) {
	var abouts []Author
	err := r.db.Find(&abouts).Error
	return abouts, err
}

func (r *repository) FindById(id string) (Author, error) {
	var about Author
	err := r.db.Where("id = ?", id).First(&about).Error
	return about, err
}
