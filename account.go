package oanda

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type AccountID string

type AccountProperties struct {
	ID           AccountID `json:"id"`
	MT4AccountID int       `json:"mt4AccountID"`
	Tags         []string  `json:"tags"`
}

func (c *Client) AccountList() ([]AccountProperties, error) {
	resp, err := c.sendGetRequest("/v3/accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer closeBody(resp)

	slog.Info(resp.Status)
	if resp.StatusCode != http.StatusOK {
		return nil, decodeErrorResponse(resp)
	}

	accountsResp := struct {
		Accounts []AccountProperties `json:"accounts"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&accountsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}
	return accountsResp.Accounts, nil
}
