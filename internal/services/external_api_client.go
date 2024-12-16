// external_api_client.go
package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// ExternalAPIClient interacts with external APIs.
type ExternalAPIClient struct {
	BaseURL string
	APIKey  string
}

// NewExternalAPIClient initializes a new ExternalAPIClient.
func NewExternalAPIClient(baseURL, apiKey string) *ExternalAPIClient {
	return &ExternalAPIClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

// Get makes a GET request to an external API.
func (c *ExternalAPIClient) Get(endpoint string) (string, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}