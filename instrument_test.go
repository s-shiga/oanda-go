package oanda

import (
	"testing"
	"time"
)

func TestClient_AccountInstruments(t *testing.T) {
	client := setupClient(t)

	t.Run("without instruments", func(t *testing.T) {
		resp, err := client.Instrument.List(t.Context())
		if err != nil {
			t.Errorf("failed to list instruments: %v", err)
		}
		debugResponse(resp)
	})

	t.Run("with instruments", func(t *testing.T) {
		resp, err := client.Instrument.List(t.Context(), "EUR_USD", "USD_JPY")
		if err != nil {
			t.Errorf("failed to list instruments: %v", err)
		}
		debugResponse(resp)
	})
}

func TestClient_Candlesticks(t *testing.T) {
	client := setupClientWithoutAccountID(t)
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
	req := NewCandlesticksRequest("USD_JPY", M1).SetFrom(from).SetTo(to)
	resp, err := client.Instrument.Candlesticks(t.Context(), req)
	if err != nil {
		t.Errorf("failed to get candlesticks: %v", err)
	}
	if len(resp.Candles) == 0 {
		t.Errorf("got no candlesticks")
	}
	debugResponse(resp)
}
