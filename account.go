package oanda

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Definitions https://developer.oanda.com/rest-live-v20/account-df/

// AccountID is the string representation of an Account Identifier.
type AccountID string

// Account is the full representation of a client's Account. This includes full open Trade,
// open Position and pending Order representation.
type Account struct {
	// ID is the Account's identifier.
	ID AccountID `json:"id"`
	// Alias is the client-assigned alias for the Account. Only provided if the Account has an alias.
	Alias string `json:"alias"`
	// Currency is the home currency of the Account.
	Currency Currency `json:"currency"`
	// CreatedByUserID is the ID of the user that created the Account.
	CreatedByUserID int `json:"createdByUserID"`
	// CreatedTime is the date/time when the Account was created.
	CreatedTime DateTime `json:"createdTime"`
	// GuaranteedStopLossOrderParameters contains the current mutability and hedging settings
	// related to guaranteed Stop Loss orders.
	GuaranteedStopLossOrderParameters GuaranteedStopLossOrderParameters `json:"guaranteedStopLossOrderParameters"`
	// GuaranteedStopLossOrderMode describes the guaranteed Stop Loss Order mode of the Account.
	GuaranteedStopLossOrderMode GuaranteedStopLossOrderMode `json:"guaranteedStopLossOrderMode"`
	// ResettablePLTime is the date/time that the Account's resettablePL was last reset.
	ResettablePLTime DateTime `json:"resettablePLTime"`
	// MarginRate is the current financing margin rate used for financing calculations for the Account.
	MarginRate DecimalNumber `json:"marginRate"`
	// OpenTradeCount is the number of Trades currently open in the Account.
	OpenTradeCount int `json:"openTradeCount"`
	// OpenPositionCount is the number of Positions currently open in the Account.
	OpenPositionCount int `json:"openPositionCount"`
	// PendingOrderCount is the number of Orders currently pending in the Account.
	PendingOrderCount int `json:"pendingOrderCount"`
	// HedgingEnabled indicates whether the Account is permitted to create hedged Trades.
	HedgingEnabled bool `json:"hedgingEnabled"`
	// UnrealizedPL is the total unrealized profit/loss for all Trades currently open in the Account.
	UnrealizedPL AccountUnits `json:"unrealizedPL"`
	// NAV is the net asset value of the Account. Equal to Account balance + unrealizedPL.
	NAV AccountUnits `json:"NAV"`
	// MarginUsed is the total amount of margin used by the Account.
	MarginUsed AccountUnits `json:"marginUsed"`
	// MarginAvailable is the total margin available for the Account.
	MarginAvailable AccountUnits `json:"marginAvailable"`
	// PositionValue is the value of the Account's open positions represented in the Account's
	// home currency.
	PositionValue AccountUnits `json:"positionValue"`
	// MarginCloseoutUnrealizedPL is the Account's margin closeout unrealized PL.
	MarginCloseoutUnrealizedPL AccountUnits `json:"marginCloseoutUnrealizedPL"`
	// MarginCloseoutNAV is the Account's margin closeout NAV.
	MarginCloseoutNAV AccountUnits `json:"marginCloseoutNAV"`
	// MarginCloseoutMarginUsed is the Account's margin closeout margin used.
	MarginCloseoutMarginUsed AccountUnits `json:"marginCloseoutMarginUsed"`
	// MarginCloseoutPercent is the Account's margin closeout percentage. When this value is 1.0
	// or above the Account is in a margin closeout situation.
	MarginCloseoutPercent DecimalNumber `json:"marginCloseoutPercent"`
	// MarginCloseoutPositionValue is the value of the Account's open positions as used for margin
	// closeout calculations represented in the Account's home currency.
	MarginCloseoutPositionValue DecimalNumber `json:"marginCloseoutPositionValue"`
	// WithdrawalLimit is the current WithdrawalLimit for the account which will be zero or a
	// positive value indicating how much can be withdrawn from the account.
	WithdrawalLimit AccountUnits `json:"withdrawalLimit"`
	// MarginCallMarginUsed is the Account's margin call margin used.
	MarginCallMarginUsed AccountUnits `json:"marginCallMarginUsed"`
	// MarginCallPercent is the Account's margin call percentage. When this value is 1.0 or above
	// the Account is in a margin call situation.
	MarginCallPercent DecimalNumber `json:"marginCallPercent"`
	// Balance is the current balance of the Account.
	Balance AccountUnits `json:"balance"`
	// PL is the total profit/loss realized over the lifetime of the Account.
	PL AccountUnits `json:"pl"`
	// ResettablePL is the total realized profit/loss for the Account since it was last reset by
	// the client.
	ResettablePL AccountUnits `json:"resettablePL"`
	// Financing is the total amount of financing paid/collected over the lifetime of the Account.
	Financing AccountUnits `json:"financing"`
	// Commission is the total amount of commission paid over the lifetime of the Account.
	Commission AccountUnits `json:"commission"`
	// DividendAdjustment is the total amount of dividend adjustment paid or collected over the
	// lifetime of the Account in the Account's home currency.
	DividendAdjustment AccountUnits `json:"dividendAdjustment"`
	// GuaranteedExecutionFees is the total amount of fees charged over the lifetime of the Account
	// for the execution of guaranteed Stop Loss Orders.
	GuaranteedExecutionFees AccountUnits `json:"guaranteedExecutionFees"`
	// MarginCallEnterTime is the date/time when the Account entered a margin call state. Only
	// provided if the Account is in a margin call.
	MarginCallEnterTime DateTime `json:"marginCallEnterTime"`
	// MarginCallExtensionCount is the number of times that the Account's current margin call was
	// extended.
	MarginCallExtensionCount int `json:"marginCallExtensionCount"`
	// LastMarginCallExtensionTime is the date/time of the Account's last margin call extension.
	LastMarginCallExtensionTime DateTime `json:"lastMarginCallExtensionTime"`
	// LastTransactionID is the ID of the last Transaction created for the Account.
	LastTransactionID TransactionID `json:"lastTransactionID"`
	// Trades is the details of the Trades currently open in the Account.
	Trades []TradeSummary `json:"trades"`
	// Positions is the details of the Positions currently open in the Account.
	Positions []Position `json:"positions"`
	// Orders is the details of the Orders currently pending in the Account.
	Orders []Order `json:"orders"`
}

