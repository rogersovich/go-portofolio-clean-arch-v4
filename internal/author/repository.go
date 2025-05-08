package author

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(params GetAllAuthorParams) ([]Author, int, error)
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

func (r *repository) FindAll(params GetAllAuthorParams) ([]Author, int, error) {
	var authors []Author
	var totalCount int

	//todo: Build the raw Count SQL query
	rawCountSQL := `
		SELECT 
			count(*)
		FROM authors
	`

	// Initialize the WHERE clause and arguments
	whereClauses := []string{"deleted_at IS NULL"}
	queryArgs := []interface{}{}

	//? field "name"
	if params.Name != "" {
		whereClauses = append(whereClauses, "(name LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Name+"%")
	}

	// Apply date range filtering for created_at if provided
	if len(params.CreatedAt) == 1 {
		// If only one date is provided, use equality
		whereClauses = append(whereClauses, "(created_at LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.CreatedAt[0]+"%")
	} else if len(params.CreatedAt) == 2 {
		// If two dates are provided, use BETWEEN
		whereClauses = append(whereClauses, "(created_at BETWEEN ? AND ?)")
		queryArgs = append(queryArgs, params.CreatedAt[0], params.CreatedAt[1])
	}

	//? Construct the WHERE clause
	whereSQL := ""
	if len(whereClauses) != 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	finalCountSQL := fmt.Sprintf(`
		%s
		%s`, rawCountSQL, whereSQL)

	// Add LIMIT and OFFSET arguments
	err := r.db.Raw(finalCountSQL, queryArgs...).Scan(&totalCount).Error

	if err != nil {
		return nil, 0, err
	}

	//todo: Build the raw SQL query
	rawSQL := `
		SELECT
			id,
			name,
			avatar_url,
			avatar_file_name,
			created_at
		FROM authors
	`

	//? Construct the ORDER BY clause
	orderBySQL := fmt.Sprintf("ORDER BY %s %s", params.Order, params.Sort)

	// Construct the final SQL query with LIMIT and OFFSET
	finalSQL := fmt.Sprintf(`
		%s
		%s
		%s
		LIMIT ? OFFSET ?`, rawSQL, whereSQL, orderBySQL)

	// Add LIMIT and OFFSET arguments
	offset := (params.Page - 1) * params.Limit
	queryArgs = append(queryArgs, params.Limit, offset)

	// Execute the raw SQL query
	err = r.db.Raw(finalSQL, queryArgs...).Scan(&authors).Error

	if err != nil {
		return nil, 0, err
	}

	return authors, totalCount, nil
}

func (r *repository) FindById(id int) (Author, error) {
	var data Author
	err := r.db.Where("id = ?", id).First(&data).Error
	if err == gorm.ErrRecordNotFound {
		errStr := fmt.Errorf("author with id %d not found", id)
		return Author{}, errStr
	}
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
