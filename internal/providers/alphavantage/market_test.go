package alphavantage

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestGetMarketStatus(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		wantErr     bool
		wantMarkets int
	}{
		{
			name: "successful response",
			response: `{
				"endpoint": "Global Market Open & Close Status",
				"markets": [
					{
						"market_type": "Equity",
            			"region": "United States",
            			"primary_exchanges": "NASDAQ, NYSE, AMEX, BATS",
            			"local_open": "09:30",
            			"local_close": "16:15",
            			"current_status": "open",
            			"notes": ""
        			},
					{
						"market_type": "Cryptocurrency",
						"region": "Global",
						"primary_exchanges": "Global",
						"local_open": "00:00",
						"local_close": "23:59",
						"current_status": "open",
						"notes": "The cryptocurrency market is open 24 hours a day"
					}
				]
			}`,
			wantErr:     false,
			wantMarkets: 2,
		},
		{
			name: "empty markets",
			response: `{
				"endpoint": "Global Market Open & Close Status",
				"markets": []
			}`,
			wantErr:     false,
			wantMarkets: 0,
		},
		{
			name:     "invalid JSON",
			response: `{"endpoint": "Global Market Open & Close Status"`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Query().Get("function") != "MARKET_STATUS" {
						t.Errorf("unexpected function parameter: %s", r.URL.Query().Get("function"))
					}
					if r.URL.Query().Get("apikey") != "test-api-key" {
						t.Errorf("unexpected API key: %s", r.URL.Query().Get("apikey"))
					}
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, tt.response)
				},
			))
			defer server.Close()

			client := NewClient("test-api-key")
			client.baseURL = server.URL
			client.httpClient = server.Client()

			result, err := client.GetMarketStatus()

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				if result != nil {
					t.Errorf("expected nil result, got %+v", result)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("expected a result, got nil")
			}

			if len(result.Markets) != tt.wantMarkets {
				t.Errorf("expected %d markets, got %d", tt.wantMarkets, len(result.Markets))
			}
		})
	}
}
