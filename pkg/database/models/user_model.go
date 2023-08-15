package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	SSOID        string `gorm:"column:ssoid" json:"ssoid" validate:"required"`
	SeeksterUser SeeksterUser
}
