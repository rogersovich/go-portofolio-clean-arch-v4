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

	// Build the raw SQL query
	rawSQL := `
		SELECT 
			id,	
			name,
			avatar_url,
			avatar_file_name,
			created_at
		FROM authors
	`
	// Initialize the WHERE clause and arguments
	whereClauses := []string{"deleted_at IS NULL"}
	queryArgs := []interface{}{}

	//? field "search"
	if params.Name != "" {
		whereClauses = append(whereClauses, "(name LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Name+"%")
	}

	//? Construct the WHERE clause
	whereSQL := ""
	if len(whereClauses) != 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

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
	err := r.db.Raw(finalSQL, queryArgs...).Scan(&authors).Error

	if err != nil {
		return nil, 0, err
	}

	rawCountSQL := `
		SELECT 
			count(*)
		FROM authors
	`

	// Initialize the WHERE clause and arguments
	whereCountClauses := []string{"deleted_at IS NULL"}
	queryCountArgs := []interface{}{}

	//? field "search"
	if params.Name != "" {
		whereCountClauses = append(whereCountClauses, "(name LIKE ?)")
		queryCountArgs = append(queryCountArgs, "%"+params.Name+"%")
	}

	//? Construct the WHERE clause
	whereCountSQL := ""
	if len(whereCountClauses) != 0 {
		whereCountSQL = "WHERE " + strings.Join(whereCountClauses, " AND ")
	}

	finalCountSQL := fmt.Sprintf(`
		%s
		%s`, rawCountSQL, whereCountSQL)

	// Add LIMIT and OFFSET arguments
	err = r.db.Raw(finalCountSQL, queryCountArgs...).Scan(&totalCount).Error

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
