package testimonial

import "fmt"

type Service interface {
	GetAllTestimonials() ([]TestimonialResponse, error)
	GetTestimonialById(id int) (TestimonialResponse, error)
	CreateTestimonial(p CreateTestimonialRequest) (TestimonialResponse, error)
	UpdateTestimonial(p UpdateTestimonialRequest) error
	DeleteTestimonial(id int) (Testimonial, error)
	ChangeStatusTestimonial(id int, isUsed string) error
	ChangeMultiStatusTestimonial(ids []int, isUsed string) error
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
	payload := CreateTestimonialDTO{
		Name:      p.Name,
		Via:       p.Via,
		Role:      p.Role,
		WorkingAt: p.WorkingAt,
		IsUsed:    false,
	}

	data, err := s.repo.CreateTestimonial(payload)
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

	payload := UpdateTestimonialDTO{
		ID:        p.ID,
		Name:      p.Name,
		Via:       p.Via,
		Role:      p.Role,
		WorkingAt: p.WorkingAt,
		IsUsed:    p.IsUsed == "Y",
	}

	err = s.repo.UpdateTestimonial(payload)
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

func (s *service) ChangeStatusTestimonial(id int, isUsed string) error {
	_, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	isUsedBool := isUsed == "Y"

	err = s.repo.ChangeStatusTestimonial(id, isUsedBool)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ChangeMultiStatusTestimonial(ids []int, isUsed string) error {
	countData, err := s.repo.FindByMultiId(ids)
	if err != nil {
		return err
	}

	if len(countData) != len(ids) {
		err := fmt.Errorf("some testimonial_ids not found in database")
		return err
	}

	isUsedBool := isUsed == "Y"

	err = s.repo.ChangeMultiStatusTestimonial(ids, isUsedBool)
	if err != nil {
		return err
	}
	return nil
}
