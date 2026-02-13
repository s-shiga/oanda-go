package oanda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// -----------------------------------------------------------------
// Definitions https://developer.oanda.com/rest-live-v20/account-df/
// -----------------------------------------------------------------

// AccountID is the string representation of an Account Identifier.
type AccountID = string

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
	GuaranteedStopLossOrderParameters *GuaranteedStopLossOrderParameters `json:"guaranteedStopLossOrderParameters,omitempty"`
	// GuaranteedStopLossOrderMode describes the guaranteed Stop Loss Order mode of the Account.
	GuaranteedStopLossOrderMode GuaranteedStopLossOrderMode `json:"guaranteedStopLossOrderMode"`
	// ResettablePLTime is the date/time that the Account's resettablePL was last reset.
	ResettablePLTime *DateTime `json:"resettablePLTime,omitempty"`
	// MarginRate is the current financing margin rate used for financing calculations for the Account.
	MarginRate DecimalNumber `json:"marginRate,omitempty"`
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
	MarginCallMarginUsed *AccountUnits `json:"marginCallMarginUsed,omitempty"`
	// MarginCallPercent is the Account's margin call percentage. When this value is 1.0 or above
	// the Account is in a margin call situation.
	MarginCallPercent *DecimalNumber `json:"marginCallPercent,omitempty"`
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
	MarginCallEnterTime *DateTime `json:"marginCallEnterTime,omitempty"`
	// MarginCallExtensionCount is the number of times that the Account's current margin call was
	// extended.
	MarginCallExtensionCount *int `json:"marginCallExtensionCount,omitempty"`
	// LastMarginCallExtensionTime is the date/time of the Account's last margin call extension.
	LastMarginCallExtensionTime *DateTime `json:"lastMarginCallExtensionTime,omitempty"`
	// LastTransactionID is the ID of the last Transaction created for the Account.
	LastTransactionID TransactionID `json:"lastTransactionID"`
	// Trades is the details of the Trades currently open in the Account.
	Trades []TradeSummary `json:"trades,omitempty"`
	// Positions is the details of the Positions currently open in the Account.
	Positions []Position `json:"positions,omitempty"`
	// Orders is the details of the Orders currently pending in the Account.
	Orders []Order `json:"orders,omitempty"`
}

