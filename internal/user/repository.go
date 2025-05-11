package user

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(params GetAllUserParams) ([]User, int, error)
	FindById(id int) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(id int) error
	CheckUniqueEmail(email string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(params GetAllUserParams) ([]User, int, error) {
	var users []User
	var totalCount int

	//todo: Build the raw Count SQL query
	rawCountSQL := `
		SELECT 
			count(*)
		FROM users
	`

	// Initialize the WHERE clause and arguments
	whereClauses := []string{"deleted_at IS NULL"}
	queryArgs := []interface{}{}

	//? field "username"
	if params.Username != "" {
		whereClauses = append(whereClauses, "(username LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Username+"%")
	}

	//? field "email"
	if params.Email != "" {
		whereClauses = append(whereClauses, "(email LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Email+"%")
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
			username,
			email,
			created_at
		FROM users
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
	err = r.db.Raw(finalSQL, queryArgs...).Scan(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

func (r *repository) FindById(id int) (User, error) {
	var data User
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) UpdateUser(user User) (User, error) {
	err := r.db.Save(&user).Error
	return user, err
}

func (r *repository) DeleteUser(id int) error {
	var data User
	err := r.db.Where("id = ?", id).Delete(&data).Error
	return err
}

func (r *repository) CheckUniqueEmail(email string) (bool, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return true, nil
	}

	return false, nil
}
