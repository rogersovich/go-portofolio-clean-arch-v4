package testimonial

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Testimonial, error)
	FindById(id int) (Testimonial, error)
	CreateTestimonial(p CreateTestimonialRequest) (Testimonial, error)
	UpdateTestimonial(p UpdateTestimonialRequest) error
	DeleteTestimonial(id int) (Testimonial, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Testimonial, error) {
	var datas []Testimonial
	err := r.db.Find(&datas).Error
	return datas, err
}

func (r *repository) FindById(id int) (Testimonial, error) {
	var data Testimonial
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) CreateTestimonial(p CreateTestimonialRequest) (Testimonial, error) {
	data := Testimonial{
		Name:      p.Name,
		Via:       p.Via,
		Role:      p.Role,
		WorkingAt: p.WorkingAt,
	}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateTestimonial(p UpdateTestimonialRequest) error {
	data := Testimonial{
		ID:        p.ID,
		Name:      p.Name,
		Via:       p.Via,
		Role:      p.Role,
		WorkingAt: p.WorkingAt,
	}
	err := r.db.Updates(&data).Error
	return err
}

func (r *repository) DeleteTestimonial(id int) (Testimonial, error) {
	var data Testimonial

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return Testimonial{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return Testimonial{}, err
	}

	// Step 3: Return the data
	return data, nil
}
