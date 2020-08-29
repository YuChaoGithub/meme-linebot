package config

import (
	"os"
	"time"
)

// Config contains a ServerConfig, LineBotConfig and a DBConfig for the configurations of the app.
type Config struct {
	AdminSecret string
	Server      ServerConfig
	LineBot     LineBotConfig
	DB          DBConfig
}

// ServerConfig defines the configurations of the webserver.
type ServerConfig struct {
	Port         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// LineBotConfig defines the configurations of the line bot client.
type LineBotConfig struct {
	ChannelSecret      string
	ChannelAccessToken string
}

// DBConfig defines the configurations of the DB connection.
type DBConfig struct {
	Dialect       string
	ConnectionURL string
}

var conf Config

// Initialize the config struct from the environment variables.
func init() {
	conf = Config{
		AdminSecret: os.Getenv("ADMIN_SECRET"),
		Server: ServerConfig{
			Port:         ":" + os.Getenv("PORT"),
			IdleTimeout:  time.Minute,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		LineBot: LineBotConfig{
			ChannelSecret:      os.Getenv("LINE_CHANNEL_SECRET"),
			ChannelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		},
		DB: DBConfig{
			Dialect:       "postgres",
			ConnectionURL: os.Getenv("DATABASE_URL"),
		},
	}
}

// GetConfig returns the initialized configuration.
func GetConfig() *Config {
	return &conf
}
