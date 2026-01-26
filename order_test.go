package oanda

import "testing"

func TestClient_OrderList(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	req := NewOrderListRequest(accountID).SetInstrument("USD_JPY")
	orders, lastTransactionID, err := client.OrderList(t.Context(), req)
	if err != nil {
		t.Errorf("failed to list orders: %v", err)
	}
	t.Logf("orders: %#v", orders)
	t.Logf("lastTransactionID: %v", lastTransactionID)
}

func TestClient_OrderListPending(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	orders, lastTransactionID, err := client.OrderListPending(t.Context(), accountID)
	if err != nil {
		t.Errorf("failed to list pending orders: %v", err)
	}
	t.Logf("orders: %#v", orders)
	t.Logf("lastTransactionID: %v", lastTransactionID)
}
