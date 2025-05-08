package user

import "fmt"

type Service interface {
	GetAllUsers(params GetAllUserParams) ([]UserResponse, int, error)
	GetUserById(id int) (UserResponse, error)
	UpdateUser(user User) (UserResponse, error)
	DeleteUser(id int) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllUsers(params GetAllUserParams) ([]UserResponse, int, error) {
	datas, total, err := s.repo.FindAll(params)
	if err != nil {
		return nil, 0, err
	}

	var result []UserResponse
	for _, p := range datas {
		result = append(result, ToUserResponse(p))
	}
	return result, total, nil
}

func (s *service) GetUserById(id int) (UserResponse, error) {
	data, err := s.repo.FindById(id)
	if err != nil {
		return UserResponse{}, err
	}
	return ToUserResponse(data), nil
}

func (s *service) UpdateUser(user User) (UserResponse, error) {
	//todo: Get User
	oldData, err := s.repo.FindById(user.ID)
	if err != nil {
		return UserResponse{}, err
	}

	if oldData.Email != user.Email {
		//todo: Check Unique Email
		is_unique, _ := s.repo.CheckUniqueEmail(user.Email)
		if !is_unique {
			return UserResponse{}, fmt.Errorf("email already exist")
		}
	}

	user.CreatedAt = oldData.CreatedAt
	user.Password = oldData.Password

	data, err := s.repo.UpdateUser(user)
	if err != nil {
		return UserResponse{}, err
	}
	return ToUserResponse(data), nil
}

func (s *service) DeleteUser(id int) error {
	err := s.repo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
