package exampleapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *Client) GetUser(id string) (User, error) {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/users/%s", c.BaseURL, id))
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return User{}, err
		}
		return User{}, fmt.Errorf("API error: %s (code %d)", errResp.Message, errResp.Code)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return User{}, err
	}

	return user, nil
}
