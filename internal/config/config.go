package config

import (
	"fmt"
	"os"

	"github.com/Georgi-Progger/task-tracker-backend/internal/config/env"
	"github.com/joho/godotenv"
)

type Config struct {
	DbConfig
	server
}

type server struct {
	Port string
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("error loading .env file: %w", err)
	}

	db, err := env.NewDbConfig()
	if err != nil {
		return Config{}, fmt.Errorf("error create config: %w", err)
	}

	appPort := os.Getenv("APP_PORT")

	cfg := Config{
		db,
		server{
			Port: appPort,
		},
	}
	return cfg, nil
}
