package repository

import (
	"context"
	"database/sql"
	"instawall/internal/domain"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
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
