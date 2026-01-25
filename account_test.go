package oanda

import (
	"os"
	"testing"
)

func setupClient(t *testing.T) *Client {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func setupAccountID(t *testing.T) AccountID {
	accountID, ok := os.LookupEnv("OANDA_ACCOUNT_ID")
	if !ok {
		t.Fatal("OANDA_ACCOUNT_ID not set")
	}
	return AccountID(accountID)
}

func TestClient_AccountsList(t *testing.T) {
	client := setupClient(t)
	accounts, err := client.AccountList()
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
	account, lastTransactionID, err := client.AccountDetails(accountID)
	if err != nil {
		t.Errorf("failed to get account details: %v", err)
	}
	t.Logf("%#v", account)
	t.Logf("lastTransactionID: %v", lastTransactionID)
}

func TestClient_AccountSummary(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	accountSummary, lastTransactionID, err := client.AccountSummary(accountID)
	if err != nil {
		t.Errorf("failed to get account summary: %v", err)
	}
	t.Logf("%#v", accountSummary)
	t.Logf("lastTransactionID: %v", lastTransactionID)
}
