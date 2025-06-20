package project

import (
	"context"
	"fmt"
	"time"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/project_content_image"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/project_technology"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/statistic"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
	"gorm.io/gorm"
)

type Service interface {
	GetAllProjects(params GetAllProjectParams) ([]ProjectResponse, int, error)
	GetProjectByIdWithRelations(id int) (ProjectRelationResponse, error)
	GetProjectById(id int) (ProjectResponse, error)
	CreateProject(p CreateProjectRequest) (ProjectResponse, error)
	UpdateProject(p UpdateProjectRequest) (ProjectUpdateResponse, error)
	UpdateProjectStatistic(p ProjectStatisticUpdateRequest) (ProjectStatisticUpdateResponse, error)
	DeleteProject(id int) (Project, error)
	ChangeStatusProject(req ProjectChangeStatusRequest) (ProjectChangeStatusResponse, error)
}

type service struct {
	projectTechService   project_technology.Service
	projectImagesService project_content_image.Service
	statisticService     statistic.Service
	projectRepo          Repository
	db                   *gorm.DB
}

func NewService(
	projectTechSvc project_technology.Service,
	projctImagesSvc project_content_image.Service,
	statisticSvc statistic.Service,
	r Repository,
	db *gorm.DB,
) Service {
	return &service{
		projectTechService:   projectTechSvc,
		projectImagesService: projctImagesSvc,
		statisticService:     statisticSvc,
		projectRepo:          r,
		db:                   db,
	}
}

func (s *service) GetAllProjects(params GetAllProjectParams) ([]ProjectResponse, int, error) {
	datas, total, err := s.projectRepo.FindAll(params)
	if err != nil {
		return nil, 0, err
	}

	var result []ProjectResponse
	for _, p := range datas {
		result = append(result, ToProjectResponse(p))
	}
	return result, total, nil
}

func (s *service) GetProjectByIdWithRelations(id int) (ProjectRelationResponse, error) {
	data, err := s.projectRepo.FindByIdWithRelations(id)
	if err != nil {
		return ProjectRelationResponse{}, err
	}

	// Mapping result
	projectMap := map[int]*ProjectRelationResponse{}

	for _, row := range data {
		projectID := int(row.ID)

		//? "Comma-ok" itu fitur spesial
		_, exists := projectMap[projectID]
		if !exists {
			var publishedAtPointer *string
			if row.PublishedAt != nil {
				formattedPublishedAt := row.PublishedAt.Format("2006-01-02 15:04:05")
				publishedAtPointer = &formattedPublishedAt
			}

			var projectStatistic *ProjectStatisticDTO
			if row.StatisticID != 0 {
				projectStatistic = &ProjectStatisticDTO{
					ID:    row.StatisticID,
					Views: row.StatisticViews,
					Likes: row.StatisticLikes,
					Type:  row.StatisticType,
				}
			}
			projectMap[projectID] = &ProjectRelationResponse{
				ID:            projectID,
				Title:         row.Title,
				Slug:          row.Slug,
				Description:   row.Description,
				ImageUrl:      row.ImageUrl,
				ImageFileName: row.ImageFileName,
				RepositoryUrl: row.RepositoryUrl,
				Summary:       row.Summary,
				Status:        row.Status,
				IsHighlight:   row.IsHighlight,
				PublishedAt:   publishedAtPointer,
				CreatedAt:     row.CreatedAt.Format("2006-01-02 15:04:05"),
				StatisticID:   row.StatisticID,
				Statistic:     projectStatistic,
				Technologies:  []ProjectTechnologiesDTO{},
				ContentImages: []ProjectContentImagesDTO{},
			}
		}

		if row.ProjectTechnologyID != 0 {
			seen := make(map[int]bool)
			for _, tech := range projectMap[projectID].Technologies {
				seen[tech.ProjectTechID] = true
			}

			if !seen[row.ProjectTechnologyID] {
				projectMap[projectID].Technologies = append(projectMap[projectID].Technologies, ProjectTechnologiesDTO{
					ProjectTechID: row.ProjectTechnologyID,
					TechID:        row.TechnologyID,
					TechName:      row.TechnologyName,
				})
			}
		}

		if row.ProjectImgID != 0 {
			seen := make(map[int]bool)
			for _, img := range projectMap[projectID].ContentImages {
				seen[img.ProjectImageID] = true
			}

			if !seen[row.ProjectImgID] {
				projectMap[projectID].ContentImages = append(projectMap[projectID].ContentImages, ProjectContentImagesDTO{
					ProjectImageID: row.ProjectImgID,
					ImageFileName:  row.ProjectImgFileName,
					ImageUrl:       row.ProjectImgUrl,
				})
			}
		}

	}

	// Convert Map to Struct
	var result ProjectRelationResponse
	for _, v := range projectMap {
		result = *v
		break
	}

	return result, nil
}

