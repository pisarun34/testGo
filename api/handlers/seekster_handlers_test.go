package handlers

import (
	"TESTGO/pkg/database"
	"TESTGO/pkg/database/mysql"
	"TESTGO/pkg/external/seekster"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSeeksterAPI struct {
	mock.Mock
}

type MockRedisClient struct {
	data map[string]string
}

func (m *MockSeeksterAPI) SignInByPhone(user mysql.User) (*seekster.SignResponse, *resty.Response, error) {
	args := m.Called(user)
	signResponse, _ := args.Get(0).(*seekster.SignResponse)
	restyResponse, _ := args.Get(1).(*resty.Response)

	return signResponse, restyResponse, args.Error(2)
}

func (m *MockSeeksterAPI) SignUp(user mysql.User) (*seekster.SignUpResponse, *resty.Response, error) {
	args := m.Called(user)
	return args.Get(0).(*seekster.SignUpResponse), args.Get(1).(*resty.Response), args.Error(2)
}

func NewMockRedisClient() database.RedisClientInterface {
	return &MockRedisClient{data: make(map[string]string)}
}

func (m *MockRedisClient) GetSeeksterToken(ctx context.Context, ssoid string) (string, error) {
	if value, exists := m.data[ssoid]; exists {
		return value, nil
	}
	return "", fmt.Errorf("redis: nil")
}

func (m *MockRedisClient) SetSeeksterToken(ctx context.Context, ssoid string, seeksterToken string) error {
	m.data[ssoid] = seeksterToken
	return nil
}

func (m *MockRedisClient) CloseRedis() {
	// ไม่ต้องทำอะไร
}

func ResettableByteReader(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}

func TestSeeksterSignIn_Success_RedisHaveToken(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/seekster/signin", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("ssoid", "22030729")
	mockClient := new(MockSeeksterAPI)

	// ตัวอย่างการ mock error จาก SignInByPhone
	mockClient.On("SignInByPhone", mock.AnythingOfType("mysql.User")).Return(nil, nil, errors.New("Mock error"))
	mockRedis := NewMockRedisClient()
	mockRedis.SetSeeksterToken(context.Background(), "22030729", "mock_token")
	// 2. Act
	SeeksterSignin(mockClient, c, mockRedis, nil) // ใช้ nil สำหรับ db ในตัวอย่าง

	// 3. Assert
	// ใช้ testify's assert เพื่อตรวจสอบ
	expectedStatusCode := http.StatusOK // กำหนดค่า expectedStatusCode ที่คาดหวัง
	assert.Equal(t, expectedStatusCode, w.Code)
	assert.Equal(t, gin.H{"code": 10001, "message": "Success"}, w.Body.String())
	// คุณสามารถเพิ่ม assertions ได้ตามที่คุณต้องการ

	//mockClient.AssertExpectations(t) // ตรวจสอบว่ามีการเรียกฟังก์ชันที่ถูก mock ตามที่คาดหวัง
}

func TestSeeksterSignIn_Success_RedisNoToken(t *testing.T) {
	godotenv.Load("../../.env")
	// 1. Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/seekster/signin", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("ssoid", "22030729")
	mockClient := new(MockSeeksterAPI)
	db := mysql.Initialize()

	signInResponse := &seekster.SignResponse{
		ID:                1,
		Type:              "devices",
		Locale:            "",
		RegistrationToken: "",
		AccessToken:       "testAccessToken",
		UUID:              "testUUID",
		User: struct {
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
		}{
			ID:                   1,
			Type:                 "customers",
			AvatarURL:            "https://seekster-company.oss-ap-southeast-1.aliyuncs.com/workforce-images/default_avatar.png",
			DisplayName:          "test",
			FirstName:            nil,
			LastName:             nil,
			Email:                nil,
			Verified:             true,
			AcceptedLatestPolicy: false,
			AcceptedLatestTerms:  true,
		},
		Client: struct {
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
		}{
			ID:                      1,
			Name:                    "testClientName",
			Slug:                    "testClientSlug",
			SecretKey:               "testSecretKey",
			Platform:                "testPlatform",
			IssuedByID:              1,
			CreatedAt:               time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:               time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			AccessType:              "testAccessType",
			LatestVersion:           "testLatestVersion",
			MinimumSupportedVersion: "testMinimumSupportedVersion",
			Official:                true,
			AccessKey:               "testAccessKey",
			Color:                   "testColor",
			WebhookURL:              "testWebhookURL",
			Refs:                    []interface{}{"testRefs"},
			FcmServerKey:            "testFcmServerKey",
			Scheme:                  "testScheme",
			ContactEmail:            "testContactEmail",
			ContactNumber:           "testContactNumber",
			TenantID:                1,
			AppID:                   1,
			CreatedByType:           "testCreatedByType",
			CreatedByID:             1,
			Opener:                  "testOpener",
		},
		Tenant: struct {
			Name           string `json:"name"`
			Slug           string `json:"slug"`
			DarkLogoURL    string `json:"dark_logo_url"`
			LightLogoURL   string `json:"light_logo_url"`
			DefaultIconURL string `json:"default_icon_url"`
			ActiveIconURL  string `json:"active_icon_url"`
			WebFaviconURL  string `json:"web_favicon_url"`
		}{
			Name:           "testTenantName",
			Slug:           "testTenantSlug",
			DarkLogoURL:    "testDarkLogoURL",
			LightLogoURL:   "testLightLogoURL",
			DefaultIconURL: "testDefaultIconURL",
			ActiveIconURL:  "testActiveIconURL",
			WebFaviconURL:  "testWebFaviconURL",
		},
	}

	// Convert the signInResponse to a JSON byte slice
	responseBody, err := json.Marshal(signInResponse)
	if err != nil {
		t.Fatal(err)
	}

	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
			Header: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
		},
	}
	//bodyBytes, _ := io.ReadAll(resp.RawResponse.Body)

	//errMsg := fmt.Sprintf("API request failed: %s - %s", resp.Status(), string(bodyBytes))

	// ตัวอย่างการ mock error จาก SignInByPhone errors.New("Mock error")
	mockClient.On("SignInByPhone", mock.AnythingOfType("mysql.User")).Return(signInResponse, resp, nil)
	mockRedis := NewMockRedisClient()
	// 2. Act
	SeeksterSignin(mockClient, c, mockRedis, db) // ใช้ nil สำหรับ db ในตัวอย่าง

	// 3. Assert
	// ใช้ testify's assert เพื่อตรวจสอบ
	expectedStatusCode := http.StatusOK // กำหนดค่า expectedStatusCode ที่คาดหวัง
	assert.Equal(t, expectedStatusCode, w.Code)
	assert.Equal(t, string("{\"code\":10001,\"message\":\"Success\"}"), w.Body.String())
	// คุณสามารถเพิ่ม assertions ได้ตามที่คุณต้องการ

	mockClient.AssertExpectations(t) // ตรวจสอบว่ามีการเรียกฟังก์ชันที่ถูก mock ตามที่คาดหวัง
}

