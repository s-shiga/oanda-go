package oanda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------
// Definitions https://developer.oanda.com/rest-live-v20/trade-df/
// ---------------------------------------------------------------

// TradeID is the unique identifier for a Trade within an Account. It is a string representation
// of the OANDA-assigned TradeID, derived from the TransactionID of the Transaction that opened
// the Trade. Example: "1523"
type TradeID = string

// TradeSpecifier is the specification of a Trade as referred to by clients. Either the Trade's
// OANDA-assigned TradeID or the Trade's client-provided ClientID prefixed by the "@" symbol.
type TradeSpecifier = string

// Trade is the specification of a Trade within an Account. This includes the full representation
// of the Trade's dependent Orders in addition to the IDs of those Orders.
type Trade struct {
	// ID is the Trade's identifier, unique within the Trade's Account.
	ID TradeID `json:"id"`
	// Instrument is the Trade's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Price is the execution price of the Trade.
	Price PriceValue `json:"price"`
	// OpenTime is the date/time when the Trade was opened.
	OpenTime DateTime `json:"openTime"`
	// State is the current state of the Trade.
	State TradeState `json:"state"`
	// InitialUnits is the initial size of the Trade. Negative values indicate a short Trade,
	// and positive values indicate a long Trade.
	InitialUnits DecimalNumber `json:"initialUnits"`
	// InitialMarginRequired is the margin required at the time the Trade was created. Note, this
	// is the 'pure' margin required, it is not the 'effective' margin used that factors in the
	// trade risk if a GSLO is attached to the trade.
	InitialMarginRequired AccountUnits `json:"initialMarginRequired"`
	// CurrentUnits is the number of units currently open for the Trade. This value is reduced to
	// 0.0 as the Trade is closed.
	CurrentUnits DecimalNumber `json:"currentUnits"`
	// RealizedPL is the total profit/loss realized on the closed portion of the Trade.
	RealizedPL AccountUnits `json:"realizedPL"`
	// UnrealizedPL is the unrealized profit/loss on the open portion of the Trade.
	UnrealizedPL AccountUnits `json:"unrealizedPL"`
	// MarginUsed is the margin currently used by the Trade.
	MarginUsed AccountUnits `json:"marginUsed"`
	// AverageClosePrice is the average closing price of the Trade. Only present if the Trade has
	// been closed or reduced at least once.
	AverageClosePrice *PriceValue `json:"averageClosePrice,omitempty"`
	// ClosingTransactionIDs is the list of Transaction IDs associated with closing portions of
	// this Trade.
	ClosingTransactionIDs []TransactionID `json:"closingTransactionIDs,omitempty"`
	// Financing is the financing paid/collected for this Trade.
	Financing AccountUnits `json:"financing"`
	// DividendAdjustment is the dividend adjustment paid or collected for this Trade.
	DividendAdjustment AccountUnits `json:"dividendAdjustment"`
	// CloseTime is the date/time when the Trade was fully closed. Only provided for Trades whose
	// state is CLOSED.
	CloseTime *DateTime `json:"closeTime,omitempty"`
	// ClientExtensions are the client extensions of the Trade.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// TakeProfitOrder is the full representation of the Trade's Take Profit Order, only provided
	// if such an Order exists.
	TakeProfitOrder *TakeProfitOrder `json:"takeProfitOrder,omitempty"`
	// StopLossOrder is the full representation of the Trade's Stop Loss Order, only provided if
	// such an Order exists.
	StopLossOrder *StopLossOrder `json:"stopLossOrder,omitempty"`
	// TrailingStopLossOrder is the full representation of the Trade's Trailing Stop Loss Order,
	// only provided if such an Order exists.
	TrailingStopLossOrder *TrailingStopLossOrder `json:"trailingStopLossOrder,omitempty"`
}

