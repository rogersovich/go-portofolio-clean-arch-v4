package blog

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Blog, error)
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
