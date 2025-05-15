package repository

import "clean_architecture_example/internal/domain" // Replace with your module path

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}