package config

type DbConfig interface {
	GetUrlDb() string
}

type RedisConfig interface {
	GetUrlRedis() string
}

type AppConfig interface {
	GetPort() string
}

type BrokerConfig interface {
	GetUrlBroker() string
}
