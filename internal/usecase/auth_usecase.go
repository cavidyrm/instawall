package usecase

import (
	"context"
	"instawall/internal/domain"
	"instawall/internal/repository"
	"instawall/pkg/jwt"
)

type AuthUseCase interface {
	Login(ctx context.Context, req *domain.LoginRequest) (string, error)
}

type authUseCase struct {
	repo repository.UserRepository
}

func NewAuthUseCase(repo repository.UserRepository) AuthUseCase {
	return &authUseCase{repo: repo}
}

func (uc *authUseCase) Login(ctx context.Context, req *domain.LoginRequest) (string, error) {
	user, err := uc.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	//if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
	//	return "", err
	//}

	return jwt.GenerateToken(user.ID)
}
