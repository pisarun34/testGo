package dtos

type SubscriptionPackages struct {
	ID                 uint   `json:"id"`
	PackageNumber      string `json:"package_number"`
	Name               string `json:"name"`
	PackageName        string `json:"package_name"`
	PackageDescription string `json:"package_description"`
	Quantity           int    `json:"quantity"`
	QuantityUnit       string `json:"quantity_unit"`
	DisplayQuantity    string `json:"display_quantity"`
	PriceID            uint   `json:"-" gorm:"foreignKey:PriceRef"`
	PackageID          uint   `json:"-" gorm:"foreignKey:PackagesRef"`
	Price              Price  `json:"price"`
}
