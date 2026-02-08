package oanda

import (
	"testing"
	"time"
)

func TestClient_TransactionList(t *testing.T) {
	client := setupClient(t)
	req := NewTransactionListRequest()
	resp, err := client.Transaction.List(t.Context(), req)
	if err != nil {
		t.Errorf("failed to list transactions: %v", err)
	}
	debugResponse(resp)
}

func TestTransactionService_Details(t *testing.T) {
	client := setupClient(t)
	transactionID := "533"
	resp, err := client.Transaction.Details(t.Context(), transactionID)
	if err != nil {
		t.Errorf("failed to get details: %v", err)
	}
	debugResponse(resp)
}

func TestTransactionService_GetByIDRange(t *testing.T) {
	client := setupClient(t)
	req := NewTransactionGetByIDRangeRequest("500", "510")
	resp, err := client.Transaction.GetByIDRange(t.Context(), req)
	if err != nil {
		t.Errorf("failed to get transactions by ID range: %v", err)
	}
	debugResponse(resp)
}

func TestTransactionStreamService_Stream(t *testing.T) {
	client := setupStreamClient(t)
	ch := make(chan TransactionStreamItem)
	done := make(chan struct{}, 1)
	go func() {
		for item := range ch {
			debugResponse(item)
		}
	}()
	go func() {
		time.Sleep(10 * time.Second)
		done <- struct{}{}
	}()
	defer close(ch)
	if err := client.Transaction.Stream(t.Context(), ch, done); err != nil {
		t.Errorf("got error: %v", err)
	}
}
