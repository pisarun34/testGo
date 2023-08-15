package seekster

import (
	"TESTGO/pkg/database/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRestyClient to mock RestyClienter methods for testing
type MockRestyClient struct {
	mock.Mock
}

func (m *MockRestyClient) R() Requester {
	call := m.Called()
	return call.Get(0).(Requester)
}

type MockRestyRequest struct {
	mock.Mock
}

func (m *MockRestyRequest) SetHeader(header, value string) Requester {
	return m.Called(header, value).Get(0).(Requester)
}

func (m *MockRestyRequest) SetResult(result interface{}) Requester {
	return m.Called(result).Get(0).(Requester)
}

func (m *MockRestyRequest) SetBody(body interface{}) Requester {
	return m.Called(body).Get(0).(Requester)
}

func (m *MockRestyRequest) Post(url string) (*resty.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*resty.Response), args.Error(1)
}

func ResettableByteReader(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}

func TestSignInByPhone(t *testing.T) {
	// Mock resty.Request behavior

	signInResponse := &SignResponse{
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
	ErrorJson := []byte(`{"error": "error"}`)
	errorReader := ResettableByteReader(ErrorJson)

	mockRequest := new(MockRestyRequest)

	tests := []struct {
		name            string
		mockFunc        func(*MockRestyRequest)
		input           models.User
		expectErr       bool
		expectRespErr   bool
		expectResp      *SignResponse
		expectErrorResp *resty.Response
		expectError     error
	}{
		{
			name: "successful sign in",
			mockFunc: func(m *MockRestyRequest) {
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
				//signInData := &SignResponse{ID: 1} // Mock your response data
				m.On("SetHeader", "Content-Type", mock.Anything).Return(mockRequest)
				m.On("SetHeader", "Authorization", mock.Anything).Return(mockRequest)
				m.On("SetResult", mock.Anything).Return(mockRequest)
				m.On("SetBody", mock.Anything).Return(mockRequest)
				m.On("Post", mock.Anything).Return(resp, nil)
			},
			input: models.User{
				SeeksterUser: models.SeeksterUser{
					PhoneNumber: "1234567890",
					Password:    "password",
					UUID:        "testUUID",
				},
			},
			expectErr:  false,
			expectResp: signInResponse,
		},
		{
			name: "request error",
			mockFunc: func(m *MockRestyRequest) {
				mockError := errors.New("mock error")
				resp := &resty.Response{}
				m.On("SetHeader", "Content-Type", mock.Anything).Return(mockRequest)
				m.On("SetHeader", "Authorization", mock.Anything).Return(mockRequest)
				m.On("SetResult", mock.Anything).Return(mockRequest)
				m.On("SetBody", mock.Anything).Return(mockRequest)
				m.On("Post", mock.Anything).Return(resp, mockError)
			},
			input: models.User{
				SeeksterUser: models.SeeksterUser{
					PhoneNumber: "1234567890",
					Password:    "password",
					UUID:        "testUUID",
				},
			},
			expectErr:       true,
			expectRespErr:   true,
			expectResp:      nil,
			expectErrorResp: nil,
			expectError:     errors.New("mock error"),
		},
		{
			name: "failed sign in",
			mockFunc: func(m *MockRestyRequest) {
				resp := &resty.Response{
					RawResponse: &http.Response{
						StatusCode: 403,
						Body:       io.NopCloser(errorReader), // mock the response body here
						Header: http.Header{
							"Content-Type": []string{"application/json; charset=utf-8"},
						},
					},
				}
				//signInData := &SignResponse{ID: 1} // Mock your response data
				m.On("SetHeader", "Content-Type", mock.Anything).Return(mockRequest)
				m.On("SetHeader", "Authorization", mock.Anything).Return(mockRequest)
				m.On("SetResult", mock.Anything).Return(mockRequest)
				m.On("SetBody", mock.Anything).Return(mockRequest)
				m.On("Post", mock.Anything).Return(resp, nil)
			},
			input: models.User{
				SeeksterUser: models.SeeksterUser{
					PhoneNumber: "1234567890",
					Password:    "password",
					UUID:        "testUUID",
				},
			},
			expectErr:     true,
			expectRespErr: false,
			expectResp:    &SignResponse{},
			expectErrorResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: 403,
					Body:       io.NopCloser(errorReader), // mock the response body here
					Header: http.Header{
						"Content-Type": []string{"application/json; charset=utf-8"},
					},
				},
			},
		},
	}

	// Mock RestyClienter behavior

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRequest = new(MockRestyRequest)
			mockRestyClient := new(MockRestyClient)
			mockRestyClient.On("R").Return(mockRequest)

			client := &Client{
				BaseURL:    "http://mockbaseurl.com",
				restClient: mockRestyClient,
			}
			tt.mockFunc(mockRequest)
			var resp *SignResponse
			_, resps, err := client.SignInByPhone(tt.input)
			reader.Seek(0, io.SeekStart)
			if resps != nil {
				bodyBytes, err3 := io.ReadAll(resps.RawResponse.Body)
				fmt.Println("bodyBytes", string(bodyBytes))
				err2 := json.Unmarshal(bodyBytes, &resp)
				if err3 != nil {
					fmt.Println("ReadAll error:", err3)
				}
				if err2 != nil {
					fmt.Println("Unmarshal error:", err2)
				}
			}
			if tt.expectErr {
				assert.Equal(t, tt.expectErrorResp, resps)
				assert.Equal(t, tt.expectResp, resp)
				if tt.expectRespErr {
					assert.Error(t, tt.expectError, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectResp, resp)
			}

			// Assert expectations
			mockRestyClient.AssertExpectations(t)
			mockRequest.AssertExpectations(t)
		})
	}

}

// TODO: Add more tests as needed
