package config

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Host string
	Port string
}
