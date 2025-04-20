package project

import (
	"context"
	"fmt"
	"time"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

type Service interface {
	GetAllProjects() ([]ProjectResponse, error)
	GetProjectByIdWithRelations(id int) (ProjectRelationResponse, error)
	GetProjectById(id int) (ProjectResponse, error)
	CreateProject(p CreateProjectRequest) (ProjectResponse, error)
	UpdateProject(p UpdateProjectRequest) (ProjectUpdateResponse, error)
	DeleteProject(id int) (Project, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllProjects() ([]ProjectResponse, error) {
	datas, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []ProjectResponse
	for _, p := range datas {
		result = append(result, ToProjectResponse(p))
	}
	return result, nil
}

func (s *service) GetProjectByIdWithRelations(id int) (ProjectRelationResponse, error) {
	data, err := s.repo.FindByIdWithRelations(id)
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

			projectMap[projectID] = &ProjectRelationResponse{
				ID:          projectID,
				StatisticID: row.StatisticID,
				Statistic: ProjectStatisticDTO{
					ID:        row.StatisticID,
					Views:     row.StatisticViews,
					Likes:     row.StatisticLikes,
					Type:      row.StatisticType,
					CreatedAt: row.CreatedAt.Format("2006-01-02 15:04:05"),
				},
				Title:         row.Title,
				Description:   row.Description,
				ImageUrl:      row.ImageUrl,
				ImageFileName: row.ImageFileName,
				RepositoryUrl: row.RepositoryUrl,
				Summary:       row.Summary,
				Status:        row.Status,
				PublishedAt:   publishedAtPointer,
				CreatedAt:     row.CreatedAt.Format("2006-01-02 15:04:05"),
			}
		}

		projectMap[projectID].Technologies = append(projectMap[projectID].Technologies, ProjectTechnologiesDTO{
			ProjectTechID: row.ProjectTechnologyID,
			TechID:        row.TechnologyID,
			TechName:      row.TechnologyName,
		})

		projectMap[projectID].ContentImages = append(projectMap[projectID].ContentImages, ProjectContentImagesDTO{
			ProjectImageID: row.ProjectImgID,
			ImageFileName:  row.ProjectImgFileName,
			ImageUrl:       row.ProjectImgUrl,
		})
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
	data, err := s.repo.FindById(id)
	if err != nil {
		return ProjectResponse{}, err
	}
	return data, nil
}

func (s *service) CheckUpdateProjectTechnologies(projectTechs []ProjectTechUpdatePayload) error {
	total, err := s.repo.CheckUpdateProjectTechnologies(projectTechs)
	if err != nil {
		return err
	}

	if total != len(projectTechs) {
		err := fmt.Errorf("some technologies not found in database")
		return err
	}
	return nil
}

func (s *service) CheckCreateProjectTechnologies(ids []string) error {
	total, err := s.repo.CheckCreateProjectTechnologies(ids)
	if err != nil {
		return err
	}

	if total != len(ids) {
		err := fmt.Errorf("some technologies not found in database")
		return err
	}
	return nil
}

func (s *service) CheckCreateProjectImages(ids []string) error {
	total, err := s.repo.CheckCreateProjectImages(ids)
	if err != nil {
		return err
	}

	if total != len(ids) {
		err := fmt.Errorf("some project_content_images not found in database")
		return err
	}
	return nil
}

func (s *service) CheckUpdateProjectImages(project_id int, projectImages []ProjectImagesUpdatePayload) error {
	total, err := s.repo.CheckUpdateProjectImages(project_id, projectImages)
	if err != nil {
		return err
	}

	if total != len(projectImages) {
		err := fmt.Errorf("some project_content_images not found in database")
		return err
	}
	return nil
}

func (s *service) CreateProject(p CreateProjectRequest) (ProjectResponse, error) {
	if err := s.CheckCreateProjectTechnologies(p.TechnologyIds); err != nil {
		return ProjectResponse{}, err
	}

	if len(p.ContentImages) > 0 {
		if err := s.CheckCreateProjectImages(p.ContentImages); err != nil {
			return ProjectResponse{}, err
		}
	}

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
		ProjectContentImages: p.ContentImages,
		TechnologyIds:        p.TechnologyIds,
		Title:                p.Title,
		Description:          p.Description,
		ImageUrl:             imageRes.FileURL,
		ImageFileName:        imageRes.FileName,
		RepositoryUrl:        p.RepositoryUrl,
		Summary:              p.Summary,
		Status:               status,
		PublishedAt:          publishedAt,
	}
	data, err := s.repo.CreateProject(payload)
	if err != nil {
		if uploadedImage != "" {
			_ = utils.DeleteFromMinio(context.Background(), uploadedImage)
		}
		return ProjectResponse{}, err
	}
	return ToProjectResponse(data), nil
}

func (s *service) UpdateProject(p UpdateProjectRequest) (ProjectUpdateResponse, error) {
	project, err := s.GetProjectById(p.Id)
	if err != nil {
		return ProjectUpdateResponse{}, err
	}

	//* set oldFileName
	oldFileName := ""
	if project.ImageFileName != "" {
		oldFileName = project.ImageFileName
	}

	if err := s.CheckUpdateProjectTechnologies(p.TechnologyIds); err != nil {
		return ProjectUpdateResponse{}, err
	}

	//* Check on Table Temp Project Images
	// if len(p.ContentImageIds) > 0 {
	// 	if err := s.CheckUpdateProjectImages(p.Id, p.ContentImageIds); err != nil {
	// 		return ProjectUpdateResponse{}, err
	// 	}
	// }

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
	if p.IsPublished == "Y" {
		now := time.Now()
		publishedAt = &now
		status = "PUBLISHED"
	} else if p.IsPublished == "N" {
		status = "UNPUBLISHED"
	}

	payload := UpdateProjectDTO{
		Id:              p.Id,
		ContentImageIds: p.ContentImageIds,
		TechnologyIds:   p.TechnologyIds,
		Title:           p.Title,
		Description:     p.Description,
		ImageUrl:        newFileURL,
		ImageFileName:   newFileName,
		RepositoryUrl:   p.RepositoryUrl,
		Summary:         p.Summary,
		Status:          status,
		PublishedAt:     publishedAt,
	}
	_ = payload

	data, err := s.repo.UpdateProject(payload)
	if err != nil {
		return ProjectUpdateResponse{}, err
	}

	//* 3. Optional: Delete old file from MinIO
	if oldFileName != newFileName {
		err = utils.DeleteFromMinio(context.Background(), oldFileName) // ignore error or handle if needed
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}

	return data, nil
}

func (s *service) DeleteProject(id int) (Project, error) {
	data, err := s.repo.DeleteProject(id)
	if err != nil {
		return Project{}, err
	}
	return data, nil
}
