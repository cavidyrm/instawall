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
	signKey string
	repo    Repository
}

func New(signKey string, repo Repository) Service {
	return Service{repo: repo, signKey: signKey}
}

type RegisterUserRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type RegisterUserResponse struct {
	User domain.User
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
}

func (service *Service) Login(username string, password string) (LoginUserResponse, error) {

	//create token from service and return it to login response
	return LoginUserResponse{
		AccessToken: "token",
	}, nil

}

func (service *Service) Register(r RegisterUserRequest) (RegisterUserResponse, error) {
	if !utils.IsPhoneNumberValid(r.Phone) {
		return RegisterUserResponse{}, errors.New("invalid phone")
	}
	return RegisterUserResponse{}, nil
}
