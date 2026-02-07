package oanda

import (
	"testing"
)

func TestClient_Order(t *testing.T) {
	client := setupClient(t)
	var orderID OrderID

	t.Run("create", func(t *testing.T) {
		req := NewLimitOrderRequest("USD_JPY", "10000", "100.00")
		resp, err := client.Order.Create(t.Context(), req)
		if err != nil {
			t.Fatalf("failed to create order: %v", err)
		}
		orderID = resp.OrderCreateTransaction.ID
		debugResponse(resp)
	})

	t.Run("list", func(t *testing.T) {
		req := NewOrderListRequest().SetInstrument("USD_JPY")
		resp, err := client.Order.List(t.Context(), req)
		if err != nil {
			t.Errorf("failed to list orders: %v", err)
		}
		found := false
		for _, order := range resp.Orders {
			if order.GetID() == orderID {
				found = true
			}
		}
		if !found {
			t.Errorf("did not find order with ID %v in OrderList", orderID)
		}
		debugResponse(resp)
	})

	t.Run("list pending", func(t *testing.T) {
		resp, err := client.Order.ListPending(t.Context())
		if err != nil {
			t.Errorf("failed to list pending orders: %v", err)
		}
		found := false
		for _, order := range resp.Orders {
			if order.GetID() == orderID {
				found = true
			}
		}
		if !found {
			t.Errorf("did not find order with ID %v in OrderListPending", orderID)
		}
		debugResponse(resp)
	})

	t.Run("details", func(t *testing.T) {
		resp, err := client.Order.Details(t.Context(), orderID)
		if err != nil {
			t.Errorf("failed to get order details: %v", err)
		}
		debugResponse(resp)
	})

	t.Run("replace", func(t *testing.T) {
		req := NewLimitOrderRequest("USD_JPY", "10000", "110.00")
		resp, err := client.OrderReplace(t.Context(), orderID, req)
		if err != nil {
			t.Errorf("failed to replace order: %v", err)
		}
		orderID = resp.OrderCreateTransaction.ID
		debugResponse(resp)
	})

	t.Run("update client extensions", func(t *testing.T) {
		req := OrderUpdateClientExtensionsRequest{
			ClientExtensions: ClientExtensions{
				ID:      "test client extension ID",
				Tag:     "test tag",
				Comment: "test comment",
			},
			TradeClientExtensions: ClientExtensions{
				ID:      "test trade client extension ID",
				Tag:     "test tag",
				Comment: "test comment",
			},
		}
		resp, err := client.Order.UpdateClientExtensions(t.Context(), orderID, req)
		if err != nil {
			t.Errorf("failed to update client extensions: %v", err)
		}
		debugResponse(resp)
	})

	t.Run("cancel", func(t *testing.T) {
		resp, err := client.Order.Cancel(t.Context(), orderID)
		if err != nil {
			t.Errorf("failed to cancel order: %v", err)
		}
		debugResponse(resp)
	})
}
