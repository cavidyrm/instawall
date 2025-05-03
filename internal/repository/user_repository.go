package repository

import (
	"errors"
	"sync"

	"c:\Users\Cavid\Desktop\instawall/internal/domain"
)

type userRepository struct {
	users map[string]*domain.User
	mutex sync.RWMutex
}

func NewUserRepository() domain.UserRepository {
	return &userRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *userRepository) Create(user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if user with the same email already exists
	for _, u := range r.users {
		if u.Email == user.Email {
			return errors.New("user with this email already exists")
		}
	}

	r.users[user.ID] = user
	return nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *userRepository) UpdatePassword(userID, password string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	user, exists := r.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	user.Password = password
	return nil
}