// TradeSummary is the summary of a Trade within an Account. This representation does not provide
// the full representation of the Trade's dependent Orders.
type TradeSummary struct {
	// ID is the Trade's identifier, unique within the Trade's Account.
	ID TradeID `json:"id"`
	// Instrument is the Trade's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Price is the execution price of the Trade.
	Price PriceValue `json:"price"`
	// OpenTime is the date/time when the Trade was opened.
	OpenTime DateTime `json:"openTime"`
	// State is the current state of the Trade.
	State TradeState `json:"state"`
	// InitialUnits is the initial size of the Trade. Negative values indicate a short Trade,
	// and positive values indicate a long Trade.
	InitialUnits DecimalNumber `json:"initialUnits"`
	// InitialMarginRequired is the margin required at the time the Trade was created. Note, this
	// is the 'pure' margin required, it is not the 'effective' margin used that factors in the
	// trade risk if a GSLO is attached to the trade.
	InitialMarginRequired AccountUnits `json:"initialMarginRequired"`
	// CurrentUnits is the number of units currently open for the Trade. This value is reduced to
	// 0.0 as the Trade is closed.
	CurrentUnits DecimalNumber `json:"currentUnits"`
	// RealizedPL is the total profit/loss realized on the closed portion of the Trade.
	RealizedPL AccountUnits `json:"realizedPL"`
	// UnrealizedPL is the unrealized profit/loss on the open portion of the Trade.
	UnrealizedPL AccountUnits `json:"unrealizedPL"`
	// MarginUsed is the margin currently used by the Trade.
	MarginUsed AccountUnits `json:"marginUsed"`
	// AverageClosePrice is the average closing price of the Trade. Only present if the Trade has
	// been closed or reduced at least once.
	AverageClosePrice PriceValue `json:"averageClosePrice"`
	// ClosingTransactionIDs is the list of Transaction IDs associated with closing portions of
	// this Trade.
	ClosingTransactionIDs []TransactionID `json:"closingTransactionIDs,omitempty"`
	// Financing is the financing paid/collected for this Trade.
	Financing AccountUnits `json:"financing"`
	// DividendAdjustment is the dividend adjustment paid or collected for this Trade.
	DividendAdjustment AccountUnits `json:"dividendAdjustment"`
	// CloseTime is the date/time when the Trade was fully closed. Only provided for Trades whose
	// state is CLOSED.
	CloseTime *DateTime `json:"closeTime,omitempty"`
	// ClientExtensions are the client extensions of the Trade.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// TakeProfitOrderID is the ID of the Trade's Take Profit Order, only provided if such an
	// Order exists.
	TakeProfitOrderID OrderID `json:"takeProfitOrderID"`
	// StopLossOrderID is the ID of the Trade's Stop Loss Order, only provided if such an Order exists.
	StopLossOrderID OrderID `json:"stopLossOrderID"`
	// GuaranteedStopLossOrderID is the ID of the Trade's Guaranteed Stop Loss Order, only provided
	// if such an Order exists.
	GuaranteedStopLossOrderID OrderID `json:"guaranteedStopLossOrderID,omitempty"`
	// TrailingStopLossOrderID is the ID of the Trade's Trailing Stop Loss Order, only provided if
	// such an Order exists.
	TrailingStopLossOrderID OrderID `json:"trailingStopLossOrderID,omitempty"`
}

// CalculatedTradeState represents the dynamic (calculated) state of an open Trade.
type CalculatedTradeState struct {
	// ID is the Trade's ID.
	ID TradeID `json:"id"`
	// UnrealizedPL is the Trade's unrealized profit/loss.
	UnrealizedPL AccountUnits `json:"unrealizedPL"`
	// MarginUsed is the margin currently used by the Trade.
	MarginUsed AccountUnits `json:"marginUsed"`
}

// TradeState represents the current state of the Trade.
type TradeState string

const (
	// TradeStateOpen means the Trade is currently open.
	TradeStateOpen TradeState = "OPEN"
	// TradeStateClosed means the Trade has been fully closed.
	TradeStateClosed TradeState = "CLOSED"
	// TradeStateCloseWhenTradeable means the Trade will be closed as soon as the trade's
	// instrument becomes tradeable.
	TradeStateCloseWhenTradeable TradeState = "CLOSE_WHEN_TRADEABLE"
)

