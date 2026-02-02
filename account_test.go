package oanda

import (
	"os"
	"testing"
)

func setupClient(t *testing.T) *Client {
	client, err := NewPracticeClient()
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func setupAccountID(t *testing.T) AccountID {
	accountID, ok := os.LookupEnv("OANDA_ACCOUNT_ID_DEMO")
	if !ok {
		t.Fatal("OANDA_ACCOUNT_ID_DEMO not set")
	}
	return AccountID(accountID)
}

func TestClient_AccountsList(t *testing.T) {
	client := setupClient(t)
	accounts, err := client.AccountList(t.Context())
	if err != nil {
		t.Errorf("failed to list accounts: %v", err)
	}
	for _, account := range accounts {
		t.Log(account)
	}
}

func TestClient_AccountDetails(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	account, lastTransactionID, err := client.AccountDetails(t.Context(), accountID)
	if err != nil {
		t.Errorf("failed to get account details: %v", err)
	}
	t.Logf("%#v", account)
	t.Logf("lastTransactionID: %v", lastTransactionID)
}

func TestClient_AccountSummary(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	accountSummary, lastTransactionID, err := client.AccountSummary(t.Context(), accountID)
	if err != nil {
		t.Errorf("failed to get account summary: %v", err)
	}
	t.Logf("%#v", accountSummary)
	t.Logf("lastTransactionID: %v", lastTransactionID)
}

func TestClient_AccountInstruments(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	accountInstruments, lastTransactionID, err := client.AccountInstruments(t.Context(), accountID, "EUR_USD", "USD_JPY")
	if err != nil {
		t.Errorf("failed to get account instruments: %v", err)
	}
	t.Logf("%#v", accountInstruments)
	t.Logf("lastTransactionID: %v", lastTransactionID)
}

func TestClient_AccountChanges(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	transactionID := "421"
	accountChanges, accountChangesState, lastTransactionID, err := client.AccountChanges(t.Context(), accountID, transactionID)
	if err != nil {
		t.Errorf("failed to get account changes: %v", err)
	}
	t.Logf("%#v", accountChanges)
	t.Logf("%#v", accountChangesState)
	t.Logf("lastTransactionID: %v", lastTransactionID)
}
