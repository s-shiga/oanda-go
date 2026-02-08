package oanda

import (
	"context"
	"fmt"
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

type PositionListResponse struct {
	Positions         []Position    `json:"positions"`
	LastTransactionID TransactionID `json:"lastTransactionId"`
}

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

type PositionListByInstrumentResponse struct {
	Position          Position      `json:"position"`
	LastTransactionID TransactionID `json:"lastTransactionID"`
}

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
