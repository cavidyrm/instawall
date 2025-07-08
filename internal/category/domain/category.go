package domain

import (
	"github.com/google/uuid"
	"time"
)

// Category represents the core Category entity in the domain layer.
type Category struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	ImageURL    string    `db:"image_url"`
	CreatedAt   time.Time `db:"created_at"`
}
