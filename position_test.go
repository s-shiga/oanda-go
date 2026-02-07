package oanda

import "testing"

func TestClient_PositionList(t *testing.T) {
	client := setupClient(t)
	resp, err := client.Position.List(t.Context())
	if err != nil {
		t.Errorf("failed to list positions: %v", err)
	}
	debugResponse(resp)
}

func TestClient_PositionListOpen(t *testing.T) {
	client := setupClient(t)
	resp, err := client.Position.ListOpen(t.Context())
	if err != nil {
		t.Errorf("failed to list positions: %v", err)
	}
	debugResponse(resp)
}

func TestClient_PositionListInstrument(t *testing.T) {
	client := setupClient(t)
	resp, err := client.Position.ListByInstrument(t.Context(), "USD_JPY")
	if err != nil {
		t.Errorf("failed to list positions: %v", err)
	}
	debugResponse(resp)
}
