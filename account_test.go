package oanda

import (
	"encoding/json"
	"log/slog"
	"os"
	"testing"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

func setupClientWithoutAccountID(t *testing.T) (client *Client) {
	client, err := NewPracticeClientWithoutAccountID()
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func setupClient(t *testing.T) *Client {
	accountID, ok := os.LookupEnv("OANDA_ACCOUNT_ID_DEMO")
	if !ok {
		t.Fatal("OANDA_ACCOUNT_ID_DEMO not set")
	}
	client, err := NewPracticeClient(accountID)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func debugResponse(resp any) {
	b, _ := json.MarshalIndent(resp, "", "  ")
	slog.Debug(string(b))
}

func TestClient_AccountsList(t *testing.T) {
	client := setupClientWithoutAccountID(t)
	resp, err := client.AccountList(t.Context())
	if err != nil {
		t.Errorf("failed to list accounts: %v", err)
	}
	debugResponse(resp)
}

func TestClient_AccountDetails(t *testing.T) {
	client := setupClient(t)
	resp, err := client.AccountDetails(t.Context())
	if err != nil {
		t.Errorf("failed to get account details: %v", err)
	}
	debugResponse(resp)
}

func TestClient_AccountSummary(t *testing.T) {
	client := setupClient(t)
	resp, err := client.AccountSummary(t.Context())
	if err != nil {
		t.Errorf("failed to get account summary: %v", err)
	}
	debugResponse(resp)
}

func TestClient_AccountInstruments(t *testing.T) {
	client := setupClient(t)
	resp, err := client.AccountInstruments(t.Context(), "EUR_USD", "USD_JPY")
	if err != nil {
		t.Errorf("failed to get account instruments: %v", err)
	}
	debugResponse(resp)
}

func TestClient_AccountConfiguration(t *testing.T) {
	client := setupClient(t)
	req := NewAccountConfigurationRequest().SetAlias("TestAlias")
	resp, err := client.AccountConfiguration(t.Context(), req)
	if err != nil {
		t.Errorf("failed to set account configuration: %v", err)
	}
	debugResponse(resp)
}

func TestClient_AccountChanges(t *testing.T) {
	client := setupClient(t)
	transactionID := "421"
	resp, err := client.AccountChanges(t.Context(), transactionID)
	if err != nil {
		t.Errorf("failed to get account changes: %v", err)
	}
	debugResponse(resp)
}
