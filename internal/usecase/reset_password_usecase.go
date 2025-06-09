package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"instawall/internal/repository"
)

type ResetPasswordUseCase interface {
	RequestReset(email string) (string, error)
	ResetPassword(token, newPassword string) error
}

type resetPasswordUC struct {
	userRepo  repository.UserRepository
	tokenRepo repository.ResetTokenRepository
	mailer    EmailSender
}

func NewResetPasswordUseCase(
	userRepo repository.UserRepository,
	tokenRepo repository.ResetTokenRepository,
	mailer EmailSender,
) ResetPasswordUseCase {
	return &resetPasswordUC{userRepo: userRepo, tokenRepo: tokenRepo, mailer: mailer}
}

func (uc *resetPasswordUC) RequestReset(email string) (string, error) {
	ctx := context.Background()
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found")
	}

	token, err := generateToken()
	if err != nil {
		return "", err
	}

	err = uc.tokenRepo.Create(repository.PasswordResetToken{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return "", err
	}

	// Send the email
	subject := "Reset Your Password"
	body := fmt.Sprintf(`Click to reset your password: <a href="https://yourapp.com/reset?token=%s">Reset</a>`, token)

	if err := uc.mailer.Send(email, subject, body); err != nil {
		return "", errors.New("could not send email")
	}

	return token, nil
}

func (uc *resetPasswordUC) ResetPassword(token, newPassword string) error {
	t, err := uc.tokenRepo.GetByToken(token)
	if err != nil || t.ExpiresAt.Before(time.Now()) {
		return errors.New("invalid or expired token")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := uc.userRepo.UpdatePassword(t.UserID, string(hashed)); err != nil {
		return err
	}

	return uc.tokenRepo.Delete(token)
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
