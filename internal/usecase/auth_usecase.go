package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"c:\Users\Cavid\Desktop\instawall/internal/domain"
	"c:\Users\Cavid\Desktop\instawall/pkg/otp"
)

type authUseCase struct {
	userRepo   domain.UserRepository
	otpRepo    domain.OTPRepository
	otpService otp.OTPService
}

func NewAuthUseCase(userRepo domain.UserRepository, otpRepo domain.OTPRepository, otpService otp.OTPService) domain.AuthUseCase {
	return &authUseCase{
		userRepo:   userRepo,
		otpRepo:    otpRepo,
		otpService: otpService,
	}
}

func (uc *authUseCase) Register(email, password string) error {
	// Check if user already exists
	_, err := uc.userRepo.GetByEmail(email)
	if err == nil {
		return errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user
	user := &domain.User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user
	if err := uc.userRepo.Create(user); err != nil {
		return err
	}

	// Generate OTP
	code := uc.otpService.Generate()

	// Create OTP record
	otpRecord := &domain.OTP{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Code:      code,
		Purpose:   "registration",
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
	}

	// Save OTP
	if err := uc.otpRepo.Create(otpRecord); err != nil {
		return err
	}

	// Send OTP
	return uc.otpService.Send(email, code)
}

func (uc *authUseCase) Login(email, password string) (string, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// In a real application, you would generate a JWT token here
	// For simplicity, we'll just return the user ID
	return user.ID, nil
}

func (uc *authUseCase) ForgotPassword(email string) error {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	// Generate OTP
	code := uc.otpService.Generate()

	// Create OTP record
	otpRecord := &domain.OTP{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Code:      code,
		Purpose:   "password-reset",
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
	}

	// Save OTP
	if err := uc.otpRepo.Create(otpRecord); err != nil {
		return err
	}

	// Send OTP
	return uc.otpService.Send(email, code)
}

func (uc *authUseCase) VerifyOTP(code string) (string, error) {
	// Get OTP by code
	otpRecord, err := uc.otpRepo.GetByCode(code)
	if err != nil {
		return "", errors.New("invalid or expired OTP")
	}

	// Delete OTP after verification
	if err := uc.otpRepo.Delete(otpRecord.ID); err != nil {
		return "", err
	}

	return otpRecord.UserID, nil
}

func (uc *authUseCase) ResetPassword(userID, password string) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	return uc.userRepo.UpdatePassword(userID, string(hashedPassword))
}