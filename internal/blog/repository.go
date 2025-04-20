package blog

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Blog, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Blog, error) {
	var datas []Blog
	err := r.db.Find(&datas).Error
	return datas, err
}
