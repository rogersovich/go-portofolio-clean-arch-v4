package statistic

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(params GetAllStatisticParams) ([]Statistic, int, error)
	FindById(id int) (Statistic, error)
	CreateStatistic(p CreateStatisticRequest) (Statistic, error)
	CreateStatisticWithTx(p CreateStatisticRequest, tx *gorm.DB) (Statistic, error)
	UpdateStatistic(p UpdateStatisticRequest) error
	DeleteStatistic(id int) (Statistic, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(params GetAllStatisticParams) ([]Statistic, int, error) {
	var statistic []Statistic
	var totalCount int

	//todo: Build the raw Count SQL query
	rawCountSQL := `
		SELECT 
			count(*)
		FROM statistics
	`

	// Initialize the WHERE clause and arguments
	whereClauses := []string{"deleted_at IS NULL"}
	queryArgs := []interface{}{}

	//? field "type"
	if params.Type != "" {
		whereClauses = append(whereClauses, "(type LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Type+"%")
	}

	//? field "min_likes"
	if params.MinLikes != "" {
		whereClauses = append(whereClauses, "(likes >= ?)")
		queryArgs = append(queryArgs, params.MinLikes)
	}

	//? field "max_likes"
	if params.MaxLikes != "" {
		whereClauses = append(whereClauses, "(likes <= ?)")
		queryArgs = append(queryArgs, params.MaxLikes)
	}

	//? field "min_views"
	if params.MinViews != "" {
		whereClauses = append(whereClauses, "(views >= ?)")
		queryArgs = append(queryArgs, params.MinViews)
	}

	//? field "max_views"
	if params.MaxViews != "" {
		whereClauses = append(whereClauses, "(views <= ?)")
		queryArgs = append(queryArgs, params.MaxViews)
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
			likes,
			views,
			type,
			created_at
		FROM statistics
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
	err = r.db.Raw(finalSQL, queryArgs...).Scan(&statistic).Error

	if err != nil {
		return nil, 0, err
	}

	return statistic, totalCount, nil
}

func (r *repository) FindById(id int) (Statistic, error) {
	var data Statistic
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) CreateStatistic(p CreateStatisticRequest) (Statistic, error) {
	data := Statistic{
		Likes: p.Likes,
		Views: p.Views,
		Type:  p.Type}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) CreateStatisticWithTx(p CreateStatisticRequest, tx *gorm.DB) (Statistic, error) {
	data := Statistic{
		Likes: p.Likes,
		Views: p.Views,
		Type:  p.Type}
	err := tx.Create(&data).Error
	return data, err
}

func (r *repository) UpdateStatistic(p UpdateStatisticRequest) error {
	data := Statistic{
		ID:    p.ID,
		Likes: p.Likes,
		Views: p.Views,
		Type:  p.Type}
	err := r.db.Updates(&data).Error
	return err
}

func (r *repository) DeleteStatistic(id int) (Statistic, error) {
	var data Statistic

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return Statistic{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return Statistic{}, err
	}

	// Step 3: Return the data
	return data, nil
}
