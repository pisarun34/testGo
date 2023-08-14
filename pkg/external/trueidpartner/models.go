package trueidpartner

import (
	"time"
)

// Subscribers สร้าง struct สำหรับเก็บข้อมูล subscription ของ TrueID
type Subscribers struct {
	Code            int    `json:"code"`
	Message         string `json:"message"`
	PlatformModule  int    `json:"platform_module"`
	ReportDashboard int    `json:"report_dashboard"`
	Data            struct {
		ID                  string `json:"id"`
		AccountID           string `json:"account_id"`
		PartnerCode         string `json:"partner_code"`
		PartnerSubscriberID string `json:"partner_subscriber_id"`
		Status              string `json:"status"`
		Products            []any  `json:"products"`
		Services            []struct {
			ID             string    `json:"id"`
			Identifier     string    `json:"identifier"`
			IdentifierType string    `json:"identifier_type"`
			DisplayName    string    `json:"display_name"`
			Type           string    `json:"type"`
			SubType        string    `json:"sub_type"`
			Financing      string    `json:"financing"`
			Status         string    `json:"status"`
			ModifyDate     time.Time `json:"modify_date"`
			CreateDate     time.Time `json:"create_date"`
		} `json:"services"`
		LastRefresh time.Time `json:"last_refresh"`
		CreateDate  time.Time `json:"create_date"`
	} `json:"data"`
}
