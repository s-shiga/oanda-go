package oanda

import "testing"

func TestClient_OrderCreate(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	req := NewLimitOrderRequest("USD_JPY", "10000", "100.00")
	resp, err := client.OrderCreate(t.Context(), accountID, req)
	if err != nil {
		t.Errorf("failed to create order: %v", err)
	}
	t.Logf("resp: %#v", resp)
}

func TestClient_OrderList(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	req := NewOrderListRequest(accountID).SetInstrument("USD_JPY")
	orders, lastTransactionID, err := client.OrderList(t.Context(), req)
	if err != nil {
		t.Errorf("failed to list orders: %v", err)
	}
	t.Logf("orders: %#v", orders)
	if len(orders) > 0 {
		for _, order := range orders {
			t.Logf("order: %#v", order)
		}
	}
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

func TestClient_OrderDetails(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	orderID := "427"
	resp, err := client.OrderDetails(t.Context(), accountID, orderID)
	if err != nil {
		t.Errorf("failed to get order details: %v", err)
	}
	t.Logf("response: %#v", resp)
}
