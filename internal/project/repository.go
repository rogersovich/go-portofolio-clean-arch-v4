package project

import (
	"time"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/statistic"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Project, error)
	FindByIdWithRelations(id int) ([]RawProjectRelationResponse, error)
	FindById(id int) (ProjectResponse, error)
	CreateProject(p CreateProjectDTO, tx *gorm.DB) (Project, error)
	UpdateProject(p UpdateProjectDTO, tx *gorm.DB) (Project, error)
	UpdateProjectStatistic(p ProjectStatisticUpdateDTO) (ProjectStatisticUpdateResponse, error)
	DeleteProject(id int) (Project, error)
	CheckUniqueSlug(slug string) (bool, error)
	ChangeStatusProject(id int, status string, project ProjectResponse) (ProjectChangeStatusResponse, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Project, error) {
	var datas []Project
	err := r.db.Find(&datas).Error
	return datas, err
}

func (r *repository) FindByIdWithRelations(id int) ([]RawProjectRelationResponse, error) {
	var data []RawProjectRelationResponse
	err := r.db.Raw(`
		SELECT 
			p.id, 
			p.title, 
			p.description, 
			p.image_url, 
			p.image_file_name, 
			p.repository_url, 
			p.summary, 
			p.status, 
			p.published_at,
			p.created_at,
			s.id as statistic_id,
			s.views as statistic_views,
			s.likes as statistic_likes,
			s.type as statistic_type,
			pt.id as project_technology_id,
			t.id as technology_id,
			t.name as technology_name,
			pci.id as project_img_id,
      pci.image_file_name as project_img_file_name,
			pci.image_url as project_img_url
		FROM projects p 
		LEFT JOIN statistics s ON p.statistic_id = s.id
		LEFT JOIN project_technologies pt ON pt.project_id = p.id
		LEFT JOIN technologies t ON t.id = pt.technology_id
		LEFT JOIN project_content_images pci ON pci.project_id = p.id
		WHERE p.id = ? AND p.deleted_at IS NULL
	`, id).Scan(&data).Error

	if err != nil {
		return nil, err // handle DB or syntax error
	}

	if len(data) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return data, err
}

func (r *repository) FindById(id int) (ProjectResponse, error) {
	var data ProjectResponse
	err := r.db.Table("projects").Where("id = ?", id).Scan(&data).Error
	if data.ID == 0 {
		return ProjectResponse{}, gorm.ErrRecordNotFound
	}
	return data, err
}

func (r *repository) CreateProject(p CreateProjectDTO, tx *gorm.DB) (Project, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	// Create Project
	data := Project{
		StatisticID:   p.StatisticID,
		Title:         p.Title,
		Description:   p.Description,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		RepositoryUrl: p.RepositoryUrl,
		Summary:       p.Summary,
		Status:        p.Status,
		Slug:          p.Slug,
		PublishedAt:   p.PublishedAt}

	err := db.Create(&data).Error

	if err != nil {
		return Project{}, err
	}

	return data, err
}

func (r *repository) UpdateProject(p UpdateProjectDTO, tx *gorm.DB) (Project, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	//todo: UPDATE PROJECT

	data := Project{
		ID:            p.Id,
		Title:         p.Title,
		Description:   p.Description,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		RepositoryUrl: p.RepositoryUrl,
		Summary:       p.Summary,
		Status:        p.Status,
		Slug:          p.Slug,
		PublishedAt:   p.PublishedAt}

	err := db.Where("ID = ?", p.Id).Updates(&data).Error

	if err != nil {
		return Project{}, err
	}

	return data, nil
}

func (r *repository) UpdateProjectStatistic(p ProjectStatisticUpdateDTO) (ProjectStatisticUpdateResponse, error) {
	data := statistic.Statistic{
		ID:    p.StatisticID,
		Likes: p.Likes,
		Views: p.Views,
		Type:  p.Type,
	}
	err := r.db.Where("ID = ?", p.StatisticID).Updates(&data).Error
	if err != nil {
		return ProjectStatisticUpdateResponse{}, err
	}

	res := ProjectStatisticUpdateResponse{
		ProjectID:    p.ProjectID,
		ProjectTitle: p.ProjectTitle,
		StatisticID:  p.StatisticID,
		Likes:        *data.Likes,
		Views:        *data.Views,
		Type:         data.Type,
	}

	return res, nil
}

func (r *repository) DeleteProject(id int) (Project, error) {
	var data Project

	// Step 1: Find by ID
	if err := r.db.First(&data, id).Error; err != nil {
		return Project{}, err // return if not found or any error
	}

	// Step 2: Delete
	if err := r.db.Delete(&data).Error; err != nil {
		return Project{}, err
	}

	// Step 3: Return the data
	return data, nil
}

func (r *repository) CheckUniqueSlug(slug string) (bool, error) {
	var data ProjectResponse

	err := r.db.Table("projects").Where("slug = ?", slug).Scan(&data).Error

	if err != nil {
		return false, err // handle DB or syntax error
	}

	if data.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (r *repository) ChangeStatusProject(id int, status string, project ProjectResponse) (ProjectChangeStatusResponse, error) {
	now := time.Now()
	var updateMap = make(map[string]interface{})
	updateMap["status"] = status
	if status == "Published" {
		updateMap["published_at"] = now
	}
	err := r.db.Model(&Project{}).Where("id = ?", id).Updates(updateMap).Error

	if err != nil {
		return ProjectChangeStatusResponse{}, err
	}

	// Return the updated data
	var publishedAtStringPtr *string
	if status == "Published" {
		publishedAtString := now.Format("2006-01-02 15:04:05")
		publishedAtStringPtr = &publishedAtString
	}
	updatedData := ProjectChangeStatusResponse{
		ID:          id,
		Title:       project.Title,
		Status:      status,
		PublishedAt: publishedAtStringPtr,
	}

	return updatedData, nil
}
