package project_technology

type Service interface {
	GetAllProjectTechnologies() ([]ProjectTechnologyResponse, error)
	GetProjectTechnologyById(id string) (ProjectTechnologyResponse, error)
	CreateProjectTechnology(p CreateProjectTechnologyRequest) (ProjectTechnologyResponse, error)
	UpdateProjectTechnology(p UpdateProjectTechnologyRequest) (ProjectTechnologyUpdateResponse, error)
	DeleteProjectTechnology(id int) (ProjectTechnology, error)
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

func (s *service) GetProjectTechnologyById(id string) (ProjectTechnologyResponse, error) {
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

func (s *service) UpdateProjectTechnology(p UpdateProjectTechnologyRequest) (ProjectTechnologyUpdateResponse, error) {
	data, err := s.repo.UpdateProjectTechnology(p)
	if err != nil {
		return ProjectTechnologyUpdateResponse{}, err
	}

	return ToProjectTechnologyUpdateResponse(data), nil
}

func (s *service) DeleteProjectTechnology(id int) (ProjectTechnology, error) {
	data, err := s.repo.DeleteProjectTechnology(id)
	if err != nil {
		return ProjectTechnology{}, err
	}
	return data, nil
}
