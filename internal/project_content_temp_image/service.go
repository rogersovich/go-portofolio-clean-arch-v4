package project_content_temp_image

type Service interface {
	GetAllProjectContentTempImgs() ([]ProjectContentTempImgResponse, error)
	GetProjectContentTempImgById(id string) (ProjectContentTempImgResponse, error)
	CreateProjectContentTempImg(p CreateProjectContentTempImgRequest) (ProjectContentTempImgResponse, error)
	UpdateProjectContentTempImg(p UpdateProjectContentTempImgRequest) (ProjectContentTempImgUpdateResponse, error)
	DeleteProjectContentTempImg(id int) (ProjectContentTempImages, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllProjectContentTempImgs() ([]ProjectContentTempImgResponse, error) {
	datas, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []ProjectContentTempImgResponse
	for _, p := range datas {
		result = append(result, ToProjectContentTempImgResponse(p))
	}
	return result, nil
}

func (s *service) GetProjectContentTempImgById(id string) (ProjectContentTempImgResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return ProjectContentTempImgResponse{}, err
	}
	return ToProjectContentTempImgResponse(data), nil
}

func (s *service) CreateProjectContentTempImg(p CreateProjectContentTempImgRequest) (ProjectContentTempImgResponse, error) {
	data, err := s.repo.CreateProjectContentTempImg(p)
	if err != nil {
		return ProjectContentTempImgResponse{}, err
	}
	return ToProjectContentTempImgResponse(data), nil
}

func (s *service) UpdateProjectContentTempImg(p UpdateProjectContentTempImgRequest) (ProjectContentTempImgUpdateResponse, error) {
	data, err := s.repo.UpdateProjectContentTempImg(p)
	if err != nil {
		return ProjectContentTempImgUpdateResponse{}, err
	}

	return ToProjectContentTempImgUpdateResponse(data), nil
}

func (s *service) DeleteProjectContentTempImg(id int) (ProjectContentTempImages, error) {
	data, err := s.repo.DeleteProjectContentTempImg(id)
	if err != nil {
		return ProjectContentTempImages{}, err
	}

	return data, nil
}
