package oanda

import (
	"testing"
	"time"
)

func TestPriceService_LatestCandlestick(t *testing.T) {
	client := setupClient(t)
	req := NewPriceLatestCandlesticksRequest().Specification("EUR_USD:S10:BM")
	resp, err := client.Price.LatestCandlesticks(t.Context(), req)
	if err != nil {
		t.Errorf("failed to get latest candlesticks: %v", err)
	}
	debugResponse(resp)
}

func TestPriceService_Information(t *testing.T) {
	client := setupClient(t)
	req := NewPriceInformationRequest().Instruments("EUR_USD")
	resp, err := client.Price.Information(t.Context(), req)
	if err != nil {
		t.Errorf("failed to get information: %v", err)
	}
	debugResponse(resp)
}

func TestPriceStreamService_Stream(t *testing.T) {
	client := setupStreamClient(t)
	req := NewPriceStreamRequest("USD_JPY")
	ch := make(chan PriceStreamItem)
	done := make(chan struct{}, 1)
	go func() {
		for priceStreamItem := range ch {
			debugResponse(priceStreamItem)
		}
	}()
	go func() {
		time.Sleep(10 * time.Second)
		done <- struct{}{}
	}()
	defer close(ch)
	if err := client.Price.Stream(t.Context(), req, ch, done); err != nil {
		t.Errorf("got error: %v", err)
	}
}
