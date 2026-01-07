package env

import "github.com/caarlos0/env/v11"

type appEnvConfig struct {
	AppPort string `env:"APP_PORT,required"`
}

type appConfig struct {
	app appEnvConfig
}

func NewAppConfig() (*appConfig, error) {
	var app appEnvConfig
	err := env.Parse(&app)
	if err != nil {
		return nil, err
	}
	return &appConfig{app: app}, nil
}

func (cfg *appConfig) GetPort() string {
	return cfg.app.AppPort
}
