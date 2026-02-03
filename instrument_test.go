package oanda

import (
	"testing"
	"time"
)

func TestClient_Candlesticks(t *testing.T) {
	client := setupClient(t)
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
	req := NewCandlesticksRequest("USD_JPY", M1).SetFrom(from).SetTo(to)
	resp, err := client.Candlesticks(t.Context(), req)
	if err != nil {
		t.Errorf("failed to get candlesticks: %v", err)
	}
	t.Logf("response: %#v", resp)
	if len(resp.Candles) == 0 {
		t.Errorf("got no candlesticks")
	}
}