func (s *service) GetProjectById(id int) (ProjectResponse, error) {
	data, err := s.projectRepo.FindById(id)
	if err != nil {
		return ProjectResponse{}, err
	}
	return data, nil
}

func (s *service) CreateProject(p CreateProjectRequest) (ProjectResponse, error) {
	//todo: Check Is Unique Slug
	slugVal := utils.StringToSlug(p.Slug)
	is_unique_slug, err := s.projectRepo.CheckUniqueSlug(slugVal)
	if err != nil {
		return ProjectResponse{}, err
	}
	if !is_unique_slug {
		err = fmt.Errorf("slug %s already exists", slugVal)
		return ProjectResponse{}, err
	}

	//todo: Check Technology Ids
	if err := s.projectTechService.CountTechnologiesByIDs(p.TechnologyIds); err != nil {
		return ProjectResponse{}, err
	}

	//todo: Check Project Images
	if len(p.ContentImages) > 0 {
		if err := s.projectImagesService.CountUnusedProjectImages(p.ContentImages); err != nil {
			return ProjectResponse{}, err
		}
	}

	tx := s.db.Begin()

	//todo: Create Statistic
	zero := 0
	statisticPayload := statistic.CreateStatisticRequest{
		Likes: &zero,
		Views: &zero,
		Type:  "Project"}

	statRes, err := s.statisticService.CreateStatisticWithTx(statisticPayload, tx)

	if err != nil {
		tx.Rollback()
		return ProjectResponse{}, err
	}

	//todo: Upload Image File to minio
	imageRes, err := utils.HandlUploadFile(p.ImageFile, "project")
	if err != nil {
		return ProjectResponse{}, err
	}

	uploadedImage := imageRes.FileURL

	var publishedAt *time.Time
	var status string
	if p.IsPublished == "Y" {
		now := time.Now()
		publishedAt = &now
		status = "PUBLISHED"
	} else if p.IsPublished == "N" {
		status = "UNPUBLISHED"
	}

	payload := CreateProjectDTO{
		StatisticID:          statRes.ID,
		ProjectContentImages: p.ContentImages,
		TechnologyIds:        p.TechnologyIds,
		Title:                p.Title,
		Description:          p.Description,
		ImageUrl:             imageRes.FileURL,
		ImageFileName:        imageRes.FileName,
		RepositoryUrl:        p.RepositoryUrl,
		Summary:              p.Summary,
		Status:               status,
		Slug:                 slugVal,
		IsHighlight:          false,
		PublishedAt:          publishedAt,
	}

	//todo: Create Project
	data, err := s.projectRepo.CreateProject(payload, tx)
	if err != nil {
		tx.Rollback()
		if uploadedImage != "" {
			_ = utils.DeleteFromMinio(context.Background(), uploadedImage)
		}
		return ProjectResponse{}, err
	}

	//todo: Bulk Create Project Technologies
	err = s.projectTechService.BulkCreateTechnologies(p.TechnologyIds, data.ID, tx)
	if err != nil {
		tx.Rollback()
		if uploadedImage != "" {
			_ = utils.DeleteFromMinio(context.Background(), uploadedImage)
		}
		return ProjectResponse{}, err
	}

	//todo: Batch Update Project Images
	err = s.projectImagesService.BatchUpdateProjectImages(p.ContentImages, data.ID, tx)
	if err != nil {
		tx.Rollback()
		if uploadedImage != "" {
			_ = utils.DeleteFromMinio(context.Background(), uploadedImage)
		}
		return ProjectResponse{}, err
	}

	//todo: Commit Transaction
	if err := tx.Commit().Error; err != nil {
		err = fmt.Errorf("error commit transaction")
		return ProjectResponse{}, err
	}

	return ToProjectResponse(data), nil
}

