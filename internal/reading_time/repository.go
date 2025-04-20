package reading_time

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]ReadingTime, error)
	FindById(id string) (ReadingTime, error)
	CreateReadingTime(p CreateReadingTimeRequest) (ReadingTime, error)
	UpdateReadingTime(p UpdateReadingTimeRequest) (ReadingTime, error)
	DeleteReadingTime(id int) (ReadingTime, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]ReadingTime, error) {
	var datas []ReadingTime
	err := r.db.Find(&datas).Error
	return datas, err
}

func (r *repository) FindById(id string) (ReadingTime, error) {
	var data ReadingTime
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) CreateReadingTime(p CreateReadingTimeRequest) (ReadingTime, error) {
	data := ReadingTime{
		Minutes:          p.Minutes,
		TextLength:       p.TextLength,
		EstimatedSeconds: p.EstimatedSeconds,
		WordCount:        p.WordCount,
		Type:             p.Type}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateReadingTime(p UpdateReadingTimeRequest) (ReadingTime, error) {
	data := ReadingTime{
		ID:               p.Id,
		Minutes:          p.Minutes,
		TextLength:       p.TextLength,
		EstimatedSeconds: p.EstimatedSeconds,
		WordCount:        p.WordCount,
		Type:             p.Type}
	err := r.db.Updates(&data).Error
	return data, err
}

func (r *repository) DeleteReadingTime(id int) (ReadingTime, error) {
	var data ReadingTime

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return ReadingTime{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return ReadingTime{}, err
	}

	// Step 3: Return the data
	return data, nil
}
