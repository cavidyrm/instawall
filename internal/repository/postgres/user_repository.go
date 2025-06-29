package postgres

import (
	"context"
	"os/user"

	"github.com/cavidyrm/instawall/internal/domain"

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
			  VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, query, u.MobileNumber, u.PasswordHash, u.Name, u.Email).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

// GetByMobileNumber retrieves a user by their mobile number.
func (r *UserRepository) GetByMobileNumber(ctx context.Context, mobileNumber string) (*user.User, error) {
	var u user.User
	query := `SELECT id, mobile_number, password_hash, name, email, created_at, updated_at
			  FROM users WHERE mobile_number = $1`
	err := r.db.GetContext(ctx, &u, query, mobileNumber)
	if err != nil {
		return nil, err
	}
	return &u, nil
}