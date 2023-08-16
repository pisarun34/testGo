package models

import (
	"time"

	"gorm.io/gorm"
)

type Packages struct {
	gorm.Model
	CoverUrl             string                 `gorm:"type:varchar(255)"`
	BannerUrl            string                 `gorm:"type:varchar(255)"`
	Name                 string                 `gorm:"type:varchar(255)"`
	DisplayPackageCard   string                 `gorm:"type:varchar(255)"`
	DisplayUserType      string                 `gorm:"type:varchar(255)"`
	Description          string                 `gorm:"type:varchar(255)"`
	DisplayPrefixPrice   string                 `gorm:"type:varchar(255)"`
	DisplayPackagePrice  string                 `gorm:"type:varchar(255)"`
	DisplayPromotion     string                 `gorm:"type:varchar(255)"`
	VatDescription       string                 `gorm:"type:varchar(255)"`
	Type                 string                 `gorm:"type:varchar(255)"`
	Terms                string                 `gorm:"type:text"`
	Status               string                 `gorm:"type:ENUM('inactive', 'active', 'suspended');default:'Active'"`
	StartDate            time.Time              `gorm:"type:datetime(3)"`
	EndDate              time.Time              `gorm:"type:datetime(3)"`
	SubscriptionPackages []SubscriptionPackages `gorm:"foreignKey:PackageID"`
}
