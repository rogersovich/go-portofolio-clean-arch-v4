package project_technology

import (
	"fmt"
	"slices"

	"gorm.io/gorm"
)

type Service interface {
	GetAllProjectTechnologies() ([]ProjectTechnologyResponse, error)
	GetProjectTechnologyById(id int) (ProjectTechnologyResponse, error)
	CreateProjectTechnology(p CreateProjectTechnologyRequest) (ProjectTechnologyResponse, error)
	UpdateProjectTechnology(p UpdateProjectTechnologyRequest) error
	DeleteProjectTechnology(id int) (ProjectTechnology, error)
	CountTechnologiesByIDs(ids []int) error
	BulkCreateTechnologies(tech_ids []int, project_id int, tx *gorm.DB) error
	BatchUpdateTechnologies(tech_ids []int, project_id int, tx *gorm.DB) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllProjectTechnologies() ([]ProjectTechnologyResponse, error) {
	datas, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []ProjectTechnologyResponse
	for _, p := range datas {
		result = append(result, ToProjectTechnologyResponse(p))
	}
	return result, nil
}

func (s *service) GetProjectTechnologyById(id int) (ProjectTechnologyResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return ProjectTechnologyResponse{}, err
	}
	return ToProjectTechnologyResponse(data), nil
}

func (s *service) CreateProjectTechnology(p CreateProjectTechnologyRequest) (ProjectTechnologyResponse, error) {
	data, err := s.repo.CreateProjectTechnology(p)
	if err != nil {
		return ProjectTechnologyResponse{}, err
	}
	return ToProjectTechnologyResponse(data), nil
}

func (s *service) UpdateProjectTechnology(p UpdateProjectTechnologyRequest) error {
	_, err := s.repo.FindById(p.ID)
	if err != nil {
		return err
	}

	err = s.repo.UpdateProjectTechnology(p)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteProjectTechnology(id int) (ProjectTechnology, error) {
	data, err := s.repo.DeleteProjectTechnology(id)
	if err != nil {
		return ProjectTechnology{}, err
	}
	return data, nil
}

func (s *service) CountTechnologiesByIDs(ids []int) error {
	total, err := s.repo.CountTechnologiesByIDs(ids)
	if err != nil {
		return err
	}

	if total != len(ids) {
		err := fmt.Errorf("some technology_ids not found in database")
		return err
	}
	return nil
}

func (s *service) BulkCreateTechnologies(tech_ids []int, project_id int, tx *gorm.DB) error {
	var technologies []ProjectTechnology

	for _, technology_id := range tech_ids {
		technologies = append(technologies, ProjectTechnology{
			ProjectID:    project_id,
			TechnologyID: technology_id,
		})
	}

	return s.repo.BulkCreateTechnologies(technologies, tx)
}

func (s *service) BatchUpdateTechnologies(tech_ids []int, project_id int, tx *gorm.DB) error {
	existingProjectTechs, err := s.repo.FindExistingProjectTechnologies(project_id)
	if err != nil {
		return err
	}

	var existing_tech_ids []int
	var isNewTech bool = true
	for _, item := range existingProjectTechs {
		if slices.Contains(tech_ids, item.TechnologyID) {
			isNewTech = false
		}
		existing_tech_ids = append(existing_tech_ids, item.TechnologyID)
	}

	if !isNewTech && len(tech_ids) == len(existing_tech_ids) {
		return nil
	}

	err = s.repo.BulkHardDeleteTechnology(existing_tech_ids, project_id, tx)
	if err != nil {
		return err
	}

	newTechIds := []ProjectTechnology{}
	for _, tech_id := range tech_ids {
		newTechIds = append(newTechIds, ProjectTechnology{
			ProjectID:    project_id,
			TechnologyID: tech_id,
		})
	}

	err = s.repo.BulkCreateTechnologies(newTechIds, tx)
	if err != nil {
		return err
	}

	return nil
}
