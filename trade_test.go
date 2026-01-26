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
	t.Logf("transactions: %#v", transactionID)
}
