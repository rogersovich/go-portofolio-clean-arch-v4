package author

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Author, error)
	FindById(id string) (Author, error)
	CreateAuthor(p CreateAuthorRequest) (Author, error)
	UpdateAuthor(p UpdateAuthorDTO) (Author, error)
	DeleteAuthor(id int) (Author, error)
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

func (r *repository) CreateAuthor(p CreateAuthorRequest) (Author, error) {
	about := Author{
		Name:           p.Name,
		AvatarUrl:      p.AvatarUrl,
		AvatarFileName: p.AvatarFileName}
	err := r.db.Create(&about).Error
	return about, err
}

func (r *repository) UpdateAuthor(p UpdateAuthorDTO) (Author, error) {
	author := Author{
		ID:             p.Id,
		Name:           p.Name,
		AvatarUrl:      p.AvatarUrl,
		AvatarFileName: p.AvatarFileName}
	err := r.db.Updates(&author).Error
	return author, err
}

func (r *repository) DeleteAuthor(id int) (Author, error) {
	var author Author

	// Step 1: Find the author by ID
	if err := r.db.First(&author, id).Error; err != nil {
		return Author{}, err // return if not found or any error
	}

	// Step 2: Delete the found author
	if err := r.db.Delete(&author).Error; err != nil {
		return Author{}, err
	}

	// Step 3: Return the deleted author's data
	return author, nil
}
