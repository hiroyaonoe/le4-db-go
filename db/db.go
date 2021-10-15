package db

import (
	"fmt"

	"github.com/hiroyaonoe/le4db-go/config"
	"github.com/jmoiron/sqlx"
)

func NewDB() (*DB, error) {
	db, err := sqlx.Connect("postgres", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL: %w", err)
	}

	return db, nil
}
