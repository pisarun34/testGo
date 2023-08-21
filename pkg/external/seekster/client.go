package seekster

import (
	"TESTGO/middlewares"
	"TESTGO/pkg/database/models"
	"TESTGO/pkg/utils"
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type Requester interface {
	SetHeader(header, value string) Requester
	SetResult(result interface{}) Requester
	SetBody(body interface{}) Requester
	Post(url string) (*resty.Response, error)
	//SetDoNotParseResponse(bool) Requester
}

func (rrw *RestyRequestWrapper) SetDoNotParseResponse(val bool) Requester {
	rrw.request.SetDoNotParseResponse(val)
	return rrw
}

// RestyClienter is an interface that wraps methods we use from resty.Client.
type RestyClienter interface {
	R() Requester
}

// RestyClientWrapper wraps resty.Client to satisfy the RestyClienter interface.
type RestyClientWrapper struct {
	client *resty.Client
}

func (rrw *RestyRequestWrapper) SetHeader(header, value string) Requester {
	rrw.request.SetHeader(header, value)
	return rrw
}

func (rrw *RestyRequestWrapper) SetResult(result interface{}) Requester {
	rrw.request.SetResult(result)
	return rrw
}

func (rrw *RestyRequestWrapper) SetBody(body interface{}) Requester {
	rrw.request.SetBody(body)
	return rrw
}

func (rrw *RestyRequestWrapper) Post(url string) (*resty.Response, error) {
	return rrw.request.Post(url)
}

type RestyRequestWrapper struct {
	request *resty.Request
}

func (rcw *RestyClientWrapper) R() Requester {
	return &RestyRequestWrapper{request: rcw.client.R()}
}

// Client คือ struct ที่ใช้สำหรับเรียกใช้งาน API ของ Seekster
type Client struct {
	BaseURL string
	//restClient *resty.Client
	restClient RestyClienter
}

// NewClient คือฟังก์ชันสร้าง Client สำหรับเรียกใช้งาน API ของ Seekster
func NewClient() *Client {
	return &Client{
		BaseURL: os.Getenv("SEEKSTER_API_BASE_URL"),
		//restClient: resty.New(),
		restClient: &RestyClientWrapper{client: resty.New()},
	}
}

// SignInByPhone คือฟังก์ชันสำหรับเรียกใช้งาน SignInByPhone API ของ Seekster
func (c *Client) SignInByPhone(seeksterUser models.User) (*SignResponse, *resty.Response, error) {

	url := fmt.Sprintf("%s/sign_in_with_phone_number_password", c.BaseURL)
	data := SignInRequest{
		DevicesAttributes: struct {
			Brand             string `json:"brand"`
			Carrier           string `json:"carrier"`
			ClientType        string `json:"client_type"`
			ClientVersion     string `json:"client_version"`
			Locale            string `json:"locale"`
			Model             string `json:"model"`
			OsVersion         string `json:"os_version"`
			UUID              string `json:"uuid"`
			RegistrationToken string `json:"registration_token"`
		}{UUID: seeksterUser.SeeksterUser.UUID},
		PhoneNumber: seeksterUser.SeeksterUser.PhoneNumber,
		Password:    seeksterUser.SeeksterUser.Password,
	}
	var signInResponse SignResponse
	// ทำ POST request โดยส่ง Headers และ Body
	resp, err := c.restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", os.Getenv("SEEKSTER_API_KEY")).
		SetResult(&signInResponse).
		SetBody(data).
		Post(url)
	//fmt.Println("resp", resp)
	if err != nil {
		return nil, nil, err
	}
	if resp.IsError() {
		errMsg := fmt.Sprintf("API request failed: %s - %s", resp.Status(), resp.String())
		return nil, resp, errors.New(errMsg)
	}

	return &signInResponse, resp, nil
}

func (c *Client) SignUp(seeksterUser models.User) (*SignUpResponse, *resty.Response, error) {

	url := fmt.Sprintf("%s/register", c.BaseURL)

	client := resty.New()

	// ตั้งค่า Headers ที่จะใช้ทั่วโปรแกรม

	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Authorization", os.Getenv("SEEKSTER_API_KEY"))

	data := SignUpRequest{
		DevicesAttributes: struct {
			Brand             string `json:"brand"`
			Carrier           string `json:"carrier"`
			OwnerType         string `json:"owner_type"`
			ClientVersion     string `json:"client_version"`
			Locale            string `json:"locale"`
			Model             string `json:"model"`
			OsVersion         string `json:"os_version"`
			RegistrationToken string `json:"registration_token"`
			UUID              string `json:"uuid"`
		}{UUID: seeksterUser.SeeksterUser.UUID},
		PhoneNumber: seeksterUser.SeeksterUser.PhoneNumber,
		Password:    seeksterUser.SeeksterUser.Password,
	}

	var user SignUpResponse
	// ทำ POST request โดยส่ง Headers และ Body
	resp, err := client.R().
		SetResult(&user).
		SetBody(data).
		Post(url)

	if err != nil {
		return nil, nil, err
	}
	if resp.IsError() {
		errMsg := fmt.Sprintf("API request failed: %s - %s", resp.Status(), resp.String())
		fmt.Println(errMsg)
		return nil, resp, errors.New(errMsg)
	}

	return &user, nil, nil
}

func (c *Client) GetServiceList(g *gin.Context) (*GetServiceListResponse, *resty.Response, error) {
	token, err := middlewares.CheckAndExtractSeeksterToken(g)
	if err != nil {
		return nil, nil, err
	}
	queryParameters := make(map[string]string)
	for key, value := range g.Request.URL.Query() {
		queryParameters[key] = value[0]
	}
	baseURL, err := utils.BuildURLWithQuery(fmt.Sprintf("%s/services", c.BaseURL), queryParameters)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("baseURL", baseURL)
	client := resty.New()

	// ตั้งค่า Headers ที่จะใช้ทั่วโปรแกรม

	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Authorization", "Bearer "+token)

	var service GetServiceListResponse
	// ทำ POST request โดยส่ง Headers และ Body
	resp, err := client.R().
		SetResult(&service).
		Get(baseURL)

	if err != nil {
		return nil, nil, err
	}
	if resp.IsError() {
		errMsg := fmt.Sprintf("API request failed: %s - %s", resp.Status(), resp.String())
		return nil, resp, errors.New(errMsg)
	}

	return &service, resp, nil
}

func (c *Client) GetServiceDetails(g *gin.Context) (*GetServiceDetailsResponse, *resty.Response, error) {
	token, err := middlewares.CheckAndExtractSeeksterToken(g)
	if err != nil {
		return nil, nil, err
	}
	id := g.Param("id")

	queryParameters := make(map[string]string)
	for key, value := range g.Request.URL.Query() {
		queryParameters[key] = value[0]
	}
	baseURL, err := utils.BuildURLWithQuery(fmt.Sprintf("%s/services/%s", c.BaseURL, id), queryParameters)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("baseURL", baseURL)
	client := resty.New()

	// ตั้งค่า Headers ที่จะใช้ทั่วโปรแกรม

	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Authorization", "Bearer "+token)

	var service GetServiceDetailsResponse
	// ทำ POST request โดยส่ง Headers และ Body
	resp, err := client.R().
		SetResult(&service).
		Get(baseURL)

	if err != nil {
		return nil, nil, err
	}
	if resp.IsError() {
		errMsg := fmt.Sprintf("API request failed: %s - %s", resp.Status(), resp.String())
		return nil, resp, errors.New(errMsg)
	}

	return &service, resp, nil
}

func (c *Client) GetSlotsQuantity(g *gin.Context) (*GetSlotsQuantityResponse, *resty.Response, error) {
	token, err := middlewares.CheckAndExtractSeeksterToken(g)
	if err != nil {
		return nil, nil, err
	}
	id := g.Param("id")

	queryParameters := make(map[string]string)
	for key, value := range g.Request.URL.Query() {
		queryParameters[key] = value[0]
	}
	baseURL, err := utils.BuildURLWithQuery(fmt.Sprintf("%s/services/%s/slots", c.BaseURL, id), queryParameters)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("baseURL", baseURL)
	client := resty.New()

	// ตั้งค่า Headers ที่จะใช้ทั่วโปรแกรม

	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Authorization", "Bearer "+token)

	var slots GetSlotsQuantityResponse
	// ทำ POST request โดยส่ง Headers และ Body
	resp, err := client.R().
		SetResult(&slots).
		Get(baseURL)

	if err != nil {
		return nil, nil, err
	}
	if resp.IsError() {
		errMsg := fmt.Sprintf("API request failed: %s - %s", resp.Status(), resp.String())
		return nil, resp, errors.New(errMsg)
	}

	return &slots, resp, nil
}

func (c *Client) BookingServiceBySlot(g *gin.Context, RequestBody *BookingServiceBySlotRequest) (*BookingServiceBySlotResponse, *resty.Response, error) {
	token, err := middlewares.CheckAndExtractSeeksterToken(g)
	if err != nil {
		return nil, nil, err
	}
	id := g.Param("id")
	// Get query parameters
	queryParameters := make(map[string]string)
	for key, value := range g.Request.URL.Query() {
		queryParameters[key] = value[0]
	}
	// Create url with query parameters
	baseURL, err := utils.BuildURLWithQuery(fmt.Sprintf("%s/services/%s/slots", c.BaseURL, id), queryParameters)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("baseURL", baseURL)
	client := resty.New()

	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Authorization", "Bearer "+token)

	var service BookingServiceBySlotResponse
	resp, err := client.R().
		SetResult(&service).
		Get(baseURL)

	if err != nil {
		return nil, nil, err
	}
	if resp.IsError() {
		errMsg := fmt.Sprintf("API request failed: %s - %s", resp.Status(), resp.String())
		return nil, resp, errors.New(errMsg)
	}

	return &service, resp, nil
}

func (c *Client) GetInquiryList(g *gin.Context) (*GetInquiryListResponse, *resty.Response, error) {
	token, err := middlewares.CheckAndExtractSeeksterToken(g)
	if err != nil {
		return nil, nil, err
	}
	// Get query parameters
	queryParameters := make(map[string]string)
	for key, value := range g.Request.URL.Query() {
		queryParameters[key] = value[0]
	}
	// Create url with query parameters
	baseURL, err := utils.BuildURLWithQuery(fmt.Sprintf("%s/inquiries", c.BaseURL), queryParameters)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("baseURL", baseURL)
	client := resty.New()

	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Authorization", "Bearer "+token)

	var inquiries GetInquiryListResponse
	resp, err := client.R().
		SetResult(&inquiries).
		Get(baseURL)

	if err != nil {
		return nil, nil, err
	}
	if resp.IsError() {
		errMsg := fmt.Sprintf("API request failed: %s - %s", resp.Status(), resp.String())
		return nil, resp, errors.New(errMsg)
	}

	return &inquiries, resp, nil
}
