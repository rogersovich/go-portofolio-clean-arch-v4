package auth

import (
	"fmt"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(req RegisterUserRequest) (RegisterResponse, error)
	LoginUser(req LoginUserRequest) (LoginResponse, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) RegisterUser(req RegisterUserRequest) (RegisterResponse, error) {
	//todo: Check Unique Email
	is_unique, _ := s.repo.CheckUniqueEmail(req.Email)
	if !is_unique {
		return RegisterResponse{}, fmt.Errorf("email already exist")
	}

	// Hash the password
	hashPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	payload := user.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashPass),
	}

	data, err := s.repo.RegisterUser(payload)
	if err != nil {
		return RegisterResponse{}, err
	}
	return data, nil
}

func (s *service) LoginUser(req LoginUserRequest) (LoginResponse, error) {
	data, err := s.repo.LoginUser(req)
	if err != nil {
		return LoginResponse{}, err
	}
	return data, nil
}
