package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID           uuid.UUID `db:"id"`
	MobileNumber string    `db:"mobile_number"`
	PasswordHash string    `db:"password_hash"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
