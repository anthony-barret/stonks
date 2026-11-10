package alphavantage

import (
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
