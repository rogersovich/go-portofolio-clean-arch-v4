package blog

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Blog, error)
	FindByIdWithRelations(id int) ([]RawBlogRelationResponse, error)
	CreateBlog(p CreateBlogDTO, tx *gorm.DB) (Blog, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Blog, error) {
	var datas []Blog
	err := r.db.Find(&datas).Error
	return datas, err
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
		JOIN authors a ON a.id = b.author_id
		JOIN reading_times rt ON rt.id = b.reading_time_id
		JOIN statistics s ON s.id = b.statistic_id
		LEFT JOIN blog_topics bt ON bt.blog_id = b.id
		JOIN topics t ON t.id = bt.topic_id
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
		PublishedAt:     p.PublishedAt,
	}

	err := db.Table("blogs").Create(&data).Error
	return data, err
}
