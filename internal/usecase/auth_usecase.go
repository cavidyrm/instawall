package usecase

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"instawall/internal/domain"
	"instawall/internal/repository"
	"instawall/pkg/jwt"
)

type AuthUseCase interface {
	Login(ctx context.Context, req *domain.LoginRequest) (string, error)
	Register(ctx context.Context, req *domain.RegisterRequest) (string, error)
}

type authUseCase struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthUseCase(repo repository.UserRepository, jwtSecret string) *authUseCase {
	return &authUseCase{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (uc *authUseCase) Login(ctx context.Context, req *domain.LoginRequest) (string, error) {
	user, err := uc.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", err
	}

	return jwt.GenerateToken(user.ID)
}

func (uc *authUseCase) Register(ctx context.Context, req *domain.RegisterRequest) (string, error) {
	exists, err := uc.repo.IsEmailTaken(ctx, req.Email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("email already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &domain.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = uc.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	return jwt.GenerateToken(user.ID)
}

//func (uc *authUseCase) ForgotPassword(ctx context.Context, email string) error {
//	user, err := uc.repo.GetByEmail(ctx, email)
//	if err != nil {
//		return fmt.Errorf("user not found")
//	}
//
//	token := uuid.NewString()
//	expiresAt := time.Now().Add(15 * time.Minute)
//
//	err = uc.repo.(ctx, user.ID, token, expiresAt)
//	if err != nil {
//		return err
//	}
//
//	// TODO: Send email. For now, just log token.
//	fmt.Printf("Password reset link: http://localhost:8080/reset-password?token=%s\n", token)
//	return nil
//}
//
//func (uc *authUseCase) ResetPassword(ctx context.Context, token, newPassword string) error {
//	userID, err := uc.repo.GetUserIDByResetToken(ctx, token)
//	if err != nil {
//		return err
//	}
//
//	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
//	if err != nil {
//		return err
//	}
//
//	if err := uc.repo.UpdateUserPassword(ctx, userID, string(hashed)); err != nil {
//		return err
//	}
//
//	return uc.repo.DeleteResetToken(ctx, token)
//}
