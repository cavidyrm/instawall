package postgres

import (
	"context"
	"github.com/cavidyrm/instawall/internal/user/domain" // <-- IMPORTANT: Replace with your actual module name
	"github.com/jmoiron/sqlx"
)

// UserRepository is a PostgreSQL implementation of the UserRepository.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user in the database.
func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO users (mobile_number, password_hash, name, email)
			  VALUES ($1, $2, $3, $4) RETURNING id, role, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, query, u.MobileNumber, u.PasswordHash, u.Name, u.Email).Scan(&u.ID, &u.Role, &u.CreatedAt, &u.UpdatedAt)
}

// GetByMobileNumber retrieves a user by their mobile number.
func (r *UserRepository) GetByMobileNumber(ctx context.Context, mobileNumber string) (*domain.User, error) {
	var u domain.User
	query := `SELECT id, mobile_number, password_hash, name, email, role, created_at, updated_at
			  FROM users WHERE mobile_number = $1`
	err := r.db.GetContext(ctx, &u, query, mobileNumber)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByID retrieves a user by their ID.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	query := `SELECT id, name, email, role, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &u, query, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
