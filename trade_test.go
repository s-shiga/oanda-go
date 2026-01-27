package oanda

import "testing"

func TestClient_TradeList(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	req := NewTradeListRequest(accountID).SetInstrument("USD_JPY")
	trades, transactionID, err := client.TradeList(t.Context(), req)
	if err != nil {
		t.Errorf("failed to list trades: %s", err)
	}
	t.Logf("trades: %#v", trades)
	t.Logf("last transaction: %#v", transactionID)
}

func TestClient_TradeListOpen(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	trades, transactionID, err := client.TradeListOpen(t.Context(), accountID)
	if err != nil {
		t.Errorf("failed to list trades: %s", err)
	}
	t.Logf("trades: %#v", trades)
	t.Logf("last transaction: %#v", transactionID)
}

func TestClient_TradeDetails(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	tradeID := "15808"
	trade, transactionID, err := client.TradeDetails(t.Context(), accountID, tradeID)
	if err != nil {
		t.Errorf("failed to list trades: %s", err)
	}
	t.Logf("trades: %#v", trade)
	t.Logf("last transaction: %#v", transactionID)
}