// AccountProperties contains properties related to an Account.
type AccountProperties struct {
	// ID is the Account's identifier.
	ID AccountID `json:"id"`
	// MT4AccountID is the Account's associated MT4 Account ID. This field will not be present
	// if the Account is not an MT4 account.
	MT4AccountID int `json:"mt4AccountID"`
	// Tags is the Account's tags. Tags are user-defined labels that can be applied to Accounts.
	Tags []string `json:"tags"`
}

// GuaranteedStopLossOrderParameters contains the current mutability and hedging settings related
// to guaranteed Stop Loss orders.
type GuaranteedStopLossOrderParameters struct {
	// MutabilityMarketOpen indicates whether or not guaranteed Stop Loss Orders can be cancelled
	// or have their price changed while the market is open.
	MutabilityMarketOpen GuaranteedStopLossOrderMutability `json:"mutabilityMarketOpen"`
	// MutabilityMarketHalted indicates whether or not guaranteed Stop Loss Orders can be cancelled
	// or have their price changed while the market is halted.
	MutabilityMarketHalted GuaranteedStopLossOrderMutability `json:"mutabilityMarketHalted"`
}

// GuaranteedStopLossOrderMode describes the guaranteed Stop Loss Order mode of an Account.
type GuaranteedStopLossOrderMode string

const (
	// GuaranteedStopLossOrderModeDisabled means the Account is not permitted to create guaranteed
	// Stop Loss Orders.
	GuaranteedStopLossOrderModeDisabled GuaranteedStopLossOrderMode = "DISABLED"
	// GuaranteedStopLossOrderModeAllowed means the Account is able, but not required to have
	// guaranteed Stop Loss Orders for open Trades.
	GuaranteedStopLossOrderModeAllowed GuaranteedStopLossOrderMode = "ALLOWED"
	// GuaranteedStopLossOrderModeRequired means the Account is required to have guaranteed Stop
	// Loss Orders for all open Trades.
	GuaranteedStopLossOrderModeRequired GuaranteedStopLossOrderMode = "REQUIRED"
)

