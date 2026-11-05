package alphavantage

import (
	"encoding/json"
	"io"
)

type DailyTimeSeries struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type TimeSeries struct {
	MetaData struct {
		Information   string `json:"1. Information"`
		Symbol        string `json:"2. Symbol"`
		LastRefreshed string `json:"3. Last Refreshed"`
		OutputSize    string `json:"4. Output Size"`
		TimeZone      string `json:"5. Time Zone"`
	} `json:"Meta Data"`
	TimeSeries map[string]DailyTimeSeries `json:"Time Series (Daily)"`
}

func (c *Client) GetTimeSeriesDailyForSymbol(symbol string) (*TimeSeries, error) {
	url := c.baseURL + "?function=TIME_SERIES_DAILY&symbol=" + symbol + "&apikey=" + c.apiKey
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	timeSeries := TimeSeries{}
	err = json.Unmarshal(body, &timeSeries)
	if err != nil {
		return nil, err
	}
	return &timeSeries, nil
}
