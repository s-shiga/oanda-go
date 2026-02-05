package oanda

import (
	"fmt"
	"log/slog"
	"testing"
)

func TestClient_Order(t *testing.T) {
	client := setupClient(t)
	var orderID OrderID

	t.Run("create", func(t *testing.T) {
		req := NewLimitOrderRequest("USD_JPY", "10000", "100.00")
		resp, err := client.OrderCreate(t.Context(), req)
		if err != nil {
			t.Fatalf("failed to create order: %v", err)
		}
		orderID = resp.OrderCreateTransaction.ID
		slog.Debug("OrderCreate:", "OrderCreateTransaction", fmt.Sprintf("%#v", resp.OrderCreateTransaction))
		slog.Debug("OrderCreate:", "OrderFillTransaction", fmt.Sprintf("%#v", resp.OrderFillTransaction))
		slog.Debug("OrderCreate:", "OrderCancelTransaction", fmt.Sprintf("%#v", resp.OrderCancelTransaction))
		slog.Debug("OrderCreate:", "OrderReissueTransaction", fmt.Sprintf("%#v", resp.OrderReissueTransaction))
		slog.Debug("OrderCreate:", "OrderReissueRejectTransaction", fmt.Sprintf("%#v", resp.OrderReissueRejectTransaction))
		slog.Debug("OrderCreate:", "RelatedTransactionIDs", resp.RelatedTransactionIDs)
		slog.Debug("OrderCreate:", "LastTransactionID", resp.LastTransactionID)
	})

	t.Run("list", func(t *testing.T) {
		req := NewOrderListRequest().SetInstrument("USD_JPY")
		resp, err := client.OrderList(t.Context(), req)
		if err != nil {
			t.Errorf("failed to list orders: %v", err)
		}
		found := false
		for _, order := range resp.Orders {
			slog.Debug("OrderList:", "Order", fmt.Sprintf("%#v", order))
			if order.GetID() == orderID {
				found = true
			}
		}
		if !found {
			t.Errorf("did not find order with ID %v in OrderList", orderID)
		}
		slog.Debug("OrderList:", "LastTransactionID", resp.LastTransactionID)
	})

	t.Run("list pending", func(t *testing.T) {
		resp, err := client.OrderListPending(t.Context())
		if err != nil {
			t.Errorf("failed to list pending orders: %v", err)
		}
		found := false
		for _, order := range resp.Orders {
			slog.Debug("OrderListPending:", "Order", fmt.Sprintf("%#v", order))
			if order.GetID() == orderID {
				found = true
			}
		}
		if !found {
			t.Errorf("did not find order with ID %v in OrderListPending", orderID)
		}
		slog.Debug("OrderListPending:", "LastTransactionID", resp.LastTransactionID)
	})

	t.Run("details", func(t *testing.T) {
		resp, err := client.OrderDetails(t.Context(), orderID)
		if err != nil {
			t.Errorf("failed to get order details: %v", err)
		}
		slog.Debug("OrderDetails:", "Order", fmt.Sprintf("%#v", resp.Order))
		slog.Debug("OrderDetails:", "LastTransactionID", resp.LastTransactionID)
	})

	t.Run("replace", func(t *testing.T) {
		req := NewLimitOrderRequest("USD_JPY", "10000", "110.00")
		resp, err := client.OrderReplace(t.Context(), orderID, req)
		if err != nil {
			t.Errorf("failed to replace order: %v", err)
		}
		orderID = resp.OrderCreateTransaction.ID
		slog.Debug("OrderReplace:", "OrderCancelTransaction", fmt.Sprintf("%#v", resp.OrderCancelTransaction))
		slog.Debug("OrderReplace:", "OrderCreateTransaction", fmt.Sprintf("%#v", resp.OrderCreateTransaction))
		slog.Debug("OrderReplace:", "OrderFillTransaction", fmt.Sprintf("%#v", resp.OrderFillTransaction))
		slog.Debug("OrderReplace:", "OrderReissueTransaction", fmt.Sprintf("%#v", resp.OrderReissueTransaction))
		slog.Debug("OrderReplace:", "OrderReissueRejectTransaction", fmt.Sprintf("%#v", resp.OrderReissueRejectTransaction))
		slog.Debug("OrderReplace:", "RelatedTransactionIDs", resp.RelatedTransactionIDs)
		slog.Debug("OrderReplace:", "LastTransactionID", resp.LastTransactionID)
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
		resp, err := client.OrderUpdateClientExtensions(t.Context(), orderID, req)
		if err != nil {
			t.Errorf("failed to update client extensions: %v", err)
		}
		slog.Debug("OrderUpdateClientExtensions:", "OrderClientExtensionsModifyTransaction", fmt.Sprintf("%#v", resp.OrderClientExtensionsModifyTransaction))
		slog.Debug("OrderUpdateClientExtensions:", "LastTransactionID", resp.LastTransactionID)
		slog.Debug("OrderUpdateClientExtensions:", "RelatedTransactionIDs", resp.RelatedTransactionIDs)
	})

	t.Run("cancel", func(t *testing.T) {
		resp, err := client.OrderCancel(t.Context(), orderID)
		if err != nil {
			t.Errorf("failed to cancel order: %v", err)
		}
		slog.Debug("OrderCancel:", "OrderCancelTransaction", fmt.Sprintf("%#v", resp.OrderCancelTransaction))
		slog.Debug("OrderCancel:", "RelatedTransactionIDs", resp.RelatedTransactionIDs)
		slog.Debug("OrderCancel:", "LastTransactionID", resp.LastTransactionID)
	})
}
