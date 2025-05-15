package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq" // PostgreSQL driver

	"instawall/internal/delivery/httpserver"
	"instawall/internal/repository/postgresql"
	"instawall/internal/usecase"
)

func main() {
	// Load configuration (e.g., from environment variables)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// Initialize database connection
	db, err := postgresql.NewDatabaseConnection(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Initialize repositories
	userRepo := postgresql.NewUserRepository(db)
	otpRepo := postgresql.NewOTPRepository(db)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, otpRepo)

	// Initialize handlers
	authHandler := httpserver.NewAuthHandler(authUseCase)

	// Setup Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Setup routes
	httpserver.NewHTTPServer(e, authHandler)

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server stopped: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}