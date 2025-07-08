package migration

import (
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

//go:embed ../../migrations/*.sql
var migrationsFS embed.FS

// Run applies all up migrations from the embedded SQL files.
func Run(db *sqlx.DB, dbName string) {
	log.Println("Starting database migration...")

	// 1. Create the source driver from the embedded filesystem.
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		log.Fatalf("could not create migration source driver: %v", err)
	}

	// 2. Create the database driver using the existing database connection.
	dbDriver, err := postgres.WithInstance(db.DB, &postgres.Config{
		DatabaseName: dbName,
	})
	if err != nil {
		log.Fatalf("could not create migration database driver: %v", err)
	}

	// 3. Create a new migrate instance using the correct function: NewWithInstance.
	// This function takes the driver instances directly.
	m, err := migrate.NewWithInstance(
		"iofs",       // A name for the source driver
		sourceDriver, // The source driver instance
		dbName,       // The database name
		dbDriver,     // The database driver instance
	)
	if err != nil {
		log.Fatalf("could not create migrate instance: %v", err)
	}

	// 4. Apply all available "up" migrations.
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Database is already up-to-date. No changes applied.")
		} else {
			log.Fatalf("failed to apply migrations: %v", err)
		}
	} else {
		log.Println("Database migration completed successfully.")
	}
}
