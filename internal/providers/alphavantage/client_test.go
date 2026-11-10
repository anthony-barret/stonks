package alphavantage

import (
	"strings"
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

func TestCheckAPIError(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantErr     bool
		wantMessage string
	}{
		{
			name:    "valid JSON without API error",
			body:    `{"bestMatches":[]}`,
			wantErr: false,
		},
		{
			name:        "rate limit error",
			body:        `{"Note":"API rate limit reached"}`,
			wantErr:     true,
			wantMessage: "API rate limit reached",
		},
		{
			name:        "API error message",
			body:        `{"Error Message":"Invalid API call"}`,
			wantErr:     true,
			wantMessage: "Invalid API call",
		},
		{
			name:        "information message",
			body:        `{"Information":"Demo API key restriction"}`,
			wantErr:     true,
			wantMessage: "Demo API key restriction",
		},
		{
			name:        "invalid JSON",
			body:        `{"Note":"invalid JSON"`,
			wantErr:     true,
			wantMessage: "invalid API response:",
		},
		{
			name:    "empty JSON object",
			body:    `{}`,
			wantErr: false,
		},
		{
			name:    "empty body",
			body:    ``,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkAPIError([]byte(tt.body))

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}

				if !strings.Contains(err.Error(), tt.wantMessage) {
					t.Errorf(
						"expected error containing %q, got %q",
						tt.wantMessage,
						err.Error(),
					)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
