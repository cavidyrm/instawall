package repository

import (
	"context"
	"time"

	"instawall/internal/domain"
)

// OTPRepository defines the methods for interacting with OTP data.
type OTPRepository interface {
	Create(ctx context.Context, otp *domain.OTP) error
	FindByUserIDAndCode(ctx context.Context, userID string, code string) (*domain.OTP, error)
	DeleteByUserID(ctx context.Context, userID string) error
}