func (s *service) UpdateProject(p UpdateProjectRequest) (ProjectUpdateResponse, error) {
	//todo: Get Project
	project, err := s.GetProjectById(p.Id)
	if err != nil {
		return ProjectUpdateResponse{}, err
	}

	//todo: Check Is Unique Slug
	slugVal := utils.StringToSlug(p.Slug)

	if project.Slug != slugVal {
		is_unique_slug, err := s.projectRepo.CheckUniqueSlug(slugVal)
		if err != nil {
			return ProjectUpdateResponse{}, err
		}
		if !is_unique_slug {
			err = fmt.Errorf("slug %s already exists", slugVal)
			return ProjectUpdateResponse{}, err
		}
	}

	//todo: set oldFileName
	oldFileName := ""
	if project.ImageFileName != "" {
		oldFileName = project.ImageFileName
	}

	//todo: Check Technology Ids
	var tech_ids []int
	for _, tech := range p.TechnologyIds {
		tech_ids = append(tech_ids, tech.TechID)
	}
	if err := s.projectTechService.CountTechnologiesByIDs(tech_ids); err != nil {
		return ProjectUpdateResponse{}, err
	}

	tx := s.db.Begin()

	//todo: Batch Update Technologies
	err = s.projectTechService.BatchUpdateTechnologies(tech_ids, p.Id, tx)
	if err != nil {
		tx.Rollback()
		return ProjectUpdateResponse{}, err
	}

	// todo: Sync Project Images
	oldProjectImages, err := s.projectImagesService.SyncProjectImages(p.ProjectImages, p.Id, tx)

	if err != nil {
		tx.Rollback()
		return ProjectUpdateResponse{}, err
	}

	var newFileURL string
	var newFileName string

	if p.ImageFile != nil {
		imageRes, err := utils.HandlUploadFile(p.ImageFile, "project")
		if err != nil {
			return ProjectUpdateResponse{}, err
		}

		newFileURL = imageRes.FileURL
		newFileName = imageRes.FileName
	} else {
		newFileURL = project.ImageUrl // keep existing if not updated
		newFileName = project.ImageFileName
	}

	var publishedAt *time.Time
	var status string
	var oldIsPublished string
	if project.Status == "Published" {
		oldIsPublished = "Y"
	} else {
		oldIsPublished = "N"
	}

	if oldIsPublished != p.IsPublished {
		if p.IsPublished == "Y" {
			now := time.Now()
			publishedAt = &now
			status = "Published"
		} else if p.IsPublished == "N" {
			status = "Unpublished"
		}
	} else {
		status = project.Status
	}

	payload := UpdateProjectDTO{
		Id:            p.Id,
		Title:         p.Title,
		Description:   p.Description,
		ImageUrl:      newFileURL,
		ImageFileName: newFileName,
		RepositoryUrl: p.RepositoryUrl,
		Summary:       p.Summary,
		Status:        status,
		Slug:          slugVal,
		IsHighlight:   p.IsHighlight,
		PublishedAt:   publishedAt,
	}

	data, err := s.projectRepo.UpdateProject(payload, tx)
	if err != nil {
		tx.Rollback()
		return ProjectUpdateResponse{}, err
	}

	//todo: Delete Old Project Images
	if len(oldProjectImages) > 0 {
		slice_image_urls := []string{}
		for _, item := range oldProjectImages {
			slice_image_urls = append(slice_image_urls, item.ImageUrl)
		}

		err = s.projectImagesService.BulkDeleteHardByImageUrls(slice_image_urls, tx)
		if err != nil {
			tx.Rollback()
			return ProjectUpdateResponse{}, err
		}
	}

	//todo: Commit Transaction
	if err := tx.Commit().Error; err != nil {
		err = fmt.Errorf("error commit transaction")
		return ProjectUpdateResponse{}, err
	}

	//todo: Delete Old Image
	if oldFileName != newFileName {
		_ = utils.DeleteFromMinio(context.Background(), oldFileName)
	}

	return ToProjectUpdateResponse(data), nil
}

func (s *service) UpdateProjectStatistic(p ProjectStatisticUpdateRequest) (ProjectStatisticUpdateResponse, error) {
	project, err := s.GetProjectById(p.ProjectID)
	if err != nil {
		return ProjectStatisticUpdateResponse{}, err
	}

	payload := ProjectStatisticUpdateDTO{
		ProjectID:    p.ProjectID,
		ProjectTitle: project.Title,
		StatisticID:  p.StatisticID,
		Likes:        p.Likes,
		Views:        p.Views,
		Type:         p.Type,
	}

	data, err := s.projectRepo.UpdateProjectStatistic(payload)
	if err != nil {
		return ProjectStatisticUpdateResponse{}, err
	}
	return data, nil
}

func (s *service) DeleteProject(id int) (Project, error) {
	data, err := s.projectRepo.DeleteProject(id)
	if err != nil {
		return Project{}, err
	}
	return data, nil
}

func (s *service) ChangeStatusProject(req ProjectChangeStatusRequest) (ProjectChangeStatusResponse, error) {
	project, err := s.GetProjectById(req.ID)
	if err != nil {
		return ProjectChangeStatusResponse{}, err
	}

	data, err := s.projectRepo.ChangeStatusProject(req.ID, req.Status, project)
	if err != nil {
		return ProjectChangeStatusResponse{}, err
	}
	return data, nil
}