// GuaranteedStopLossOrderMutability describes the mutability of guaranteed Stop Loss Orders.
type GuaranteedStopLossOrderMutability string

const (
	// GuaranteedStopLossOrderMutabilityFixed means the guaranteed Stop Loss Order cannot be replaced
	// or cancelled.
	GuaranteedStopLossOrderMutabilityFixed GuaranteedStopLossOrderMutability = "FIXED"
	// GuaranteedStopLossOrderMutabilityReplaceable means the guaranteed Stop Loss Order can only
	// be replaced, not cancelled.
	GuaranteedStopLossOrderMutabilityReplaceable GuaranteedStopLossOrderMutability = "REPLACEABLE"
	// GuaranteedStopLossOrderMutabilityCancelable means the guaranteed Stop Loss Order can be
	// either cancelled or replaced.
	GuaranteedStopLossOrderMutabilityCancelable GuaranteedStopLossOrderMutability = "CANCELABLE"
	// GuaranteedStopLossOrderMutabilityPriceWidenOnly means the guaranteed Stop Loss Order can
	// only be replaced to widen the gap from the current price.
	GuaranteedStopLossOrderMutabilityPriceWidenOnly GuaranteedStopLossOrderMutability = "PRICE_WIDEN_ONLY"
)

// AccountSummary is a summary representation of a client's Account. The AccountSummary does not
// provide to full specification of pending Orders, open Trades and Positions.
type AccountSummary struct {
	// ID is the Account's identifier.
	ID AccountID `json:"id"`
	// Alias is the client-assigned alias for the Account. Only provided if the Account has an alias.
	Alias string `json:"alias"`
	// Currency is the home currency of the Account.
	Currency Currency `json:"currency"`
	// CreatedByUserID is the ID of the user that created the Account.
	CreatedByUserID int `json:"createdByUserID"`
	// CreatedTime is the date/time when the Account was created.
	CreatedTime DateTime `json:"createdTime"`
	// GuaranteedStopLossOrderParameters contains the current mutability and hedging settings
	// related to guaranteed Stop Loss orders.
	GuaranteedStopLossOrderParameters GuaranteedStopLossOrderParameters `json:"guaranteedStopLossOrderParameters"`
	// GuaranteedStopLossOrderMode describes the guaranteed Stop Loss Order mode of the Account.
	GuaranteedStopLossOrderMode GuaranteedStopLossOrderMode `json:"guaranteedStopLossOrderMode"`
	// ResettablePLTime is the date/time that the Account's resettablePL was last reset.
	ResettablePLTime DateTime `json:"resettablePLTime"`
	// MarginRate is the current financing margin rate used for financing calculations for the Account.
	MarginRate DecimalNumber `json:"marginRate"`
	// OpenTradeCount is the number of Trades currently open in the Account.
	OpenTradeCount int `json:"openTradeCount"`
	// OpenPositionCount is the number of Positions currently open in the Account.
	OpenPositionCount int `json:"openPositionCount"`
	// PendingOrderCount is the number of Orders currently pending in the Account.
	PendingOrderCount int `json:"pendingOrderCount"`
	// HedgingEnabled indicates whether the Account is permitted to create hedged Trades.
	HedgingEnabled bool `json:"hedgingEnabled"`
	// UnrealizedPL is the total unrealized profit/loss for all Trades currently open in the Account.
	UnrealizedPL AccountUnits `json:"unrealizedPL"`
	// NAV is the net asset value of the Account. Equal to Account balance + unrealizedPL.
	NAV AccountUnits `json:"NAV"`
	// MarginUsed is the total amount of margin used by the Account.
	MarginUsed AccountUnits `json:"marginUsed"`
	// MarginAvailable is the total margin available for the Account.
	MarginAvailable AccountUnits `json:"marginAvailable"`
	// PositionValue is the value of the Account's open positions represented in the Account's
	// home currency.
	PositionValue AccountUnits `json:"positionValue"`
	// MarginCloseoutUnrealizedPL is the Account's margin closeout unrealized PL.
	MarginCloseoutUnrealizedPL AccountUnits `json:"marginCloseoutUnrealizedPL"`
	// MarginCloseoutNAV is the Account's margin closeout NAV.
	MarginCloseoutNAV AccountUnits `json:"marginCloseoutNAV"`
	// MarginCloseoutMarginUsed is the Account's margin closeout margin used.
	MarginCloseoutMarginUsed AccountUnits `json:"marginCloseoutMarginUsed"`
	// MarginCloseoutPercent is the Account's margin closeout percentage. When this value is 1.0
	// or above the Account is in a margin closeout situation.
	MarginCloseoutPercent DecimalNumber `json:"marginCloseoutPercent"`
	// MarginCloseoutPositionValue is the value of the Account's open positions as used for margin
	// closeout calculations represented in the Account's home currency.
	MarginCloseoutPositionValue DecimalNumber `json:"marginCloseoutPositionValue"`
	// WithdrawalLimit is the current WithdrawalLimit for the account which will be zero or a
	// positive value indicating how much can be withdrawn from the account.
	WithdrawalLimit AccountUnits `json:"withdrawalLimit"`
	// MarginCallMarginUsed is the Account's margin call margin used.
	MarginCallMarginUsed AccountUnits `json:"marginCallMarginUsed"`
	// MarginCallPercent is the Account's margin call percentage. When this value is 1.0 or above
	// the Account is in a margin call situation.
	MarginCallPercent DecimalNumber `json:"marginCallPercent"`
	// Balance is the current balance of the Account.
	Balance AccountUnits `json:"balance"`
	// PL is the total profit/loss realized over the lifetime of the Account.
	PL AccountUnits `json:"pl"`
	// ResettablePL is the total realized profit/loss for the Account since it was last reset by
	// the client.
	ResettablePL AccountUnits `json:"resettablePL"`
	// Financing is the total amount of financing paid/collected over the lifetime of the Account.
	Financing AccountUnits `json:"financing"`
	// Commission is the total amount of commission paid over the lifetime of the Account.
	Commission AccountUnits `json:"commission"`
	// DividendAdjustment is the total amount of dividend adjustment paid or collected over the
	// lifetime of the Account in the Account's home currency.
	DividendAdjustment AccountUnits `json:"dividendAdjustment"`
	// GuaranteedExecutionFees is the total amount of fees charged over the lifetime of the Account
	// for the execution of guaranteed Stop Loss Orders.
	GuaranteedExecutionFees AccountUnits `json:"guaranteedExecutionFees"`
	// MarginCallEnterTime is the date/time when the Account entered a margin call state. Only
	// provided if the Account is in a margin call.
	MarginCallEnterTime DateTime `json:"marginCallEnterTime"`
	// MarginCallExtensionCount is the number of times that the Account's current margin call was
	// extended.
	MarginCallExtensionCount int `json:"marginCallExtensionCount"`
	// LastMarginCallExtensionTime is the date/time of the Account's last margin call extension.
	LastMarginCallExtensionTime DateTime `json:"lastMarginCallExtensionTime"`
	// LastTransactionID is the ID of the last Transaction created for the Account.
	LastTransactionID TransactionID `json:"lastTransactionID"`
}

