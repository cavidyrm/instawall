package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"instawall/internal/repository"
	"log"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"instawall/config"
)

func main() {
	cfg := config.LoadConfig()

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Migration driver error: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Make sure path is correct relative to project root
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("Migration setup failed: %v", err)
	}

	err = m.Up()
	switch err {
	case nil:
		fmt.Println("✅ Migrations applied successfully.")
	case migrate.ErrNoChange:
		fmt.Println("ℹ️ No migrations to apply.")
	default:
		log.Fatalf("Migration failed: %v", err)
	}
}
