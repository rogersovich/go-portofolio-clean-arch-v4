package project_content_temp_image

import (
	"context"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]ProjectContentTempImages, error)
	FindById(id string) (ProjectContentTempImages, error)
	CreateProjectContentTempImg(p CreateProjectContentTempImgRequest) (ProjectContentTempImages, error)
	UpdateProjectContentTempImg(p UpdateProjectContentTempImgRequest) (ProjectContentTempImages, error)
	DeleteProjectContentTempImg(id int) (ProjectContentTempImages, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]ProjectContentTempImages, error) {
	var datas []ProjectContentTempImages
	err := r.db.Find(&datas).Error
	return datas, err
}

func (r *repository) FindById(id string) (ProjectContentTempImages, error) {
	var data ProjectContentTempImages
	err := r.db.Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *repository) CreateProjectContentTempImg(p CreateProjectContentTempImgRequest) (ProjectContentTempImages, error) {
	data := ProjectContentTempImages{
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		IsUsed:        p.IsUsed == "Y"}
	err := r.db.Create(&data).Error
	return data, err
}

func (r *repository) UpdateProjectContentTempImg(p UpdateProjectContentTempImgRequest) (ProjectContentTempImages, error) {
	data := ProjectContentTempImages{
		ID:            p.Id,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		IsUsed:        p.IsUsed == "Y"}
	err := r.db.Updates(&data).Error
	return data, err
}

func (r *repository) DeleteProjectContentTempImg(id int) (ProjectContentTempImages, error) {
	var data ProjectContentTempImages

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return ProjectContentTempImages{}, err // return if not found or any error
	}

	oldFileName := data.ImageFileName

	// Step 2: Permanently Delete
	if err := r.db.Unscoped().Delete(&data).Error; err != nil {
		return ProjectContentTempImages{}, err
	}

	err := utils.DeleteFromMinio(context.Background(), oldFileName)
	if err != nil {
		utils.Logger.Error(err.Error())
	}

	// Step 3: Return the data
	return data, nil
}
