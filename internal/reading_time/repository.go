package reading_time

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(params GetAllReadingTimeParams) ([]ReadingTime, int, error)
	FindById(id int) (ReadingTime, error)
	CreateReadingTime(p CreateReadingTimeRequest, tx *gorm.DB) (ReadingTime, error)
	UpdateReadingTime(p UpdateReadingTimeRequest, tx *gorm.DB) error
	DeleteReadingTime(id int) (ReadingTime, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(params GetAllReadingTimeParams) ([]ReadingTime, int, error) {
	var reading_time []ReadingTime
	var totalCount int

	//todo: Build the raw Count SQL query
	rawCountSQL := `
		SELECT 
			count(*)
		FROM reading_times
	`

	// Initialize the WHERE clause and arguments
	whereClauses := []string{"deleted_at IS NULL"}
	queryArgs := []interface{}{}

	//? field "min_minutes"
	if params.MinMinutes != "" {
		whereClauses = append(whereClauses, "(minutes >= ?)")
		queryArgs = append(queryArgs, params.MinMinutes)
	}

	//? field "max_minutes"
	if params.MaxMinutes != "" {
		whereClauses = append(whereClauses, "(minutes <= ?)")
		queryArgs = append(queryArgs, params.MaxMinutes)
	}

	//? field "min_estimates"
	if params.MinEstimates != "" {
		whereClauses = append(whereClauses, "(estimated_seconds >= ?)")
		queryArgs = append(queryArgs, params.MinEstimates)
	}

	//? field "max_estimates"
	if params.MaxEstimates != "" {
		whereClauses = append(whereClauses, "(estimated_seconds <= ?)")
		queryArgs = append(queryArgs, params.MaxEstimates)
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
			minutes,
			text_length,
			estimated_seconds,
			word_count,
			type,
			created_at
		FROM reading_times
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
	err = r.db.Raw(finalSQL, queryArgs...).Scan(&reading_time).Error

	if err != nil {
		return nil, 0, err
	}

	return reading_time, totalCount, nil
}

func (r *repository) FindById(id int) (ReadingTime, error) {
	var data ReadingTime
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) CreateReadingTime(p CreateReadingTimeRequest, tx *gorm.DB) (ReadingTime, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	data := ReadingTime{
		Minutes:          p.Minutes,
		TextLength:       p.TextLength,
		EstimatedSeconds: p.EstimatedSeconds,
		WordCount:        p.WordCount,
		Type:             p.Type}
	err := db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateReadingTime(p UpdateReadingTimeRequest, tx *gorm.DB) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	data := ReadingTime{
		ID:               p.ID,
		Minutes:          p.Minutes,
		TextLength:       p.TextLength,
		EstimatedSeconds: p.EstimatedSeconds,
		WordCount:        p.WordCount,
		Type:             p.Type}
	err := db.Updates(&data).Error
	return err
}

func (r *repository) DeleteReadingTime(id int) (ReadingTime, error) {
	var data ReadingTime

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return ReadingTime{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return ReadingTime{}, err
	}

	// Step 3: Return the data
	return data, nil
}