func TestSeeksterSignIn_DbNotFoundUser(t *testing.T) {
	godotenv.Load("../../.env")
	// 1. Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/seekster/signin", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("ssoid", "22030730")
	mockClient := new(MockSeeksterAPI)
	db := mysql.Initialize()

	signInResponse := &seekster.SignResponse{
		ID:                1,
		Type:              "devices",
		Locale:            "",
		RegistrationToken: "",
		AccessToken:       "testAccessToken",
		UUID:              "testUUID",
		User: struct {
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
		}{
			ID:                   1,
			Type:                 "customers",
			AvatarURL:            "https://seekster-company.oss-ap-southeast-1.aliyuncs.com/workforce-images/default_avatar.png",
			DisplayName:          "test",
			FirstName:            nil,
			LastName:             nil,
			Email:                nil,
			Verified:             true,
			AcceptedLatestPolicy: false,
			AcceptedLatestTerms:  true,
		},
		Client: struct {
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
		}{
			ID:                      1,
			Name:                    "testClientName",
			Slug:                    "testClientSlug",
			SecretKey:               "testSecretKey",
			Platform:                "testPlatform",
			IssuedByID:              1,
			CreatedAt:               time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:               time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			AccessType:              "testAccessType",
			LatestVersion:           "testLatestVersion",
			MinimumSupportedVersion: "testMinimumSupportedVersion",
			Official:                true,
			AccessKey:               "testAccessKey",
			Color:                   "testColor",
			WebhookURL:              "testWebhookURL",
			Refs:                    []interface{}{"testRefs"},
			FcmServerKey:            "testFcmServerKey",
			Scheme:                  "testScheme",
			ContactEmail:            "testContactEmail",
			ContactNumber:           "testContactNumber",
			TenantID:                1,
			AppID:                   1,
			CreatedByType:           "testCreatedByType",
			CreatedByID:             1,
			Opener:                  "testOpener",
		},
		Tenant: struct {
			Name           string `json:"name"`
			Slug           string `json:"slug"`
			DarkLogoURL    string `json:"dark_logo_url"`
			LightLogoURL   string `json:"light_logo_url"`
			DefaultIconURL string `json:"default_icon_url"`
			ActiveIconURL  string `json:"active_icon_url"`
			WebFaviconURL  string `json:"web_favicon_url"`
		}{
			Name:           "testTenantName",
			Slug:           "testTenantSlug",
			DarkLogoURL:    "testDarkLogoURL",
			LightLogoURL:   "testLightLogoURL",
			DefaultIconURL: "testDefaultIconURL",
			ActiveIconURL:  "testActiveIconURL",
			WebFaviconURL:  "testWebFaviconURL",
		},
	}

	// Convert the signInResponse to a JSON byte slice
	responseBody, err := json.Marshal(signInResponse)
	if err != nil {
		t.Fatal(err)
	}

	jsonContent := []byte(responseBody)
	reader := ResettableByteReader(jsonContent)

	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Body:       io.NopCloser(reader), // mock the response body here
			Header: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
		},
	}

	// ตัวอย่างการ mock error จาก SignInByPhone errors.New("Mock error")
	mockClient.On("SignInByPhone", mock.AnythingOfType("mysql.User")).Return(signInResponse, resp, nil)
	mockRedis := NewMockRedisClient()

	// 2. Act
	SeeksterSignin(mockClient, c, mockRedis, db) // ใช้ nil สำหรับ db ในตัวอย่าง

	// 3. Assert
	// ใช้ testify's assert เพื่อตรวจสอบ
	expectedStatusCode := http.StatusInternalServerError // กำหนดค่า expectedStatusCode ที่คาดหวัง
	assert.Equal(t, expectedStatusCode, w.Code)
	assert.Equal(t, string("{\"code\":10888,\"message\":\"Internal server error\"}"), w.Body.String())
	// คุณสามารถเพิ่ม assertions ได้ตามที่คุณต้องการ

	//mockClient.AssertExpectations(t) // ตรวจสอบว่ามีการเรียกฟังก์ชันที่ถูก mock ตามที่คาดหวัง
}

