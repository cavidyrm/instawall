package postgresql

import (
	"database/sql"
	"errors"
	"time"

	"instawall/internal/domain" // Replace with your module name
)

type OTPRepository struct {
	db *sql.DB
}

func NewOTPRepository(db *sql.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) Create(otp *domain.OTP) error {
	query := `
		INSERT INTO otps (user_id, code, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := r.db.QueryRow(query, otp.UserID, otp.Code, otp.ExpiresAt, otp.CreatedAt).Scan(&otp.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *OTPRepository) FindByUserIDAndCode(userID int, code string) (*domain.OTP, error) {
	otp := &domain.OTP{}
	query := `
		SELECT id, user_id, code, expires_at, created_at
		FROM otps
		WHERE user_id = $1 AND code = $2
	`
	err := r.db.QueryRow(query, userID, code).Scan(
		&otp.ID,
		&otp.UserID,
		&otp.Code,
		&otp.ExpiresAt,
		&otp.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // OTP not found
		}
		return nil, err
	}
	return otp, nil
}

func (r *OTPRepository) Delete(userID int) error {
	query := `
		DELETE FROM otps
		WHERE user_id = $1
	`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
func (r *OTPRepository) DeleteExpired(userID int) error {
	query := `
		DELETE FROM otps
		WHERE user_id = $1 AND expires_at < $2
	`
	_, err := r.db.Exec(query, userID, time.Now())
	if err != nil {
		return err
	}
	return nil
}