package oanda

import "testing"

func TestPositionService(t *testing.T) {
	client := setupClient(t)

	t.Run("create position", func(t *testing.T) {
		req := NewMarketOrderRequest("USD_JPY", "10000")
		resp, err := client.Order.Create(t.Context(), req)
		if err != nil {
			t.Fatalf("failed to create market order: %v", err)
		}
		debugResponse(resp)
	})

	t.Run("list", func(t *testing.T) {
		resp, err := client.Position.List(t.Context())
		if err != nil {
			t.Errorf("failed to list positions: %v", err)
		}
		debugResponse(resp)
	})

	t.Run("list open", func(t *testing.T) {
		resp, err := client.Position.ListOpen(t.Context())
		if err != nil {
			t.Errorf("failed to list positions: %v", err)
		}
		debugResponse(resp)
	})

	t.Run("list by instrument", func(t *testing.T) {
		resp, err := client.Position.ListByInstrument(t.Context(), "USD_JPY")
		if err != nil {
			t.Errorf("failed to list positions: %v", err)
		}
		debugResponse(resp)
	})

	t.Run("close", func(t *testing.T) {
		req := NewPositionCloseRequest().SetLongAll()
		resp, err := client.Position.Close(t.Context(), "USD_JPY", req)
		if err != nil {
			t.Errorf("failed to close position: %v", err)
		}
		debugResponse(resp)
	})
}
