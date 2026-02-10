package oanda

import "testing"

func TestPositionService_List(t *testing.T) {
	client := setupClient(t)
	resp, err := client.Position.List(t.Context())
	if err != nil {
		t.Errorf("failed to list positions: %v", err)
	}
	debugResponse(resp)
}

func TestPositionService_ListOpen(t *testing.T) {
	client := setupClient(t)
	resp, err := client.Position.ListOpen(t.Context())
	if err != nil {
		t.Errorf("failed to list positions: %v", err)
	}
	debugResponse(resp)
}

func TestPositionService_ListByInstrument(t *testing.T) {
	client := setupClient(t)
	resp, err := client.Position.ListByInstrument(t.Context(), "USD_JPY")
	if err != nil {
		t.Errorf("failed to list positions: %v", err)
	}
	debugResponse(resp)
}
