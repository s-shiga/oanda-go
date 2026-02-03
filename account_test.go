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
	resp, err := client.AccountList(t.Context())
	if err != nil {
		t.Errorf("failed to list accounts: %v", err)
	}
	for _, account := range resp.Accounts {
		t.Log(account)
	}
}

func TestClient_AccountDetails(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	resp, err := client.AccountDetails(t.Context(), accountID)
	if err != nil {
		t.Errorf("failed to get account details: %v", err)
	}
	t.Logf("response: %#v", resp)
}

func TestClient_AccountSummary(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	resp, err := client.AccountSummary(t.Context(), accountID)
	if err != nil {
		t.Errorf("failed to get account summary: %v", err)
	}
	t.Logf("response: %#v", resp)
}

func TestClient_AccountInstruments(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	resp, err := client.AccountInstruments(t.Context(), accountID, "EUR_USD", "USD_JPY")
	if err != nil {
		t.Errorf("failed to get account instruments: %v", err)
	}
	t.Logf("response: %#v", resp)
}

func TestClient_AccountConfiguration(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	req := NewAccountConfigurationRequest().SetAlias("TestAlias")
	resp, err := client.AccountConfiguration(t.Context(), accountID, req)
	if err != nil {
		t.Errorf("failed to set account configuration: %v", err)
	}
	t.Logf("%#v", resp)
}

func TestClient_AccountChanges(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	transactionID := "421"
	resp, err := client.AccountChanges(t.Context(), accountID, transactionID)
	if err != nil {
		t.Errorf("failed to get account changes: %v", err)
	}
	t.Logf("response: %#v", resp)
}
