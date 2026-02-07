package oanda

import "testing"

func TestPriceService_LatestCandlestick(t *testing.T) {
	client := setupClient(t)
	req := NewPriceLatestCandlesticksRequest().Specification("EUR_USD:S10:BM")
	resp, err := client.Price.LatestCandlesticks(t.Context(), req)
	if err != nil {
		t.Errorf("failed to get latest candlesticks: %v", err)
	}
	debugResponse(resp)
}
