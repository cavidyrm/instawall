package userservice

import (
	"errors"
	"instawall/internal/domain"
	"instawall/utils"
)

type Repository interface {
	IsPhoneNumberExist(phone string) bool
	Register(u domain.User) error
}
type Service struct {
	repo Repository
}

type RegisterUserRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type RegisterUserResponse struct {
	User domain.User
}

func (service *Service) Login(username string, password string) (string, error) {
	return "", nil
}

func (service *Service) Register(r RegisterUserRequest) (RegisterUserResponse, error) {
	if !utils.IsPhoneNumberValid(r.Phone) {
		return RegisterUserResponse{}, errors.New("invalid phone")
	}
	return RegisterUserResponse{}, nil
}
