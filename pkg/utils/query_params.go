package utils

import "net/url"

func BuildURLWithQuery(baseURL string, params map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	query := u.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	u.RawQuery = query.Encode()
	return u.String(), nil
}
