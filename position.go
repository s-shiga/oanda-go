package oanda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// ------------------------------------------------------------------
// Definitions https://developer.oanda.com/rest-live-v20/position-df/
// ------------------------------------------------------------------

// Position is the specification of a Position within an Account.
type Position struct {
	// Instrument of the Position.
	Instrument InstrumentName `json:"instrument"`
	// PL (profit/loss) realized by the Position over the lifetime of the Account.
	PL AccountUnits `json:"pl"`
	// UnrealizedPL of all open Trades that contribute to this Position.
	UnrealizedPL AccountUnits `json:"unrealizedPL"`
	// MarginUsed is margin currently used by the Position.
	MarginUsed AccountUnits `json:"marginUsed"`
	// ResettablePL is Profit/loss realized by the Position since the Account’s resettablePL was
	// last reset by the client.
	ResettablePL AccountUnits `json:"resettablePL"`
	// Financing is the total amount of financing paid/collected over the lifetime of the Position.
	Financing AccountUnits `json:"financing"`
	// Commission is the total amount of commission paid over the lifetime of the Position.
	Commission AccountUnits `json:"commission"`
	// DividendAdjustment is the total amount of dividend adjustments paid or collected over the
	// lifetime of the Position in the Account's home currency.
	DividendAdjustment AccountUnits `json:"dividendAdjustment"`
	// GuaranteedExecutionsFees is the total amount of fees charged over the lifetime of the Account
	// for the execution of guaranteed Stop Loss Orders attached to Trades for this Position.
	GuaranteedExecutionsFees AccountUnits `json:"guaranteedExecutionsFees"`
	// Long is the details of the long side of the Position.
	Long PositionSide `json:"long"`
	// Short is the details of the short side of the Position.
	Short PositionSide `json:"short"`
}

// PositionSide represents the holdings in a single direction (long or short) for a Position.
type PositionSide struct {
	// Units is the number of units in the position (a positive number for a long position,
	// and a negative number for a short position).
	Units DecimalNumber `json:"units"`
	// AveragePrice is the volume-weighted average of the underlying Trade open prices for
	// the Position.
	AveragePrice PriceValue `json:"averagePrice"`
	// TradeIDs is the list of the open Trade IDs which contribute to the open Position.
	TradeIDs []TradeID `json:"tradeIDs,omitempty"`
	// PL is the profit/loss realized by the PositionSide over the lifetime of the Account.
	PL AccountUnits `json:"pl"`
	// UnrealizedPL is the unrealized profit/loss of all open Trades that contribute to this PositionSide.
	UnrealizedPL AccountUnits `json:"unrealizedPL"`
	// ResettablePL is the profit/loss realized by the PositionSide since the Account's resettablePL
	// was last reset by the client.
	ResettablePL AccountUnits `json:"resettablePL"`
	// Financing is the total amount of financing paid/collected for this PositionSide over the
	// lifetime of the Account.
	Financing AccountUnits `json:"financing"`
	// DividendAdjustment is the total amount of dividend adjustments paid or collected for this
	// PositionSide over the lifetime of the Account.
	DividendAdjustment AccountUnits `json:"dividendAdjustment"`
	// GuaranteedExecutionFees is the total amount of fees charged over the lifetime of the Account
	// for the execution of guaranteed Stop Loss Orders attached to Trades for this PositionSide.
	GuaranteedExecutionFees AccountUnits `json:"guaranteedExecutionFees"`
}

// CalculatedPositionState represents the dynamic (calculated) state of a Position.
type CalculatedPositionState struct {
	// Instrument is the Position's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// NetUnrealizedPL is the Position's net unrealized profit/loss.
	NetUnrealizedPL AccountUnits `json:"netUnrealizedPL"`
	// LongUnrealizedPL is the unrealized profit/loss of the Position's long side.
	LongUnrealizedPL AccountUnits `json:"longUnrealizedPL"`
	// ShortUnrealizedPL is the unrealized profit/loss of the Position's short side.
	ShortUnrealizedPL AccountUnits `json:"shortUnrealizedPL"`
	// MarginUsed is the margin currently used by the Position.
	MarginUsed AccountUnits `json:"marginUsed"`
}

// ----------------------------------------------------------------
// Endpoints https://developer.oanda.com/rest-live-v20/position-ep/
// ----------------------------------------------------------------

type positionService struct {
	client *Client
}

func newPositionService(client *Client) *positionService {
	return &positionService{client}
}

// PositionListResponse is the response returned by [positionService.List] and [positionService.ListOpen].
type PositionListResponse struct {
	Positions         []Position    `json:"positions"`
	LastTransactionID TransactionID `json:"lastTransactionId"`
}

