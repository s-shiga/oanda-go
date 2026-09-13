package oanda

import "testing"

func TestPricingService(t *testing.T) {
	runEndpointTests(t, []endpointTest{
		{
			name:     `latest candlesticks`,
			method:   `GET`,
			path:     testAccountPath + "/candles/latest",
			query:    `candleSpecifications=USD_JPY:M1:M`,
			response: `{"latestCandles":[{"instrument":"USD_JPY","granularity":"M1","candles":[{"time":"2025-01-01T00:00:00Z","mid":{"o":"150.000","h":"150.100","l":"149.900","c":"150.050"},"volume":12,"complete":true}]}]}`,
			want:     `[{"instrument":"USD_JPY","granularity":"M1","candles":[{"time":"2025-01-01T00:00:00Z","mid":{"o":"150.000","h":"150.100","l":"149.900","c":"150.050"},"volume":12,"complete":true}]}]`,
			call: func(c *Client) (any, error) {
				return c.Price.LatestCandlesticks(t.Context(), NewPriceLatestCandlesticksRequest().AddSpecifications("USD_JPY:M1:M"))
			},
		},
		{
			name:     `information`,
			method:   `GET`,
			path:     testAccountPath + "/pricing",
			query:    `instruments=USD_JPY`,
			response: `{"prices":[{"type":"PRICE","instrument":"USD_JPY","time":"2025-01-01T00:00:00Z","tradeable":true,"bids":[{"price":"150.000","liquidity":1000000}],"asks":[{"price":"150.010","liquidity":1000000}]}],"time":"2025-01-01T00:00:00Z"}`,
			call: func(c *Client) (any, error) {
				return c.Price.Information(t.Context(), NewPriceInformationRequest().AddInstruments("USD_JPY"))
			},
		},
		{
			name:     `account candlesticks`,
			method:   `GET`,
			path:     testAccountPath + "/instruments/USD_JPY/candles",
			query:    `granularity=M1&count=1&units=10&price=M`,
			response: `{"instrument":"USD_JPY","granularity":"M1","candles":[{"time":"2025-01-01T00:00:00Z","mid":{"o":"150.000","h":"150.100","l":"149.900","c":"150.050"},"volume":12,"complete":true}]}`,
			call: func(c *Client) (any, error) {
				return c.Price.Candlesticks(t.Context(), NewPriceCandlesticksRequest("USD_JPY", M1).SetCount(1).SetUnits(10).Mid())
			},
		},
	})
}
