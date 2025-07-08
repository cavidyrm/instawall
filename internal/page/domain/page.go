package domain

import (
	"github.com/google/uuid"
	"time"
)

// Page represents the core Page entity in the domain layer.
type Page struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	ImageURL    string    `db:"image_url"`
	Link        string    `db:"link"`
	HasIssue    bool      `db:"has_issue"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
