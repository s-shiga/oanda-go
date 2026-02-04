package oanda

import (
	"fmt"
	"log/slog"
	"testing"
)

func TestClient_OrderCreate(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	req := NewLimitOrderRequest("USD_JPY", "10000", "100.00")
	resp, err := client.OrderCreate(t.Context(), accountID, req)
	if err != nil {
		t.Errorf("failed to create order: %v", err)
	}
	slog.Debug("OrderCreate:", "resp", fmt.Sprintf("%#v", resp))
}

func TestClient_OrderList(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	req := NewOrderListRequest(accountID).SetInstrument("USD_JPY")
	resp, err := client.OrderList(t.Context(), req)
	if err != nil {
		t.Errorf("failed to list orders: %v", err)
	}
	for _, order := range resp.Orders {
		slog.Debug("OrderList:", "Order", fmt.Sprintf("%#v", order))
	}
	slog.Debug("OrderList:", "LastTransactionID", resp.LastTransactionID)
}

func TestClient_OrderListPending(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	resp, err := client.OrderListPending(t.Context(), accountID)
	if err != nil {
		t.Errorf("failed to list pending orders: %v", err)
	}
	for _, order := range resp.Orders {
		slog.Debug("OrderListPending:", "Order", fmt.Sprintf("%#v", order))
	}
	slog.Debug("OrderListPending:", "LastTransactionID", resp.LastTransactionID)
}

func TestClient_OrderDetails(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	orderID := "427"
	resp, err := client.OrderDetails(t.Context(), accountID, orderID)
	if err != nil {
		t.Errorf("failed to get order details: %v", err)
	}
	slog.Debug("OrderDetails:", "Order", fmt.Sprintf("%#v", resp.Order))
	slog.Debug("OrderDetails:", "LastTransactionID", resp.LastTransactionID)
}