// TradeStateFilter is used to filter Trades by their state.
type TradeStateFilter string

const (
	// TradeStateFilterOpen filters for Trades in the OPEN state.
	TradeStateFilterOpen TradeStateFilter = "OPEN"
	// TradeStateFilterClosed filters for Trades in the CLOSED state.
	TradeStateFilterClosed TradeStateFilter = "CLOSED"
	// TradeStateFilterCloseWhenTradeable filters for Trades in the CLOSE_WHEN_TRADEABLE state.
	TradeStateFilterCloseWhenTradeable TradeStateFilter = "CLOSE_WHEN_TRADEABLE"
	// TradeStateFilterAll selects all Trades regardless of their state.
	TradeStateFilterAll TradeStateFilter = "ALL"
)

// TradePL classifies a Trade's profit/loss.
type TradePL string

const (
	// TradePLPositive means the Trade's profit/loss is positive (profitable).
	TradePLPositive TradePL = "POSITIVE"
	// TradePLNegative means the Trade's profit/loss is negative (losing).
	TradePLNegative TradePL = "NEGATIVE"
	// TradePLZero means the Trade's profit/loss is zero (break-even).
	TradePLZero TradePL = "ZERO"
)

// -------------------------------------------------------------
// Endpoints https://developer.oanda.com/rest-live-v20/trade-ep/
// -------------------------------------------------------------

type tradeService struct {
	client *Client
}

func newTradeService(client *Client) *tradeService {
	return &tradeService{client}
}

// TradeListRequest represents a request to list Trades for an Account.
type TradeListRequest struct {
	// IDs is a list of Trade IDs to retrieve. If specified, only Trades with these IDs will
	// be returned.
	IDs []TradeID
	// State filters Trades by their state. Default is OPEN.
	State *TradeStateFilter
	// Instrument filters Trades by the instrument they are for.
	Instrument *InstrumentName
	// Count is the maximum number of Trades to return. Maximum value is 500.
	Count *int
	// BeforeID returns only Trades that were opened before this Trade ID.
	BeforeID *TradeID
}

// NewTradeListRequest creates a new TradeListRequest for the given account.
func NewTradeListRequest() *TradeListRequest {
	return &TradeListRequest{
		IDs: make([]TradeID, 0),
	}
}

// AddIDs adds Trade IDs to filter the results.
func (r *TradeListRequest) AddIDs(id ...TradeID) *TradeListRequest {
	r.IDs = append(r.IDs, id...)
	return r
}

// SetStateFilter sets the state filter for the Trades to return.
func (r *TradeListRequest) SetStateFilter(filter TradeStateFilter) *TradeListRequest {
	r.State = &filter
	return r
}

// SetInstrument filters the Trades by the specified instrument.
func (r *TradeListRequest) SetInstrument(instrument InstrumentName) *TradeListRequest {
	r.Instrument = &instrument
	return r
}

// SetCount sets the maximum number of Trades to return. Maximum value is 500.
func (r *TradeListRequest) SetCount(count int) *TradeListRequest {
	r.Count = &count
	return r
}

// SetBeforeID returns only Trades that were opened before this Trade ID.
func (r *TradeListRequest) SetBeforeID(beforeID TradeID) *TradeListRequest {
	r.BeforeID = &beforeID
	return r
}

// validate checks that the request parameters are valid.
func (r *TradeListRequest) validate() error {
	if r.Count != nil {
		if *r.Count < 0 {
			return errors.New("count must be greater than or equal to 0")
		}
		if *r.Count > 500 {
			return errors.New("count must be less than or equal to 500")
		}
	}
	return nil
}

