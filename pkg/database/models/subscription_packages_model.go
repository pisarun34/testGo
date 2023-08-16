package models

import "gorm.io/gorm"

type SubscriptionPackages struct {
	gorm.Model
	PackageNumber      string `gorm:"type:varchar(255)"`
	Name               string `gorm:"type:varchar(255)"`
	PackageName        string `gorm:"type:varchar(255)"`
	PackageDescription string `gorm:"type:varchar(255)"`
	Quantity           int    `gorm:"type:integer"`
	QuantityUnit       string `gorm:"type:varchar(255)"`
	DisplayQuantity    string `gorm:"type:varchar(255)"`
	PriceID            uint   `gorm:"foreignKey:PriceRef"`
	PackageID          uint   `gorm:"foreignKey:PackagesRef"`
	Status             string `gorm:"type:ENUM('inactive', 'active', 'suspended');default:'Active'"`
	Price              Price
}
