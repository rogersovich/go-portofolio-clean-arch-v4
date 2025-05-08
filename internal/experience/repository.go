package experience

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(params GetAllExperienceParams) ([]Experience, int, error)
	FindById(id int) (Experience, error)
	CreateExperience(p CreateExperienceDTO) (Experience, error)
	UpdateExperience(p UpdateExperienceDTO) error
	DeleteExperience(id int) (Experience, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(params GetAllExperienceParams) ([]Experience, int, error) {
	var experience []Experience
	var totalCount int

	//todo: Build the raw Count SQL query
	rawCountSQL := `
		SELECT 
			count(*)
		FROM experiences
	`

	// Initialize the WHERE clause and arguments
	whereClauses := []string{"deleted_at IS NULL"}
	queryArgs := []interface{}{}

	//? field "position"
	if params.Position != "" {
		whereClauses = append(whereClauses, "(position LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Position+"%")
	}

	//? field "company_name"
	if params.CompanyName != "" {
		whereClauses = append(whereClauses, "(company_name LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.CompanyName+"%")
	}

	//? field "work_type"
	if params.WorkType != "" {
		whereClauses = append(whereClauses, "(work_type LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.WorkType+"%")
	}

	//? field "company_name"
	if params.CompanyName != "" {
		whereClauses = append(whereClauses, "(company_name LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.CompanyName+"%")
	}

	//? field "country"
	if params.Country != "" {
		whereClauses = append(whereClauses, "(country LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Country+"%")
	}

	//? field "city"
	if params.City != "" {
		whereClauses = append(whereClauses, "(city LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.City+"%")
	}

	//? field "summary_html"
	if params.SummaryHTML != "" {
		whereClauses = append(whereClauses, "(summary_html LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.SummaryHTML+"%")
	}

	//? field "is_current"
	if params.IsCurrent != "" {
		is_current := params.IsCurrent == "Y"
		whereClauses = append(whereClauses, "(is_current = ?)")
		queryArgs = append(queryArgs, is_current)
	}

	//? field "from_date"
	if len(params.FromDate) == 1 {
		// If only one date is provided, use equality
		whereClauses = append(whereClauses, "(from_date LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.FromDate[0]+"%")
	} else if len(params.FromDate) == 2 {
		// If two dates are provided, use BETWEEN
		whereClauses = append(whereClauses, "(from_date BETWEEN ? AND ?)")
		queryArgs = append(queryArgs, params.FromDate[0], params.FromDate[1])
	}

	//? field "to_date"
	if len(params.ToDate) == 1 {
		// If only one date is provided, use equality
		whereClauses = append(whereClauses, "(to_date LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.ToDate[0]+"%")
	} else if len(params.ToDate) == 2 {
		// If two dates are provided, use BETWEEN
		whereClauses = append(whereClauses, "(to_date BETWEEN ? AND ?)")
		queryArgs = append(queryArgs, params.ToDate[0], params.ToDate[1])
	}

	//? field "created_at"
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

	fmt.Println(finalCountSQL)
	fmt.Println(queryArgs)

	// Add LIMIT and OFFSET arguments
	err := r.db.Raw(finalCountSQL, queryArgs...).Scan(&totalCount).Error

	if err != nil {
		return nil, 0, err
	}

	//todo: Build the raw SQL query
	rawSQL := `
		SELECT
			id,
			position,
			company_name,
			work_type,
			country,
			city,
			summary_html,
			from_date,
			to_date,
			comp_image_url,
			comp_image_file_name,
			comp_website_url,
			is_current,
			created_at
		FROM experiences
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
	err = r.db.Raw(finalSQL, queryArgs...).Scan(&experience).Error

	if err != nil {
		return nil, 0, err
	}

	return experience, totalCount, nil
}

func (r *repository) FindById(id int) (Experience, error) {
	var data Experience
	err := r.db.Where("id = ?", id).First(&data).Error
	if err == gorm.ErrRecordNotFound {
		return Experience{}, gorm.ErrRecordNotFound
	}
	return data, err
}

func (r *repository) CreateExperience(p CreateExperienceDTO) (Experience, error) {
	data := Experience{
		Position:          p.Position,
		CompanyName:       p.CompanyName,
		WorkType:          p.WorkType,
		Country:           p.Country,
		City:              p.City,
		SummaryHTML:       p.SummaryHTML,
		FromDate:          p.FromDate,
		ToDate:            p.ToDate,
		CompImageUrl:      p.CompImageUrl,
		CompImageFileName: p.CompImageFileName,
		CompWebsiteUrl:    p.CompWebsiteUrl,
		IsCurrent:         p.IsCurrent}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateExperience(p UpdateExperienceDTO) error {
	// Create a map with only the fields that are non-zero
	updateMap := map[string]interface{}{
		"position":             p.Position,
		"company_name":         p.CompanyName,
		"work_type":            p.WorkType,
		"country":              p.Country,
		"city":                 p.City,
		"summary_html":         p.SummaryHTML,
		"from_date":            p.FromDate,
		"to_date":              p.ToDate,
		"comp_image_url":       p.CompImageUrl,
		"comp_image_file_name": p.CompImageFileName,
		"comp_website_url":     p.CompWebsiteUrl,
		"is_current":           p.IsCurrent,
	}

	err := r.db.Table("experiences").Where("id = ?", p.ID).Updates(updateMap).Error
	if err != nil {
		return err
	}

	return err
}

func (r *repository) DeleteExperience(id int) (Experience, error) {
	var data Experience

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return Experience{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return Experience{}, err
	}

	// Step 3: Return the data
	return data, nil
}