// List retrieves all Positions for the Account configured via [WithAccountID].
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/positions
//
// Reference: https://developer.oanda.com/rest-live-v20/position-ep/#collapse_endpoint_1
func (s *positionService) List(ctx context.Context) (*PositionListResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/positions", s.client.accountID)
	httpResp, err := s.client.sendGetRequest(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var resp PositionListResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &resp, nil
}

// ListOpen retrieves all open Positions for the Account configured via [WithAccountID].
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/openPositions
//
// Reference: https://developer.oanda.com/rest-live-v20/position-ep/#collapse_endpoint_2
func (s *positionService) ListOpen(ctx context.Context) (*PositionListResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/openPositions", s.client.accountID)
	httpResp, err := s.client.sendGetRequest(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var resp PositionListResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &resp, nil
}

// PositionListByInstrumentResponse is the response returned by [positionService.ListByInstrument].
type PositionListByInstrumentResponse struct {
	Position          Position      `json:"position"`
	LastTransactionID TransactionID `json:"lastTransactionID"`
}

// ListByInstrument retrieves the Position for a specific instrument in the Account configured via [WithAccountID].
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/positions/{instrument}
//
// Reference: https://developer.oanda.com/rest-live-v20/position-ep/#collapse_endpoint_3
func (s *positionService) ListByInstrument(ctx context.Context, instrument InstrumentName) (*PositionListByInstrumentResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/positions/%v", s.client.accountID, instrument)
	httpResp, err := s.client.sendGetRequest(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var resp PositionListByInstrumentResponse
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &resp, nil
}

// PositionCloseRequest represents a request to close (fully or partially) a Position's long and/or short side.
type PositionCloseRequest struct {
	LongUnits             *string           `json:"longUnits,omitempty"`
	LongClientExtensions  *ClientExtensions `json:"longClientExtensions,omitempty"`
	ShortUnits            *string           `json:"shortUnits,omitempty"`
	ShortClientExtensions *ClientExtensions `json:"shortClientExtensions,omitempty"`
}

// NewPositionCloseRequest creates a new request with empty fields.
func NewPositionCloseRequest() *PositionCloseRequest {
	return &PositionCloseRequest{}
}

// SetLongAll sets the request to close all units of the long side.
func (r *PositionCloseRequest) SetLongAll() *PositionCloseRequest {
	v := "ALL"
	r.LongUnits = &v
	return r
}

// SetLongUnits sets the number of long units to close.
func (r *PositionCloseRequest) SetLongUnits(units uint) *PositionCloseRequest {
	v := strconv.FormatUint(uint64(units), 10)
	r.LongUnits = &v
	return r
}

// SetLongClientExtensions sets the client extensions for the long side close.
func (r *PositionCloseRequest) SetLongClientExtensions(extensions *ClientExtensions) *PositionCloseRequest {
	r.LongClientExtensions = extensions
	return r
}

// SetShortAll sets the request to close all units of the short side.
func (r *PositionCloseRequest) SetShortAll() *PositionCloseRequest {
	v := "ALL"
	r.ShortUnits = &v
	return r
}

// SetShortUnits sets the number of short units to close.
func (r *PositionCloseRequest) SetShortUnits(units uint) *PositionCloseRequest {
	v := strconv.FormatUint(uint64(units), 10)
	r.ShortUnits = &v
	return r
}

// SetShortClientExtensions sets the client extensions for the short side close.
func (r *PositionCloseRequest) SetShortClientExtensions(extensions *ClientExtensions) *PositionCloseRequest {
	r.ShortClientExtensions = extensions
	return r
}

func (r *PositionCloseRequest) body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

// PositionCloseResponse is the successful response returned by [positionService.Close].
type PositionCloseResponse struct {
	LongOrderCreateTransaction  *MarketOrderTransaction `json:"longOrderCreateTransaction,omitempty"`
	LongOrderFillTransaction    *OrderFillTransaction   `json:"longOrderFillTransaction,omitempty"`
	LongOrderCancelTransaction  *OrderCancelTransaction `json:"longOrderCancelTransaction,omitempty"`
	ShortOrderCreateTransaction *MarketOrderTransaction `json:"shortOrderCreateTransaction,omitempty"`
	ShortOrderFillTransaction   *OrderFillTransaction   `json:"shortOrderFillTransaction,omitempty"`
	ShortOrderCancelTransaction *OrderCancelTransaction `json:"shortOrderCancelTransaction,omitempty"`
	RelatedTransactionIDs       []TransactionID         `json:"relatedTransactionIDs"`
	LastTransactionID           TransactionID           `json:"lastTransactionID"`
}

// PositionCloseErrorResponse is the error response returned by [positionService.Close].
type PositionCloseErrorResponse struct {
	LongOrderRejectTransaction  *MarketOrderRejectTransaction `json:"longOrderRejectTransaction,omitempty"`
	ShortOrderRejectTransaction *MarketOrderRejectTransaction `json:"shortOrderRejectTransaction,omitempty"`
	RelatedTransactionIDs       []TransactionID               `json:"relatedTransactionIDs"`
	LastTransactionID           TransactionID                 `json:"lastTransactionID"`
	ErrorCode                   string                        `json:"errorCode"`
	ErrorMessage                string                        `json:"errorMessage"`
}

// Error implements the error interface.
func (r PositionCloseErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", r.ErrorCode, r.ErrorMessage)
}

// Close closes (fully or partially) the long and/or short side of a Position for a specific instrument.
//
// This corresponds to the OANDA API endpoint: PUT /v3/accounts/{accountID}/positions/{instrument}/close
//
// Reference: https://developer.oanda.com/rest-live-v20/position-ep/#collapse_endpoint_4
func (s *positionService) Close(ctx context.Context, instrument InstrumentName, req *PositionCloseRequest) (*PositionCloseResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/positions/%v/close", s.client.accountID, instrument)
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
		var resp PositionCloseResponse
		if err := decodeResponse(httpResp, &resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &resp, nil
	case http.StatusBadRequest:
		var resp PositionCloseErrorResponse
		if err := decodeResponse(httpResp, &resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return nil, BadRequest{HTTPError{StatusCode: httpResp.StatusCode, Message: "bad request", Err: resp}}
	case http.StatusNotFound:
		var resp PositionCloseErrorResponse
		if err := decodeResponse(httpResp, &resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return nil, NotFoundError{HTTPError{StatusCode: httpResp.StatusCode, Message: "not found", Err: resp}}
	default:
		return nil, decodeErrorResponse(httpResp)
	}
}
