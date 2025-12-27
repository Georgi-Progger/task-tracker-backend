package datasource

import (
	"log/slog"
	"os"

	"github.com/Georgi-Progger/task-tracker-backend/internal/config"
	"github.com/jmoiron/sqlx"
)

func NewDb(cfg config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.GetUrlDb())
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}

	return db, nil
}