// AccountChangesState represents the price-dependent state of an Account.
type AccountChangesState struct {
	// UnrealizedPL is the total unrealized profit/loss for all Trades currently open in the Account.
	UnrealizedPL AccountUnits `json:"unrealizedPL"`
	// NAV is the net asset value of the Account. Equal to Account balance + unrealizedPL.
	NAV AccountUnits `json:"NAV"`
	// MarginUsed is the total amount of margin used by the Account.
	MarginUsed AccountUnits `json:"marginUsed"`
	// MarginAvailable is the total margin available for the Account.
	MarginAvailable AccountUnits `json:"marginAvailable"`
	// PositionValue is the value of the Account's open positions represented in the Account's
	// home currency.
	PositionValue AccountUnits `json:"positionValue"`
	// MarginCloseoutUnrealizedPL is the Account's margin closeout unrealized PL.
	MarginCloseoutUnrealizedPL AccountUnits `json:"marginCloseoutUnrealizedPL"`
	// MarginCloseoutNAV is the Account's margin closeout NAV.
	MarginCloseoutNAV AccountUnits `json:"marginCloseoutNAV"`
	// MarginCloseoutMarginUsed is the Account's margin closeout margin used.
	MarginCloseoutMarginUsed AccountUnits `json:"marginCloseoutMarginUsed"`
	// MarginCloseoutPercent is the Account's margin closeout percentage. When this value is 1.0
	// or above the Account is in a margin closeout situation.
	MarginCloseoutPercent DecimalNumber `json:"marginCloseoutPercent"`
	// MarginCloseoutPositionValue is the value of the Account's open positions as used for margin
	// closeout calculations represented in the Account's home currency.
	MarginCloseoutPositionValue DecimalNumber `json:"marginCloseoutPositionValue"`
	// WithdrawalLimit is the current WithdrawalLimit for the account which will be zero or a
	// positive value indicating how much can be withdrawn from the account.
	WithdrawalLimit AccountUnits `json:"withdrawalLimit"`
	// MarginCallMarginUsed is the Account's margin call margin used.
	MarginCallMarginUsed AccountUnits `json:"marginCallMarginUsed"`
	// MarginCallPercent is the Account's margin call percentage. When this value is 1.0 or above
	// the Account is in a margin call situation.
	MarginCallPercent DecimalNumber `json:"marginCallPercent"`
	// Balance is the current balance of the Account.
	Balance AccountUnits `json:"balance"`
	// PL is the total profit/loss realized over the lifetime of the Account.
	PL AccountUnits `json:"pl"`
	// ResettablePL is the total realized profit/loss for the Account since it was last reset by
	// the client.
	ResettablePL AccountUnits `json:"resettablePL"`
	// Financing is the total amount of financing paid/collected over the lifetime of the Account.
	Financing AccountUnits `json:"financing"`
	// Commission is the total amount of commission paid over the lifetime of the Account.
	Commission AccountUnits `json:"commission"`
	// DividendAdjustment is the total amount of dividend adjustment paid or collected over the
	// lifetime of the Account in the Account's home currency.
	DividendAdjustment AccountUnits `json:"dividendAdjustment"`
	// GuaranteedExecutionFees is the total amount of fees charged over the lifetime of the Account
	// for the execution of guaranteed Stop Loss Orders.
	GuaranteedExecutionFees AccountUnits `json:"guaranteedExecutionFees"`
}

