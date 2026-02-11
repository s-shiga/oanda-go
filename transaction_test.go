package oanda

import (
	"strconv"
	"testing"
	"time"
)

func TestTransactionService(t *testing.T) {
	client := setupClient(t)
	var lastTransactionID TransactionID

	t.Run("list", func(t *testing.T) {
		req := NewTransactionListRequest()
		resp, err := client.Transaction.List(t.Context(), req)
		if err != nil {
			t.Errorf("failed to list transactions: %v", err)
		}
		debugResponse(resp)
		lastTransactionID = resp.LastTransactionID
	})

	t.Run("details", func(t *testing.T) {
		resp, err := client.Transaction.Details(t.Context(), lastTransactionID)
		if err != nil {
			t.Errorf("failed to get details: %v", err)
		}
		debugResponse(resp)
	})

	t.Run("get by ID range", func(t *testing.T) {
		id, err := strconv.Atoi(lastTransactionID)
		if err != nil {
			t.Errorf("failed to convert last transaction id to int: %v", err)
		}
		from := strconv.Itoa(id - 20)
		to := strconv.Itoa(id - 10)
		req := NewTransactionGetByIDRangeRequest(from, to)
		resp, err := client.Transaction.GetByIDRange(t.Context(), req)
		if err != nil {
			t.Errorf("failed to get transactions by ID range: %v", err)
		}
		debugResponse(resp)
	})

	t.Run("get by since ID", func(t *testing.T) {
		id, err := strconv.Atoi(lastTransactionID)
		if err != nil {
			t.Errorf("failed to convert last transaction id to int: %v", err)
		}
		since := strconv.Itoa(id - 10)
		req := NewTransactionGetBySinceIDRequest(since)
		resp, err := client.Transaction.GetBySinceID(t.Context(), req)
		if err != nil {
			t.Errorf("failed to get transactions by since ID: %v", err)
		}
		debugResponse(resp)
	})
}

func TestStreamClient_Transaction(t *testing.T) {
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
	if err := client.Transaction(t.Context(), ch, done); err != nil {
		t.Errorf("got error: %v", err)
	}
}
