package oanda

import "testing"

func TestInstrumentService(t *testing.T) {
	runEndpointTests(t, []endpointTest{
		{
			name:     "candlestick time range",
			method:   "GET",
			path:     "/v3/instruments/USD_JPY/candles",
			query:    "granularity=M1&from=2025-01-01T00:00:00Z&to=2025-01-02T00:00:00Z",
			response: `{"instrument":"USD_JPY","granularity":"M1","candles":[{"time":"2025-01-01T00:00:00Z","volume":12,"complete":true}]}`,
			call: func(c *Client) (any, error) {
				return c.Instrument.Candlesticks(t.Context(), NewCandlesticksRequest("USD_JPY", M1).
					SetFrom(mustTime(t, "2025-01-01T00:00:00Z")).SetTo(mustTime(t, "2025-01-02T00:00:00Z")))
			},
		},
		{
			name:     `list`,
			method:   `GET`,
			path:     testAccountPath + "/instruments",
			response: `{"instruments":[{"name":"EUR_USD","pipLocation":-4,"displayPrecision":5},{"name":"USD_JPY","pipLocation":-2,"displayPrecision":3}],"lastTransactionID":"50"}`,
			call:     func(c *Client) (any, error) { return c.Instrument.List(t.Context()) },
		},
		{
			name:     `filter instruments`,
			method:   `GET`,
			path:     testAccountPath + "/instruments",
			query:    `instruments=EUR_USD,USD_JPY`,
			response: `{"instruments":[{"name":"EUR_USD","pipLocation":-4,"displayPrecision":5},{"name":"USD_JPY","pipLocation":-2,"displayPrecision":3}],"lastTransactionID":"50"}`,
			call:     func(c *Client) (any, error) { return c.Instrument.List(t.Context(), "EUR_USD", "USD_JPY") },
		},
		{
			name:     `candlesticks`,
			method:   `GET`,
			path:     "/v3/instruments/USD_JPY/candles",
			query:    `granularity=M1&count=1&price=M`,
			response: `{"instrument":"USD_JPY","granularity":"M1","candles":[{"time":"2025-01-01T00:00:00Z","mid":{"o":"150.000","h":"150.100","l":"149.900","c":"150.050"},"volume":12,"complete":true}]}`,
			call: func(c *Client) (any, error) {
				return c.Instrument.Candlesticks(t.Context(), NewCandlesticksRequest("USD_JPY", M1).SetCount(1).Mid())
			},
		},
	})
}
