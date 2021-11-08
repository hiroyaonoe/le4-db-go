package db

import (
	"fmt"

	"github.com/hiroyaonoe/le4-db-go/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL: %w", err)
	}

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	return db, nil
}