// values validates parameters and returns url.Values for the request.
func (r *TradeListRequest) values() (url.Values, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if len(r.IDs) > 0 {
		v.Set("ids", strings.Join(r.IDs, ","))
	}
	if r.State != nil {
		v.Set("state", string(*r.State))
	}
	if r.Instrument != nil {
		v.Set("instrument", *r.Instrument)
	}
	if r.Count != nil {
		v.Set("count", strconv.Itoa(*r.Count))
	}
	if r.BeforeID != nil {
		v.Set("beforeID", *r.BeforeID)
	}
	return v, nil
}

// TradeListResponse represents the response from a Trade list request.
type TradeListResponse struct {
	// Trades is the list of Trade details.
	Trades []Trade `json:"trades"`
	// LastTransactionID is the ID of the most recent Transaction created for the Account.
	LastTransactionID TransactionID `json:"lastTransactionID"`
}

// List retrieves a list of Trades for the Account configured via [WithAccountID].
// Use [NewTradeListRequest] to create and configure filter parameters.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/trades
//
// Reference: https://developer.oanda.com/rest-live-v20/trade-ep/#collapse_endpoint_2
func (s *tradeService) List(ctx context.Context, req *TradeListRequest) (*TradeListResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%s/trades", s.client.accountID)
	v, err := req.values()
	if err != nil {
		return nil, err
	}
	httpResp, err := s.client.sendGetRequest(ctx, path, v)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var resp TradeListResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &resp, nil
}

