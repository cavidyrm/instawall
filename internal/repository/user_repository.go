package repository

import (
	"context"
	"database/sql"
	"fmt"
	"instawall/internal/domain"
	"time"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user *domain.User) error
	SaveResetToken(ctx context.Context, userID uint64, token string, expiresAt time.Time) error
	UpdatePassword(userID uint64, newHashedPassword string) error
}

type userRepo struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{DB: db}
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := domain.User{}
	err := r.DB.QueryRowContext(ctx, "SELECT id, email, password FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.DB.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	return exists, err
}

func (r *userRepo) CreateUser(ctx context.Context, user *domain.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	return r.DB.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&user.ID)
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.DB.QueryRowContext(ctx, "SELECT id, email FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Email)
	return user, err
}

func (r *userRepo) SaveResetToken(ctx context.Context, userID uint64, token string, expiresAt time.Time) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO password_reset_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, token, expiresAt)
	return err
}

func (r *userRepo) GetUserIDByResetToken(ctx context.Context, token string) (int, error) {
	var userID int
	var expiresAt time.Time
	err := r.DB.QueryRowContext(ctx, "SELECT user_id, expires_at FROM password_reset_tokens WHERE token = $1", token).
		Scan(&userID, &expiresAt)

	if err != nil {
		return 0, err
	}

	if time.Now().After(expiresAt) {
		return 0, fmt.Errorf("reset token expired")
	}

	return userID, nil
}

func (r *userRepo) UpdateUserPassword(ctx context.Context, userID uint64, newHashedPassword string) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE users SET password = $1 WHERE id = $2", newHashedPassword, userID)
	return err
}

func (r *userRepo) DeleteResetToken(ctx context.Context, token string) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM password_reset_tokens WHERE token = $1", token)
	return err
}

func (r *userRepo) UpdatePassword(userID uint64, newHashedPassword string) error {
	_, err := r.DB.Exec(`UPDATE users SET password = $1 WHERE id = $2`, newHashedPassword, userID)
	return err
}
