package alphavantage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://www.alphavantage.co/query",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type APIError struct {
	ErrorMessage string `json:"Error Message"`
	Information  string `json:"Information"`
	Note         string `json:"Note"`
}

func checkAPIError(body []byte) error {
	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return fmt.Errorf("invalid API response: %w", err)
	}

	switch {
	case apiErr.Note != "":
		return fmt.Errorf("%s", apiErr.Note)
	case apiErr.ErrorMessage != "":
		return fmt.Errorf("%s", apiErr.ErrorMessage)
	case apiErr.Information != "":
		return fmt.Errorf("%s", apiErr.Information)
	}

	return nil
}
