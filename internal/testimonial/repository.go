package testimonial

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(params GetAllTestimonialParams) ([]Testimonial, int, error)
	FindById(id int) (Testimonial, error)
	FindByMultiId(ids []int) ([]Testimonial, error)
	CreateTestimonial(p CreateTestimonialDTO) (Testimonial, error)
	UpdateTestimonial(p UpdateTestimonialDTO) error
	DeleteTestimonial(id int) (Testimonial, error)
	ChangeStatusTestimonial(id int, isUsed bool) error
	ChangeMultiStatusTestimonial(ids []int, isUsed bool) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(params GetAllTestimonialParams) ([]Testimonial, int, error) {
	var testimonial []Testimonial
	var totalCount int

	//todo: Build the raw Count SQL query
	rawCountSQL := `
		SELECT 
			count(*)
		FROM testimonials
	`

	// Initialize the WHERE clause and arguments
	whereClauses := []string{"deleted_at IS NULL"}
	queryArgs := []interface{}{}

	//? field "name"
	if params.Name != "" {
		whereClauses = append(whereClauses, "(name LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Name+"%")
	}

	//? field "role"
	if params.Role != "" {
		whereClauses = append(whereClauses, "(role LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Role+"%")
	}

	//? field "working_at"
	if params.WorkingAt != "" {
		whereClauses = append(whereClauses, "(working_at LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.WorkingAt+"%")
	}

	//? field "is_used"
	if params.IsUsed != "" {
		is_used := params.IsUsed == "Y"
		whereClauses = append(whereClauses, "(is_used = ?)")
		queryArgs = append(queryArgs, is_used)
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
			via,
			role,
			working_at,
			is_used,
			created_at
		FROM testimonials
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
	err = r.db.Raw(finalSQL, queryArgs...).Scan(&testimonial).Error

	if err != nil {
		return nil, 0, err
	}

	return testimonial, totalCount, nil
}

func (r *repository) FindById(id int) (Testimonial, error) {
	var data Testimonial
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) FindByMultiId(ids []int) ([]Testimonial, error) {
	var datas []Testimonial
	err := r.db.Table("testimonials").Where("id IN ?", ids).Scan(&datas).Error
	return datas, err
}

func (r *repository) CreateTestimonial(p CreateTestimonialDTO) (Testimonial, error) {
	data := Testimonial{
		Name:      p.Name,
		Via:       p.Via,
		Role:      p.Role,
		WorkingAt: p.WorkingAt,
		IsUsed:    p.IsUsed,
	}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateTestimonial(p UpdateTestimonialDTO) error {
	updates := map[string]interface{}{
		"id":         p.ID,
		"name":       p.Name,
		"via":        p.Via,
		"role":       p.Role,
		"working_at": p.WorkingAt,
		"is_used":    p.IsUsed,
	}
	err := r.db.Table("testimonials").Where("id = ?", p.ID).Updates(updates).Error
	if err != nil {
		return err
	}
	return err
}

func (r *repository) DeleteTestimonial(id int) (Testimonial, error) {
	var data Testimonial

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return Testimonial{}, err // return if not found or any error
	}

	// Step 2: Hard Delete
	if err := r.db.Unscoped().Delete(&data).Error; err != nil {
		return Testimonial{}, err
	}

	// Step 3: Return the data
	return data, nil
}

func (r *repository) ChangeStatusTestimonial(id int, isUsed bool) error {
	updates := map[string]interface{}{
		"is_used": isUsed,
	}
	err := r.db.Table("testimonials").Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return err
	}
	return err
}

func (r *repository) ChangeMultiStatusTestimonial(ids []int, isUsed bool) error {
	updates := map[string]interface{}{
		"is_used": isUsed,
	}
	err := r.db.Table("testimonials").Where("id IN ?", ids).Updates(updates).Error
	if err != nil {
		return err
	}
	return err
}
