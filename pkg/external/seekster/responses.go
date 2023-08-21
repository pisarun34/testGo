package seekster

import (
	"time"
)

// RegisterUser คือโครงสร้างข้อมูล Register ที่ได้จากการเรียกใช้งาน API ของ Seekster
type SignResponse struct {
	ID                int    `json:"id"`
	Type              string `json:"type"`
	Locale            string `json:"locale"`
	RegistrationToken string `json:"registration_token"`
	AccessToken       string `json:"access_token"`
	UUID              string `json:"uuid"`
	User              struct {
		ID                   int         `json:"id"`
		Type                 string      `json:"type"`
		AvatarURL            string      `json:"avatar_url"`
		DisplayName          string      `json:"display_name"`
		FirstName            interface{} `json:"first_name"`
		LastName             interface{} `json:"last_name"`
		Email                interface{} `json:"email"`
		Verified             bool        `json:"verified"`
		AcceptedLatestPolicy bool        `json:"accepted_latest_policy"`
		AcceptedLatestTerms  bool        `json:"accepted_latest_terms"`
	} `json:"user"`
	Client struct {
		ID                      int           `json:"id"`
		Name                    string        `json:"name"`
		Slug                    string        `json:"slug"`
		SecretKey               string        `json:"secret_key"`
		Platform                string        `json:"platform"`
		IssuedByID              int           `json:"issued_by_id"`
		CreatedAt               time.Time     `json:"created_at"`
		UpdatedAt               time.Time     `json:"updated_at"`
		AccessType              string        `json:"access_type"`
		LatestVersion           string        `json:"latest_version"`
		MinimumSupportedVersion string        `json:"minimum_supported_version"`
		Official                bool          `json:"official"`
		AccessKey               string        `json:"access_key"`
		Color                   string        `json:"color"`
		WebhookURL              string        `json:"webhook_url"`
		Refs                    []interface{} `json:"refs"`
		FcmServerKey            string        `json:"fcm_server_key"`
		Scheme                  string        `json:"scheme"`
		ContactEmail            interface{}   `json:"contact_email"`
		ContactNumber           interface{}   `json:"contact_number"`
		TenantID                int           `json:"tenant_id"`
		AppID                   int           `json:"app_id"`
		CreatedByType           interface{}   `json:"created_by_type"`
		CreatedByID             int           `json:"created_by_id"`
		Opener                  string        `json:"opener"`
	} `json:"client"`
	Tenant struct {
		Name           string `json:"name"`
		Slug           string `json:"slug"`
		DarkLogoURL    string `json:"dark_logo_url"`
		LightLogoURL   string `json:"light_logo_url"`
		DefaultIconURL string `json:"default_icon_url"`
		ActiveIconURL  string `json:"active_icon_url"`
		WebFaviconURL  string `json:"web_favicon_url"`
	} `json:"tenant"`
}

type SignUpResponse struct {
	ID                int    `json:"id"`
	Type              string `json:"type"`
	Locale            string `json:"locale"`
	RegistrationToken string `json:"registration_token"`
	AccessToken       string `json:"access_token"`
	UUID              string `json:"uuid"`
	User              struct {
		ID                   int    `json:"id"`
		Type                 string `json:"type"`
		AvatarURL            string `json:"avatar_url"`
		DisplayName          string `json:"display_name"`
		FirstName            string `json:"first_name"`
		LastName             string `json:"last_name"`
		Email                any    `json:"email"`
		Verified             bool   `json:"verified"`
		AcceptedLatestPolicy bool   `json:"accepted_latest_policy"`
		AcceptedLatestTerms  bool   `json:"accepted_latest_terms"`
	} `json:"user"`
	Client struct {
		ID                      int       `json:"id"`
		Name                    string    `json:"name"`
		Slug                    string    `json:"slug"`
		SecretKey               string    `json:"secret_key"`
		Platform                string    `json:"platform"`
		IssuedByID              int       `json:"issued_by_id"`
		CreatedAt               time.Time `json:"created_at"`
		UpdatedAt               time.Time `json:"updated_at"`
		AccessType              string    `json:"access_type"`
		LatestVersion           string    `json:"latest_version"`
		MinimumSupportedVersion string    `json:"minimum_supported_version"`
		Official                bool      `json:"official"`
		AccessKey               string    `json:"access_key"`
		Color                   string    `json:"color"`
		WebhookURL              string    `json:"webhook_url"`
		Refs                    []any     `json:"refs"`
		FcmServerKey            string    `json:"fcm_server_key"`
		Scheme                  string    `json:"scheme"`
		ContactEmail            any       `json:"contact_email"`
		ContactNumber           any       `json:"contact_number"`
		TenantID                int       `json:"tenant_id"`
		AppID                   any       `json:"app_id"`
		CreatedByType           any       `json:"created_by_type"`
		CreatedByID             any       `json:"created_by_id"`
		Opener                  string    `json:"opener"`
	} `json:"client"`
	Tenant struct {
		Name           string `json:"name"`
		Slug           string `json:"slug"`
		DarkLogoURL    string `json:"dark_logo_url"`
		LightLogoURL   string `json:"light_logo_url"`
		DefaultIconURL string `json:"default_icon_url"`
		ActiveIconURL  string `json:"active_icon_url"`
		WebFaviconURL  string `json:"web_favicon_url"`
	} `json:"tenant"`
}

