package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"instawall/config"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	return sql.Open("postgres", cfg.GetDBConnStr())
}
