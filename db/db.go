package db

import (
	"os"

	"github.com/hiroyaonoe/le4-db-go/config"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

var (
	db *sqlx.DB
)

// InitDB はDBへのコネクションを構成する
func InitDB() {
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false},
	)
	dbrow := sqldblogger.OpenDriver(
		config.DSN(),
		pq.Driver{},
		zerologadapter.New(logger),
	)
	db = sqlx.NewDb(dbrow, "postgres")

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)
}

// GetDB は構成済みのDBへのコネクションを返す
func GetDB() *sqlx.DB {
	return db
}
