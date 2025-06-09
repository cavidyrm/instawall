package repository

import (
	"database/sql"
	"time"
)

type PasswordResetToken struct {
	Token     string
	UserID    uint64
	ExpiresAt time.Time
}

type ResetTokenRepository interface {
	Create(token PasswordResetToken) error
	GetByToken(token string) (*PasswordResetToken, error)
	Delete(token string) error
}

type resetTokenRepo struct {
	db *sql.DB
}

func NewResetTokenRepository(db *sql.DB) ResetTokenRepository {
	return &resetTokenRepo{db: db}
}

func (r *resetTokenRepo) Create(t PasswordResetToken) error {
	_, err := r.db.Exec(`INSERT INTO password_reset_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`,
		t.UserID, t.Token, t.ExpiresAt)
	return err
}

func (r *resetTokenRepo) GetByToken(token string) (*PasswordResetToken, error) {
	row := r.db.QueryRow(`SELECT user_id, expires_at FROM password_reset_tokens WHERE token = $1`, token)
	var t PasswordResetToken
	t.Token = token
	if err := row.Scan(&t.UserID, &t.ExpiresAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *resetTokenRepo) Delete(token string) error {
	_, err := r.db.Exec(`DELETE FROM password_reset_tokens WHERE token = $1`, token)
	return err
}
