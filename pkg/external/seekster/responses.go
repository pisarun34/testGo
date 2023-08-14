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
