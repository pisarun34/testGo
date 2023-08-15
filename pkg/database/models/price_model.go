package models

import "gorm.io/gorm"

type Price struct {
	gorm.Model
	Fractional   int    `gorm:"type:integer"`
	Decimal      string `gorm:"type:varchar(255)"`
	Tax          int    `gorm:"type:integer"`
	TaxDecimal   string `gorm:"type:varchar(255)"`
	DisplayValue string `gorm:"type:varchar(255)"`
	FullDisplay  string `gorm:"type:varchar(255)"`
	Currency     string `gorm:"type:varchar(255)"`
}
