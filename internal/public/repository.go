package public

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/statistic"
	"gorm.io/gorm"
)

type Repository interface {
	GetTechnologiesPublic() ([]TechnologyProfilePublicResponse, error)
	GetAboutPublic() (AboutPublicResponse, error)
	GetCurrentWork() (CurrentWorkPublicResponse, error)
	GetExperiencesPublic() ([]ExperiencesPublicResponse, error)
	GetRawPublicPaginateBlogs(params BlogPublicParams) ([]BlogPaginatePublicRaw, int, error)
	GetRawPublicBlogs(params BlogPublicParams, uniquePaginateBlogIDs []int) ([]BlogPublicRaw, error)
	GetPublicBlogBySlug(slug string) ([]SingleBlogPublicRaw, error)
	GetRawPublicBlogTopics(params BlogPublicParams, uniqueBlogIDs []int) ([]BlogTopicPublicRaw, error)
	GetPublicTestimonials() ([]TestimonialPublicResponse, error)
	GetPublicTopics() ([]TopicPublicResponse, error)
	GetRawPublicPaginateProjects(params ProjectPublicParams) ([]ProjectPaginatePublicRaw, int, error)
	GetRawPublicProjectTechnologies(params ProjectPublicParams, uniqueProjectIDs []int) ([]ProjectTechnologyPublicRaw, error)
	GetPublicProjectBySlug(slug string) ([]SingleProjectPublicRaw, error)
	GetPublicTechnologies() ([]TechnologyPublicResponse, error)
	GetPublicAuthors() ([]AuthorPublicResponse, error)
	FindProjectById(id int) (ProjectByIdResponse, error)
	UpdatePublicProjectStatistic(p ProjectStatisticUpdatePublicDTO) (ProjectStatisticUpdatePubblicResponse, error)
	FindBlogById(id int) (BlogByIdResponse, error)
	UpdatePublicBlogStatistic(p BlogStatisticUpdatePublicDTO) (BlogStatisticUpdatePubblicResponse, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetTechnologiesPublic() ([]TechnologyProfilePublicResponse, error) {
	var data []TechnologyProfilePublicResponse

	rawQuery := `
		SELECT id, name, logo_url, logo_file_name
		FROM technologies
		WHERE is_major = ? AND deleted_at IS NULL
	`
	err := r.db.Raw(rawQuery, 1).Scan(&data).Error

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

	rawQuery := `
		SELECT id, position, company_name, work_type, country, city, comp_website_url
		FROM experiences
		WHERE is_current = ? AND deleted_at IS NULL
		LIMIT 1
	`
	err := r.db.Raw(rawQuery, 1).Scan(&data).Error
	if err != nil {
		return CurrentWorkPublicResponse{}, err
	}

	return data, nil
}

func (r *repository) GetExperiencesPublic() ([]ExperiencesPublicResponse, error) {
	var data []ExperiencesPublicResponse

	// Build the query
	query := r.db.Table("experiences").Where("deleted_at IS NULL")

	sort := "DESC"
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

func (r *repository) GetRawPublicPaginateBlogs(params BlogPublicParams) ([]BlogPaginatePublicRaw, int, error) {
	var datas []BlogPaginatePublicRaw
	var totalCount int

	//todo: Build the raw Count SQL query
	rawCountSQL := `
		SELECT 
			count(*)
		FROM blogs b
	`

	// Build the raw SQL query
	rawSQL := `
		SELECT 
			b.id, 
			b.title,
			s.id as statistic_id,
			s.likes as statistic_likes,
			s.views as statistic_views,
			s.type as statistic_type
		FROM blogs b
		LEFT JOIN statistics s ON s.id = b.statistic_id
	`

	// Initialize the WHERE clause and arguments
	whereClauses := []string{"b.deleted_at IS NULL"}
	queryArgs := []interface{}{}

	//? field "status"
	whereClauses = append(whereClauses, "b.status = ?")
	queryArgs = append(queryArgs, "Published")

	//? field "search"
	if params.Search != "" {
		whereClauses = append(whereClauses, "(b.title LIKE ? OR b.summary LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Search+"%", "%"+params.Search+"%")
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
		return []BlogPaginatePublicRaw{}, 0, nil
	}

	//! Build Query Paginate

	//? Construct the ORDER BY clause
	orderBySQL := fmt.Sprintf("ORDER BY b.%s %s", params.Order, params.Sort)

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
	err = r.db.Raw(finalSQL, queryArgs...).Scan(&datas).Error

	if err != nil {
		return []BlogPaginatePublicRaw{}, 0, err
	}

	return datas, totalCount, nil
}

func (r *repository) GetRawPublicBlogs(params BlogPublicParams, uniquePaginateBlogIDs []int) ([]BlogPublicRaw, error) {
	var datas []BlogPublicRaw

	rawTopicSQL := `
		SELECT 
			b.id, 
			b.title,
			b.summary,
			b.description_html,
			b.banner_url,
			b.banner_file_name,
			b.published_at,
			b.status,
			b.slug,
			b.is_highlight,
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
			t.name as topic_name
		FROM blogs b
		LEFT JOIN authors a ON a.id = b.author_id
		LEFT JOIN reading_times rt ON rt.id = b.reading_time_id
		LEFT JOIN statistics s ON s.id = b.statistic_id
		LEFT JOIN blog_topics bt on bt.blog_id = b.id
		LEFT JOIN topics t on t.id = bt.topic_id
	`

	whereClauses := []string{}
	queryArgs := []interface{}{}

	if len(uniquePaginateBlogIDs) > 0 {
		whereClauses = append(whereClauses, "b.id IN (?)")
		queryArgs = append(queryArgs, uniquePaginateBlogIDs)
	}

	if len(params.Topics) > 0 {
		whereClauses = append(whereClauses, "t.id IN (?)")
		queryArgs = append(queryArgs, params.Topics)
	}

	//? Construct the WHERE clause
	whereSQL := ""
	if len(whereClauses) != 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	//? Construct the ORDER BY clause
	orderBySQL := fmt.Sprintf("ORDER BY b.%s %s", params.Order, params.Sort)

	// Construct the final SQL query with LIMIT and OFFSET
	finalSQL := fmt.Sprintf(`
		%s
		%s
		%s`, rawTopicSQL, whereSQL, orderBySQL)

	err := r.db.Raw(finalSQL, queryArgs...).Scan(&datas).Error
	if err != nil {
		return []BlogPublicRaw{}, err
	}

	return datas, nil
}

func (r *repository) GetRawPublicBlogTopics(params BlogPublicParams, uniqueBlogIDs []int) ([]BlogTopicPublicRaw, error) {
	var datas []BlogTopicPublicRaw

	rawTopicSQL := `
		SELECT 
			b.id as blog_id, 
			t.id as topic_id,
			t.name as topic_name
		FROM blogs b
		LEFT JOIN statistics s ON s.id = b.statistic_id
		LEFT JOIN blog_topics bt on bt.blog_id = b.id
		LEFT JOIN topics t on t.id = bt.topic_id
	`

	whereClauses := []string{}
	queryArgs := []interface{}{}

	if len(uniqueBlogIDs) > 0 {
		whereClauses = append(whereClauses, "b.id IN (?)")
		queryArgs = append(queryArgs, uniqueBlogIDs)
	}

	//? Construct the WHERE clause
	whereSQL := ""
	if len(whereClauses) != 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	//? Construct the ORDER BY clause
	orderBySQL := fmt.Sprintf("ORDER BY b.%s %s", params.Order, params.Sort)

	// Construct the final SQL query with LIMIT and OFFSET
	finalSQL := fmt.Sprintf(`
		%s
		%s
		%s`, rawTopicSQL, whereSQL, orderBySQL)

	err := r.db.Raw(finalSQL, queryArgs...).Scan(&datas).Error
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
			b.is_highlight,
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
		return []SingleBlogPublicRaw{}, errors.New("data not found")
	}

	return datas, nil
}

func (r *repository) GetPublicTestimonials() ([]TestimonialPublicResponse, error) {
	var datas []TestimonialPublicResponse
	err := r.db.Table("testimonials").Where("deleted_at IS NULL AND is_used = ?", 1).Order("updated_at DESC").Scan(&datas).Error
	return datas, err
}

func (r *repository) GetPublicTopics() ([]TopicPublicResponse, error) {
	var datas []TopicPublicResponse
	err := r.db.Table("topics").Where("deleted_at IS NULL").Order("created_at DESC").Scan(&datas).Error
	return datas, err
}

func (r *repository) GetRawPublicPaginateProjects(params ProjectPublicParams) ([]ProjectPaginatePublicRaw, int, error) {
	var datas []ProjectPaginatePublicRaw
	var totalCount int

	//todo: Build the raw Count SQL query
	rawCountSQL := `
		SELECT 
			count(*)
		FROM projects p
	`

	// Build the raw SQL query
	rawSQL := `
		SELECT 
			p.id,
			p.title,
			p.summary,
			p.image_url,
			p.image_file_name,
			p.repository_url,
			p.published_at,
			p.slug,
			p.is_highlight
		FROM projects p
	`

	// Initialize the WHERE clause and arguments
	whereClauses := []string{"p.deleted_at IS NULL"}
	queryArgs := []interface{}{}

	//? field "status"
	whereClauses = append(whereClauses, "p.status = ?")
	queryArgs = append(queryArgs, "Published")

	//? field "search"
	if params.Search != "" {
		whereClauses = append(whereClauses, "(p.title LIKE ? OR p.summary LIKE ?)")
		queryArgs = append(queryArgs, "%"+params.Search+"%", "%"+params.Search+"%")
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
		return []ProjectPaginatePublicRaw{}, 0, nil
	}

	//! Build the raw SQL query

	//? Construct the ORDER BY clause
	orderBySQL := fmt.Sprintf("ORDER BY p.%s %s", params.Order, params.Sort)

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
	err = r.db.Raw(finalSQL, queryArgs...).Scan(&datas).Error

	if err != nil {
		return []ProjectPaginatePublicRaw{}, 0, err
	}

	return datas, totalCount, nil
}

func (r *repository) GetRawPublicProjectTechnologies(params ProjectPublicParams, uniqueProjectIDs []int) ([]ProjectTechnologyPublicRaw, error) {
	var datas []ProjectTechnologyPublicRaw

	rawSQL := `
		SELECT 
			p.id AS project_id,
			t.id AS tech_id,
			t.name AS tech_name,
			t.logo_url AS tech_logo_url,
			t.logo_file_name AS tech_logo_file_name,
			t.link AS tech_link
		FROM projects p
		LEFT JOIN project_technologies pt ON pt.project_id = p.id
		LEFT JOIN technologies t ON t.id = pt.technology_id
	`

	whereClauses := []string{}
	queryArgs := []interface{}{}

	if len(uniqueProjectIDs) > 0 {
		whereClauses = append(whereClauses, "p.id IN (?)")
		queryArgs = append(queryArgs, uniqueProjectIDs)
	}

	//? Construct the WHERE clause
	whereSQL := ""
	if len(whereClauses) != 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	//? Construct the ORDER BY clause
	orderBySQL := fmt.Sprintf("ORDER BY p.%s %s", params.Order, params.Sort)

	//? Construct the final SQL query with LIMIT and OFFSET
	finalSQL := fmt.Sprintf(`
		%s
		%s
		%s`, rawSQL, whereSQL, orderBySQL)

	err := r.db.Raw(finalSQL, queryArgs...).Scan(&datas).Error
	if err != nil {
		return []ProjectTechnologyPublicRaw{}, err
	}

	return datas, nil
}

func (r *repository) GetPublicProjectBySlug(slug string) ([]SingleProjectPublicRaw, error) {
	var datas []SingleProjectPublicRaw

	// Build the raw SQL query
	rawSQL := `
		SELECT 
			p.id, 
			p.title,
			p.description,
			p.image_url,
			p.image_file_name,
			p.repository_url,
			p.summary,
			p.status,
			p.slug,
			p.published_at,
			p.is_highlight,
			s.id as statistic_id,
			s.likes as statistic_likes,
			s.views as statistic_views,
			s.type as statistic_type,
			pci.id as content_image_id,
      pci.image_url as content_image_url,
      pci.image_file_name as content_image_file_name,
      t.id as tech_id,
      t.name as tech_name,
			t.logo_url as tech_logo_url,
			t.link as tech_link
		FROM projects p
		LEFT JOIN statistics s ON s.id = p.statistic_id
		LEFT JOIN project_content_images pci ON pci.project_id = p.id
    LEFT JOIN project_technologies pt ON pt.project_id = p.id
    LEFT JOIN technologies t ON t.id = pt.technology_id
		WHERE 
			p.deleted_at IS NULL AND 
			p.slug = ? AND
			p.status = ?
	`

	// Execute the raw SQL query
	err := r.db.Raw(rawSQL, slug, "Published").Scan(&datas).Error

	if err != nil {
		return []SingleProjectPublicRaw{}, err
	}

	if len(datas) == 0 {
		return []SingleProjectPublicRaw{}, errors.New("data not found")
	}

	return datas, nil
}

func (r *repository) GetPublicTechnologies() ([]TechnologyPublicResponse, error) {
	var datas []TechnologyPublicResponse
	err := r.db.Table("technologies").Where("deleted_at IS NULL").Order("updated_at DESC").Scan(&datas).Error
	return datas, err
}

func (r *repository) GetPublicAuthors() ([]AuthorPublicResponse, error) {
	var datas []AuthorPublicResponse
	err := r.db.Table("authors").Where("deleted_at IS NULL").Order("updated_at DESC").Scan(&datas).Error
	return datas, err
}

func (r *repository) FindProjectById(id int) (ProjectByIdResponse, error) {
	var data ProjectByIdResponse
	err := r.db.Table("projects").Where("id = ?", id).Scan(&data).Error
	if data.ID == 0 {
		return ProjectByIdResponse{}, gorm.ErrRecordNotFound
	}
	return data, err
}

func (r *repository) UpdatePublicProjectStatistic(p ProjectStatisticUpdatePublicDTO) (ProjectStatisticUpdatePubblicResponse, error) {
	data := statistic.Statistic{
		ID:    p.StatisticID,
		Likes: p.Likes,
		Views: p.Views,
		Type:  p.Type,
	}
	err := r.db.Where("ID = ?", p.StatisticID).Updates(&data).Error
	if err != nil {
		return ProjectStatisticUpdatePubblicResponse{}, err
	}

	res := ProjectStatisticUpdatePubblicResponse{
		ProjectID:    p.ProjectID,
		ProjectTitle: p.ProjectTitle,
		StatisticID:  p.StatisticID,
		Likes:        *data.Likes,
		Views:        *data.Views,
		Type:         data.Type,
	}

	return res, nil
}

func (r *repository) FindBlogById(id int) (BlogByIdResponse, error) {
	var data BlogByIdResponse
	err := r.db.Table("blogs").Where("id = ?", id).Scan(&data).Error
	if data.ID == 0 {
		return BlogByIdResponse{}, gorm.ErrRecordNotFound
	}
	return data, err
}

func (r *repository) UpdatePublicBlogStatistic(p BlogStatisticUpdatePublicDTO) (BlogStatisticUpdatePubblicResponse, error) {
	data := statistic.Statistic{
		ID:    p.StatisticID,
		Likes: p.Likes,
		Views: p.Views,
		Type:  p.Type,
	}
	err := r.db.Where("ID = ?", p.StatisticID).Updates(&data).Error
	if err != nil {
		return BlogStatisticUpdatePubblicResponse{}, err
	}

	res := BlogStatisticUpdatePubblicResponse{
		BlogID:      p.BlogID,
		Title:       p.Title,
		StatisticID: p.StatisticID,
		Likes:       *data.Likes,
		Views:       *data.Views,
		Type:        data.Type,
	}

	return res, nil
}
