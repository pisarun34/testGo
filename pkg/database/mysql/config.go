package mysql

import (
	"os"
)

type Config struct {
	User string
	Pass string
	Name string
	Host string
}

func NewConfig() *Config {
	return &Config{
		User: os.Getenv("DB_USER"),
		Pass: os.Getenv("DB_PASS"),
		Name: os.Getenv("DB_NAME"),
		Host: os.Getenv("DB_HOST"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
