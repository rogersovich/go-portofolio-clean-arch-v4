package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]User, error)
	FindById(id int) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(id int) error
	CheckUniqueEmail(email string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]User, error) {
	var datas []User
	err := r.db.Find(&datas).Error
	return datas, err
}

func (r *repository) FindById(id int) (User, error) {
	var data User
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) UpdateUser(user User) (User, error) {
	err := r.db.Save(&user).Error
	return user, err
}

func (r *repository) DeleteUser(id int) error {
	var data User
	err := r.db.Where("id = ?", id).Delete(&data).Error
	return err
}

func (r *repository) CheckUniqueEmail(email string) (bool, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return true, nil
	}

	return false, nil
}
