package technology

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(params GetAllTechnologyParams) ([]Technology, int, error)
	FindById(id int) (Technology, error)
	CreateTechnology(p CreateTechnologyDTO) (Technology, error)
	UpdateTechnology(p UpdateTechnologyDTO) error
	DeleteTechnology(id int) (Technology, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(params GetAllTechnologyParams) ([]Technology, int, error) {
	var technology []Technology
	var totalCount int

	//todo: Build the raw Count SQL query
	rawCountSQL := `
		SELECT 
			count(*)
		FROM technologies
	`

	// Initialize the WHERE clause and arguments
	whereClauses := []string{"deleted_at IS NULL"}
	queryArgs := []interface{}{}

	//? field "name"
	if params.Name != "" {
		whereClauses = append(whereClauses, "(name LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Name+"%")
	}

	//? field "description_html"
	if params.DescriptionHTML != "" {
		whereClauses = append(whereClauses, "(description_html LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.DescriptionHTML+"%")
	}

	//? field "is_major"
	if params.IsMajor != "" {
		is_used := params.IsMajor == "Y"
		whereClauses = append(whereClauses, "(is_major = ?)")
		queryArgs = append(queryArgs, is_used)
	}

	//? field "created_at"
	if len(params.CreatedAt) == 1 {
		whereClauses = append(whereClauses, "(created_at LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.CreatedAt[0]+"%")
	} else if len(params.CreatedAt) == 2 {
		// Parse the dates and adjust the time for the range
		startDate, err := time.Parse("2006-01-02", params.CreatedAt[0])
		if err != nil {
			return nil, 0, err
		}
		endDate, err := time.Parse("2006-01-02", params.CreatedAt[1])
		if err != nil {
			return nil, 0, err
		}

		startDate = startDate.Truncate(24 * time.Hour)        // Start at 00:00:00
		endDate = endDate.Add(24*time.Hour - time.Nanosecond) // End at 23:59:59.999
		whereClauses = append(whereClauses, "(created_at BETWEEN ? AND ?)")
		queryArgs = append(queryArgs, startDate, endDate)
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
			description_html,
			logo_url,
			logo_file_name,
			is_major,
			link,
			created_at
		FROM technologies
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
	err = r.db.Raw(finalSQL, queryArgs...).Scan(&technology).Error

	if err != nil {
		return nil, 0, err
	}

	return technology, totalCount, nil
}

func (r *repository) FindById(id int) (Technology, error) {
	var data Technology
	err := r.db.Where("id = ?", id).First(&data).Error
	if err == gorm.ErrRecordNotFound {
		return Technology{}, gorm.ErrRecordNotFound
	}
	return data, err
}

func (r *repository) CreateTechnology(p CreateTechnologyDTO) (Technology, error) {
	data := Technology{
		Name:            p.Name,
		DescriptionHTML: p.DescriptionHTML,
		LogoUrl:         p.LogoUrl,
		LogoFileName:    p.LogoFileName,
		IsMajor:         p.IsMajor,
		Link:            p.Link,
	}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateTechnology(p UpdateTechnologyDTO) error {

	updateMap := map[string]interface{}{
		"id":               p.ID,
		"name":             p.Name,
		"description_html": p.DescriptionHTML,
		"logo_url":         p.LogoUrl,
		"logo_file_name":   p.LogoFileName,
		"is_major":         p.IsMajor,
		"link":             p.Link,
		"updated_at":       time.Now(),
	}
	err := r.db.Table("technologies").Where("id = ?", p.ID).Updates(updateMap).Error
	if err != nil {
		return err
	}

	return err
}

func (r *repository) DeleteTechnology(id int) (Technology, error) {
	var data Technology

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return Technology{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return Technology{}, err
	}

	// Step 3: Return the data
	return data, nil
}
