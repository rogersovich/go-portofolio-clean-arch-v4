package testimonial

type Service interface {
	GetAllTestimonials() ([]TestimonialResponse, error)
	GetTestimonialById(id int) (TestimonialResponse, error)
	CreateTestimonial(p CreateTestimonialRequest) (TestimonialResponse, error)
	UpdateTestimonial(p UpdateTestimonialRequest) error
	DeleteTestimonial(id int) (Testimonial, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllTestimonials() ([]TestimonialResponse, error) {
	datas, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []TestimonialResponse
	for _, p := range datas {
		result = append(result, ToTestimonialResponse(p))
	}
	return result, nil
}

func (s *service) GetTestimonialById(id int) (TestimonialResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return TestimonialResponse{}, err
	}
	return ToTestimonialResponse(data), nil
}

func (s *service) CreateTestimonial(p CreateTestimonialRequest) (TestimonialResponse, error) {
	data, err := s.repo.CreateTestimonial(p)
	if err != nil {
		return TestimonialResponse{}, err
	}
	return ToTestimonialResponse(data), nil
}

func (s *service) UpdateTestimonial(p UpdateTestimonialRequest) error {
	_, err := s.repo.FindById(p.ID)
	if err != nil {
		return err
	}

	err = s.repo.UpdateTestimonial(p)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteTestimonial(id int) (Testimonial, error) {
	data, err := s.repo.DeleteTestimonial(id)
	if err != nil {
		return Testimonial{}, err
	}
	return data, nil
}
