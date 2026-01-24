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

func (c *Client) AccountsList() ([]AccountProperties, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.URL+"/v3/accounts", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+c.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	slog.Info(resp.Status)
	if resp.StatusCode != http.StatusOK {
		errResp := struct {
			Message string `json:"errorMessage"`
		}{}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	accountsResp := struct {
		Accounts []AccountProperties `json:"accounts"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&accountsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}
	return accountsResp.Accounts, nil
}
