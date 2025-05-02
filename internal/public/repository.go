package public

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAllAuthors(params AuthorPublicParams) ([]AuthorPublicResponse, error)
	GetProfile() (ProfilePublicResponse, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAllAuthors(params AuthorPublicParams) ([]AuthorPublicResponse, error) {
	var datas []AuthorPublicResponse

	// Build the query
	query := r.db.Table("authors").Where("deleted_at IS NULL")

	// Filter by 'name' if provided
	if params.Name != "" {
		query = query.Where("name LIKE ?", "%"+params.Name+"%")
	}

	// Apply sorting if provided
	if params.Sort != "" && params.Order != "" {
		query = query.Order(params.Order + " " + params.Sort) // Dynamically apply sorting (e.g., id ASC)
	}

	// Apply pagination (LIMIT and OFFSET)
	if params.Limit > 0 {
		query = query.Offset((params.Page - 1) * params.Limit).Limit(params.Limit)
	}

	// Execute the query
	err := query.Find(&datas).Error
	return datas, err
}

func (r *repository) GetProfile() (ProfilePublicResponse, error) {
	var about AboutPublicResponse

	err := r.db.Table("abouts").Where("is_used = ? AND deleted_at IS NULL", 1).Scan(&about).Error
	if err != nil {
		return ProfilePublicResponse{}, err
	}

	var technologies []TechnologyProfilePublicResponse
	err = r.db.Table("technologies").Where("is_major = ? AND deleted_at IS NULL", 1).Scan(&technologies).Error
	if err != nil {
		return ProfilePublicResponse{}, err
	}

	data := ProfilePublicResponse{
		About:        about,
		Technologies: technologies,
	}
	return data, nil
}
