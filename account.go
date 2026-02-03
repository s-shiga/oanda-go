package oanda

import (
	"bytes"
	"context"
	"encoding/json"
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

type AccountListResponse struct {
	Accounts []AccountProperties `json:"accounts"`
}

// AccountList retrieves the list of accounts authorized for the provided token.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts
//
// Returns:
//   - []AccountProperties: A slice of account properties for all authorized accounts.
//   - error: An error if the request fails or response cannot be decoded.
//
// Reference: https://developer.oanda.com/rest-live-v20/account-ep/#collapse_endpoint_1
func (c *Client) AccountList(ctx context.Context) (*AccountListResponse, error) {
	httpResp, err := c.sendGetRequest(ctx, "/v3/accounts", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var resp AccountListResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type AccountDetailsResponse struct {
	Account           Account       `json:"account"`
	LastTransactionID TransactionID `json:"lastTransactionID"`
}

// AccountDetails retrieves the full details for a single Account.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}
//
// Parameters:
//   - ctx: Context for the request.
//   - id: The Account identifier to retrieve details for.
//
// Returns:
//   - *Account: Full account details including open trades, positions, and pending orders.
//   - TransactionID: The ID of the most recent transaction created for the account.
//   - error: An error if the request fails or response cannot be decoded.
//
// Reference: https://developer.oanda.com/rest-live-v20/account-ep/#collapse_endpoint_2
func (c *Client) AccountDetails(ctx context.Context, id AccountID) (*AccountDetailsResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v", id)
	httpResp, err := c.sendGetRequest(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var resp AccountDetailsResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type AccountSummaryResponse struct {
	Account           AccountSummary `json:"account"`
	LastTransactionID TransactionID  `json:"lastTransactionID"`
}

// AccountSummary retrieves a summary for a single Account.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/summary
//
// Unlike AccountDetails, this method does not include full pending Order, open Trade,
// and Position representations, making it more lightweight for cases where only
// account-level summary information is needed.
//
// Parameters:
//   - ctx: Context for the request.
//   - id: The Account identifier to retrieve the summary for.
//
// Returns:
//   - *AccountSummary: Summary account information without full trade/position details.
//   - TransactionID: The ID of the most recent transaction created for the account.
//   - error: An error if the request fails or response cannot be decoded.
//
// Reference: https://developer.oanda.com/rest-live-v20/account-ep/#collapse_endpoint_3
func (c *Client) AccountSummary(ctx context.Context, id AccountID) (*AccountSummaryResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/summary", id)
	httpResp, err := c.sendGetRequest(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var resp AccountSummaryResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type AccountInstrumentsResponse struct {
	Instruments       []Instrument  `json:"instruments"`
	LastTransactionID TransactionID `json:"lastTransactionID"`
}

// AccountInstruments retrieves the list of tradeable instruments for the given Account.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/instruments
//
// The instruments returned are those that can be traded using the specified account.
// You can optionally filter the results by providing specific instrument names.
//
// Parameters:
//   - ctx: Context for the request.
//   - id: The Account identifier.
//   - instruments: Optional list of instrument names to filter. If empty, all tradeable
//     instruments are returned.
//
// Returns:
//   - []Instrument: A slice of instruments available for trading.
//   - TransactionID: The ID of the most recent transaction created for the account.
//   - error: An error if the request fails or response cannot be decoded.
//
// Reference: https://developer.oanda.com/rest-live-v20/account-ep/#collapse_endpoint_4
func (c *Client) AccountInstruments(ctx context.Context, id AccountID, instruments ...InstrumentName) (*AccountInstrumentsResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/instruments", id)
	v := url.Values{}
	if len(instruments) != 0 {
		v.Set("instruments", strings.Join(instruments, ","))
	}
	httpResp, err := c.sendGetRequest(ctx, path, v)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var resp AccountInstrumentsResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type AccountConfigurationRequest struct {
	Alias      string        `json:"alias"`
	MarginRate DecimalNumber `json:"marginRate"`
}

func NewAccountConfigurationRequest() *AccountConfigurationRequest {
	return &AccountConfigurationRequest{}
}

func (r *AccountConfigurationRequest) SetAlias(alias string) *AccountConfigurationRequest {
	r.Alias = alias
	return r
}

func (r *AccountConfigurationRequest) SetMarginRate(marginRate DecimalNumber) *AccountConfigurationRequest {
	r.MarginRate = marginRate
	return r
}

type AccountConfigurationResponse struct {
	ClientConfigureTransaction ClientConfigureTransaction `json:"clientConfigureTransaction"`
	LastTransactionID          TransactionID              `json:"lastTransactionID"`
}

func (c *Client) AccountConfiguration(ctx context.Context, accountID AccountID, req *AccountConfigurationRequest) (*AccountConfigurationResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/configuration", accountID)
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	reqBody := bytes.NewBuffer(jsonReq)
	resp, err := c.sendPatchRequest(ctx, path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to send PATCH request: %w", err)
	}
	var accountConfigurationResp AccountConfigurationResponse
	if err := decodeResponse(resp, &accountConfigurationResp); err != nil {
		return nil, err
	}
	return &accountConfigurationResp, nil
}

type AccountChangesResponse struct {
	Changes           AccountChanges      `json:"changes"`
	State             AccountChangesState `json:"state"`
	LastTransactionID TransactionID       `json:"lastTransactionID"`
}

// AccountChanges retrieves the changes to an Account's Orders, Trades, and Positions
// since a specified TransactionID.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/changes
//
// This endpoint is useful for polling-based synchronization. By tracking the
// LastTransactionID from each response, you can efficiently fetch only the changes
// that have occurred since your last request.
//
// Parameters:
//   - ctx: Context for the request.
//   - id: The Account identifier.
//   - since: The TransactionID to get Account changes since. This should typically
//     be the LastTransactionID returned from a previous call.
//
// Returns:
//   - *AccountChanges: The changes to orders, trades, and positions since the given transaction.
//   - *AccountChangesState: The current price-dependent state of the account.
//   - TransactionID: The ID of the most recent transaction created for the account.
//   - error: An error if the request fails or response cannot be decoded.
//
// Reference: https://developer.oanda.com/rest-live-v20/account-ep/#collapse_endpoint_6
func (c *Client) AccountChanges(ctx context.Context, id AccountID, since TransactionID) (*AccountChangesResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/changes", id)
	v := url.Values{}
	v.Set("sinceTransactionID", since)
	httpResp, err := c.sendGetRequest(ctx, path, v)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var resp AccountChangesResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
