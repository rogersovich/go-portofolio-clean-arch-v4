package topic

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Topic, error)
	FindById(id int) (Topic, error)
	CreateTopic(p CreateTopicRequest) (Topic, error)
	UpdateTopic(p UpdateTopicRequest) error
	DeleteTopic(id int) (Topic, error)
	CheckTopicIds(ids []int) ([]Topic, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Topic, error) {
	var datas []Topic
	err := r.db.Find(&datas).Error
	return datas, err
}

func (r *repository) FindById(id int) (Topic, error) {
	var data Topic
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) CreateTopic(p CreateTopicRequest) (Topic, error) {
	data := Topic{
		Name: p.Name}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateTopic(p UpdateTopicRequest) error {
	data := Topic{
		ID:   p.ID,
		Name: p.Name}
	err := r.db.Updates(&data).Error
	return err
}

func (r *repository) DeleteTopic(id int) (Topic, error) {
	var data Topic

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return Topic{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return Topic{}, err
	}

	// Step 3: Return the data
	return data, nil
}

func (r *repository) CheckTopicIds(ids []int) ([]Topic, error) {
	var data []Topic
	err := r.db.Where("id in (?)", ids).Select("id", "name").Find(&data).Error
	return data, err
}
