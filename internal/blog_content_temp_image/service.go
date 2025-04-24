package blog_content_temp_image

import "fmt"

type Service interface {
	GetAllBlogContentTempImgs() ([]BlogContentTempImgResponse, error)
	GetBlogContentTempImgById(id string) (BlogContentTempImgResponse, error)
	CreateBlogContentTempImg(p CreateBlogContentTempImgRequest) (BlogContentTempImgResponse, error)
	UpdateBlogContentTempImg(p UpdateBlogContentTempImgRequest) (BlogContentTempImgUpdateResponse, error)
	DeleteBlogContentTempImg(id int) (BlogContentTempImages, error)
	CountTempImages(tempImages []CountTempImagesDTO) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllBlogContentTempImgs() ([]BlogContentTempImgResponse, error) {
	datas, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []BlogContentTempImgResponse
	for _, p := range datas {
		result = append(result, ToBlogContentTempImgResponse(p))
	}
	return result, nil
}

func (s *service) GetBlogContentTempImgById(id string) (BlogContentTempImgResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return BlogContentTempImgResponse{}, err
	}
	return ToBlogContentTempImgResponse(data), nil
}

func (s *service) CreateBlogContentTempImg(p CreateBlogContentTempImgRequest) (BlogContentTempImgResponse, error) {
	data, err := s.repo.CreateBlogContentTempImg(p)
	if err != nil {
		return BlogContentTempImgResponse{}, err
	}
	return ToBlogContentTempImgResponse(data), nil
}

func (s *service) UpdateBlogContentTempImg(p UpdateBlogContentTempImgRequest) (BlogContentTempImgUpdateResponse, error) {
	data, err := s.repo.UpdateBlogContentTempImg(p)
	if err != nil {
		return BlogContentTempImgUpdateResponse{}, err
	}

	return ToBlogContentTempImgUpdateResponse(data), nil
}

func (s *service) DeleteBlogContentTempImg(id int) (BlogContentTempImages, error) {
	data, err := s.repo.DeleteBlogContentTempImg(id)
	if err != nil {
		return BlogContentTempImages{}, err
	}

	return data, nil
}

func (s *service) CountTempImages(tempImages []CountTempImagesDTO) error {
	total, err := s.repo.CountTempImages(tempImages)
	if err != nil {
		return err
	}

	if total != len(tempImages) {
		err := fmt.Errorf("some blog_content_temp_images not found in database")
		return err
	}
	return nil
}
