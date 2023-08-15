package dtos

type Packages struct {
	ID                   uint                   `json:"id"`
	CoverUrl             string                 `json:"cover_url"`
	BannerUrl            string                 `json:"banner_url"`
	Name                 string                 `json:"name"`
	DisplayPackageCard   string                 `json:"display_package_card"`
	DisplayUserType      string                 `json:"display_user_type"`
	Description          string                 `json:"description"`
	DisplayPrefixPrice   string                 `json:"display_prefix_price"`
	DisplayPackagePrice  string                 `json:"display_package_price"`
	DisplayPromotion     string                 `json:"display_promotion"`
	VatDescription       string                 `json:"vat_description"`
	Type                 string                 `json:"type"`
	Terms                string                 `json:"terms"`
	SubscriptionPackages []SubscriptionPackages `json:"packages" gorm:"foreignKey:PackageID"`
}
