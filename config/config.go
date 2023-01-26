package config

type Job struct {
	Company_name string   `json:"company_name"`
	Date_posted  string   `json:"date_posted"`
	Job_title    string   `json:"job_title"`
	Job_id       string   `json:"_id"`
	Job_link     string   `json:"job_link"`
	Job_location string   `json:"job_location"`
	Keywords     []string `json:"keywords"`
	Source       string   `json:"source"`
}

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
	Host                   string
	Port                   string
	DeleteOldData          bool `mapstructure:"delete_old_data"`
	DeleteOldDataAfterDays int  `mapstructure:"delete_old_data_after"`
}

type FirebaseConfig struct {
	URL         string
	Credentials string
}

type MongoConfig struct {
	URL       string
	DockerURL string
}
