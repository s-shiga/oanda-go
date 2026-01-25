package oanda

import (
	"os"
	"testing"
)

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

func TestClient_AccountDetails(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	accountID, ok := os.LookupEnv("OANDA_ACCOUNT_ID")
	if !ok {
		t.Skip("OANDA_ACCOUNT_ID environment variable not set")
	}
	account, lastTransactionID, err := client.AccountDetails(AccountID(accountID))
	if err != nil {
		t.Errorf("failed to get account details: %v", err)
	}
	t.Logf("%#v", account)
	t.Logf("lastTransactionID: %v", lastTransactionID)
}
