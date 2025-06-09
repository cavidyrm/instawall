package main

import (
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"instawall/config"
	"instawall/internal/validator"
	"instawall/pkg/email"
	"log"

	"instawall/internal/delivery"
	"instawall/internal/repository"
	"instawall/internal/usecase"
)

func main() {
	cfg := config.LoadConfig()
	mailer := email.NewSMTPSender(
		cfg.SMTPFrom,
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPass,
	)
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	authUC := usecase.NewAuthUseCase(userRepo, cfg.JWTSecret)

	e := echo.New()
	e.Validator = validator.NewValidator()
	delivery.NewAuthHandler(e, authUC)
	resetTokenRepo := repository.NewResetTokenRepository(db)

	resetUC := usecase.NewResetPasswordUseCase(userRepo, resetTokenRepo, mailer)
	delivery.NewResetPasswordHandler(e, resetUC)

	e.Logger.Fatal(e.Start(":8080"))
}
