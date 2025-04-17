package statistic

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Statistic, error)
	FindById(id string) (Statistic, error)
	CreateStatistic(p CreateStatisticRequest) (Statistic, error)
	UpdateStatistic(p UpdateStatisticRequest) (Statistic, error)
	DeleteStatistic(id int) (Statistic, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Statistic, error) {
	var datas []Statistic
	err := r.db.Find(&datas).Error
	return datas, err
}

func (r *repository) FindById(id string) (Statistic, error) {
	var data Statistic
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) CreateStatistic(p CreateStatisticRequest) (Statistic, error) {
	data := Statistic{
		Likes: *p.Likes,
		Views: *p.Views,
		Type:  p.Type}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateStatistic(p UpdateStatisticRequest) (Statistic, error) {
	data := Statistic{
		ID:    p.Id,
		Likes: *p.Likes,
		Views: *p.Views,
		Type:  p.Type}
	err := r.db.Updates(&data).Error
	return data, err
}

func (r *repository) DeleteStatistic(id int) (Statistic, error) {
	var data Statistic

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return Statistic{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return Statistic{}, err
	}

	// Step 3: Return the data
	return data, nil
}