type GetServiceListResponse []struct {
	ID               int    `json:"id"`
	Type             string `json:"type"`
	Name             string `json:"name"`
	Slug             string `json:"slug"`
	Code             string `json:"code"`
	RegionalSlug     string `json:"regional_slug"`
	IconURL          string `json:"icon_url"`
	BannerURL        string `json:"banner_url"`
	ThumbnailURL     string `json:"thumbnail_url"`
	IconWebpURL      string `json:"icon_webp_url"`
	BannerWebpURL    string `json:"banner_webp_url"`
	ThumbnailWebpURL string `json:"thumbnail_webp_url"`
	Position         int    `json:"position"`
	Packages         []struct {
		ID                    int       `json:"id"`
		ServiceID             int       `json:"service_id"`
		PriceSatangs          int       `json:"price_satangs"`
		PriceCurrency         string    `json:"price_currency"`
		UpdatedAt             time.Time `json:"updated_at"`
		CreatedAt             time.Time `json:"created_at"`
		CommissionRate        float64   `json:"commission_rate"`
		Image                 any       `json:"image"`
		DeletedAt             any       `json:"deleted_at"`
		Position              int       `json:"position"`
		Duration              int       `json:"duration"`
		Code                  string    `json:"code"`
		PayoutSatangs         int       `json:"payout_satangs"`
		PayoutCurrency        string    `json:"payout_currency"`
		ParentID              any       `json:"parent_id"`
		MinimumAmount         int       `json:"minimum_amount"`
		MaximumAmount         int       `json:"maximum_amount"`
		MinimumPriceSatangs   int       `json:"minimum_price_satangs"`
		MinimumPriceCurrency  string    `json:"minimum_price_currency"`
		ItemType              string    `json:"item_type"`
		MinimumPayoutSatangs  int       `json:"minimum_payout_satangs"`
		MinimumPayoutCurrency string    `json:"minimum_payout_currency"`
		TenantID              int       `json:"tenant_id"`
		Subscribable          bool      `json:"subscribable"`
		IsSpecial             bool      `json:"is_special"`
		Name                  string    `json:"name"`
		Description           string    `json:"description"`
		Unit                  string    `json:"unit"`
	} `json:"packages"`
	Title           string `json:"title"`
	Subscribable    bool   `json:"subscribable"`
	Description     string `json:"description"`
	Terms           string `json:"terms"`
	NameEn          string `json:"name_en"`
	NameTh          string `json:"name_th"`
	TitleEn         any    `json:"title_en"`
	TitleTh         string `json:"title_th"`
	DescriptionEn   string `json:"description_en"`
	DescriptionTh   string `json:"description_th"`
	TermsEn         string `json:"terms_en"`
	TermsTh         string `json:"terms_th"`
	CheapestPackage struct {
		ID            int    `json:"id"`
		Type          string `json:"type"`
		Name          string `json:"name"`
		Code          string `json:"code"`
		Description   string `json:"description"`
		Image         any    `json:"image"`
		MinimumAmount int    `json:"minimum_amount"`
		MaximumAmount int    `json:"maximum_amount"`
		ItemType      string `json:"item_type"`
		Unit          string `json:"unit"`
		Duration      int    `json:"duration"`
		Subscribable  bool   `json:"subscribable"`
		Price         struct {
			Fractional   int     `json:"fractional"`
			Decimal      float64 `json:"decimal"`
			DisplayValue string  `json:"display_value"`
			FullDisplay  string  `json:"full_display"`
			Currency     string  `json:"currency"`
		} `json:"price"`
		MinimumPrice struct {
			Fractional   int     `json:"fractional"`
			Decimal      float64 `json:"decimal"`
			DisplayValue string  `json:"display_value"`
			FullDisplay  string  `json:"full_display"`
			Currency     string  `json:"currency"`
		} `json:"minimum_price"`
		ServiceID             int       `json:"service_id"`
		PriceSatangs          int       `json:"price_satangs"`
		PriceCurrency         string    `json:"price_currency"`
		UpdatedAt             time.Time `json:"updated_at"`
		CreatedAt             time.Time `json:"created_at"`
		CommissionRate        float64   `json:"commission_rate"`
		DeletedAt             any       `json:"deleted_at"`
		Position              int       `json:"position"`
		PayoutSatangs         int       `json:"payout_satangs"`
		PayoutCurrency        string    `json:"payout_currency"`
		ParentID              any       `json:"parent_id"`
		MinimumPriceSatangs   int       `json:"minimum_price_satangs"`
		MinimumPriceCurrency  string    `json:"minimum_price_currency"`
		MinimumPayoutSatangs  int       `json:"minimum_payout_satangs"`
		MinimumPayoutCurrency string    `json:"minimum_payout_currency"`
		TenantID              int       `json:"tenant_id"`
		IsSpecial             bool      `json:"is_special"`
		NameEn                string    `json:"name_en"`
		NameTh                string    `json:"name_th"`
		DescriptionEn         string    `json:"description_en"`
		DescriptionTh         string    `json:"description_th"`
		UnitEn                string    `json:"unit_en"`
		UnitTh                string    `json:"unit_th"`
	} `json:"cheapest_package"`
}
