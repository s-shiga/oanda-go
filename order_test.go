package oanda

import (
	"testing"
)

func TestOrderService(t *testing.T) {
	t.Run("market order", testMarketOrder)
	t.Run("limit order", testLimitOrder)
}

func checkOrders(t *testing.T, orders []Order, takeProfitOrderID, stopLossOrderID OrderID) {
	foundTakeProfit := false
	foundStopLoss := false
	for _, order := range orders {
		if order.GetID() == takeProfitOrderID {
			foundTakeProfit = true
		}
		if order.GetID() == stopLossOrderID {
			foundStopLoss = true
		}
	}
	if !foundTakeProfit {
		t.Errorf("did not find take profit order with ID %v in OrderListPending", takeProfitOrderID)
	}
	if !foundStopLoss {
		t.Errorf("did not find stop loss order with ID %v in OrderListPending", stopLossOrderID)
	}
}

func testMarketOrder(t *testing.T) {
	client := setupClient(t)
	var tradeID TradeID
	var takeProfitOrderID OrderID
	var stopLossOrderID OrderID

	t.Run("create market order", func(t *testing.T) {
		req := NewMarketOrderRequest("USD_JPY", "10000")
		resp, err := client.Order.Create(t.Context(), req)
		if err != nil {
			t.Fatalf("failed to create market order: %v", err)
		}
		tradeID = resp.OrderFillTransaction.TradeOpened.TradeID
		debugResponse(resp)
	})

	t.Run("create take profit order", func(t *testing.T) {
		req := NewTakeProfitOrderRequest(tradeID, "170.00")
		resp, err := client.Order.Create(t.Context(), req)
		if err != nil {
			t.Errorf("failed to create take profit order: %v", err)
		}
		if tp := resp.OrderCreateTransaction.GetType(); tp != TransactionTypeTakeProfitOrder {
			t.Errorf("wrong transaction type: %v", tp)
		}
		takeProfitOrderID = resp.OrderCreateTransaction.GetID()
		debugResponse(resp)
	})

	t.Run("create stop loss order", func(t *testing.T) {
		req := NewStopLossOrderRequest(tradeID).SetDistance("10.000")
		resp, err := client.Order.Create(t.Context(), req)
		if err != nil {
			t.Errorf("failed to create stop loss order: %v", err)
		}
		if tp := resp.OrderCreateTransaction.GetType(); tp != TransactionTypeStopLossOrder {
			t.Errorf("wrong transaction type: %v", tp)
		}
		stopLossOrderID = resp.OrderCreateTransaction.GetID()
		debugResponse(resp)
	})

	t.Run("list", func(t *testing.T) {
		req := NewOrderListRequest().SetInstrument("USD_JPY")
		resp, err := client.Order.List(t.Context(), req)
		if err != nil {
			t.Errorf("failed to list orders: %v", err)
		}
		checkOrders(t, resp.Orders, takeProfitOrderID, stopLossOrderID)
		debugResponse(resp)
	})

	t.Run("list pending", func(t *testing.T) {
		resp, err := client.Order.ListPending(t.Context())
		if err != nil {
			t.Errorf("failed to list pending orders: %v", err)
		}
		checkOrders(t, resp.Orders, takeProfitOrderID, stopLossOrderID)
		debugResponse(resp)
	})

	t.Run("close trade", func(t *testing.T) {
		req := NewTradeCloseALLRequest()
		resp, err := client.Trade.Close(t.Context(), tradeID, req)
		if err != nil {
			t.Errorf("failed to close trade: %v", err)
		}
		debugResponse(resp)
	})
}

func testLimitOrder(t *testing.T) {
	client := setupClient(t)
	var orderID OrderID

	t.Run("create", func(t *testing.T) {
		req := NewLimitOrderRequest("USD_JPY", "10000", "100.00")
		resp, err := client.Order.Create(t.Context(), req)
		if err != nil {
			t.Fatalf("failed to create order: %v", err)
		}
		orderID = resp.OrderCreateTransaction.GetID()
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
		orderID = resp.OrderCreateTransaction.GetID()
		debugResponse(resp)
	})

	t.Run("update client extensions", func(t *testing.T) {
		req := OrderUpdateClientExtensionsRequest{
			ClientExtensions:      NewClientExtensions().SetID("test client extension ID").SetTag("test tag"),
			TradeClientExtensions: NewClientExtensions().SetID("test trade client extension ID").SetComment("test comment"),
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

func TestOrderList(t *testing.T) {
	client := setupClient(t)
	req := NewOrderListRequest().SetInstrument("USD_JPY")
	resp, err := client.Order.List(t.Context(), req)
	if err != nil {
		t.Errorf("failed to list orders: %v", err)
	}
	debugResponse(resp)
}
