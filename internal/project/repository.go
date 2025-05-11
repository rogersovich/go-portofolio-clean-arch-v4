package project

import (
	"fmt"
	"strings"
	"time"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/statistic"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(params GetAllProjectParams) ([]Project, int, error)
	FindByIdWithRelations(id int) ([]RawProjectRelationResponse, error)
	FindById(id int) (ProjectResponse, error)
	CreateProject(p CreateProjectDTO, tx *gorm.DB) (Project, error)
	UpdateProject(p UpdateProjectDTO, tx *gorm.DB) (Project, error)
	UpdateProjectStatistic(p ProjectStatisticUpdateDTO) (ProjectStatisticUpdateResponse, error)
	DeleteProject(id int) (Project, error)
	CheckUniqueSlug(slug string) (bool, error)
	ChangeStatusProject(id int, status string, project ProjectResponse) (ProjectChangeStatusResponse, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(params GetAllProjectParams) ([]Project, int, error) {
	var project []Project
	var totalCount int

	//todo: Build the raw Count SQL query
	rawCountSQL := `
		SELECT 
			count(*)
		FROM projects
	`

	// Initialize the WHERE clause and arguments
	whereClauses := []string{"deleted_at IS NULL"}
	queryArgs := []interface{}{}

	//? field "title"
	if params.Title != "" {
		whereClauses = append(whereClauses, "(title LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Title+"%")
	}

	//? field "status"
	if params.Status != "" {
		whereClauses = append(whereClauses, "(status LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Status+"%")
	}

	//? field "published_at"
	if len(params.PublishedAt) == 1 {
		// If only one date is provided, use equality
		whereClauses = append(whereClauses, "(published_at LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.PublishedAt[0]+"%")
	} else if len(params.PublishedAt) == 2 {
		// Parse the dates and adjust the time for the range
		startDate, err := time.Parse("2006-01-02", params.PublishedAt[0])
		if err != nil {
			return nil, 0, err
		}
		endDate, err := time.Parse("2006-01-02", params.PublishedAt[1])
		if err != nil {
			return nil, 0, err
		}

		startDate = startDate.Truncate(24 * time.Hour)        // Start at 00:00:00
		endDate = endDate.Add(24*time.Hour - time.Nanosecond) // End at 23:59:59.999
		whereClauses = append(whereClauses, "(from_date BETWEEN ? AND ?)")
		queryArgs = append(queryArgs, startDate, endDate)
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

	if totalCount == 0 {
		return project, totalCount, nil
	}

	//todo: Build the raw SQL query
	rawSQL := `
		SELECT
			id,
			title,
			description,
			image_url,
			image_file_name,
			repository_url,
			summary,
			status,
			slug,
			published_at,
			created_at
		FROM projects
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
	err = r.db.Raw(finalSQL, queryArgs...).Scan(&project).Error

	if err != nil {
		return nil, 0, err
	}

	return project, totalCount, nil
}

func (r *repository) FindByIdWithRelations(id int) ([]RawProjectRelationResponse, error) {
	var data []RawProjectRelationResponse
	err := r.db.Raw(`
		SELECT 
			p.id, 
			p.title, 
			p.description, 
			p.image_url, 
			p.image_file_name, 
			p.repository_url, 
			p.summary, 
			p.status, 
			p.published_at,
			p.created_at,
			s.id as statistic_id,
			s.views as statistic_views,
			s.likes as statistic_likes,
			s.type as statistic_type,
			pt.id as project_technology_id,
			t.id as technology_id,
			t.name as technology_name,
			pci.id as project_img_id,
      pci.image_file_name as project_img_file_name,
			pci.image_url as project_img_url
		FROM projects p 
		LEFT JOIN statistics s ON p.statistic_id = s.id
		LEFT JOIN project_technologies pt ON pt.project_id = p.id
		LEFT JOIN technologies t ON t.id = pt.technology_id
		LEFT JOIN project_content_images pci ON pci.project_id = p.id
		WHERE p.id = ? AND p.deleted_at IS NULL
	`, id).Scan(&data).Error

	if err != nil {
		return nil, err // handle DB or syntax error
	}

	if len(data) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return data, err
}

func (r *repository) FindById(id int) (ProjectResponse, error) {
	var data ProjectResponse
	err := r.db.Table("projects").Where("id = ?", id).Scan(&data).Error
	if data.ID == 0 {
		return ProjectResponse{}, gorm.ErrRecordNotFound
	}
	return data, err
}

func (r *repository) CreateProject(p CreateProjectDTO, tx *gorm.DB) (Project, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	// Create Project
	data := Project{
		StatisticID:   p.StatisticID,
		Title:         p.Title,
		Description:   p.Description,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		RepositoryUrl: p.RepositoryUrl,
		Summary:       p.Summary,
		Status:        p.Status,
		Slug:          p.Slug,
		PublishedAt:   p.PublishedAt}

	err := db.Create(&data).Error

	if err != nil {
		return Project{}, err
	}

	return data, err
}

func (r *repository) UpdateProject(p UpdateProjectDTO, tx *gorm.DB) (Project, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	//todo: UPDATE PROJECT

	data := Project{
		ID:            p.Id,
		Title:         p.Title,
		Description:   p.Description,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		RepositoryUrl: p.RepositoryUrl,
		Summary:       p.Summary,
		Status:        p.Status,
		Slug:          p.Slug,
		PublishedAt:   p.PublishedAt}

	err := db.Where("ID = ?", p.Id).Updates(&data).Error

	if err != nil {
		return Project{}, err
	}

	return data, nil
}

func (r *repository) UpdateProjectStatistic(p ProjectStatisticUpdateDTO) (ProjectStatisticUpdateResponse, error) {
	data := statistic.Statistic{
		ID:    p.StatisticID,
		Likes: p.Likes,
		Views: p.Views,
		Type:  p.Type,
	}
	err := r.db.Where("ID = ?", p.StatisticID).Updates(&data).Error
	if err != nil {
		return ProjectStatisticUpdateResponse{}, err
	}

	res := ProjectStatisticUpdateResponse{
		ProjectID:    p.ProjectID,
		ProjectTitle: p.ProjectTitle,
		StatisticID:  p.StatisticID,
		Likes:        *data.Likes,
		Views:        *data.Views,
		Type:         data.Type,
	}

	return res, nil
}

func (r *repository) DeleteProject(id int) (Project, error) {
	var data Project

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return Project{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return Project{}, err
	}

	// Step 3: Return the data
	return data, nil
}

func (r *repository) CheckUniqueSlug(slug string) (bool, error) {
	var data ProjectResponse

	err := r.db.Table("projects").Where("slug = ?", slug).Scan(&data).Error

	if err != nil {
		return false, err // handle DB or syntax error
	}

	if data.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (r *repository) ChangeStatusProject(id int, status string, project ProjectResponse) (ProjectChangeStatusResponse, error) {
	now := time.Now()
	var updateMap = make(map[string]interface{})
	updateMap["status"] = status
	if status == "Published" {
		updateMap["published_at"] = now
	}
	err := r.db.Model(&Project{}).Where("id = ?", id).Updates(updateMap).Error

	if err != nil {
		return ProjectChangeStatusResponse{}, err
	}

	// Return the updated data
	var publishedAtStringPtr *string
	if status == "Published" {
		publishedAtString := now.Format("2006-01-02 15:04:05")
		publishedAtStringPtr = &publishedAtString
	}
	updatedData := ProjectChangeStatusResponse{
		ID:          id,
		Title:       project.Title,
		Status:      status,
		PublishedAt: publishedAtStringPtr,
	}

	return updatedData, nil
}
