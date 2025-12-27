package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type dbEnvConfig struct {
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	DbName   string `env:"DB_NAME,required"`
	PortDb   string `env:"DB_PORT,required"`
	Host     string `env:"DB_HOST,required"`
}

type dbConfig struct {
	db dbEnvConfig
}

func NewDbConfig() (*dbConfig, error) {
	var db dbEnvConfig
	err := env.Parse(&db)
	if err != nil {
		return nil, err
	}
	return &dbConfig{db: db}, nil
}

func (cfg *dbConfig) GetUrlDb() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.db.User, cfg.db.Password, cfg.db.DbName, cfg.db.Host, cfg.db.PortDb)
}
