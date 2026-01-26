package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type redisEnvConfig struct {
	PortDb string `env:"REDIS_PORT,required"`
	Host   string `env:"REDIS_HOST,required"`
}

type redisConfig struct {
	redis redisEnvConfig
}

func NewRedisConfig() (*redisConfig, error) {
	var redis redisEnvConfig
	err := env.Parse(&redis)
	if err != nil {
		return nil, err
	}
	return &redisConfig{redis: redis}, nil
}

func (cfg *redisConfig) GetUrlRedis() string {
	return fmt.Sprintf("%s:%s", cfg.redis.Host, cfg.redis.PortDb)
}
