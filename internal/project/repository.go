package project

import (
	"fmt"
	"strings"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/statistic"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Project, error)
	FindByIdWithRelations(id int) ([]RawProjectRelationResponse, error)
	FindById(id int) (ProjectResponse, error)
	CreateProject(p CreateProjectDTO, tx *gorm.DB) (Project, error)
	CheckUpdateProjectTechnologies(projectTechs []ProjectTechUpdatePayload) (int, error)
	CheckUpdateProjectImages(projectImages []ProjectImagesUpdatePayload) (int, error)
	UpdateProject(p UpdateProjectDTO) (ProjectUpdateResponse, error)
	UpdateProjectStatistic(p ProjectStatisticUpdateDTO) (ProjectStatisticUpdateResponse, error)
	DeleteProject(id int) (Project, error)
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
		JOIN statistics s ON p.statistic_id = s.id
		JOIN project_technologies pt ON pt.project_id = p.id
		JOIN technologies t ON t.id = pt.technology_id
		JOIN project_content_images pci ON pci.project_id = p.id
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

func (r *repository) CheckUpdateProjectTechnologies(projectTechs []ProjectTechUpdatePayload) (total int, err error) {
	var ids []int
	for _, v := range projectTechs {
		ids = append(ids, v.TechID)
	}
	err = r.db.Raw(`
		SELECT COUNT(*) FROM technologies 
		WHERE id IN ? AND
		deleted_at IS NULL
	`, ids).Scan(&total).Error
	return total, err
}

func (r *repository) CheckUpdateProjectImages(projectImages []ProjectImagesUpdatePayload) (total int, err error) {
	var image_urls []string
	for _, v := range projectImages {
		image_urls = append(image_urls, v.ImageUrl)
	}
	err = r.db.Raw(`
		SELECT COUNT(*) FROM project_content_temp_images 
		WHERE image_url IN ? AND
		deleted_at IS NULL
	`, image_urls).Scan(&total).Error
	return total, err
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
		PublishedAt:   p.PublishedAt}

	err := db.Create(&data).Error

	if err != nil {
		return Project{}, err
	}

	return data, err
}

func (r *repository) BulkUpdateTechIDsByID(payload []ProjectTechUpdatePayload, tx *gorm.DB) error {
	var ids []interface{}
	var values []interface{}
	var caseStmt strings.Builder

	caseStmt.WriteString("CASE id")
	for _, item := range payload {
		caseStmt.WriteString(" WHEN ? THEN ?")
		values = append(values, item.ID, item.TechID)
		ids = append(ids, item.ID)
	}
	caseStmt.WriteString(" END")

	// Construct IN clause
	inClause := strings.TrimRight(strings.Repeat("?,", len(ids)), ",")

	// Combine all parameters: CASE params + IN list
	values = append(values, ids...)

	query := fmt.Sprintf(`
		UPDATE project_technologies 
		SET technology_id = %s 
		WHERE id IN (%s)
	`, caseStmt.String(), inClause)

	//todo: Output
	// UPDATE project_technologies
	// SET tech_id = CASE id
	// 	WHEN 1 THEN 2001
	// 	WHEN 2 THEN 2002
	// END
	// WHERE id IN (1, 2);

	if err := tx.Exec(query, values...).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) BulkUpdateContentImgIDsByID(payload []ProjectImagesUpdatePayload, tx *gorm.DB) error {
	var ids []interface{}
	var args []interface{}
	var imgUrlCase, ImgFileNameCase strings.Builder

	imgUrlCase.WriteString("CASE id")
	ImgFileNameCase.WriteString("CASE id")

	for _, p := range payload {
		imgUrlCase.WriteString(" WHEN ? THEN ?")
		args = append(args, p.ID, p.ImageUrl)

		ImgFileNameCase.WriteString(" WHEN ? THEN ?")
		args = append(args, p.ID, p.ImageFileName)

		ids = append(ids, p.ID)
	}

	imgUrlCase.WriteString(" END")
	ImgFileNameCase.WriteString(" END")

	// Add the IDs again for the WHERE clause
	placeholders := strings.TrimRight(strings.Repeat("?,", len(ids)), ",")
	args = append(args, utils.ToInterfaceSlice(ids)...)

	query := fmt.Sprintf(`
		UPDATE project_content_images 
		SET 
			image_url = %s,
			image_file_name = %s
		WHERE id IN (%s)
	`, imgUrlCase.String(), ImgFileNameCase.String(), placeholders)

	//todo: Output
	// UPDATE project_content_images
	// SET image_url = CASE id
	// 	WHEN 1 THEN 'https://example.com/image1.jpg'
	// 	WHEN 2 THEN 'https://example.com/image2.jpg'
	// END,
	// image_file_name = CASE id
	// 	WHEN 1 THEN 'image1.jpg'
	// 	WHEN 2 THEN 'image2.jpg'
	// END
	// WHERE id IN (1, 2);

	if err := tx.Exec(query, args...).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateProject(p UpdateProjectDTO) (ProjectUpdateResponse, error) {
	tx := r.db.Begin()

	//todo BEGIN: UPDATE PROJECT TECH

	err := r.BulkUpdateTechIDsByID(p.TechnologyIds, tx)
	if err != nil {
		tx.Rollback()
		return ProjectUpdateResponse{}, err
	}

	//todo BEGIN: UPDATE PROJECT CONTENT IMAGES

	err = r.BulkUpdateContentImgIDsByID(p.ContentImageIds, tx)
	if err != nil {
		tx.Rollback()
		return ProjectUpdateResponse{}, err
	}

	//todo BEGIN: UPDATE PROJECT

	data := Project{
		Title:         p.Title,
		Description:   p.Description,
		ImageUrl:      p.ImageUrl,
		ImageFileName: p.ImageFileName,
		RepositoryUrl: p.RepositoryUrl,
		Summary:       p.Summary,
		Status:        p.Status,
		PublishedAt:   p.PublishedAt}
	err = tx.Where("ID = ?", p.Id).Updates(&data).Error

	if err != nil {
		tx.Rollback()
		return ProjectUpdateResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return ProjectUpdateResponse{}, err
	}

	var publishedAt *string

	if p.PublishedAt != nil {
		publishedAtFormatted := p.PublishedAt.Format("2006-01-02")
		publishedAt = &publishedAtFormatted
	}
	res := ProjectUpdateResponse{
		ID:            p.Id,
		Title:         data.Title,
		Description:   data.Description,
		ImageUrl:      data.ImageUrl,
		ImageFileName: data.ImageFileName,
		RepositoryUrl: data.RepositoryUrl,
		Summary:       data.Summary,
		Status:        data.Status,
		PublishedAt:   publishedAt,
	}

	return res, nil
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
