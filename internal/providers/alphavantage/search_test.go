package alphavantage

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTyckerSearchNetworkError(t *testing.T) {
	client := NewClient("test-api-key")

	client.baseURL = "http://127.0.0.1:1"
	client.httpClient = &http.Client{
		Timeout: time.Second,
	}

	result, err := client.TyckerSearch("test")

	if err == nil {
		t.Fatal("expected a network error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
}

func TestTyckerSearch(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		symbol      string
		wantErr     bool
		wantMatches int
	}{
		{
			name:   "successful response",
			symbol: "tesco",
			response: `{
				"bestMatches": [
					{
						"1. symbol": "TSCO.LON",
						"2. name": "Tesco PLC",
						"3. type": "Equity",
						"4. region": "United Kingdom",
						"5. marketOpen": "08:00",
						"6. marketClose": "16:30",
						"7. timezone": "UTC+01",
						"8. currency": "GBX",
						"9. matchScore": "0.7273"
					},
					{
						"1. symbol": "TSCDF",
						"2. name": "Tesco plc",
						"3. type": "Equity",
						"4. region": "United States",
						"5. marketOpen": "09:30",
						"6. marketClose": "16:00",
						"7. timezone": "UTC-04",
						"8. currency": "USD",
						"9. matchScore": "0.7143"
					}
				]
			}`,
			wantErr:     false,
			wantMatches: 2,
		},
		{
			name:   "empty matches",
			symbol: "tescoaaaaaaaaaaaaaaaaaaaaaaaa",
			response: `{
				"bestMatches": []
			}`,
			wantErr:     false,
			wantMatches: 0,
		},
		{
			name:   "empty tycker",
			symbol: "",
			response: `{
				"Error Message": "Invalid API call. Please retry or visit the documentation (https://www.alphavantage.co/documentation/) for SYMBOL_SEARCH."
			}`,
			wantErr: true,
		},
		{
			name:     "invalid JSON",
			symbol:   "tesco",
			response: `{"bestMatches": [`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Query().Get("function") != "SYMBOL_SEARCH" {
						t.Errorf("unexpected function parameter: %s", r.URL.Query().Get("function"))
					}
					if r.URL.Query().Get("apikey") != "test-api-key" {
						t.Errorf("unexpected API key: %s", r.URL.Query().Get("apikey"))
					}
					if r.URL.Query().Get("keywords") != tt.symbol {
						t.Errorf("expected keywords %q, got %q", tt.symbol, r.URL.Query().Get("keywords"))
					}
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, tt.response)
				},
			))
			defer server.Close()

			client := NewClient("test-api-key")
			client.baseURL = server.URL
			client.httpClient = server.Client()

			result, err := client.TyckerSearch(tt.symbol)

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

			if len(result.BestMatches) != tt.wantMatches {
				t.Errorf("expected %d markets, got %d", tt.wantMatches, len(result.BestMatches))
			}
		})
	}
}
