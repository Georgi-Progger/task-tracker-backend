package config

import (
	"fmt"

	"github.com/Georgi-Progger/task-tracker-backend/internal/config/env"
	"github.com/joho/godotenv"
)

type Config struct {
	DbConfig
	AppConfig
	BrokerConfig
	RedisConfig
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("error loading .env file: %w", err)
	}

	db, err := env.NewDbConfig()
	if err != nil {
		return Config{}, fmt.Errorf("error create config: %w", err)
	}

	app, err := env.NewAppConfig()
	if err != nil {
		return Config{}, fmt.Errorf("error create config: %w", err)
	}

	broker, err := env.NewBrokerConfig()
	if err != nil {
		return Config{}, fmt.Errorf("error create config: %w", err)
	}

	redis, err := env.NewRedisConfig()
	if err != nil {
		return Config{}, fmt.Errorf("error create config: %w", err)
	}

	cfg := Config{
		DbConfig:     db,
		AppConfig:    app,
		BrokerConfig: broker,
		RedisConfig:  redis,
	}
	return cfg, nil
}
