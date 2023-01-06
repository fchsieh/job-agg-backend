package config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Firebase FirebaseConfig
	Mongo    MongoConfig
}

type DatabaseConfig struct {
	Database   string
	Collection string
	Document   string
}

type ServerConfig struct {
	Host string
	Port string
}

type FirebaseConfig struct {
	URL         string
	Credentials string
}

type MongoConfig struct {
	URL string
}
