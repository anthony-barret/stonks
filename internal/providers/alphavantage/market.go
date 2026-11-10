package alphavantage

import (
	"encoding/json"
	"io"
)

type Market struct {
	MarketType       string `json:"market_type"`
	Region           string `json:"region"`
	PrimaryExchanges string `json:"primary_exchanges"`
	LocalOpen        string `json:"local_open"`
	LocalClose       string `json:"local_close"`
	CurrentStatus    string `json:"current_status"`
	Notes            string `json:"notes"`
}

type MarketStatus struct {
	Endpoint string   `json:"endpoint"`
	Markets  []Market `json:"markets"`
}

func (c *Client) GetMarketStatus() (*MarketStatus, error) {
	url := c.baseURL + "?function=MARKET_STATUS&apikey=" + c.apiKey
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = checkAPIError(body)
	if err != nil {
		return nil, err
	}
	market := MarketStatus{}
	err = json.Unmarshal(body, &market)
	if err != nil {
		return nil, err
	}
	return &market, nil
}
