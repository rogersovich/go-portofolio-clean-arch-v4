package blog_topic

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]BlogTopic, error)
	FindById(id string) (BlogTopic, error)
	CreateBlogTopic(p CreateBlogTopicRequest) (BlogTopic, error)
	BulkCreateBlogTopic(topic_ids []int, blog_id int, tx *gorm.DB) error
	UpdateBlogTopic(p UpdateBlogTopicRequest) (BlogTopic, error)
	DeleteBlogTopic(id int) (BlogTopic, error)
	BulkDeleteHard(topic_ids []int, tx *gorm.DB) error
	FindExistingBlogTopics(blog_id int) ([]BlogTopicExistingResponse, error)
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

func (r *repository) FindExistingBlogTopics(blog_id int) ([]BlogTopicExistingResponse, error) {
	var data []BlogTopicExistingResponse
	err := r.db.Table("blog_topics").
		Where("blog_id = ?", blog_id).
		Select("id, blog_id, topic_id").
		Find(&data).Error
	return data, err
}

func (r *repository) CreateBlogTopic(p CreateBlogTopicRequest) (BlogTopic, error) {
	data := BlogTopic{
		BlogID:  p.BlogID,
		TopicID: p.TopicID}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) BulkCreateBlogTopic(topic_ids []int, blog_id int, tx *gorm.DB) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	var blog_topics []BlogTopic

	for _, topic_id := range topic_ids {
		blog_topics = append(blog_topics, BlogTopic{
			BlogID:  blog_id,
			TopicID: topic_id,
		})
	}

	if err := db.Create(&blog_topics).Error; err != nil {
		return err
	}

	return nil
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

func (r *repository) BulkDeleteHard(topic_ids []int, tx *gorm.DB) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	// Create a raw SQL query to delete records with IDs in the slice
	query := "DELETE FROM blog_topics WHERE topic_id IN ?"

	// Execute the raw query
	if err := db.Exec(query, topic_ids).Error; err != nil {
		return err
	}

	return nil
}
