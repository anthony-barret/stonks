package alphavantage

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key")

	if client.apiKey != "test-api-key" {
		t.Errorf("expected API key %q, got %q", "test-api-key", client.apiKey)
	}

	if client.baseURL != "https://www.alphavantage.co/query" {
		t.Errorf("unexpected base URL: %s", client.baseURL)
	}

	if client.httpClient == nil {
		t.Fatal("expected HTTP client, got nil")
	}

	if client.httpClient.Timeout == 0 {
		t.Error("expected HTTP client timeout to be set")
	}
}
