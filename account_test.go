package oanda

import (
	"fmt"
	"testing"
)

func TestClient_AccountsList(t *testing.T) {
	client := setupClientWithoutAccountID(t)
	resp, err := client.Account.List(t.Context())
	if err != nil {
		t.Errorf("failed to list accounts: %v", err)
	}
	debugResponse(resp)
}

func TestClient_AccountDetails(t *testing.T) {
	client := setupClient(t)
	resp, err := client.Account.Details(t.Context())
	if err != nil {
		t.Errorf("failed to get account details: %v", err)
	}
	fmt.Printf("%+v\n", resp)
	debugResponse(resp.Account)
}

func TestClient_AccountSummary(t *testing.T) {
	client := setupClient(t)
	resp, err := client.Account.Summary(t.Context())
	if err != nil {
		t.Errorf("failed to get account summary: %v", err)
	}
	debugResponse(resp)
}

func TestClient_AccountConfiguration(t *testing.T) {
	client := setupClient(t)
	req := NewAccountConfigureRequest().SetAlias("TestAlias")
	resp, err := client.Account.Configure(t.Context(), req)
	if err != nil {
		t.Errorf("failed to set account configuration: %v", err)
	}
	debugResponse(resp)
}

func TestClient_AccountChanges(t *testing.T) {
	client := setupClient(t)
	transactionID := "500"
	resp, err := client.Account.Changes(t.Context(), transactionID)
	if err != nil {
		t.Errorf("failed to get account changes: %v", err)
	}
	debugResponse(resp)
}
