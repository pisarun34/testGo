package models

import "gorm.io/gorm"

type SeeksterUser struct {
	gorm.Model
	PhoneNumber string `gorm:"type:varchar(15)" validate:"required,numeric,len=10:15"`
	Password    string `json:"password" validate:"required"`
	UUID        string `json:"uuid" validate:"required"`
	UserID      uint   `gorm:"foreignKey:ID"`
}
