package oanda

import (
	"strconv"
	"testing"
)

func TestAccountService(t *testing.T) {
	client := setupClient(t)
	var lastTransactionID TransactionID

	t.Run("list", func(t *testing.T) {
		client := setupClientWithoutAccountID(t)
		resp, err := client.Account.List(t.Context())
		if err != nil {
			t.Errorf("failed to list accounts: %v", err)
		}
		debugResponse(resp)
	})

	t.Run("details", func(t *testing.T) {
		resp, err := client.Account.Details(t.Context())
		if err != nil {
			t.Errorf("failed to get account details: %v", err)
		}
		lastTransactionID = resp.LastTransactionID
		debugResponse(resp.Account)
	})

	t.Run("summary", func(t *testing.T) {
		resp, err := client.Account.Summary(t.Context())
		if err != nil {
			t.Errorf("failed to get account summary: %v", err)
		}
		debugResponse(resp)
	})

	t.Run("configure", func(t *testing.T) {
		req := NewAccountConfigureRequest().SetAlias("TestAlias")
		resp, err := client.Account.Configure(t.Context(), req)
		if err != nil {
			t.Errorf("failed to set account configuration: %v", err)
		}
		debugResponse(resp)
	})

	t.Run("changes", func(t *testing.T) {
		var transactionID TransactionID
		id, err := strconv.Atoi(lastTransactionID)
		if err != nil {
			t.Errorf("failed to parse last transaction id: %v", err)
		}
		if id > 10 {
			transactionID = strconv.Itoa(id - 10)
		} else {
			transactionID = lastTransactionID
		}
		resp, err := client.Account.Changes(t.Context(), transactionID)
		if err != nil {
			t.Errorf("failed to get account changes: %v", err)
		}
		debugResponse(resp)
	})
}
