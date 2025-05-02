package public

import (
	"errors"
	"fmt"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
	"gorm.io/gorm"
)

type Repository interface {
	FindAllAuthors(params AuthorPublicParams) ([]AuthorPublicResponse, error)
	GetTechnologiesPublic() ([]TechnologyProfilePublicResponse, error)
	GetAboutPublic() (AboutPublicResponse, error)
	GetCurrentWork() (CurrentWorkPublicResponse, error)
	GetExperiencesPublic() ([]ExperiencesPublicResponse, error)
	GetPublicBlogs(params BlogPublicParams) ([]BlogPublicRaw, error)
	GetPublicBlogTopics(params BlogPublicParams, uniqueBlogIDs []int) ([]BlogTopicPublicRaw, error)
	GetPublicBlogBySlug(slug string) ([]SingleBlogPublicRaw, error)
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
			b.summary,
			b.banner_url,
			b.banner_file_name,
			b.published_at,
			b.status,
			b.slug,
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
			s.type as statistic_type
		FROM blogs b
		LEFT JOIN authors a ON a.id = b.author_id
		LEFT JOIN reading_times rt ON rt.id = b.reading_time_id
		LEFT JOIN statistics s ON s.id = b.statistic_id
		WHERE 
			b.deleted_at IS NULL AND 
			b.status = ?
	`

	blogArgs := []interface{}{"Published"}

	if params.Search != "" {
		rawSQL += " AND (b.title LIKE ? OR b.summary LIKE ?)"
		blogArgs = append(blogArgs, "%"+params.Search+"%", "%"+params.Search+"%")
	}

	// Apply sorting if provided
	if params.Sort != "" && params.Order != "" {
		rawSQL += " ORDER BY " + params.Sort + " " + params.Order
	}

	// Apply pagination (LIMIT and OFFSET)
	if params.Limit > 0 {
		rawSQL += " LIMIT ? OFFSET ?"
		offset := (params.Page - 1) * params.Limit
		blogArgs = append(blogArgs, params.Limit, offset)
	}

	// Execute the raw SQL query
	err := r.db.Raw(rawSQL, blogArgs...).Scan(&datas).Error

	if err != nil {
		return []BlogPublicRaw{}, err
	}

	return datas, nil
}

func (r *repository) GetPublicBlogTopics(params BlogPublicParams, uniqueBlogIDs []int) ([]BlogTopicPublicRaw, error) {
	var datas []BlogTopicPublicRaw

	// Create a slice of "?" placeholders equal to the length of the input slice
	placeholders := utils.SliceIntToPlaceholder(uniqueBlogIDs)

	rawTopicSQL := fmt.Sprintf(`
		SELECT 
			b.id as blog_id,
			t.id as topic_id,
			t.name as topic_name
		FROM blogs b
		JOIN blog_topics bt on bt.blog_id = b.id
		JOIN topics t on t.id = bt.topic_id
		JOIN statistics s on s.id = b.statistic_id
		WHERE b.id IN (%s)
	`, placeholders)

	// Convert uniqueBlogIDs into a slice of interfaces for GORM query
	blogTopicArgs := make([]interface{}, len(uniqueBlogIDs))
	for i, id := range uniqueBlogIDs {
		blogTopicArgs[i] = id
	}

	// Apply sorting if provided
	if params.Sort != "" && params.Order != "" {
		rawTopicSQL += " ORDER BY " + params.Sort + " " + params.Order
	}

	err := r.db.Raw(rawTopicSQL, blogTopicArgs...).Scan(&datas).Error
	if err != nil {
		return []BlogTopicPublicRaw{}, err
	}

	return datas, nil
}

func (r *repository) GetPublicBlogBySlug(slug string) ([]SingleBlogPublicRaw, error) {
	var datas []SingleBlogPublicRaw

	// Build the raw SQL query
	rawSQL := `
		SELECT 
			b.id, 
			b.title,
			b.summary,
			b.banner_url,
			b.banner_file_name,
			b.description_html,
			b.published_at,
			b.status,
			b.slug,
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
			bct.id as content_image_id,
      bct.image_url as content_image_url,
      bct.image_file_name as content_image_file_name,
      t.id as topic_id,
      t.name as topic_name
		FROM blogs b
		LEFT JOIN authors a ON a.id = b.author_id
		LEFT JOIN reading_times rt ON rt.id = b.reading_time_id
		LEFT JOIN statistics s ON s.id = b.statistic_id
		LEFT JOIN blog_content_images bct ON bct.blog_id = b.id
    LEFT JOIN blog_topics bt ON bt.blog_id = b.id
    LEFT JOIN topics t ON t.id = bt.topic_id
		WHERE 
			b.deleted_at IS NULL AND 
			b.slug = ? AND
			b.status = ?
	`

	// Execute the raw SQL query
	err := r.db.Raw(rawSQL, slug, "Published").Scan(&datas).Error

	if err != nil {
		return []SingleBlogPublicRaw{}, err
	}

	if len(datas) == 0 {
		return []SingleBlogPublicRaw{}, errors.New("blog not found")
	}

	return datas, nil
}
