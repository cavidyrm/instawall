package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"instawall/internal/domain" // Replace with your module path
)

type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository.
// It takes a *sql.DB connection as a dependency.
// This function is part of the repository layer, specifically the PostgreSQL implementation.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *domain.User) error {
	query := `INSERT INTO users (id, username, email, password, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// Check for errors during password hashing.
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)
	// Generate a new UUID for the user ID.
	user.ID = uuid.New()

	// Execute the insert query and scan the returned ID into the user struct.
	err = r.db.QueryRow(query, user.ID, user.Username, user.Email, user.Password, user.CreatedAt).Scan(&user.ID)
	// Check for errors during database insertion.
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = $1`

	user := &domain.User{}
	// Execute the select query and scan the results into the user struct.
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		// If no rows are returned, indicate that the user was not found.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return user, nil
}

func (r *UserRepository) FindUserByID(id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE id = $1`

	user := &domain.User{}
	// Execute the select query and scan the results into the user struct.
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		// If no rows are returned, indicate that the user was not found.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(user *domain.User) error {
    query := `UPDATE users SET username = $1, email = $2, password = $3 WHERE id = $4`

	// Execute the update query.
    _, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.ID)
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }

    return nil
}

func (r *UserRepository) DeleteUser(id uuid.UUID) error {
    query := `DELETE FROM users WHERE id = $1`

	// Execute the delete query.
    _, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }

    return nil
}