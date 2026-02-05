package oanda

import (
	"fmt"
	"testing"
)

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
	fmt.Printf("%+v\n", resp)
	debugResponse(resp.Account)
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
