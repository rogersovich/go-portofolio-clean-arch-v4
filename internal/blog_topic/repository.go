package blog_topic

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]BlogTopic, error)
	FindById(id string) (BlogTopic, error)
	CreateBlogTopic(p CreateBlogTopicRequest) (BlogTopic, error)
	UpdateBlogTopic(p UpdateBlogTopicRequest) (BlogTopic, error)
	DeleteBlogTopic(id int) (BlogTopic, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]BlogTopic, error) {
	var datas []BlogTopic
	err := r.db.Find(&datas).Error
	return datas, err
}

func (r *repository) FindById(id string) (BlogTopic, error) {
	var data BlogTopic
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) CreateBlogTopic(p CreateBlogTopicRequest) (BlogTopic, error) {
	data := BlogTopic{
		BlogID:  p.BlogID,
		TopicID: p.TopicID}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateBlogTopic(p UpdateBlogTopicRequest) (BlogTopic, error) {
	data := BlogTopic{
		ID:      p.Id,
		BlogID:  p.BlogID,
		TopicID: p.TopicID}
	err := r.db.Updates(&data).Error
	return data, err
}

func (r *repository) DeleteBlogTopic(id int) (BlogTopic, error) {
	var data BlogTopic

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return BlogTopic{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return BlogTopic{}, err
	}

	// Step 3: Return the data
	return data, nil
}
