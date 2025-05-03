package repository

import (
	"errors"
	"sync"
	"time"

	"c:\Users\Cavid\Desktop\instawall/internal/domain"
)

type otpRepository struct {
	otps  map[string]*domain.OTP
	mutex sync.RWMutex
}

func NewOTPRepository() domain.OTPRepository {
	return &otpRepository{
		otps: make(map[string]*domain.OTP),
	}
}

func (r *otpRepository) Create(otp *domain.OTP) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.otps[otp.ID] = otp
	return nil
}

func (r *otpRepository) GetByCode(code string) (*domain.OTP, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, otp := range r.otps {
		if otp.Code == code {
			// Check if OTP is expired
			if time.Now().After(otp.ExpiresAt) {
				return nil, errors.New("otp expired")
			}
			return otp, nil
		}
	}

	return nil, errors.New("otp not found")
}

func (r *otpRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.otps, id)
	return nil
}