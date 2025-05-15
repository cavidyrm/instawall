package usecase

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"instawall/internal/domain"
	"instawall/internal/repository"
)

// AuthUseCase defines the interface for authentication-related business logic.
type AuthUseCase interface {
	Register(username, email, password string) (*domain.User, error)
	Login(email, password string) (*domain.User, error)
	SendPasswordResetOTP(email string) error
	ValidatePasswordResetOTP(email, otp string) (*domain.User, error)
}

// authUseCase implements the AuthUseCase interface.
type authUseCase struct {
	userRepo repository.UserRepository
	otpRepo  repository.OTPRepository
}

// NewAuthUseCase creates a new instance of AuthUseCase.
func NewAuthUseCase(userRepo repository.UserRepository, otpRepo repository.OTPRepository) AuthUseCase {
	return &authUseCase{
		userRepo: userRepo,
		otpRepo:  otpRepo,
	}
}

// Register creates a new user.
func (uc *authUseCase) Register(username, email, password string) (*domain.User, error) {
	// Check if a user with the same email already exists
	existingUser, err := uc.userRepo.FindByEmail(email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create the new user
	newUser := &domain.User{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// Save the user to the database
	createdUser, err := uc.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

// Login authenticates a user.
func (uc *authUseCase) Login(email, password string) (*domain.User, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// SendPasswordResetOTP generates and sends an OTP for password reset.
func (uc *authUseCase) SendPasswordResetOTP(email string) error {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Generate a random OTP (you'll need a function for this)
	otpCode := "123456" // Placeholder for OTP generation

	// Set OTP expiry time (e.g., 15 minutes)
	expiresAt := time.Now().Add(15 * time.Minute)

	// Delete any existing OTPs for this user
	uc.otpRepo.DeleteByUserID(user.ID) // Delete any existing OTPs for this user

	// Create and save the new OTP
	newOTP := &domain.OTP{
		UserID:    user.ID,
		Code:      otpCode,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	_, err = uc.otpRepo.Create(newOTP)
	if err != nil {
		return err
	}

	// TODO: Implement sending the OTP via email or SMS

	return nil
}

// ValidatePasswordResetOTP validates the provided OTP for password reset.
func (uc *authUseCase) ValidatePasswordResetOTP(email, otp string) (*domain.User, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	otpEntry, err := uc.otpRepo.FindByUserIDAndCode(user.ID, otp)
	if err != nil {
		if errors.Is(err, repository.ErrOTPNotFound) {
			return nil, errors.New("invalid OTP")
		}
		return nil, err
	}

	// Check if the OTP has expired
	if time.Now().After(otpEntry.ExpiresAt) {
		// Optionally delete the expired OTP
		uc.otpRepo.Delete(otpEntry.ID) // Ignore error
		return nil, errors.New("OTP has expired")
	}

	// Optionally delete the OTP after successful validation to prevent reuse
	// uc.otpRepo.Delete(otpEntry.ID) // Ignore error

	return user, nil
}