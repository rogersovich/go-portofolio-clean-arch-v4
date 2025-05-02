package author

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Author, error)
	FindById(id int) (Author, error)
	CreateAuthor(p CreateAuthorDTO) (Author, error)
	UpdateAuthor(p UpdateAuthorDTO) error
	DeleteAuthor(id int) error
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

func (r *repository) FindById(id int) (Author, error) {
	var data Author
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) CreateAuthor(p CreateAuthorDTO) (Author, error) {
	about := Author{
		Name:           p.Name,
		AvatarUrl:      p.AvatarUrl,
		AvatarFileName: p.AvatarFileName}
	err := r.db.Create(&about).Error
	return about, err
}

func (r *repository) UpdateAuthor(p UpdateAuthorDTO) error {
	author := Author{
		ID:             p.ID,
		Name:           p.Name,
		AvatarUrl:      p.AvatarUrl,
		AvatarFileName: p.AvatarFileName}
	err := r.db.Updates(&author).Error
	return err
}

func (r *repository) DeleteAuthor(id int) error {
	// Hard Delete
	if err := r.db.Unscoped().Where("id = ?", id).Delete(&Author{}).Error; err != nil {
		return err
	}

	// Return the data
	return nil
}
