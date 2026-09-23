package oanda

import (
	"encoding/json"
	"testing"
)

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
			response: `{"prices":[{"type":"PRICE","instrument":"USD_JPY","time":"2025-01-01T00:00:00Z","tradeable":true,"bids":[{"price":"150.000","liquidity":1000000}],"asks":[{"price":"150.010","liquidity":"250000.5"}]}],"time":"2025-01-01T00:00:00Z"}`,
			want:     `{"prices":[{"type":"PRICE","instrument":"USD_JPY","time":"2025-01-01T00:00:00Z","tradeable":true,"bids":[{"price":"150.000","liquidity":"1000000"}],"asks":[{"price":"150.010","liquidity":"250000.5"}]}],"time":"2025-01-01T00:00:00Z"}`,
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
		{
			name:     `account candlesticks with fractional units`,
			method:   `GET`,
			path:     testAccountPath + "/instruments/USD_JPY/candles",
			query:    `granularity=M1&units=0.125`,
			response: `{"instrument":"USD_JPY","granularity":"M1","candles":[]}`,
			call: func(c *Client) (any, error) {
				return c.Price.Candlesticks(t.Context(), NewPriceCandlesticksRequest("USD_JPY", M1).SetUnitsDecimal("0.125"))
			},
		},
	})
}

func TestPriceCandlesticksInvalidUnits(t *testing.T) {
	for _, units := range []DecimalNumber{"", "0", "-1.5", "1/2", "1e2"} {
		t.Run(string(units), func(t *testing.T) {
			_, err := NewPriceCandlesticksRequest("EUR_USD", M1).SetUnitsDecimal(units).values()
			if err == nil {
				t.Errorf("units %q should be rejected", units)
			}
		})
	}
}

func TestPriceBucketLiquidity(t *testing.T) {
	cases := []struct {
		name string
		json string
		want DecimalNumber
	}{
		{"integer number", `{"price":"1.1","liquidity":1000000}`, "1000000"},
		{"decimal number", `{"price":"1.1","liquidity":0.5}`, "0.5"},
		{"integer string", `{"price":"1.1","liquidity":"1000000"}`, "1000000"},
		{"decimal string", `{"price":"1.1","liquidity":"250000.5"}`, "250000.5"},
		{"null", `{"price":"1.1","liquidity":null}`, ""},
		{"missing", `{"price":"1.1"}`, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var bucket PriceBucket
			if err := json.Unmarshal([]byte(tc.json), &bucket); err != nil {
				t.Fatal(err)
			}
			if bucket.Price != "1.1" || bucket.Liquidity != tc.want {
				t.Errorf("bucket = %#v, want price 1.1 and liquidity %q", bucket, tc.want)
			}
		})
	}
	var bucket PriceBucket
	if err := json.Unmarshal([]byte(`{"price":"1.1","liquidity":"lots"}`), &bucket); err == nil {
		t.Error("want error for non-numeric liquidity, got nil")
	}
}
