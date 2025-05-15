package domain

import (
	"time"
)

type OTP struct {
	ID int `json:"id" db:"id"`
	UserID int `json:"user_id" db:"user_id"`
	Code string `json:"code" db:"code"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}