package main

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"instawall/internal/delivery"
	"instawall/internal/repository"
	"instawall/internal/usecase"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/test?sslmode=disable")
	if err != nil {
		panic(err)
	}

	e := echo.New()

	userRepo := repository.NewUserRepository(db)
	authUC := usecase.NewAuthUseCase(userRepo)
	delivery.NewAuthHandler(e, authUC)

	e.Logger.Fatal(e.Start(":8080"))
}
