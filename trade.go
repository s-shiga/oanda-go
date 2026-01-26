package oanda

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Definitions https://developer.oanda.com/rest-live-v20/trade-df/

// TradeID is the unique identifier for a Trade within an Account. It is a string representation
// of the OANDA-assigned TradeID, derived from the TransactionID of the Transaction that opened
// the Trade. Example: "1523"
type TradeID = string

// TradeSpecifier is the specification of a Trade as referred to by clients. Either the Trade's
// OANDA-assigned TradeID or the Trade's client-provided ClientID prefixed by the "@" symbol.
type TradeSpecifier string

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
	AverageClosePrice PriceValue `json:"averageClosePrice"`
	// ClosingTransactionIDs is the list of Transaction IDs associated with closing portions of
	// this Trade.
	ClosingTransactionIDs []TransactionID `json:"closingTransactionIDs"`
	// Financing is the financing paid/collected for this Trade.
	Financing AccountUnits `json:"financing"`
	// DividendAdjustment is the dividend adjustment paid or collected for this Trade.
	DividendAdjustment AccountUnits `json:"dividendAdjustment"`
	// CloseTime is the date/time when the Trade was fully closed. Only provided for Trades whose
	// state is CLOSED.
	CloseTime DateTime `json:"closeTime"`
	// ClientExtensions are the client extensions of the Trade.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
	// TakeProfitOrder is the full representation of the Trade's Take Profit Order, only provided
	// if such an Order exists.
	TakeProfitOrder TakeProfitOrder `json:"takeProfitOrder"`
	// StopLossOrder is the full representation of the Trade's Stop Loss Order, only provided if
	// such an Order exists.
	StopLossOrder StopLossOrder `json:"stopLossOrder"`
	// TrailingStopLossOrder is the full representation of the Trade's Trailing Stop Loss Order,
	// only provided if such an Order exists.
	TrailingStopLossOrder TrailingStopLossOrder `json:"trailingStopLossOrder"`
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
	ClosingTransactionIDs []TransactionID `json:"closingTransactionIDs"`
	// Financing is the financing paid/collected for this Trade.
	Financing AccountUnits `json:"financing"`
	// DividendAdjustment is the dividend adjustment paid or collected for this Trade.
	DividendAdjustment AccountUnits `json:"dividendAdjustment"`
	// CloseTime is the date/time when the Trade was fully closed. Only provided for Trades whose
	// state is CLOSED.
	CloseTime DateTime `json:"closeTime"`
	// ClientExtensions are the client extensions of the Trade.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
	// TakeProfitOrderID is the ID of the Trade's Take Profit Order, only provided if such an
	// Order exists.
	TakeProfitOrderID OrderID `json:"takeProfitOrderID"`
	// StopLossOrderID is the ID of the Trade's Stop Loss Order, only provided if such an Order exists.
	StopLossOrderID OrderID `json:"stopLossOrderID"`
	// GuaranteedStopLossOrderID is the ID of the Trade's Guaranteed Stop Loss Order, only provided
	// if such an Order exists.
	GuaranteedStopLossOrderID OrderID `json:"guaranteedStopLossOrderID"`
	// TrailingStopLossOrderID is the ID of the Trade's Trailing Stop Loss Order, only provided if
	// such an Order exists.
	TrailingStopLossOrderID OrderID `json:"trailingStopLossOrderID"`
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

// Endpoints https://developer.oanda.com/rest-live-v20/trade-ep/

type TradeListRequest struct {
	AccountID  AccountID
	IDs        []TradeID
	State      *TradeStateFilter
	Instrument *InstrumentName
	Count      *int
	BeforeID   *TradeID
}

func NewTradeListRequest(accountID AccountID) *TradeListRequest {
	return &TradeListRequest{
		AccountID: accountID,
		IDs:       make([]TradeID, 0),
	}
}

func (r *TradeListRequest) AddIDs(id ...TradeID) *TradeListRequest {
	r.IDs = append(r.IDs, id...)
	return r
}

func (r *TradeListRequest) SetStateFilter(filter TradeStateFilter) *TradeListRequest {
	r.State = &filter
	return r
}

func (r *TradeListRequest) SetInstrument(instrument InstrumentName) *TradeListRequest {
	r.Instrument = &instrument
	return r
}

func (r *TradeListRequest) SetCount(count int) *TradeListRequest {
	r.Count = &count
	return r
}

func (r *TradeListRequest) SetBeforeID(beforeID TradeID) *TradeListRequest {
	r.BeforeID = &beforeID
	return r
}

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

func (c *Client) TradeList(ctx context.Context, req *TradeListRequest) ([]Trade, TransactionID, error) {
	path := fmt.Sprintf("/v3/accounts/%s/trades", req.AccountID)
	v, err := req.values()
	if err != nil {
		return nil, "", err
	}
	resp, err := c.sendGetRequest(ctx, path, v)
	if err != nil {
		return nil, "", fmt.Errorf("failed to send request: %w", err)
	}
	tradeListResp := struct {
		Trades             []Trade       `json:"trades"`
		LastTransactionsID TransactionID `json:"lastTransactionsID"`
	}{}
	if err := decodeResponse(resp, &tradeListResp); err != nil {
		return nil, "", fmt.Errorf("failed to decode response: %w", err)
	}
	return tradeListResp.Trades, tradeListResp.LastTransactionsID, nil
}
