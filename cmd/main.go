package main

import (
	"log"

	userHandler "github.com/cavidyrm/instawall/internal/user/delivery"
	"github.com/cavidyrm/instawall/internal/user/repository/postgres"
	"github.com/cavidyrm/instawall/internal/user/usecase"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection (use environment variables in production)
	dbURL := "postgres://youruser:yourpassword@localhost:5432/yourdb?sslmode=disable"
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Initialize Echo
	e := echo.New()

	// Standard Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize dependencies
	userRepo := postgres.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Register handlers (which in turn register routes)
	userHandler.NewUserHandler(e, userUsecase)

	// Start server
	log.Println("Starting server on :8080")
	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
	}
}