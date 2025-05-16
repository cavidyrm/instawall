package main

import (
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"instawall/config"
	"instawall/internal/validator"
	"log"

	"instawall/internal/delivery"
	"instawall/internal/repository"
	"instawall/internal/usecase"
)

func main() {
	cfg := config.LoadConfig()

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	authUC := usecase.NewAuthUseCase(userRepo, cfg.JWTSecret)

	e := echo.New()
	e.Validator = validator.NewValidator()
	delivery.NewAuthHandler(e, authUC)

	e.Logger.Fatal(e.Start(":8080"))
}
