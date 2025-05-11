package topic

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(params GetAllTopicParams) ([]Topic, int, error)
	FindById(id int) (Topic, error)
	CreateTopic(p CreateTopicRequest) (Topic, error)
	UpdateTopic(p UpdateTopicRequest) error
	DeleteTopic(id int) (Topic, error)
	CheckTopicIds(ids []int) ([]Topic, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(params GetAllTopicParams) ([]Topic, int, error) {
	var topics []Topic
	var totalCount int

	//todo: Build the raw Count SQL query
	rawCountSQL := `
		SELECT 
			count(*)
		FROM topics
	`

	// Initialize the WHERE clause and arguments
	whereClauses := []string{"deleted_at IS NULL"}
	queryArgs := []interface{}{}

	//? field "name"
	if params.Name != "" {
		whereClauses = append(whereClauses, "(name LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Name+"%")
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
			created_at
		FROM topics
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
	err = r.db.Raw(finalSQL, queryArgs...).Scan(&topics).Error

	if err != nil {
		return nil, 0, err
	}

	return topics, totalCount, nil
}

func (r *repository) FindById(id int) (Topic, error) {
	var data Topic
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) CreateTopic(p CreateTopicRequest) (Topic, error) {
	data := Topic{
		Name: p.Name}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateTopic(p UpdateTopicRequest) error {
	data := Topic{
		ID:   p.ID,
		Name: p.Name}
	err := r.db.Updates(&data).Error
	return err
}

func (r *repository) DeleteTopic(id int) (Topic, error) {
	var data Topic

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return Topic{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Unscoped().Delete(&data).Error; err != nil {
		return Topic{}, err
	}

	// Step 3: Return the data
	return data, nil
}

func (r *repository) CheckTopicIds(ids []int) ([]Topic, error) {
	var data []Topic
	err := r.db.Where("id in (?)", ids).Select("id", "name").Find(&data).Error
	return data, err
}
