package oanda

import (
	"testing"
)

func TestClient_Trade(t *testing.T) {
	client := setupClient(t)
	var tradeID TradeID

	t.Run("open order", func(t *testing.T) {
		req := NewMarketOrderRequest("USD_JPY", "10000")
		resp, err := client.OrderCreate(t.Context(), req)
		if err != nil {
			t.Fatalf("failed to create order: %v", err)
		}
		tradeID = resp.OrderFillTransaction.TradeOpened.TradeID
		debugResponse(resp)
	})

	if tradeID == "" {
		t.Fatal("no trade ID")
	}

	t.Run("list", func(t *testing.T) {
		req := NewTradeListRequest().SetInstrument("USD_JPY")
		resp, err := client.TradeList(t.Context(), req)
		if err != nil {
			t.Errorf("failed to list trades: %s", err)
		}
		found := false
		for _, trade := range resp.Trades {
			if trade.ID == tradeID {
				found = true
			}
		}
		if !found {
			t.Errorf("failed to find trade: %s", tradeID)
		}
		debugResponse(resp)
	})

	t.Run("list open", func(t *testing.T) {
		resp, err := client.TradeListOpen(t.Context())
		if err != nil {
			t.Errorf("failed to list open trades: %s", err)
		}
		found := false
		for _, trade := range resp.Trades {
			if trade.ID == tradeID {
				found = true
			}
		}
		if !found {
			t.Errorf("failed to find trade: %s", tradeID)
		}
		debugResponse(resp)
	})

	t.Run("details", func(t *testing.T) {
		resp, err := client.TradeDetails(t.Context(), tradeID)
		if err != nil {
			t.Errorf("failed to get trade details: %s", err)
		}
		debugResponse(resp)
	})

	t.Run("update client extensions", func(t *testing.T) {
		req := TradeUpdateClientExtensionsRequest{
			ClientExtensions{
				ID:      "test ID",
				Tag:     "test Tag",
				Comment: "test Comment",
			},
		}
		resp, err := client.TradeUpdateClientExtensions(t.Context(), tradeID, req)
		if err != nil {
			t.Errorf("failed to update client extensions: %s", err)
		}
		debugResponse(resp)
	})

	t.Run("update orders", func(t *testing.T) {
		req := &TradeUpdateOrdersRequest{
			TakeProfit: NewTakeProfitDetails("200.00"),
			StopLoss:   NewStopLossDetails().SetDistance("10.00"),
		}
		resp, err := client.TradeUpdateOrders(t.Context(), tradeID, req)
		if err != nil {
			t.Errorf("failed to update orders: %s", err)
		}
		debugResponse(resp)
	})

	t.Run("close", func(t *testing.T) {
		req := NewTradeCloseALLRequest()
		resp, err := client.TradeClose(t.Context(), tradeID, req)
		if err != nil {
			t.Errorf("failed to close trade: %s", err)
		}
		debugResponse(resp)
	})
}
