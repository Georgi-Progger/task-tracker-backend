package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type brokerEnvConfig struct {
	BrokerPort string `env:"BROKER_APP_PORT,required"`
	BrokerHost string `env:"BROKER_APP_HOST,required"`
}

type brokerConfig struct {
	broker brokerEnvConfig
}

func NewBrokerConfig() (*brokerConfig, error) {
	var broker brokerEnvConfig
	err := env.Parse(&broker)
	if err != nil {
		return nil, err
	}
	return &brokerConfig{broker: broker}, nil
}

func (cfg *brokerConfig) GetUrlBroker() string {
	return fmt.Sprintf("%s:%s", cfg.broker.BrokerHost, cfg.broker.BrokerPort)
}