// AccountChanges represents the changes to an Account's Orders, Trades and Positions since a
// specified Account TransactionID.
type AccountChanges struct {
	// OrdersCreated is the array of Orders created since the specified transaction ID.
	OrdersCreated []Order `json:"ordersCreated"`
	// OrdersCancelled is the array of Orders cancelled since the specified transaction ID.
	OrdersCancelled []Order `json:"ordersCancelled"`
	// OrdersFilled is the array of Orders filled since the specified transaction ID.
	OrdersFilled []Order `json:"ordersFilled"`
	// OrdersTriggered is the array of Orders triggered since the specified transaction ID.
	OrdersTriggered []Order `json:"ordersTriggered"`
	// TradesOpened is the array of Trades opened since the specified transaction ID.
	TradesOpened []TradeSummary `json:"tradesOpened"`
	// TradesReduced is the array of Trades reduced since the specified transaction ID.
	TradesReduced []TradeSummary `json:"tradesReduced"`
	// TradesClosed is the array of Trades closed since the specified transaction ID.
	TradesClosed []TradeSummary `json:"tradesClosed"`
	// Positions is the array of Positions that have changed since the specified transaction ID.
	Positions []Position `json:"positions"`
	// Transactions is the array of Transactions that have been generated since the specified
	// transaction ID.
	Transactions []Transaction `json:"transactions"`
}

// AccountFinancingMode describes the financing mode of an Account.
type AccountFinancingMode string

