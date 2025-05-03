package auth

import (
	"fmt"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/user"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repository interface {
	RegisterUser(user user.User) (RegisterResponse, error)
	LoginUser(user LoginUserRequest) (LoginResponse, error)
	CheckUniqueEmail(email string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) RegisterUser(payload user.User) (RegisterResponse, error) {
	err := r.db.Create(&payload).Error

	if err != nil {
		return RegisterResponse{}, err
	}

	var user user.User
	err = r.db.Where("email = ?", payload.Email).First(&user).Error

	if err != nil {
		return RegisterResponse{}, err
	}

	response := RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return response, err
}

func (r *repository) LoginUser(payload LoginUserRequest) (LoginResponse, error) {
	var user user.User
	if err := r.db.Where("email = ? AND deleted_at IS NULL", payload.Email).First(&user).Error; err != nil {
		err = fmt.Errorf("email not found")
		return LoginResponse{}, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		err = fmt.Errorf("invalid username or password")
		return LoginResponse{}, err
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		err = fmt.Errorf("error generating token")
		return LoginResponse{}, err
	}

	response := LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}

	return response, nil
}

func (r *repository) CheckUniqueEmail(email string) (bool, error) {
	var user user.User
	err := r.db.Where("email = ?", email).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return true, nil
	}

	return false, nil
}
