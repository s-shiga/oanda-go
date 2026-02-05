package oanda

import "testing"

func TestClient_TradeList(t *testing.T) {
	client := setupClient(t)
	req := NewTradeListRequest().SetInstrument("USD_JPY")
	trades, transactionID, err := client.TradeList(t.Context(), req)
	if err != nil {
		t.Errorf("failed to list trades: %s", err)
	}
	t.Logf("trades: %#v", trades)
	t.Logf("last transaction: %#v", transactionID)
}

func TestClient_TradeListOpen(t *testing.T) {
	client := setupClient(t)
	trades, transactionID, err := client.TradeListOpen(t.Context())
	if err != nil {
		t.Errorf("failed to list trades: %s", err)
	}
	t.Logf("trades: %#v", trades)
	t.Logf("last transaction: %#v", transactionID)
}

func TestClient_TradeDetails(t *testing.T) {
	client := setupClient(t)
	tradeID := "15808"
	trade, transactionID, err := client.TradeDetails(t.Context(), tradeID)
	if err != nil {
		t.Errorf("failed to list trades: %s", err)
	}
	t.Logf("trades: %#v", trade)
	t.Logf("last transaction: %#v", transactionID)
}

//func TestClient_Trade(t *testing.T) {
//	client := setupClient(t)
//	accountID := setupAccountID(t)
//	var tradeID TradeID
//
//	t.Run("open order", func(t *testing.T) {
//		req := NewMarketOrderRequest("USD_JPY", "10000")
//		resp, err := client.OrderCreate(t.Context(), accountID, req)
//		if err != nil {
//			t.Fatalf("failed to create order: %v", err)
//		}
//		tradeID = resp.OrderFillTransaction.TradeOpened.TradeID
//	})
//
//	t.Run("list", func(t *testing.T) {
//		req := NewTradeLIst
//		resp, err := client.TradeList(t.Context())
//	})
//}
