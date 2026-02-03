package oanda

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

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
		slog.Debug("AccountList:", "Account", account)
	}
}

func TestClient_AccountDetails(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	resp, err := client.AccountDetails(t.Context(), accountID)
	if err != nil {
		t.Errorf("failed to get account details: %v", err)
	}
	slog.Debug("AccountDetails:", "Account", fmt.Sprintf("%#v", resp.Account))
	slog.Debug("AccountDetails:", "LastTransactionID", resp.LastTransactionID)
}

func TestClient_AccountSummary(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	resp, err := client.AccountSummary(t.Context(), accountID)
	if err != nil {
		t.Errorf("failed to get account summary: %v", err)
	}
	slog.Debug("AccountSummary:", "Account", fmt.Sprintf("%#v", resp.Account))
}

func TestClient_AccountInstruments(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	resp, err := client.AccountInstruments(t.Context(), accountID, "EUR_USD", "USD_JPY")
	if err != nil {
		t.Errorf("failed to get account instruments: %v", err)
	}
	for _, instrument := range resp.Instruments {
		slog.Debug("AccountInstruments:", "instrument", fmt.Sprintf("%#v", instrument))
	}
	slog.Debug("AccountInstruments:", "LastTransactionID", resp.LastTransactionID)
}

func TestClient_AccountConfiguration(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	req := NewAccountConfigurationRequest().SetAlias("TestAlias")
	resp, err := client.AccountConfiguration(t.Context(), accountID, req)
	if err != nil {
		t.Errorf("failed to set account configuration: %v", err)
	}
	slog.Debug("AccountConfiguration:", "AccountConfiguration", fmt.Sprintf("%#v", resp.ClientConfigureTransaction))
	slog.Debug("AccountConfiguration:", "LastTransactionID", resp.LastTransactionID)
}

func TestClient_AccountChanges(t *testing.T) {
	client := setupClient(t)
	accountID := setupAccountID(t)
	transactionID := "421"
	resp, err := client.AccountChanges(t.Context(), accountID, transactionID)
	if err != nil {
		t.Errorf("failed to get account changes: %v", err)
	}
	for _, order := range resp.Changes.OrdersCreated {
		slog.Debug("AccountChanges:", "OrderCreated", fmt.Sprintf("%#v", order))
	}
	for _, order := range resp.Changes.OrdersCancelled {
		slog.Debug("AccountChanges:", "OrderCancelled", fmt.Sprintf("%#v", order))
	}
	for _, order := range resp.Changes.OrdersFilled {
		slog.Debug("AccountChanges:", "OrderFilled", fmt.Sprintf("%#v", order))
	}
	for _, order := range resp.Changes.OrdersTriggered {
		slog.Debug("AccountChanges:", "OrderTriggered", fmt.Sprintf("%#v", order))
	}
	for _, trade := range resp.Changes.TradesOpened {
		slog.Debug("AccountChanges:", "TradeOpened", fmt.Sprintf("%#v", trade))
	}
	for _, trade := range resp.Changes.TradesReduced {
		slog.Debug("AccountChanges:", "TradeReduced", fmt.Sprintf("%#v", trade))
	}
	for _, trade := range resp.Changes.TradesClosed {
		slog.Debug("AccountChanges:", "TradeClosed", fmt.Sprintf("%#v", trade))
	}
	for _, position := range resp.Changes.Positions {
		slog.Debug("AccountChanges:", "Position", fmt.Sprintf("%#v", position))
	}
	for _, transaction := range resp.Changes.Transactions {
		slog.Debug("AccountChanges:", "Transaction", fmt.Sprintf("%#v", transaction))
	}
}
