package mysql

import (
	"os"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	SSOID        string `gorm:"column:ssoid" json:"ssoid" validate:"required"`
	SeeksterUser SeeksterUser
}

type SeeksterUser struct {
	gorm.Model
	PhoneNumber string `gorm:"type:varchar(15)" validate:"required,numeric,len=10:15"`
	Password    string `json:"password" validate:"required"`
	UUID        string `json:"uuid" validate:"required"`
	UserID      uint   `gorm:"foreignKey:ID"`
}

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
