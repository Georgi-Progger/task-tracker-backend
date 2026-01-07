package config

type DbConfig interface {
	GetUrlDb() string
}

type AppConfig interface {
	GetPort() string
}

type BrokerConfig interface {
	GetUrlBroker() string
}
