package trueidpartner

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

// Client คือ struct ที่ใช้สำหรับเรียกใช้งาน API ของ TrueID
type Client struct {
	BaseURL    string
	restClient *resty.Client
}

// NewClient คือฟังก์ชันสร้าง Client สำหรับเรียกใช้งาน API ของ TrueID
func NewClient() *Client {
	return &Client{
		BaseURL:    "https://partner-internal-api.trueid-preprod.net/partner/v1",
		restClient: resty.New(),
	}
}

// GetSubscribers คือ struct ที่ใช้สำหรับเก็บข้อมูลผู้ใช้งานที่สมัครสมาชิก TrueID
func (c *Client) GetSubscribers(ssoid string) (*Subscribers, *resty.Response, error) {

	url := fmt.Sprintf("%s/partner/v1/accounts/%s/subscribers?partner_code=TRUE&status=active", c.BaseURL, ssoid)

	client := resty.New()

	// ตั้งค่า Headers ที่จะใช้ทั่วโปรแกรม

	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Authorization", os.Getenv("TRUEID_PARTNER_API_KEY"))

	var subscribers Subscribers
	// ทำ POST request โดยส่ง Headers และ Body
	resp, err := client.R().
		SetResult(&subscribers).
		Get(url)

	if err != nil {
		return nil, resp, err
	}
	if resp.IsError() {
		return nil, resp, fmt.Errorf("API request failed: %s", resp.Status())
	}

	return &subscribers, resp, nil
}