const (
	// AccountFinancingModeNoFinancing means no financing is paid/charged for open Trades in the Account.
	AccountFinancingModeNoFinancing AccountFinancingMode = "NO_FINANCING"
	// AccountFinancingModeSecondBySecond means second-by-second financing is paid/charged for open
	// Trades in the Account, both daily and when the Trade is closed.
	AccountFinancingModeSecondBySecond AccountFinancingMode = "SECOND_BY_SECOND"
	// AccountFinancingModeDaily means a full day's worth of financing is paid/charged for open
	// Trades in the Account daily at 5pm New York time.
	AccountFinancingModeDaily AccountFinancingMode = "DAILY"
)

// PositionAggregationMode describes how Positions are aggregated for margin closeout purposes.
type PositionAggregationMode string

const (
	// PositionAggregationModeAbsoluteSum means the Position value or margin for each side (long
	// and short) of the Position are computed independently and added together.
	PositionAggregationModeAbsoluteSum PositionAggregationMode = "ABSOLUTE_SUM"
	// PositionAggregationModeMaximalSide means the Position value or margin for each side (long
	// and short) of the Position are computed independently. The Position value or margin chosen
	// is the maximal absolute value of the two.
	PositionAggregationModeMaximalSide PositionAggregationMode = "MAXIMAL_SIDE"
	// PositionAggregationModeNetSum means the units for each side (long and short) of the Position
	// are netted together and the resulting value (long or short) is used to compute the Position
	// value or margin.
	PositionAggregationModeNetSum PositionAggregationMode = "NET_SUM"
)

// Endpoints https://developer.oanda.com/rest-live-v20/account-ep/

func (c *Client) AccountList(ctx context.Context) ([]AccountProperties, error) {
	resp, err := c.sendGetRequest(ctx, "/v3/accounts", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	accountsResp := struct {
		Accounts []AccountProperties `json:"accounts"`
	}{}
	if err := decodeResponse(resp, &accountsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return accountsResp.Accounts, nil
}

func (c *Client) AccountDetails(ctx context.Context, id AccountID) (*Account, TransactionID, error) {
	resp, err := c.sendGetRequest(ctx, fmt.Sprintf("/v3/accounts/%v", id), nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to send request: %w", err)
	}
	accountsDetailsResp := struct {
		Account           Account       `json:"account"`
		LastTransactionID TransactionID `json:"lastTransactionID"`
	}{}
	if err := decodeResponse(resp, &accountsDetailsResp); err != nil {
		return nil, "", fmt.Errorf("failed to decode response body: %w", err)
	}
	return &accountsDetailsResp.Account, accountsDetailsResp.LastTransactionID, nil
}

func (c *Client) AccountSummary(ctx context.Context, id AccountID) (*AccountSummary, TransactionID, error) {
	resp, err := c.sendGetRequest(ctx, fmt.Sprintf("/v3/accounts/%v/summary", id), nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to send request: %w", err)
	}
	accountsSummaryResp := struct {
		Account           AccountSummary `json:"account"`
		LastTransactionID TransactionID  `json:"lastTransactionID"`
	}{}
	if err := decodeResponse(resp, &accountsSummaryResp); err != nil {
		return nil, "", fmt.Errorf("failed to decode response body: %w", err)
	}
	return &accountsSummaryResp.Account, accountsSummaryResp.LastTransactionID, nil
}

func (c *Client) AccountInstruments(ctx context.Context, id AccountID, instruments ...InstrumentName) ([]Instrument, TransactionID, error) {
	v := url.Values{}
	if len(instruments) != 0 {
		v.Set("instruments", strings.Join(instruments, ","))
	}
	resp, err := c.sendGetRequest(ctx, fmt.Sprintf("/v3/accounts/%v/instruments", id), v)
	if err != nil {
		return nil, "", fmt.Errorf("failed to send request: %w", err)
	}
	println(resp.StatusCode)
	accountsInstrumentsResp := struct {
		Instruments       []Instrument  `json:"instruments"`
		LastTransactionID TransactionID `json:"lastTransactionID"`
	}{}
	if err := decodeResponse(resp, &accountsInstrumentsResp); err != nil {
		return nil, "", fmt.Errorf("failed to decode response body: %w", err)
	}
	return accountsInstrumentsResp.Instruments, accountsInstrumentsResp.LastTransactionID, nil
}
