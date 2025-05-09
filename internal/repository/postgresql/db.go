package postgresql

import (
	"database/sql"
	"fmt"
	"log"
)

type DB struct {
	db *sql.DB
}

func New() *DB {
	//postgres://%s:%s@%s:%s/%s
	connStr := "postgres://user:password@localhost:5432/test"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Errorf("failed to connect to postgresql database: %w", err))
	}
	log.Println("connected to postgresql database")
	return &DB{db: db}
}
