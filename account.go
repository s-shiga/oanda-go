package oanda

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type AccountID string

type Account struct {
	ID              AccountID `json:"id"`
	Alias           string    `json:"alias"`
	Currency        `json:"currency"`
	CreatedByUserID int      `json:"createdByUserID"`
	CreatedTime     DateTime `json:"createdTime"`
	// Deprecated: Will be removed in a future API update
	GuaranteedStopLossOrderParameters bool          `json:"guaranteedStopLossOrderParameters"`
	ResettablePLTime                  DateTime      `json:"resettablePLTime"`
	MarginRate                        DecimalNumber `json:"marginRate"`
	OpenTradeCount                    int           `json:"openTradeCount"`
	OpenPositionCount                 int           `json:"openPositionCount"`
	PendingOrderCount                 int           `json:"pendingOrderCount"`
	HedgingEnabled                    bool          `json:"hedgingEnabled"`
	UnrealizedPL                      AccountUnits  `json:"unrealizedPL"`
	NAV                               AccountUnits  `json:"NAV"`
	MarginUsed                        AccountUnits  `json:"marginUsed"`
	MarginAvailable                   AccountUnits  `json:"marginAvailable"`
	PositionValue                     AccountUnits  `json:"positionValue"`
	MarginCloseoutUnrealizedPL        AccountUnits  `json:"marginCloseoutUnrealizedPL"`
}

type AccountProperties struct {
	ID           AccountID `json:"id"`
	MT4AccountID int       `json:"mt4AccountID"`
	Tags         []string  `json:"tags"`
}

type GuaranteedStopLossOrderParameters struct {
	MutabilityMarketOpen   GuaranteedStopLossOrderMutability `json:"mutabilityMarketOpen"`
	MutabilityMarketHalted GuaranteedStopLossOrderMutability `json:"mutabilityMarketHalted"`
}

type GuaranteedStopLossOrderMode string

const (
	GuaranteedStopLossOrderModeDisabled GuaranteedStopLossOrderMode = "DISABLED"
	GuaranteedStopLossOrderModeAllowed  GuaranteedStopLossOrderMode = "ALLOWED"
	GuaranteedStopLossOrderModeRequired GuaranteedStopLossOrderMode = "REQUIRED"
)

type GuaranteedStopLossOrderMutability string

const (
	GuaranteedStopLossOrderMutabilityFixed          GuaranteedStopLossOrderMutability = "FIXED"
	GuaranteedStopLossOrderMutabilityReplaceable    GuaranteedStopLossOrderMutability = "REPLACEABLE"
	GuaranteedStopLossOrderMutabilityCancelable     GuaranteedStopLossOrderMutability = "CANCELABLE"
	GuaranteedStopLossOrderMutabilityPriceWidenOnly GuaranteedStopLossOrderMutability = "PRICE_WIDEN_ONLY"
)

type AccountSummary struct {
	ID                                AccountID `json:"id"`
	Alias                             string    `json:"alias"`
	Currency                          Currency  `json:"currency"`
	CreatedByUserID                   int       `json:"createdByUserID"`
	CreatedTime                       DateTime  `json:"createdTime"`
	GuaranteedStopLossOrderParameters `json:"guaranteedStopLossOrderParameters"`
	GuaranteedStopLossOrderMode       `json:"guaranteedStopLossOrderMode"`
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

func (c *Client) AccountDetails(id AccountID) (*Account, string, error) {
	resp, err := c.sendGetRequest(fmt.Sprintf("/v3/accounts/%v", id))
	if err != nil {
		return nil, "", fmt.Errorf("failed to send request: %w", err)
	}
	defer closeBody(resp)

	slog.Info(resp.Status)
	if resp.StatusCode != http.StatusOK {
		return nil, "", decodeErrorResponse(resp)
	}

	accountsDetailsResp := struct {
		Account           Account `json:"account"`
		LastTransactionID string  `json:"lastTransactionID"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&accountsDetailsResp); err != nil {
		return nil, "", fmt.Errorf("failed to decode response body: %w", err)
	}
	return &accountsDetailsResp.Account, accountsDetailsResp.LastTransactionID, nil
}
