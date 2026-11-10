package alphavantage

import (
	"net/http"
	"testing"
	"time"
)

func TestGetMarketStatusNetworkError(t *testing.T) {
	client := NewClient("test-api-key")

	client.baseURL = "http://127.0.0.1:1"
	client.httpClient = &http.Client{
		Timeout: time.Second,
	}

	result, err := client.GetMarketStatus()

	if err == nil {
		t.Fatal("expected a network error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
}
