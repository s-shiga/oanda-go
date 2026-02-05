package oanda

import "testing"

func TestClient_PositionList(t *testing.T) {
	client := setupClient(t)
	positions, transactionID, err := client.PositionList(t.Context())
	if err != nil {
		t.Errorf("failed to list positions: %v", err)
	}
	t.Logf("positions: %#v", positions)
	t.Logf("transactionID: %#v", transactionID)
}

func TestClient_PositionListOpen(t *testing.T) {
	client := setupClient(t)
	positions, transactionID, err := client.PositionListOpen(t.Context())
	if err != nil {
		t.Errorf("failed to list positions: %v", err)
	}
	t.Logf("positions: %#v", positions)
	t.Logf("transactionID: %#v", transactionID)
}

func TestClient_PositionListInstrument(t *testing.T) {
	client := setupClient(t)
	position, transactionID, err := client.PositionListInstrument(t.Context(), "USD_JPY")
	if err != nil {
		t.Errorf("failed to list positions: %v", err)
	}
	t.Logf("position: %#v", position)
	t.Logf("transactionID: %#v", transactionID)
}
