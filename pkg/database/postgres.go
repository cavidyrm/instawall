// pkg/database/postgres.go
package database

import (
	"fmt"
	"github.com/cavidyrm/instawall/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPostgresDB creates a new PostgreSQL database connection.
func NewPostgresDB(cfg config.PostgresConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection is alive.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
