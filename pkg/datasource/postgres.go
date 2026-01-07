package datasource

import (
	"os"

	"github.com/Georgi-Progger/task-tracker-backend/pkg/logger"
	"github.com/jmoiron/sqlx"
)

func NewDb(dsn string, logger logger.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error(err, "Failed to connect to database")
		os.Exit(1)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		logger.Error(err, "Failed to ping database")
		os.Exit(1)
	}

	return db, nil
}
