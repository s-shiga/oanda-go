package oanda

import "testing"

func TestClient_AccountsList(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	accounts, err := client.AccountList()
	if err != nil {
		t.Errorf("failed to list accounts: %v", err)
	}
	for _, account := range accounts {
		t.Log(account)
	}
}
