package public

import "gorm.io/gorm"

type Repository interface {
	FindAllAuthors() ([]AuthorPublicResponse, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAllAuthors() ([]AuthorPublicResponse, error) {
	var datas []AuthorPublicResponse
	err := r.db.Table("authors").Where("deleted_at IS NULL").Find(&datas).Error
	return datas, err
}