// ListOpen retrieves all currently open Trades for the Account configured via [WithAccountID].
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/openTrades
//
// Reference: https://developer.oanda.com/rest-live-v20/trade-ep/#collapse_endpoint_3
func (s *tradeService) ListOpen(ctx context.Context) (*TradeListResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%s/openTrades", s.client.accountID)
	httpResp, err := s.client.sendGetRequest(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var resp TradeListResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &resp, nil
}

// TradeDetailsResponse is the response returned by [tradeService.Details].
type TradeDetailsResponse struct {
	Trade             Trade         `json:"trade"`
	LastTransactionID TransactionID `json:"lastTransactionID"`
}

// Details retrieves the details of a specific Trade for the Account configured via [WithAccountID].
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/trades/{tradeSpecifier}
//
// Reference: https://developer.oanda.com/rest-live-v20/trade-ep/#collapse_endpoint_4
func (s *tradeService) Details(ctx context.Context, specifier TradeSpecifier) (*TradeDetailsResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%s/trades/%s", s.client.accountID, specifier)
	httpResp, err := s.client.sendGetRequest(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var resp TradeDetailsResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &resp, nil
}

// TradeCloseRequest represents a request to close (fully or partially) a Trade.
// Use [NewTradeCloseRequest] or [NewTradeCloseALLRequest] to create one.
type TradeCloseRequest struct {
	// Units is indication of how much of the Trade to close. Either the string “ALL”
	// (indicating that all of the Trade should be closed), or a DecimalNumber
	// representing the number of units of the open Trade to Close using a
	// TradeClose MarketOrder. The units specified must always be positive, and
	// the magnitude of the value cannot exceed the magnitude of the Trade’s
	// open units.
	Units DecimalNumber `json:"units"`
}

func (r TradeCloseRequest) body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

// NewTradeCloseRequest creates a request to partially close a Trade by the specified number of units.
func NewTradeCloseRequest(units DecimalNumber) TradeCloseRequest {
	return TradeCloseRequest{Units: units}
}

// NewTradeCloseALLRequest creates a request to fully close a Trade.
func NewTradeCloseALLRequest() TradeCloseRequest {
	return TradeCloseRequest{Units: "ALL"}
}

// TradeCloseResponse is the successful response returned by [Client.TradeClose].
type TradeCloseResponse struct {
	OrderCreateTransaction MarketOrderTransaction  `json:"orderCreateTransaction"`
	OrderFillTransaction   OrderFillTransaction    `json:"orderFillTransaction"`
	OrderCancelTransaction *OrderCancelTransaction `json:"orderCancelTransaction,omitempty"`
	RelatedTransactionIDs  []TransactionID         `json:"relatedTransactionIDs"`
	LastTransactionID      TransactionID           `json:"lastTransactionID"`
}

// TradeCloseBadRequestResponse is the error response returned by [Client.TradeClose] on a 400 status.
type TradeCloseBadRequestResponse struct {
	OrderRejectTransaction MarketOrderRejectTransaction `json:"orderRejectTransaction"`
	ErrorCode              string                       `json:"errorCode"`
	ErrorMessage           string                       `json:"errorMessage"`
}

// Error implements the error interface.
func (r TradeCloseBadRequestResponse) Error() string {
	return fmt.Sprintf("%s: %s", r.ErrorCode, r.ErrorMessage)
}

// TradeCloseNotFoundResponse is the error response returned by [Client.TradeClose] on a 404 status.
type TradeCloseNotFoundResponse struct {
	OrderRejectTransaction MarketOrderRejectTransaction `json:"orderRejectTransaction"`
	LastTransactionID      TransactionID                `json:"lastTransactionID"`
	RelatedTransactionIDs  []TransactionID              `json:"relatedTransactionIDs"`
	ErrorCode              string                       `json:"errorCode"`
	ErrorMessage           string                       `json:"errorMessage"`
}

// Error implements the error interface.
func (r TradeCloseNotFoundResponse) Error() string {
	return fmt.Sprintf("%s: %s", r.ErrorCode, r.ErrorMessage)
}

// Close closes (fully or partially) a specific Trade for the Account configured via [WithAccountID].
//
// This corresponds to the OANDA API endpoint: PUT /v3/accounts/{accountID}/trades/{tradeSpecifier}/close
//
// Reference: https://developer.oanda.com/rest-live-v20/trade-ep/#collapse_endpoint_5
func (s *tradeService) Close(ctx context.Context, specifier TradeSpecifier, req TradeCloseRequest) (*TradeCloseResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%s/trades/%s/close", s.client.accountID, specifier)
	body, err := req.body()
	if err != nil {
		return nil, err
	}
	httpResp, err := s.client.sendPutRequest(ctx, path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to send PUT request: %w", err)
	}
	defer closeBody(httpResp)
	switch httpResp.StatusCode {
	case http.StatusOK:
		var resp TradeCloseResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &resp, nil
	case http.StatusBadRequest:
		var resp TradeCloseBadRequestResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return nil, BadRequest{HTTPError{StatusCode: httpResp.StatusCode, Message: "bad request", Err: resp}}
	case http.StatusNotFound:
		var resp TradeCloseNotFoundResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return nil, NotFoundError{HTTPError{StatusCode: httpResp.StatusCode, Message: "not found", Err: resp}}
	default:
		return nil, decodeErrorResponse(httpResp)
	}
}

// TradeUpdateClientExtensionsRequest is the request body for updating client extensions on a Trade.
type TradeUpdateClientExtensionsRequest struct {
	ClientExtensions ClientExtensions `json:"clientExtensions"`
}

func (r TradeUpdateClientExtensionsRequest) body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r.ClientExtensions)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

// TradeUpdateClientExtensionsResponse is the successful response returned by [tradeService.UpdateClientExtensions].
type TradeUpdateClientExtensionsResponse struct {
	TradeClientExtensionsModifyTransaction TradeClientExtensionsModifyTransaction `json:"tradeClientExtensionsModifyTransaction"`
	RelatedTransactionIDs                  []TransactionID                        `json:"relatedTransactionIDs"`
	LastTransactionID                      TransactionID                          `json:"lastTransactionID"`
}

// TradeUpdateClientExtensionsErrorResponse is the error response returned by [tradeService.UpdateClientExtensions].
type TradeUpdateClientExtensionsErrorResponse struct {
	TradeClientExtensionsModifyRejectTransaction TradeClientExtensionsModifyRejectTransaction `json:"tradeClientExtensionsModifyRejectTransaction"`
	LastTransactionID                            TransactionID                                `json:"lastTransactionID"`
	RelatedTransactionIDs                        []TransactionID                              `json:"relatedTransactionIDs"`
	ErrorCode                                    string                                       `json:"errorCode"`
	ErrorMessage                                 string                                       `json:"errorMessage"`
}

// Error implements the error interface.
func (r TradeUpdateClientExtensionsErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", r.ErrorCode, r.ErrorMessage)
}

// UpdateClientExtensions updates the client extensions for a Trade.
//
// This corresponds to the OANDA API endpoint: PUT /v3/accounts/{accountID}/trades/{tradeSpecifier}/clientExtensions
//
// Reference: https://developer.oanda.com/rest-live-v20/trade-ep/#collapse_endpoint_6
func (s *tradeService) UpdateClientExtensions(ctx context.Context, specifier TradeSpecifier, req TradeUpdateClientExtensionsRequest) (*TradeUpdateClientExtensionsResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%s/trades/%s/clientExtensions", s.client.accountID, specifier)
	body, err := req.body()
	if err != nil {
		return nil, err
	}
	httpResp, err := s.client.sendPutRequest(ctx, path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to send PUT request: %w", err)
	}
	defer closeBody(httpResp)
	switch httpResp.StatusCode {
	case http.StatusOK:
		var resp TradeUpdateClientExtensionsResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &resp, nil
	case http.StatusBadRequest:
		var resp TradeUpdateClientExtensionsErrorResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return nil, BadRequest{HTTPError{StatusCode: httpResp.StatusCode, Message: "bad request", Err: resp}}
	case http.StatusNotFound:
		var resp TradeUpdateClientExtensionsErrorResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return nil, NotFoundError{HTTPError{StatusCode: httpResp.StatusCode, Message: "not found", Err: resp}}
	default:
		return nil, decodeErrorResponse(httpResp)
	}
}

// TradeUpdateOrdersRequest is the request body for creating, replacing, or cancelling
// a Trade's dependent Orders (Take Profit, Stop Loss, Trailing Stop Loss, Guaranteed Stop Loss).
type TradeUpdateOrdersRequest struct {
	TakeProfit         *TakeProfitDetails         `json:"takeProfit,omitempty"`
	StopLoss           *StopLossDetails           `json:"stopLoss,omitempty"`
	TrailingStopLoss   *TrailingStopLossDetails   `json:"trailingStopLoss,omitempty"`
	GuaranteedStopLoss *GuaranteedStopLossDetails `json:"guaranteedStopLoss,omitempty"`
}

func (r TradeUpdateOrdersRequest) body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

// TradeUpdateOrdersResponse is the successful response returned by [tradeService.UpdateOrders].
type TradeUpdateOrdersResponse struct {
	TakeProfitOrderCancelTransaction         *OrderCancelTransaction             `json:"takeProfitOrderCancelTransaction,omitempty"`
	TakeProfitOrderTransaction               *TakeProfitOrderTransaction         `json:"takeProfitOrderTransaction,omitempty"`
	TakeProfitOrderFillTransaction           *OrderFillTransaction               `json:"takeProfitOrderFillTransaction,omitempty"`
	TakeProfitOrderCreatedCancelTransaction  *OrderCancelTransaction             `json:"takeProfitOrderCreatedCancelTransaction,omitempty"`
	StopLossOrderCancelTransaction           *OrderCancelTransaction             `json:"stopLossOrderCancelTransaction,omitempty"`
	StopLossOrderTransaction                 *StopLossOrderTransaction           `json:"stopLossOrderTransaction,omitempty"`
	StopLossOrderFillTransaction             *OrderFillTransaction               `json:"stopLossOrderFillTransaction,omitempty"`
	StopLossOrderCreatedCancelTransaction    *OrderCancelTransaction             `json:"stopLossOrderCreatedCancelTransaction,omitempty"`
	TrailingStopLossOrderCancelTransaction   *OrderCancelTransaction             `json:"trailingStopLossOrderCancelTransaction,omitempty"`
	TrailingStopLossOrderTransaction         *TrailingStopLossOrderTransaction   `json:"trailingStopLossOrderTransaction,omitempty"`
	GuaranteedStopLossOrderCancelTransaction *OrderCancelTransaction             `json:"guaranteedStopLossOrderCancelTransaction,omitempty"`
	GuaranteedStopLossOrderTransaction       *GuaranteedStopLossOrderTransaction `json:"guaranteedStopLossOrderTransaction,omitempty"`
	RelatedTransactionIDs                    []TransactionID                     `json:"relatedTransactionIDs"`
	LastTransactionID                        TransactionID                       `json:"lastTransactionID"`
}

// TradeUpdateOrdersErrorResponse is the error response returned by [tradeService.UpdateOrders].
type TradeUpdateOrdersErrorResponse struct {
	TakeProfitOrderCancelRejectTransaction         *OrderCancelRejectTransaction             `json:"takeProfitOrderCancelRejectTransaction,omitempty"`
	TakeProfitOrderRejectTransaction               *TakeProfitOrderRejectTransaction         `json:"takeProfitOrderRejectTransaction,omitempty"`
	StopLossOrderCancelRejectTransaction           *OrderCancelRejectTransaction             `json:"stopLossOrderCancelRejectTransaction,omitempty"`
	StopLossOrderRejectTransaction                 *StopLossOrderRejectTransaction           `json:"stopLossOrderRejectTransaction,omitempty"`
	TrailingStopLossOrderCancelRejectTransaction   *OrderCancelRejectTransaction             `json:"trailingStopLossOrderCancelRejectTransaction,omitempty"`
	TrailingStopLossOrderRejectTransaction         *TrailingStopLossOrderRejectTransaction   `json:"trailingStopLossOrderRejectTransaction,omitempty"`
	GuaranteedStopLossOrderCancelRejectTransaction *OrderCancelRejectTransaction             `json:"guaranteedStopLossOrderCancelRejectTransaction,omitempty"`
	GuaranteedStopLossOrderRejectTransaction       *GuaranteedStopLossOrderRejectTransaction `json:"guaranteedStopLossOrderRejectTransaction,omitempty"`
	LastTransactionID                              TransactionID                             `json:"lastTransactionID"`
	RelatedTransactionIDs                          []TransactionID                           `json:"relatedTransactionIDs"`
	ErrorCode                                      string                                    `json:"errorCode"`
	ErrorMessage                                   string                                    `json:"errorMessage"`
}

// Error implements the error interface.
func (r TradeUpdateOrdersErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", r.ErrorCode, r.ErrorMessage)
}

// UpdateOrders creates, replaces, or cancels a Trade's dependent Orders
// (Take Profit, Stop Loss, Trailing Stop Loss, Guaranteed Stop Loss).
//
// This corresponds to the OANDA API endpoint: PUT /v3/accounts/{accountID}/trades/{tradeSpecifier}/orders
//
// Reference: https://developer.oanda.com/rest-live-v20/trade-ep/#collapse_endpoint_7
func (s *tradeService) UpdateOrders(ctx context.Context, specifier TradeSpecifier, req *TradeUpdateOrdersRequest) (*TradeUpdateOrdersResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%s/trades/%s/orders", s.client.accountID, specifier)
	body, err := req.body()
	if err != nil {
		return nil, err
	}
	httpResp, err := s.client.sendPutRequest(ctx, path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to send PUT request: %w", err)
	}
	defer closeBody(httpResp)
	switch httpResp.StatusCode {
	case http.StatusOK:
		var resp TradeUpdateOrdersResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &resp, nil
	case http.StatusBadRequest:
		var resp TradeUpdateOrdersErrorResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return nil, BadRequest{HTTPError{StatusCode: httpResp.StatusCode, Message: "bad request", Err: resp}}
	default:
		return nil, decodeErrorResponse(httpResp)
	}
}
