package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// OTPRepository handles OTP storage and retrieval in Redis.
type OTPRepository struct {
	rdb *redis.Client
}

// NewOTPRepository creates a new OTPRepository.
func NewOTPRepository(rdb *redis.Client) *OTPRepository {
	return &OTPRepository{rdb: rdb}
}

// StoreOTP saves the OTP for a given mobile number with a 3-minute expiration.
func (r *OTPRepository) StoreOTP(ctx context.Context, mobile, otp string) error {
	return r.rdb.Set(ctx, "otp:"+mobile, otp, 3*time.Minute).Err()
}

// GetOTP retrieves the OTP for a given mobile number. It returns an error if not found.
func (r *OTPRepository) GetOTP(ctx context.Context, mobile string) (string, error) {
	return r.rdb.Get(ctx, "otp:"+mobile).Result()
}
