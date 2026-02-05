package oanda

import "testing"

func TestClient_TransactionList(t *testing.T) {
	client := setupClient(t)
	req := NewTransactionListRequest()
	resp, err := client.TransactionList(t.Context(), req)
	if err != nil {
		t.Errorf("failed to list transactions: %v", err)
	}
	debugResponse(resp)
}
