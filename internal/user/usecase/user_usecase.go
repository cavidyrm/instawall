package usecase

import (
	"context"

	"github.com/cavidyrm/instawall/internal/middleware"
	"github.com/cavidyrm/instawall/internal/user/domain"
	"golang.org/x/crypto/bcrypt"
)

// UserUsecase provides user-related use cases.
type UserUsecase struct {
	userRepo UserRepository
}

// UserRepository interface remains the same.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByMobileNumber(ctx context.Context, mobileNumber string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error) // Add this for profile retrieval
}

// NewUserUsecase creates a new UserUsecase.
func NewUserUsecase(userRepo UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

// Register remains the same.
func (uc *UserUsecase) Register(ctx context.Context, mobileNumber, password, name, email string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		MobileNumber: mobileNumber,
		PasswordHash: string(hashedPassword),
		Name:         name,
		Email:        email,
	}

	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// Login authenticates a user and returns a JWT.
func (uc *UserUsecase) Login(ctx context.Context, mobileNumber, password string) (string, error) {
	existingUser, err := uc.userRepo.GetByMobileNumber(ctx, mobileNumber)
	if err != nil {
		return "", err // Consider custom error types for "user not found"
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(password)); err != nil {
		return "", err // Invalid password
	}

	// Generate JWT
	token, err := middleware.GenerateToken(existingUser.ID.String(), existingUser.Name)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetProfile retrieves a user's profile.
func (uc *UserUsecase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	return uc.userRepo.GetByID(ctx, userID)
}