package testimonial

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Testimonial, error)
	FindById(id int) (Testimonial, error)
	FindByMultiId(ids []int) ([]Testimonial, error)
	CreateTestimonial(p CreateTestimonialDTO) (Testimonial, error)
	UpdateTestimonial(p UpdateTestimonialDTO) error
	DeleteTestimonial(id int) (Testimonial, error)
	ChangeStatusTestimonial(id int, isUsed bool) error
	ChangeMultiStatusTestimonial(ids []int, isUsed bool) error
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

func (r *repository) FindByMultiId(ids []int) ([]Testimonial, error) {
	var datas []Testimonial
	err := r.db.Table("testimonials").Where("id IN ?", ids).Scan(&datas).Error
	return datas, err
}

func (r *repository) CreateTestimonial(p CreateTestimonialDTO) (Testimonial, error) {
	data := Testimonial{
		Name:      p.Name,
		Via:       p.Via,
		Role:      p.Role,
		WorkingAt: p.WorkingAt,
		IsUsed:    p.IsUsed,
	}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateTestimonial(p UpdateTestimonialDTO) error {
	updates := map[string]interface{}{
		"id":         p.ID,
		"name":       p.Name,
		"via":        p.Via,
		"role":       p.Role,
		"working_at": p.WorkingAt,
		"is_used":    p.IsUsed,
	}
	err := r.db.Table("testimonials").Where("id = ?", p.ID).Updates(updates).Error
	if err != nil {
		return err
	}
	return err
}

func (r *repository) DeleteTestimonial(id int) (Testimonial, error) {
	var data Testimonial

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return Testimonial{}, err // return if not found or any error
	}

	// Step 2: Hard Delete
	if err := r.db.Unscoped().Delete(&data).Error; err != nil {
		return Testimonial{}, err
	}

	// Step 3: Return the data
	return data, nil
}

func (r *repository) ChangeStatusTestimonial(id int, isUsed bool) error {
	updates := map[string]interface{}{
		"is_used": isUsed,
	}
	err := r.db.Table("testimonials").Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return err
	}
	return err
}

func (r *repository) ChangeMultiStatusTestimonial(ids []int, isUsed bool) error {
	updates := map[string]interface{}{
		"is_used": isUsed,
	}
	err := r.db.Table("testimonials").Where("id IN ?", ids).Updates(updates).Error
	if err != nil {
		return err
	}
	return err
}