func TestSeeksterSignIn_Response_failed(t *testing.T) {
	godotenv.Load("../../.env")
	// 1. Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/seekster/signin", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("ssoid", "22030729")
	mockClient := new(MockSeeksterAPI)
	db := mysql.Initialize()

	type ErrorResponse struct {
		Error   string      `json:"error"`
		Message interface{} `json:"message"`
		Details interface{} `json:"details"`
	}

	signInResponse := ErrorResponse{
		Error:   "invalid_credentials",
		Message: nil,
		Details: nil,
	}

	// Convert the signInResponse to a JSON byte slice
	responseBody, err := json.Marshal(signInResponse)
	if err != nil {
		t.Fatal(err)
	}

	reader := ResettableByteReader(responseBody)
	reader.Seek(0, io.SeekStart)
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: 400,
			Body:       io.NopCloser(reader), // mock the response body here
			Header: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
		},
	}

	errMsg := fmt.Sprintf("API request failed: %s - %s", resp.Status(), resp.String())
	// ตัวอย่างการ mock error จาก SignInByPhone errors.New("Mock error")
	mockClient.On("SignInByPhone", mock.AnythingOfType("mysql.User")).Return(nil, resp, errors.New(errMsg))
	mockRedis := NewMockRedisClient()

	// 2. Act
	SeeksterSignin(mockClient, c, mockRedis, db) // ใช้ nil สำหรับ db ในตัวอย่าง

	// 3. Assert
	// ใช้ testify's assert เพื่อตรวจสอบ
	expectedStatusCode := http.StatusBadRequest // กำหนดค่า expectedStatusCode ที่คาดหวัง
	assert.Equal(t, expectedStatusCode, w.Code)
	assert.Equal(t, string("{\"details\":null,\"error\":\"invalid_credentials\",\"message\":null}"), w.Body.String())
	// คุณสามารถเพิ่ม assertions ได้ตามที่คุณต้องการ

	//mockClient.AssertExpectations(t) // ตรวจสอบว่ามีการเรียกฟังก์ชันที่ถูก mock ตามที่คาดหวัง
}
