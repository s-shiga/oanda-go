package oanda

import (
	"encoding/json"
	"log/slog"
	"os"
	"testing"
)

func getAPIKey(t *testing.T) string {
	t.Helper()
	apiKey, ok := os.LookupEnv("OANDA_API_KEY_DEMO")
	if !ok {
		t.Skip("OANDA_API_KEY_DEMO not set; skipping integration test")
	}
	return apiKey
}

func getAccountID(t *testing.T) string {
	t.Helper()
	accountID, ok := os.LookupEnv("OANDA_ACCOUNT_ID_DEMO")
	if !ok {
		t.Skip("OANDA_ACCOUNT_ID_DEMO not set; skipping integration test")
	}
	return accountID
}

func setupClientWithoutAccountID(t *testing.T) *Client {
	t.Helper()
	apiKey := getAPIKey(t)
	return NewDemoClient(apiKey)
}

func setupClient(t *testing.T) *Client {
	t.Helper()
	apiKey := getAPIKey(t)
	accountID := getAccountID(t)
	client := NewDemoClient(apiKey, WithAccountID(accountID))
	return client
}

func setupStreamClient(t *testing.T) *StreamClient {
	t.Helper()
	apiKey := getAPIKey(t)
	accountID := getAccountID(t)
	streamClient := NewDemoStreamClient(apiKey, WithAccountID(accountID))
	return streamClient
}

func debugResponse(resp any) {
	b, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Debug(string(b))
}
