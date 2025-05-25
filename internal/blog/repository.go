package blog

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(params GetAllBlogParams) ([]Blog, int, error)
	FindByIdWithRelations(id int) ([]RawBlogRelationResponse, error)
	FindById(id int) (BlogResponse, error)
	CreateBlog(p CreateBlogDTO, tx *gorm.DB) (Blog, error)
	UpdateBlog(p UpdateBlogDTO, tx *gorm.DB) (Blog, error)
	DeleteBlog(id int) (Blog, error)
	ChangeStatusBlog(id int, status string, blog BlogResponse) (BlogChangeStatusResponse, error)
	CheckUniqueSlug(slug string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(params GetAllBlogParams) ([]Blog, int, error) {
	var blog []Blog
	var totalCount int

	//todo: Build the raw Count SQL query
	rawCountSQL := `
		SELECT 
			count(*)
		FROM blogs
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
		whereClauses = append(whereClauses, "(status = ?)")
		queryArgs = append(queryArgs, params.Status)
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
		return blog, totalCount, nil
	}

	//todo: Build the raw SQL query
	rawSQL := `
		SELECT
			id,
			author_id,
			statistic_id,
			reading_time_id,
			title,
			description_html,
			summary,
			banner_url,
			banner_file_name,
			status,
			slug,
			is_highlight,
			published_at,
			created_at
		FROM blogs
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
	err = r.db.Raw(finalSQL, queryArgs...).Scan(&blog).Error

	if err != nil {
		return nil, 0, err
	}

	return blog, totalCount, nil
}

func (r *repository) FindById(id int) (BlogResponse, error) {
	var data BlogResponse
	err := r.db.Table("blogs").Where("id = ?", id).Scan(&data).Error
	if data.ID == 0 {
		return BlogResponse{}, gorm.ErrRecordNotFound
	}
	return data, err
}

func (r *repository) FindByIdWithRelations(id int) ([]RawBlogRelationResponse, error) {
	var data []RawBlogRelationResponse
	err := r.db.Raw(`
		SELECT 
			b.id, 
			b.title,
			b.description_html,
			b.summary,
			b.banner_url,
			b.banner_file_name,
			b.published_at,
			b.status,
			b.slug,
			b.is_highlight,
			b.created_at,
			a.id as author_id,
			a.name as author_name,
			rt.id as reading_time_id,
			rt.minutes as reading_time_minutes,
			rt.text_length as reading_time_text_length,
			rt.estimated_seconds as reading_time_estimated_seconds,
			rt.word_count as reading_time_word_count,
			rt.type as reading_time_type,
			s.id as statistic_id,
			s.likes as statistic_likes,
			s.views as statistic_views,
			s.type as statistic_type,
			t.id as topic_id,
			t.name as topic_name,
			bct.id as blog_content_image_id,
			bct.image_file_name as blog_content_image_file_name,
			bct.image_url as blog_content_image_url
		FROM blogs b
		LEFT JOIN authors a ON a.id = b.author_id
		LEFT JOIN reading_times rt ON rt.id = b.reading_time_id
		LEFT JOIN statistics s ON s.id = b.statistic_id
		LEFT JOIN blog_topics bt ON bt.blog_id = b.id
		LEFT JOIN topics t ON t.id = bt.topic_id
		LEFT JOIN blog_content_images bct ON bct.blog_id = b.id
		WHERE 
			b.id = ? AND 
			b.deleted_at IS NULL
	`, id).Scan(&data).Error

	if err != nil {
		return nil, err // handle DB or syntax error
	}

	if len(data) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return data, err
}

func (r *repository) CreateBlog(p CreateBlogDTO, tx *gorm.DB) (Blog, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	data := Blog{
		AuthorID:        p.AuthorID,
		StatisticID:     p.StatisticID,
		ReadingTimeID:   p.ReadingTimeID,
		Title:           p.Title,
		DescriptionHTML: p.DescriptionHTML,
		BannerUrl:       p.BannerUrl,
		BannerFileName:  p.BannerFileName,
		Summary:         p.Summary,
		Status:          p.Status,
		Slug:            p.Slug,
		PublishedAt:     p.PublishedAt,
		IsHighlight:     false,
	}

	err := db.Table("blogs").Create(&data).Error
	return data, err
}

func (r *repository) UpdateBlog(p UpdateBlogDTO, tx *gorm.DB) (Blog, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	updateMap := map[string]interface{}{
		"statistic_id":     p.StatisticID,
		"reading_time_id":  p.ReadingTimeID,
		"author_id":        p.AuthorID,
		"title":            p.Title,
		"description_html": p.DescriptionHTML,
		"banner_url":       p.BannerUrl,
		"banner_file_name": p.BannerFileName,
		"summary":          p.Summary,
		"status":           p.Status,
		"slug":             p.Slug,
		"is_highlight":     p.IsHighlight == "Y",
		"published_at":     p.PublishedAt,
		"updated_at":       time.Now(),
	}

	err := db.Table("blogs").Where("id = ?", p.ID).Updates(updateMap).Error

	data := Blog{
		ID:              p.ID,
		StatisticID:     p.StatisticID,
		ReadingTimeID:   p.ReadingTimeID,
		AuthorID:        p.AuthorID,
		Title:           p.Title,
		DescriptionHTML: p.DescriptionHTML,
		BannerUrl:       p.BannerUrl,
		BannerFileName:  p.BannerFileName,
		Summary:         p.Summary,
		Status:          p.Status,
		Slug:            p.Slug,
		IsHighlight:     p.IsHighlight == "Y",
		PublishedAt:     p.PublishedAt,
		UpdatedAt:       time.Now(),
	}
	return data, err
}

func (r *repository) DeleteBlog(id int) (Blog, error) {
	var data Blog

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return Blog{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return Blog{}, err
	}

	// Step 3: Return the data
	return data, nil
}

func (r *repository) ChangeStatusBlog(id int, status string, blog BlogResponse) (BlogChangeStatusResponse, error) {
	now := time.Now()
	var updateMap = make(map[string]interface{})
	updateMap["status"] = status
	if status == "Published" {
		updateMap["published_at"] = now
	}
	err := r.db.Model(&Blog{}).Where("id = ?", id).Updates(updateMap).Error

	if err != nil {
		return BlogChangeStatusResponse{}, err
	}

	// Return the updated data
	var publishedAtStringPtr *string
	if status == "Published" {
		publishedAtString := now.Format("2006-01-02 15:04:05")
		publishedAtStringPtr = &publishedAtString
	}
	updatedData := BlogChangeStatusResponse{
		ID:          id,
		Title:       blog.Title,
		Status:      status,
		PublishedAt: publishedAtStringPtr,
	}

	return updatedData, nil
}

func (r *repository) CheckUniqueSlug(slug string) (bool, error) {
	var data BlogResponse

	err := r.db.Table("blogs").Where("slug = ?", slug).Scan(&data).Error

	if err != nil {
		return false, err // handle DB or syntax error
	}

	if data.ID == 0 {
		return true, nil
	}

	return false, nil
}
