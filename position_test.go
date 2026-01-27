package oanda

import "testing"

func TestClient_PositionList(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	positions, transactionID, err := client.PositionList(t.Context(), accountID)
	if err != nil {
		t.Errorf("failed to list positions: %v", err)
	}
	t.Logf("positions: %#v", positions)
	t.Logf("transactionID: %#v", transactionID)
}
