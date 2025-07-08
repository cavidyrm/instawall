package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/cavidyrm/instawall/internal/middleware"
	"github.com/cavidyrm/instawall/internal/user/domain"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
)

// UserRepository defines the interface for user data storage.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByMobileNumber(ctx context.Context, mobileNumber string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

// OTPRepository defines the interface for OTP storage.
type OTPRepository interface {
	StoreOTP(ctx context.Context, mobile, otp string) error
	GetOTP(ctx context.Context, mobile string) (string, error)
}

// UserUsecase provides all user-related business logic.
type UserUsecase struct {
	userRepo UserRepository
	otpRepo  OTPRepository
}

// NewUserUsecase creates a new UserUsecase.
func NewUserUsecase(userRepo UserRepository, otpRepo OTPRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo, otpRepo: otpRepo}
}

// SendOTP generates, stores, and "sends" an OTP.
func (uc *UserUsecase) SendOTP(ctx context.Context, mobileNumber string) error {
	otp := generateOTP(6)
	log.Printf("Generated OTP for %s: %s\n", mobileNumber, otp) // For development
	return uc.otpRepo.StoreOTP(ctx, mobileNumber, otp)
}

// VerifyOTP checks the OTP and returns a temporary registration token if valid.
func (uc *UserUsecase) VerifyOTP(ctx context.Context, mobileNumber, otp string) (string, error) {
	storedOTP, err := uc.otpRepo.GetOTP(ctx, mobileNumber)
	if err != nil {
		return "", fmt.Errorf("OTP expired or not found")
	}
	if storedOTP != otp {
		return "", fmt.Errorf("invalid OTP")
	}
	return middleware.GenerateRegistrationToken(mobileNumber)
}

// CompleteRegistration creates the user after OTP has been verified.
func (uc *UserUsecase) CompleteRegistration(ctx context.Context, mobileNumber, password, name, email string) (*domain.User, error) {
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

// Login authenticates a user and returns a standard JWT.
func (uc *UserUsecase) Login(ctx context.Context, mobileNumber, password string) (string, error) {
	existingUser, err := uc.userRepo.GetByMobileNumber(ctx, mobileNumber)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(password)); err != nil {
		return "", err
	}
	return middleware.GenerateToken(existingUser.ID.String(), existingUser.Name, existingUser.Role)
}

// GetProfile retrieves a user's public profile.
func (uc *UserUsecase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	return uc.userRepo.GetByID(ctx, userID)
}

// generateOTP creates a random n-digit string.
func generateOTP(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
