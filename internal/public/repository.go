package public

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAllAuthors(params AuthorPublicParams) ([]AuthorPublicResponse, error)
	GetTechnologiesPublic() ([]TechnologyProfilePublicResponse, error)
	GetAboutPublic() (AboutPublicResponse, error)
	GetCurrentWork() (CurrentWorkPublicResponse, error)
	GetExperiencesPublic() ([]ExperiencesPublicResponse, error)
	GetPublicBlogs(params BlogPublicParams) ([]BlogPublicRaw, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAllAuthors(params AuthorPublicParams) ([]AuthorPublicResponse, error) {
	var datas []AuthorPublicResponse

	// Build the query
	query := r.db.Table("authors").Where("deleted_at IS NULL")

	// Filter by 'name' if provided
	if params.Name != "" {
		query = query.Where("name LIKE ?", "%"+params.Name+"%")
	}

	// Apply sorting if provided
	if params.Sort != "" && params.Order != "" {
		query = query.Order(params.Order + " " + params.Sort) // Dynamically apply sorting (e.g., id ASC)
	}

	// Apply pagination (LIMIT and OFFSET)
	if params.Limit > 0 {
		query = query.Offset((params.Page - 1) * params.Limit).Limit(params.Limit)
	}

	// Execute the query
	err := query.Find(&datas).Error
	return datas, err
}

func (r *repository) GetTechnologiesPublic() ([]TechnologyProfilePublicResponse, error) {
	var data []TechnologyProfilePublicResponse

	err := r.db.Table("technologies").Where("is_major = ? AND deleted_at IS NULL", 1).Scan(&data).Error

	if err != nil {
		return []TechnologyProfilePublicResponse{}, err
	}
	return data, nil
}

func (r *repository) GetAboutPublic() (AboutPublicResponse, error) {
	var data AboutPublicResponse

	err := r.db.Table("abouts").Where("is_used = ? AND deleted_at IS NULL", 1).First(&data).Error

	if err != nil {
		return AboutPublicResponse{}, err
	}
	return data, nil
}

func (r *repository) GetCurrentWork() (CurrentWorkPublicResponse, error) {
	var data CurrentWorkPublicResponse

	err := r.db.Table("experiences").Where("is_current = ? AND deleted_at IS NULL", 1).First(&data).Error
	if err != nil {
		return CurrentWorkPublicResponse{}, err
	}

	return data, nil
}

func (r *repository) GetExperiencesPublic() ([]ExperiencesPublicResponse, error) {
	var data []ExperiencesPublicResponse

	// Build the query
	query := r.db.Table("experiences").Where("deleted_at IS NULL")

	sort := "ASC"
	order := "from_date"

	// Apply sorting if provided
	if sort != "" && order != "" {
		query = query.Order(order + " " + sort)
	}

	// Execute the query
	err := query.Find(&data).Error

	if err != nil {
		return []ExperiencesPublicResponse{}, err
	}
	return data, nil
}

func (r *repository) GetPublicBlogs(params BlogPublicParams) ([]BlogPublicRaw, error) {
	var datas []BlogPublicRaw

	// Build the raw SQL query
	rawSQL := `
		SELECT 
			b.id, 
			b.title,
			b.description_html,
			b.summary,
			b.banner_url,
			b.banner_file_name,
			b.published_at,
			b.status,
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
			bct.id as content_image_id,
			bct.image_file_name as content_image_file_name,
			bct.image_url as content_image_url
		FROM blogs b
		LEFT JOIN authors a ON a.id = b.author_id
		LEFT JOIN reading_times rt ON rt.id = b.reading_time_id
		LEFT JOIN statistics s ON s.id = b.statistic_id
		LEFT JOIN blog_topics bt ON bt.blog_id = b.id
		LEFT JOIN topics t ON t.id = bt.topic_id
		LEFT JOIN blog_content_images bct ON bct.blog_id = b.id
		WHERE 
			b.deleted_at IS NULL AND 
			b.status = ?
	`

	// Apply sorting if provided
	if params.Sort != "" && params.Order != "" {
		rawSQL += " ORDER BY " + params.Order + " " + params.Sort
	}

	// Apply pagination (LIMIT and OFFSET)
	if params.Limit > 0 {
		rawSQL += " LIMIT ? OFFSET ?"
	}

	offset := (params.Page - 1) * params.Limit

	args := []interface{}{
		"Published",
		params.Limit,
		offset,
	}

	// Execute the raw SQL query
	err := r.db.Raw(rawSQL, args...).Scan(&datas).Error
	return datas, err
}