func (a *Account) UnmarshalJSON(b []byte) error {
	type Alias Account

	aux := &struct {
		*Alias
		Orders []json.RawMessage `json:"orders"`
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(b, aux); err != nil {
		return err
	}

	a.Orders = make([]Order, 0, len(a.Orders))

	for _, rawOrder := range aux.Orders {
		order, err := unmarshalOrder(rawOrder)
		if err != nil {
			return err
		}
		a.Orders = append(a.Orders, order)
	}
	return nil
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
	MutabilityMarketOpen *GuaranteedStopLossOrderMutability `json:"mutabilityMarketOpen,omitempty"`
	// MutabilityMarketHalted indicates whether or not guaranteed Stop Loss Orders can be cancelled
	// or have their price changed while the market is halted.
	MutabilityMarketHalted *GuaranteedStopLossOrderMutability `json:"mutabilityMarketHalted,omitempty"`
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
	Alias string `json:"alias,omitempty"`
	// Currency is the home currency of the Account.
	Currency Currency `json:"currency"`
	// CreatedByUserID is the ID of the user that created the Account.
	CreatedByUserID int `json:"createdByUserID"`
	// CreatedTime is the date/time when the Account was created.
	CreatedTime DateTime `json:"createdTime"`
	// GuaranteedStopLossOrderParameters contains the current mutability and hedging settings
	// related to guaranteed Stop Loss orders.
	GuaranteedStopLossOrderParameters *GuaranteedStopLossOrderParameters `json:"guaranteedStopLossOrderParameters,omitempty"`
	// GuaranteedStopLossOrderMode describes the guaranteed Stop Loss Order mode of the Account.
	GuaranteedStopLossOrderMode GuaranteedStopLossOrderMode `json:"guaranteedStopLossOrderMode"`
	// ResettablePLTime is the date/time that the Account's resettablePL was last reset.
	ResettablePLTime *DateTime `json:"resettablePLTime,omitempty"`
	// MarginRate is the current financing margin rate used for financing calculations for the Account.
	MarginRate DecimalNumber `json:"marginRate,omitempty"`
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
	MarginCallMarginUsed AccountUnits `json:"marginCallMarginUsed,omitempty"`
	// MarginCallPercent is the Account's margin call percentage. When this value is 1.0 or above
	// the Account is in a margin call situation.
	MarginCallPercent DecimalNumber `json:"marginCallPercent,omitempty"`
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
	MarginCallEnterTime *DateTime `json:"marginCallEnterTime,omitempty"`
	// MarginCallExtensionCount is the number of times that the Account's current margin call was
	// extended.
	MarginCallExtensionCount int `json:"marginCallExtensionCount"`
	// LastMarginCallExtensionTime is the date/time of the Account's last margin call extension.
	LastMarginCallExtensionTime *DateTime `json:"lastMarginCallExtensionTime,omitempty"`
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

func (c *AccountChanges) UnmarshalJSON(data []byte) error {
	type Alias AccountChanges
	aux := &struct {
		*Alias
		OrdersCreated   []json.RawMessage `json:"ordersCreated"`
		OrdersCancelled []json.RawMessage `json:"ordersCancelled"`
		OrdersFilled    []json.RawMessage `json:"ordersFilled"`
		OrdersTriggered []json.RawMessage `json:"ordersTriggered"`
		Transactions    []json.RawMessage `json:"transactions"`
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	ordersCreated, err := unmarshalOrders(aux.OrdersCreated)
	if err != nil {
		return err
	}
	c.OrdersCreated = ordersCreated
	ordersCancelled, err := unmarshalOrders(aux.OrdersCancelled)
	if err != nil {
		return err
	}
	c.OrdersCancelled = ordersCancelled
	ordersFilled, err := unmarshalOrders(aux.OrdersFilled)
	if err != nil {
		return err
	}
	c.OrdersFilled = ordersFilled
	ordersTriggered, err := unmarshalOrders(aux.OrdersTriggered)
	if err != nil {
		return err
	}
	c.OrdersTriggered = ordersTriggered
	transactions, err := unmarshalTransactions(aux.Transactions)
	if err != nil {
		return err
	}
	c.Transactions = transactions

	return nil
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

// ---------------------------------------------------------------
// Endpoints https://developer.oanda.com/rest-live-v20/account-ep/
// ---------------------------------------------------------------

// AccountService handles communication with the Account related endpoints of the
// OANDA v20 REST API.
type AccountService struct {
	client *Client
}

func newAccountService(client *Client) *AccountService {
	return &AccountService{client}
}

// AccountListResponse is the response returned by [AccountService.List].
type AccountListResponse struct {
	Accounts []AccountProperties `json:"accounts"`
}

// List retrieves the list of accounts authorized for the provided token.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts
//
// Returns:
//   - []AccountProperties: A slice of account properties for all authorized accounts.
//   - error: An error if the request fails or response cannot be decoded.
//
// Reference: https://developer.oanda.com/rest-live-v20/account-ep/#collapse_endpoint_1
func (s *AccountService) List(ctx context.Context) (*AccountListResponse, error) {
	return doGet[AccountListResponse](s.client, ctx, "/v3/accounts", nil)
}

// AccountDetailsResponse is the response returned by [AccountService.Details].
type AccountDetailsResponse struct {
	Account           Account       `json:"account"`
	LastTransactionID TransactionID `json:"lastTransactionID"`
}

// Details retrieves the full details for the Account configured via [WithAccountID].
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}
//
// Reference: https://developer.oanda.com/rest-live-v20/account-ep/#collapse_endpoint_2
func (s *AccountService) Details(ctx context.Context) (*AccountDetailsResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v", s.client.accountID)
	return doGet[AccountDetailsResponse](s.client, ctx, path, nil)
}

// AccountSummaryResponse is the response returned by [AccountService.Summary].
type AccountSummaryResponse struct {
	Account           AccountSummary `json:"account"`
	LastTransactionID TransactionID  `json:"lastTransactionID"`
}

// Summary retrieves a summary for the Account configured via [WithAccountID].
//
// Unlike [AccountService.Details], this method does not include full pending Order,
// open Trade, and Position representations, making it more lightweight.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/summary
//
// Reference: https://developer.oanda.com/rest-live-v20/account-ep/#collapse_endpoint_3
func (s *AccountService) Summary(ctx context.Context) (*AccountSummaryResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/summary", s.client.accountID)
	return doGet[AccountSummaryResponse](s.client, ctx, path, nil)
}

// AccountConfigureRequest represents a request to update Account configuration.
// Use [NewAccountConfigureRequest] to create one, then chain setters.
type AccountConfigureRequest struct {
	Alias      string        `json:"alias"`
	MarginRate DecimalNumber `json:"marginRate"`
}

func (r *AccountConfigureRequest) body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

// NewAccountConfigureRequest creates a new empty [AccountConfigureRequest].
func NewAccountConfigureRequest() *AccountConfigureRequest {
	return &AccountConfigureRequest{}
}

// SetAlias sets the client-assigned alias for the Account.
func (r *AccountConfigureRequest) SetAlias(alias string) *AccountConfigureRequest {
	r.Alias = alias
	return r
}

// SetMarginRate sets the margin rate for the Account.
func (r *AccountConfigureRequest) SetMarginRate(marginRate DecimalNumber) *AccountConfigureRequest {
	r.MarginRate = marginRate
	return r
}

// AccountConfigureResponse is the successful response returned by [AccountService.Configure].
type AccountConfigureResponse struct {
	ClientConfigureTransaction ClientConfigureTransaction `json:"clientConfigureTransaction"`
	LastTransactionID          TransactionID              `json:"lastTransactionID"`
}

// AccountConfigureErrorResponse is the error response returned by [AccountService.Configure]
// when the request is rejected (400 or 403).
type AccountConfigureErrorResponse struct {
	ClientConfigureRejectTransaction ClientConfigureRejectTransaction `json:"clientConfigureRejectTransaction"`
	LastTransactionID                TransactionID                    `json:"lastTransactionID"`
	ErrorCode                        string                           `json:"errorCode"`
	ErrorMessage                     string                           `json:"errorMessage"`
}

// Error implements the error interface.
func (r AccountConfigureErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", r.ErrorCode, r.ErrorMessage)
}

// Configure sets the client-configurable portions of an Account (alias and margin rate).
//
// This corresponds to the OANDA API endpoint: PATCH /v3/accounts/{accountID}/configuration
//
// Reference: https://developer.oanda.com/rest-live-v20/account-ep/#collapse_endpoint_5
func (s *AccountService) Configure(ctx context.Context, req *AccountConfigureRequest) (*AccountConfigureResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/configuration", s.client.accountID)
	var body io.Reader
	var err error
	if req != nil {
		body, err = req.body()
		if err != nil {
			return nil, err
		}
	}
	httpResp, err := s.client.sendPatchRequest(ctx, path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to send PATCH request: %w", err)
	}
	switch httpResp.StatusCode {
	case http.StatusOK:
		var resp AccountConfigureResponse
		if err := decodeResponse(httpResp, &resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &resp, nil
	case http.StatusBadRequest:
		var resp AccountConfigureErrorResponse
		if err := decodeResponse(httpResp, &resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return nil, BadRequest{HTTPError{StatusCode: httpResp.StatusCode, Message: "bad request", Err: resp}}
	case http.StatusForbidden:
		var resp AccountConfigureErrorResponse
		if err := decodeResponse(httpResp, &resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return nil, Forbidden{HTTPError{StatusCode: httpResp.StatusCode, Message: "forbidden", Err: resp}}
	default:
		return nil, decodeErrorResponse(httpResp)
	}
}

// AccountChangesResponse is the response returned by [AccountService.Changes].
type AccountChangesResponse struct {
	Changes           AccountChanges      `json:"changes"`
	State             AccountChangesState `json:"state"`
	LastTransactionID TransactionID       `json:"lastTransactionID"`
}

// Changes retrieves the changes to an Account's Orders, Trades, and Positions
// since a specified TransactionID. The Account is determined by [WithAccountID].
//
// This endpoint is useful for polling-based synchronization. By tracking the
// LastTransactionID from each response, you can efficiently fetch only the
// changes that have occurred since your last request.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/changes
//
// Reference: https://developer.oanda.com/rest-live-v20/account-ep/#collapse_endpoint_6
func (s *AccountService) Changes(ctx context.Context, since TransactionID) (*AccountChangesResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/changes", s.client.accountID)
	v := url.Values{}
	v.Set("sinceTransactionID", since)
	return doGet[AccountChangesResponse](s.client, ctx, path, v)
}
