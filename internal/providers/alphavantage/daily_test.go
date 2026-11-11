package alphavantage

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetTimeSeriesDailyForSymbolNetworkError(t *testing.T) {
	client := NewClient("test-api-key")

	client.baseURL = "http://127.0.0.1:1"
	client.httpClient = &http.Client{
		Timeout: time.Second,
	}

	result, err := client.GetTimeSeriesDailyForSymbol("test")

	if err == nil {
		t.Fatal("expected a network error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
}

func TestGetTimeSeriesDailyForSymbol(t *testing.T) {
	tests := []struct {
		name           string
		response       string
		symbol         string
		wantErr        bool
		wantTimeSeries int
	}{
		{
			name:   "successful response",
			symbol: "IBM",
			response: `{
				"Meta Data": {
					"1. Information": "Daily Prices (open, high, low, close) and Volumes",
					"2. Symbol": "IBM",
					"3. Last Refreshed": "2026-11-10",
					"4. Output Size": "Compact",
					"5. Time Zone": "US/Eastern"
				},
				"Time Series (Daily)": {
					"2026-11-10": {
						"1. open": "244.3800",
						"2. high": "246.0000",
						"3. low": "242.7000",
						"4. close": "244.5300",
						"5. volume": "2971839"
					},
					"2026-11-09": {
						"1. open": "246.0000",
						"2. high": "248.8100",
						"3. low": "245.1957",
						"4. close": "246.1100",
						"5. volume": "4629816"
					}
				}
			}`,
			wantErr:        false,
			wantTimeSeries: 2,
		},
		{
			name:   "invalid tycker",
			symbol: "IBMMMMMMMMMMMMMM",
			response: `{
				"Error Message": "Invalid API call. Please retry or visit the documentation (https://www.alphavantage.co/documentation/) for TIME_SERIES_DAILY."
			}`,
			wantErr: true,
		},
		{
			name:   "empty tycker",
			symbol: "",
			response: `{
				"Error Message": "Invalid API call. Please retry or visit the documentation (https://www.alphavantage.co/documentation/) for TIME_SERIES_DAILY."
			}`,
			wantErr: true,
		},
		{
			name:     "invalid JSON",
			symbol:   "IBM",
			response: `{"Meta Data": {"1. Inf`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Query().Get("function") != "TIME_SERIES_DAILY" {
						t.Errorf("unexpected function parameter: %s", r.URL.Query().Get("function"))
					}
					if r.URL.Query().Get("apikey") != "test-api-key" {
						t.Errorf("unexpected API key: %s", r.URL.Query().Get("apikey"))
					}
					if r.URL.Query().Get("symbol") != tt.symbol {
						t.Errorf("expected symbol %q, got %q", tt.symbol, r.URL.Query().Get("symbol"))
					}
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, tt.response)
				},
			))
			defer server.Close()

			client := NewClient("test-api-key")
			client.baseURL = server.URL
			client.httpClient = server.Client()

			result, err := client.GetTimeSeriesDailyForSymbol(tt.symbol)

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

			if len(result.TimeSeries) != tt.wantTimeSeries {
				t.Errorf("expected %d markets, got %d", tt.wantTimeSeries, len(result.TimeSeries))
			}
		})
	}
}
