package oanda

import (
	"testing"
	"time"
)

func TestPriceService_LatestCandlestick(t *testing.T) {
	client := setupClient(t)
	req := NewPriceLatestCandlesticksRequest().AddSpecifications("EUR_USD:S10:BM")
	resp, err := client.Price.LatestCandlesticks(t.Context(), req)
	if err != nil {
		t.Errorf("failed to get latest candlesticks: %v", err)
	}
	debugResponse(resp)
}

func TestPriceService_Information(t *testing.T) {
	client := setupClient(t)
	req := NewPriceInformationRequest().AddInstruments("EUR_USD")
	resp, err := client.Price.Information(t.Context(), req)
	if err != nil {
		t.Errorf("failed to get information: %v", err)
	}
	debugResponse(resp)
}

func TestStreamClient_Price(t *testing.T) {
	client := setupStreamClient(t)
	req := NewPriceStreamRequest("USD_JPY")
	ch := make(chan PriceStreamItem)
	done := make(chan struct{})
	go func() {
		for priceStreamItem := range ch {
			debugResponse(priceStreamItem)
		}
	}()
	timer := time.AfterFunc(10*time.Second, func() { close(done) })
	defer timer.Stop()
	defer close(ch)
	if err := client.Price(t.Context(), req, ch, done); err != nil {
		t.Errorf("got error: %v", err)
	}
}
