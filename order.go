package oanda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

//
// Definitions https://developer.oanda.com/rest-live-v20/order-df/
//

// Orders

type Order interface {
	getType() OrderType
}

// OrderBase is the base specification of an Order as referred to by clients.
type OrderBase struct {
	// ID is the Order's identifier, unique within the Order's Account.
	ID OrderID `json:"id"`
	// CreateTime is the time when the Order was created.
	CreateTime DateTime `json:"createTime"`
	// State is the current state of the Order.
	State OrderState `json:"state"`
	// ClientExtensions are the client extensions for the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// Type is the type of the Order.
	Type OrderType `json:"type"`
}

type TradeClosingDetails struct {
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed when the price threshold is breached.
	ClientTradeID ClientID `json:"clientTradeID"`
}

type OrdersOnFill struct {
	// TakeProfitOnFill specifies the details of a Take Profit Order to be
	// created on behalf of a client. This may happen when an Order is filled
	// that opens a Trade requiring a Take Profit, or when a Trade’s dependent
	// Take Profit Order is modified directly through the Trade.
	TakeProfitOnFill *TakeProfitDetails `json:"takeProfitOnFill,omitempty"`
	// StopLossOnFill specifies the details of a Stop Loss Order to be created
	// on behalf of a client. This may happen when an Order is filled that opens
	// a Trade requiring a Stop Loss, or when a Trade’s dependent Stop Loss
	// Order is modified directly through the Trade.
	StopLossOnFill *StopLossDetails `json:"stopLossOnFill,omitempty"`
	// GuaranteedStopLossOnFill specifies the details of a Guaranteed Stop Loss
	// Order to be created on behalf of a client. This may happen when an Order
	// is filled that opens a Trade requiring a Guaranteed Stop Loss, or when a
	// Trade’s dependent Guaranteed Stop Loss Order is modified directly through
	// the Trade.
	GuaranteedStopLossOnFill *GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill,omitempty"`
	// TrailingStopLossOnFill specifies the details of a Trailing Stop Loss
	// Order to be created on behalf of a client. This may happen when an Order
	// is filled that opens a Trade requiring a Trailing Stop Loss, or when a
	// Trade’s dependent Trailing Stop Loss Order is modified directly through
	// the Trade.
	TrailingStopLossOnFill *TrailingStopLossDetails `json:"trailingStopLossOnFill,omitempty"`
	// TradeClientExtensions to add to the Trade created when the Order is filled
	// (if such a Trade is created). Do not set, modify, or delete
	// tradeClientExtensions if your account is associated with MT4.
	TradeClientExtensions *ClientExtensions `json:"tradeClientExtensions,omitempty"`
}

type PositionClosingDetails struct {
	// TradeClose is details of the Trade requested to be closed, only provided when the
	// MarketOrder is being used to explicitly close a Trade.
	TradeClose MarketOrderTradeClose `json:"tradeClose"`
	// LongPositionCloseout details the long Position to closeout when the Order is filled and
	// whether the long Position should be fully closed or only partially closed.
	LongPositionCloseout MarketOrderPositionCloseout `json:"longPositionCloseout"`
	// ShortPositionCloseout details the short Position to closeout when the Order is filled and
	// whether the short Position should be fully closed or only partially closed.
	ShortPositionCloseout MarketOrderPositionCloseout `json:"shortPositionCloseout"`
	// MarginCloseout details the Margin Closeout that this Market Order was created for.
	MarginCloseout MarketOrderMarginCloseout `json:"marginCloseout"`
	// DelayedTradeClose details the delayed Trade close that this Market Order was created for.
	DelayedTradeClose MarketOrderDelayedTradeClose `json:"delayedTradeClose"`
}

type FillingDetails struct {
	// FillingTransactionID is the ID of the Transaction that filled this Order (only provided when
	// the Order's state is FILLED).
	FillingTransactionID *TransactionID `json:"fillingTransactionID,omitempty"`
	// FilledTime is the date/time when the Order was filled (only provided when the Order's state is FILLED).
	FilledTime *DateTime `json:"filledTime,omitempty"`
}

type CancellingDetails struct {
	// CancellingTransactionID is the ID of the Transaction that cancelled the Order (only provided
	// when the Order's state is CANCELLED).
	CancellingTransactionID *TransactionID `json:"cancellingTransactionID,omitempty"`
	// CancelledTime is the date/time when the Order was cancelled (only provided when the Order's
	// state is CANCELLED).
	CancelledTime *DateTime `json:"cancelledTime,omitempty"`
}

type RelatedTradeIDs struct {
	// TradeOpenedID is the ID of the Trade opened when the Order was filled (only provided when the
	// Order's state is FILLED and a Trade was opened as a result of the fill).
	TradeOpenedID *TradeID `json:"tradeOpenedID,omitempty"`
	// TradeReducedID is the ID of the Trade reduced when the Order was filled (only provided when the
	// Order's state is FILLED and a Trade was reduced as a result of the fill).
	TradeReducedID *TradeID `json:"tradeReducedID,omitempty"`
	// TradeClosedIDs are the IDs of the Trades closed when the Order was filled (only provided when
	// the Order's state is FILLED and one or more Trades were closed as a result of the fill).
	TradeClosedIDs []TradeID `json:"tradeClosedIDs,omitempty"`
}

type ReplaceDetails struct {
	// ReplacesOrderID is the ID of the Order that was replaced by this Order (only provided if this
	// Order was created as part of a cancel/replace).
	ReplacesOrderID *OrderID `json:"replacesOrderID,omitempty"`
	// ReplacedByOrderID is the ID of the Order that replaced this Order (only provided if this Order
	// was cancelled as part of a cancel/replace).
	ReplacedByOrderID *OrderID `json:"replacedByOrderID,omitempty"`
}

// MarketOrder is an order that is filled immediately upon creation using the current market price.
type MarketOrder struct {
	OrderBase
	// Instrument is the name of the instrument of the Order.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Order. A positive
	// number of units results in a long Order, and a negative number of units
	// results in a short Order.
	Units DecimalNumber `json:"units"`
	// The TimeInForce requested for the Market Order. Restricted to FOK or
	// IOC for a MarketOrder
	TimeInForce TimeInForce `json:"timeInForce"`
	// PriceBound is the worst price that the client is willing to have the MarketOrder filled at.
	PriceBound PriceValue `json:"priceBound"`
	PositionClosingDetails
	OrdersOnFill
	FillingDetails
	RelatedTradeIDs
	CancellingDetails
}

func (m MarketOrder) getType() OrderType {
	return m.Type
}

// FixedPriceOrder is an order that is filled immediately upon creation using a fixed price.
type FixedPriceOrder struct {
	OrderBase
	// Instrument is the name of the instrument of the Order.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Order. A positive
	// number of units results in a long Order, and a negative number of units
	// results in a short Order.
	Units DecimalNumber `json:"units"`
	// Price specified for the Fixed Price Order. This price is the exact
	// price that the Fixed Price Order will be filled at.
	Price PriceValue `json:"price"`
	// PositionFill is specification of how Positions in the Account are modified when the Order
	// is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TradeState is the state that the trade resulting from the Fixed Price Order should be set to.
	TradeState string `json:"tradeState"`
	OrdersOnFill
	FillingDetails
	RelatedTradeIDs
	CancellingDetails
}

func (m FixedPriceOrder) getType() OrderType {
	return m.Type
}

// LimitOrder is an order that is created with a price threshold, and will only be filled by a price
// that is equal to or better than the threshold.
type LimitOrder struct {
	OrderBase
	// Instrument is the name of the instrument of the Order.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Order. A positive
	// number of units results in a long Order, and a negative number of units
	// results in a short Order.
	Units DecimalNumber `json:"units"`
	// The Price threshold specified for the Limit Order. The Limit Order will
	// only be filled by a market price that is equal to or better than this Price.
	Price PriceValue `json:"price"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// PositionFill is specification of how Positions in the Account are modified when the Order
	// is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	OrdersOnFill
	FillingDetails
	RelatedTradeIDs
	ReplaceDetails
}

func (m LimitOrder) getType() OrderType {
	return m.Type
}

// StopOrder is an order that is created with a price threshold, and will only be filled by a price
// that is equal to or worse than the threshold.
type StopOrder struct {
	OrderBase
	// Instrument is the name of the instrument of the Order.
	Instrument InstrumentName `json:"instrument"`
	// Price is the price threshold specified for the Stop Order. The Stop Order will only be filled
	// by a market price that is equal to or worse than this price.
	Price PriceValue `json:"price"`
	// PriceBound is the worst market price that may be used to fill this Stop Order. If the market
	// gaps and crosses through both the price and the priceBound, the Stop Order will be cancelled
	// instead of being filled.
	PriceBound PriceValue `json:"priceBound"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	OrdersOnFill
	FillingDetails
	RelatedTradeIDs
	CancellingDetails
	ReplaceDetails
}

func (m StopOrder) getType() OrderType {
	return m.Type
}

// MarketIfTouchedOrder is an order that is created with a price threshold, and will only be filled
// by a market price that touches or crosses the threshold.
type MarketIfTouchedOrder struct {
	OrderBase
	// Instrument is the Market If Touched Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Market If Touched Order. A positive number
	// of units results in a long Order, and a negative number of units results in a short Order.
	Units DecimalNumber `json:"units"`
	// Price is the price threshold specified for the Market If Touched Order. The Market If Touched
	// Order will only be filled by a market price that crosses this price from the direction of the
	// market price at the time when the Order was created (the initialMarketPrice). Depending on the
	// value of the Order's price and initialMarketPrice, the MarketIfTouchedOrder will behave like a
	// Limit or a Stop Order.
	Price PriceValue `json:"price"`
	// PriceBound is the worst market price that may be used to fill this Market If Touched Order.
	PriceBound PriceValue `json:"priceBound"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// InitialMarketPrice is the Market price at the time when the MarketIfTouched Order was created.
	InitialMarketPrice PriceValue `json:"initialMarketPrice"`
	OrdersOnFill
	FillingDetails
	RelatedTradeIDs
	CancellingDetails
	ReplaceDetails
}

func (m MarketIfTouchedOrder) getType() OrderType {
	return m.Type
}

// TakeProfitOrder is an order that is linked to an open Trade and created with a price threshold.
// The Order will be filled (closing the Trade) by the first price that is equal to or better than
// the threshold. A Take Profit Order cannot be used to open a new Position.
type TakeProfitOrder struct {
	OrderBase
	TradeClosingDetails
	// Price is the price threshold specified for the Take Profit Order. The associated Trade will be
	// closed by a market price that is equal to or better than this threshold.
	Price PriceValue `json:"price"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	FillingDetails
	RelatedTradeIDs
	CancellingDetails
	ReplaceDetails
}

func (t TakeProfitOrder) getType() OrderType {
	return t.Type
}

// StopLossOrder is an order that is linked to an open Trade and created with a price threshold.
// The Order will be filled (closing the Trade) by the first price that is equal to or worse than
// the threshold. A Stop Loss Order cannot be used to open a new Position.
type StopLossOrder struct {
	OrderBase
	TradeClosingDetails
	// Price is the price threshold specified for the Stop Loss Order. The associated Trade will be
	// closed by a market price that is equal to or worse than this threshold.
	Price PriceValue `json:"price"`
	// Distance specifies the distance (in price units) from the Account's current price to use as
	// the Stop Loss Order price. If the Trade is long the Order's price will be the bid price minus
	// the distance. If the Trade is short the Order's price will be the ask price plus the distance.
	Distance DecimalNumber `json:"distance"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	FillingDetails
	RelatedTradeIDs
	CancellingDetails
	ReplaceDetails
}

func (s StopLossOrder) getType() OrderType {
	return s.Type
}

// GuaranteedStopLossOrder is an order that is linked to an open Trade and created with a price threshold
// which is guaranteed against slippage that may occur as the market crosses the price set for that order.
// The Order will be filled (closing the Trade) by the first price that is equal to or worse than the threshold.
// The price level specified for the Guaranteed Stop Loss Order must be at least the configured minimum
// distance (in price units) away from the entry price for the traded instrument. A Guaranteed Stop Loss
// Order cannot be used to open new Positions.
type GuaranteedStopLossOrder struct {
	OrderBase
	// GuaranteedExecutionPremium is the premium that will be charged if the Guaranteed Stop Loss Order
	// is filled at the guaranteed price. It is in price units and is charged for each unit of the Trade.
	GuaranteedExecutionPremium DecimalNumber `json:"guaranteedExecutionPremium"`
	TradeClosingDetails
	// Price is the price threshold specified for the Guaranteed Stop Loss Order. The associated Trade
	// will be closed at this price.
	Price PriceValue `json:"price"`
	// Distance specifies the distance (in price units) from the Account's current price to use as
	// the Guaranteed Stop Loss Order price. If the Trade is long the Order's price will be the bid
	// price minus the distance. If the Trade is short the Order's price will be the ask price plus
	// the distance.
	Distance DecimalNumber `json:"distance"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component. This field
	// is restricted to "DEFAULT" for long Trades (BID price used), and "INVERSE" for short Trades
	// (ASK price used).
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	FillingDetails
	RelatedTradeIDs
	CancellingDetails
	ReplaceDetails
}

func (s GuaranteedStopLossOrder) getType() OrderType {
	return s.Type
}

// TrailingStopLossOrder is an order that is linked to an open Trade and created with a price distance.
// The price distance is used to calculate a trailing stop value for the order that is in the losing
// direction from the market price at the time of the order's creation. The trailing stop value will
// follow the market price as it moves in the winning direction, and the order will be filled (closing
// the Trade) by the first price that is equal to or worse than the trailing stop value. A Trailing
// Stop Loss Order cannot be used to open new Positions.
type TrailingStopLossOrder struct {
	OrderBase
	TradeClosingDetails
	// Distance is the price distance (in price units) specified for the Trailing Stop Loss Order.
	Distance DecimalNumber `json:"distance"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// TrailingStopValue is the trigger price for the Trailing Stop Loss Order. The trailing stop
	// value will trail (follow) the market price by the TSL order's configured distance as the market
	// price moves in the winning direction. If the market price moves to a level that is equal to or
	// worse than the trailing stop value, the order will be filled and the Trade will be closed.
	TrailingStopValue PriceValue `json:"trailingStopValue"`
	FillingDetails
	RelatedTradeIDs
	CancellingDetails
	ReplaceDetails
}

func (s TrailingStopLossOrder) getType() OrderType {
	return s.Type
}

// Order Requests

type OrderRequest interface {
	body() (*bytes.Buffer, error)
}

// MarketOrderRequest is used to create a Market Order.
type MarketOrderRequest struct {
	// Type is the type of the Order to Create. Must be set to "MARKET" when creating a Market Order.
	Type OrderType `json:"type"`
	// Instrument is the Market Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Market Order. A positive number of units
	// results in a long Order, and a negative number of units results in a short Order.
	Units DecimalNumber `json:"units"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. FOK or IOC are the only valid options for Market Orders.
	// Default is FOK.
	TimeInForce TimeInForce `json:"timeInForce"`
	// PriceBound is the worst price that the client is willing to have the Market Order filled at.
	PriceBound PriceValue `json:"priceBound,omitempty"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	// Default is DEFAULT.
	PositionFill OrderPositionFill `json:"positionFill"`
	// ClientExtensions are the client extensions to add to the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// TakeProfitOnFill specifies the details of a Take Profit Order to be created on behalf of a
	// client. This may happen when an Order is filled that opens a Trade requiring a Take Profit.
	TakeProfitOnFill *TakeProfitDetails `json:"takeProfitOnFill,omitempty"`
	// StopLossOnFill specifies the details of a Stop Loss Order to be created on behalf of a client.
	// This may happen when an Order is filled that opens a Trade requiring a Stop Loss.
	StopLossOnFill *StopLossDetails `json:"stopLossOnFill,omitempty"`
	// GuaranteedStopLossOnFill specifies the details of a Guaranteed Stop Loss Order to be created
	// on behalf of a client. This may happen when an Order is filled that opens a Trade requiring
	// a Guaranteed Stop Loss.
	GuaranteedStopLossOnFill *GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill,omitempty"`
	// TrailingStopLossOnFill specifies the details of a Trailing Stop Loss Order to be created on
	// behalf of a client. This may happen when an Order is filled that opens a Trade requiring a
	// Trailing Stop Loss.
	TrailingStopLossOnFill *TrailingStopLossDetails `json:"trailingStopLossOnFill,omitempty"`
	// TradeClientExtensions are the client extensions to add to the Trade created when the Order is filled.
	// Do not set, modify, or delete tradeClientExtensions if your account is associated with MT4.
	TradeClientExtensions *ClientExtensions `json:"tradeClientExtensions,omitempty"`
}

func (r *MarketOrderRequest) body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

func NewMarketOrderRequest(instrument InstrumentName, units DecimalNumber) *MarketOrderRequest {
	return &MarketOrderRequest{
		Type:         OrderTypeMarket,
		Instrument:   instrument,
		Units:        units,
		TimeInForce:  TimeInForceFOK,
		PositionFill: OrderPositionFillDefault,
	}
}

func (r *MarketOrderRequest) SetIOC() *MarketOrderRequest {
	r.TimeInForce = TimeInForceIOC
	return r
}

func (r *MarketOrderRequest) SetPriceBound(priceBound PriceValue) *MarketOrderRequest {
	r.PriceBound = priceBound
	return r
}

func (r *MarketOrderRequest) SetPositionFill(positionFill OrderPositionFill) *MarketOrderRequest {
	r.PositionFill = positionFill
	return r
}

func (r *MarketOrderRequest) SetClientExtensions(clientExtensions *ClientExtensions) *MarketOrderRequest {
	r.ClientExtensions = clientExtensions
	return r
}

func (r *MarketOrderRequest) SetTakeProfitOnFill(details *TakeProfitDetails) *MarketOrderRequest {
	r.TakeProfitOnFill = details
	return r
}

func (r *MarketOrderRequest) SetStopLossOnFill(details *StopLossDetails) *MarketOrderRequest {
	r.StopLossOnFill = details
	return r
}

func (r *MarketOrderRequest) SetGuaranteedStopLossOnFill(details *GuaranteedStopLossDetails) *MarketOrderRequest {
	r.GuaranteedStopLossOnFill = details
	return r
}

func (r *MarketOrderRequest) SetTrailingStopLossOnFill(details *TrailingStopLossDetails) *MarketOrderRequest {
	r.TrailingStopLossOnFill = details
	return r
}

func (r *MarketOrderRequest) SetTradeClientExtensions(clientExtensions *ClientExtensions) *MarketOrderRequest {
	r.ClientExtensions = clientExtensions
	return r
}

// LimitOrderRequest is used to create a Limit Order.
type LimitOrderRequest struct {
	// Type is the type of the Order to Create. Must be set to "LIMIT" when creating a Limit Order.
	Type OrderType `json:"type"`
	// Instrument is the Limit Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Limit Order. A positive number of units
	// results in a long Order, and a negative number of units results in a short Order.
	Units DecimalNumber `json:"units"`
	// Price is the price threshold specified for the Limit Order. The Limit Order will only be
	// filled by a market price that is equal to or better than this price.
	Price PriceValue `json:"price"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	// Default is DEFAULT.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component. Default is DEFAULT.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// ClientExtensions are the client extensions to add to the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// TakeProfitOnFill specifies the details of a Take Profit Order to be created on behalf of a
	// client. This may happen when an Order is filled that opens a Trade requiring a Take Profit.
	TakeProfitOnFill *TakeProfitDetails `json:"takeProfitOnFill,omitempty"`
	// StopLossOnFill specifies the details of a Stop Loss Order to be created on behalf of a client.
	// This may happen when an Order is filled that opens a Trade requiring a Stop Loss.
	StopLossOnFill *StopLossDetails `json:"stopLossOnFill,omitempty"`
	// GuaranteedStopLossOnFill specifies the details of a Guaranteed Stop Loss Order to be created
	// on behalf of a client. This may happen when an Order is filled that opens a Trade requiring
	// a Guaranteed Stop Loss.
	GuaranteedStopLossOnFill *GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill,omitempty"`
	// TrailingStopLossOnFill specifies the details of a Trailing Stop Loss Order to be created on
	// behalf of a client. This may happen when an Order is filled that opens a Trade requiring a
	// Trailing Stop Loss.
	TrailingStopLossOnFill *TrailingStopLossDetails `json:"trailingStopLossOnFill,omitempty"`
	// TradeClientExtensions are the client extensions to add to the Trade created when the Order is filled.
	// Do not set, modify, or delete tradeClientExtensions if your account is associated with MT4.
	TradeClientExtensions *ClientExtensions `json:"tradeClientExtensions,omitempty"`
}

func (r *LimitOrderRequest) body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

func NewLimitOrderRequest(instrument InstrumentName, units DecimalNumber, price PriceValue) *LimitOrderRequest {
	return &LimitOrderRequest{
		Type:             OrderTypeLimit,
		Instrument:       instrument,
		Units:            units,
		Price:            price,
		TimeInForce:      TimeInForceGTC,
		PositionFill:     OrderPositionFillDefault,
		TriggerCondition: OrderTriggerConditionDefault,
	}
}

func (r *LimitOrderRequest) SetGTD(date DateTime) *LimitOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = &date
	return r
}

func (r *LimitOrderRequest) SetGFD() *LimitOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

func (r *LimitOrderRequest) SetPositionFill(positionFill OrderPositionFill) *LimitOrderRequest {
	r.PositionFill = positionFill
	return r
}

func (r *LimitOrderRequest) SetTriggerCondition(triggerCondition OrderTriggerCondition) *LimitOrderRequest {
	r.TriggerCondition = triggerCondition
	return r
}

func (r *LimitOrderRequest) SetClientExtensions(clientExtensions *ClientExtensions) *LimitOrderRequest {
	r.ClientExtensions = clientExtensions
	return r
}

func (r *LimitOrderRequest) SetTakeProfitOnFill(details *TakeProfitDetails) *LimitOrderRequest {
	r.TakeProfitOnFill = details
	return r
}

func (r *LimitOrderRequest) SetStopLossOnFill(details *StopLossDetails) *LimitOrderRequest {
	r.StopLossOnFill = details
	return r
}

func (r *LimitOrderRequest) SetGuaranteedStopLossOnFill(details *GuaranteedStopLossDetails) *LimitOrderRequest {
	r.GuaranteedStopLossOnFill = details
	return r
}

func (r *LimitOrderRequest) SetTrailingStopLossOnFill(details *TrailingStopLossDetails) *LimitOrderRequest {
	r.TrailingStopLossOnFill = details
	return r
}

func (r *LimitOrderRequest) SetTradeClientExtensions(clientExtensions *ClientExtensions) *LimitOrderRequest {
	r.TradeClientExtensions = clientExtensions
	return r
}

// StopOrderRequest is used to create a Stop Order.
type StopOrderRequest struct {
	// Type is the type of the Order to Create. Must be set to "STOP" when creating a Stop Order.
	Type OrderType `json:"type"`
	// Instrument is the Stop Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Stop Order. A positive number of units
	// results in a long Order, and a negative number of units results in a short Order.
	Units DecimalNumber `json:"units"`
	// Price is the price threshold specified for the Stop Order. The Stop Order will only be filled
	// by a market price that is equal to or worse than this price.
	Price PriceValue `json:"price"`
	// PriceBound is the worst market price that may be used to fill this Stop Order. If the market
	// gaps and crosses through both the price and the priceBound, the Stop Order will be cancelled
	// instead of being filled.
	PriceBound PriceValue `json:"priceBound,omitempty"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime,omitempty"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	// Default is DEFAULT.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component. Default is DEFAULT.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// ClientExtensions are the client extensions to add to the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// TakeProfitOnFill specifies the details of a Take Profit Order to be created on behalf of a
	// client. This may happen when an Order is filled that opens a Trade requiring a Take Profit.
	TakeProfitOnFill *TakeProfitDetails `json:"takeProfitOnFill,omitempty"`
	// StopLossOnFill specifies the details of a Stop Loss Order to be created on behalf of a client.
	// This may happen when an Order is filled that opens a Trade requiring a Stop Loss.
	StopLossOnFill *StopLossDetails `json:"stopLossOnFill,omitempty"`
	// GuaranteedStopLossOnFill specifies the details of a Guaranteed Stop Loss Order to be created
	// on behalf of a client. This may happen when an Order is filled that opens a Trade requiring
	// a Guaranteed Stop Loss.
	GuaranteedStopLossOnFill *GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill,omitempty"`
	// TrailingStopLossOnFill specifies the details of a Trailing Stop Loss Order to be created on
	// behalf of a client. This may happen when an Order is filled that opens a Trade requiring a
	// Trailing Stop Loss.
	TrailingStopLossOnFill *TrailingStopLossDetails `json:"trailingStopLossOnFill,omitempty"`
	// TradeClientExtensions are the client extensions to add to the Trade created when the Order is filled.
	// Do not set, modify, or delete tradeClientExtensions if your account is associated with MT4.
	TradeClientExtensions *ClientExtensions `json:"tradeClientExtensions,omitempty"`
}

func (r *StopOrderRequest) body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

func NewStopOrderRequest(instrument InstrumentName, units DecimalNumber, price PriceValue) *StopOrderRequest {
	return &StopOrderRequest{
		Type:             OrderTypeStop,
		Instrument:       instrument,
		Units:            units,
		Price:            price,
		TimeInForce:      TimeInForceGTC,
		PositionFill:     OrderPositionFillDefault,
		TriggerCondition: OrderTriggerConditionDefault,
	}
}

func (r *StopOrderRequest) SetPriceBound(priceBound PriceValue) *StopOrderRequest {
	r.PriceBound = priceBound
	return r
}

func (r *StopOrderRequest) SetGTD(date DateTime) *StopOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = date
	return r
}

func (r *StopOrderRequest) SetGFD() *StopOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

func (r *StopOrderRequest) SetPositionFill(positionFill OrderPositionFill) *StopOrderRequest {
	r.PositionFill = positionFill
	return r
}

func (r *StopOrderRequest) SetTriggerCondition(triggerCondition OrderTriggerCondition) *StopOrderRequest {
	r.TriggerCondition = triggerCondition
	return r
}

func (r *StopOrderRequest) SetClientExtensions(clientExtensions *ClientExtensions) *StopOrderRequest {
	r.ClientExtensions = clientExtensions
	return r
}

func (r *StopOrderRequest) SetTakeProfitOnFill(details *TakeProfitDetails) *StopOrderRequest {
	r.TakeProfitOnFill = details
	return r
}

func (r *StopOrderRequest) SetStopLossOnFill(details *StopLossDetails) *StopOrderRequest {
	r.StopLossOnFill = details
	return r
}

func (r *StopOrderRequest) SetGuaranteedStopLossOnFill(details *GuaranteedStopLossDetails) *StopOrderRequest {
	r.GuaranteedStopLossOnFill = details
	return r
}

func (r *StopOrderRequest) SetTrailingStopLossOnFill(details *TrailingStopLossDetails) *StopOrderRequest {
	r.TrailingStopLossOnFill = details
	return r
}

func (r *StopOrderRequest) SetTradeClientExtensions(clientExtensions *ClientExtensions) *StopOrderRequest {
	r.TradeClientExtensions = clientExtensions
	return r
}

// MarketIfTouchedOrderRequest is used to create a Market If Touched Order.
type MarketIfTouchedOrderRequest struct {
	// Type is the type of the Order to Create. Must be set to "MARKET_IF_TOUCHED" when creating a
	// Market If Touched Order.
	Type OrderType `json:"type"`
	// Instrument is the Market If Touched Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Market If Touched Order. A positive number
	// of units results in a long Order, and a negative number of units results in a short Order.
	Units DecimalNumber `json:"units"`
	// Price is the price threshold specified for the Market If Touched Order. The Market If Touched
	// Order will only be filled by a market price that crosses this price from the direction of the
	// market price at the time when the Order was created (the initialMarketPrice). Depending on the
	// value of the Order's price and initialMarketPrice, the MarketIfTouchedOrder will behave like a
	// Limit or a Stop Order.
	Price PriceValue `json:"price"`
	// PriceBound is the worst market price that may be used to fill this Market If Touched Order.
	PriceBound PriceValue `json:"priceBound,omitempty"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime,omitempty"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	// Default is DEFAULT.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component. Default is DEFAULT.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// ClientExtensions are the client extensions to add to the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// TakeProfitOnFill specifies the details of a Take Profit Order to be created on behalf of a
	// client. This may happen when an Order is filled that opens a Trade requiring a Take Profit.
	TakeProfitOnFill *TakeProfitDetails `json:"takeProfitOnFill,omitempty"`
	// StopLossOnFill specifies the details of a Stop Loss Order to be created on behalf of a client.
	// This may happen when an Order is filled that opens a Trade requiring a Stop Loss.
	StopLossOnFill *StopLossDetails `json:"stopLossOnFill,omitempty"`
	// GuaranteedStopLossOnFill specifies the details of a Guaranteed Stop Loss Order to be created
	// on behalf of a client. This may happen when an Order is filled that opens a Trade requiring
	// a Guaranteed Stop Loss.
	GuaranteedStopLossOnFill *GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill,omitempty"`
	// TrailingStopLossOnFill specifies the details of a Trailing Stop Loss Order to be created on
	// behalf of a client. This may happen when an Order is filled that opens a Trade requiring a
	// Trailing Stop Loss.
	TrailingStopLossOnFill *TrailingStopLossDetails `json:"trailingStopLossOnFill,omitempty"`
	// TradeClientExtensions are the client extensions to add to the Trade created when the Order is filled.
	// Do not set, modify, or delete tradeClientExtensions if your account is associated with MT4.
	TradeClientExtensions *ClientExtensions `json:"tradeClientExtensions,omitempty"`
}

func (r *MarketIfTouchedOrderRequest) body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

func NewMarketIfTouchedOrderRequest(instrument InstrumentName, units DecimalNumber, price PriceValue) *MarketIfTouchedOrderRequest {
	return &MarketIfTouchedOrderRequest{
		Type:             OrderTypeMarketIfTouched,
		Instrument:       instrument,
		Units:            units,
		Price:            price,
		TimeInForce:      TimeInForceGTC,
		PositionFill:     OrderPositionFillDefault,
		TriggerCondition: OrderTriggerConditionDefault,
	}
}

func (r *MarketIfTouchedOrderRequest) SetPriceBound(priceBound PriceValue) *MarketIfTouchedOrderRequest {
	r.PriceBound = priceBound
	return r
}

func (r *MarketIfTouchedOrderRequest) SetGTD(date DateTime) *MarketIfTouchedOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = date
	return r
}

func (r *MarketIfTouchedOrderRequest) SetGFD() *MarketIfTouchedOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

func (r *MarketIfTouchedOrderRequest) SetOpenOnly() *MarketIfTouchedOrderRequest {
	r.PositionFill = OrderPositionFillOpenOnly
	return r
}

func (r *MarketIfTouchedOrderRequest) SetReduceFirst() *MarketIfTouchedOrderRequest {
	r.PositionFill = OrderPositionFillReduceFirst
	return r
}

func (r *MarketIfTouchedOrderRequest) SetReduceOnly() *MarketIfTouchedOrderRequest {
	r.PositionFill = OrderPositionFillReduceOnly
	return r
}

func (r *MarketIfTouchedOrderRequest) SetTriggerCondition(triggerCondition OrderTriggerCondition) *MarketIfTouchedOrderRequest {
	r.TriggerCondition = triggerCondition
	return r
}

func (r *MarketIfTouchedOrderRequest) SetClientExtensions(clientExtensions *ClientExtensions) *MarketIfTouchedOrderRequest {
	r.ClientExtensions = clientExtensions
	return r
}

func (r *MarketIfTouchedOrderRequest) SetTakeProfitOnFill(details *TakeProfitDetails) *MarketIfTouchedOrderRequest {
	r.TakeProfitOnFill = details
	return r
}

func (r *MarketIfTouchedOrderRequest) SetStopLossOnFill(details *StopLossDetails) *MarketIfTouchedOrderRequest {
	r.StopLossOnFill = details
	return r
}

func (r *MarketIfTouchedOrderRequest) SetGuaranteedStopLossOnFill(details *GuaranteedStopLossDetails) *MarketIfTouchedOrderRequest {
	r.GuaranteedStopLossOnFill = details
	return r
}

func (r *MarketIfTouchedOrderRequest) SetTrailingStopLossOnFill(details *TrailingStopLossDetails) *MarketIfTouchedOrderRequest {
	r.TrailingStopLossOnFill = details
	return r
}

func (r *MarketIfTouchedOrderRequest) SetTradeClientExtensions(clientExtensions *ClientExtensions) *MarketIfTouchedOrderRequest {
	r.TradeClientExtensions = clientExtensions
	return r
}

// TakeProfitOrderRequest is used to create a Take Profit Order.
type TakeProfitOrderRequest struct {
	// Type is the type of the Order to Create. Must be set to "TAKE_PROFIT" when creating a Take Profit Order.
	Type OrderType `json:"type"`
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed when the price threshold is breached.
	ClientTradeID ClientID `json:"clientTradeID,omitempty"`
	// Price is the price threshold specified for the Take Profit Order. The associated Trade will be
	// closed by a market price that is equal to or better than this threshold.
	Price PriceValue `json:"price"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component. Default is DEFAULT.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// ClientExtensions are the client extensions to add to the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
}

func (r *TakeProfitOrderRequest) body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

func NewTakeProfitOrderRequest(tradeID TradeID, price PriceValue) *TakeProfitOrderRequest {
	return &TakeProfitOrderRequest{
		Type:             OrderTypeTakeProfit,
		TradeID:          tradeID,
		Price:            price,
		TimeInForce:      TimeInForceGTC,
		TriggerCondition: OrderTriggerConditionDefault,
	}
}

func (r *TakeProfitOrderRequest) SetClientTradeID(clientID ClientID) *TakeProfitOrderRequest {
	r.ClientTradeID = clientID
	return r
}

func (r *TakeProfitOrderRequest) SetGTD(date DateTime) *TakeProfitOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = date
	return r
}

func (r *TakeProfitOrderRequest) SetGFD() *TakeProfitOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

func (r *TakeProfitOrderRequest) SetTriggerCondition(triggerCondition OrderTriggerCondition) *TakeProfitOrderRequest {
	r.TriggerCondition = triggerCondition
	return r
}

func (r *TakeProfitOrderRequest) SetClientExtensions(clientExtensions *ClientExtensions) *TakeProfitOrderRequest {
	r.ClientExtensions = clientExtensions
	return r
}

// StopLossOrderRequest is used to create a Stop Loss Order.
type StopLossOrderRequest struct {
	// Type is the type of the Order to Create. Must be set to "STOP_LOSS" when creating a Stop Loss Order.
	Type OrderType `json:"type"`
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed when the price threshold is breached.
	ClientTradeID ClientID `json:"clientTradeID,omitempty"`
	// Price is the price threshold specified for the Stop Loss Order. The associated Trade will be
	// closed by a market price that is equal to or worse than this threshold. Either price or distance
	// may be specified, but not both.
	Price PriceValue `json:"price,omitempty"`
	// Distance specifies the distance (in price units) from the Account's current price to use as
	// the Stop Loss Order price. If the Trade is long the Order's price will be the bid price minus
	// the distance. If the Trade is short the Order's price will be the ask price plus the distance.
	// Either price or distance may be specified, but not both.
	Distance DecimalNumber `json:"distance,omitempty"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component. Default is DEFAULT.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// ClientExtensions are the client extensions to add to the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
}

func (r *StopLossOrderRequest) body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

func NewStopLossOrderRequest(tradeID TradeID, price PriceValue) *StopLossOrderRequest {
	return &StopLossOrderRequest{
		Type:             OrderTypeStopLoss,
		TradeID:          tradeID,
		Price:            price,
		TimeInForce:      TimeInForceGTC,
		TriggerCondition: OrderTriggerConditionDefault,
	}
}

func (r *StopLossOrderRequest) SetClientTradeID(clientID ClientID) *StopLossOrderRequest {
	r.ClientTradeID = clientID
	return r
}

func (r *StopLossOrderRequest) SetGTD(date DateTime) *StopLossOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = date
	return r
}

func (r *StopLossOrderRequest) SetGFD() *StopLossOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

func (r *StopLossOrderRequest) SetTriggerCondition(triggerCondition OrderTriggerCondition) *StopLossOrderRequest {
	r.TriggerCondition = triggerCondition
	return r
}

func (r *StopLossOrderRequest) SetClientExtensions(clientExtensions *ClientExtensions) *StopLossOrderRequest {
	r.ClientExtensions = clientExtensions
	return r
}

// GuaranteedStopLossOrderRequest is used to create a Guaranteed Stop Loss Order.
type GuaranteedStopLossOrderRequest struct {
	// Type is the type of the Order to Create. Must be set to "GUARANTEED_STOP_LOSS" when creating
	// a Guaranteed Stop Loss Order.
	Type OrderType `json:"type"`
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed when the price threshold is breached.
	ClientTradeID ClientID `json:"clientTradeID,omitempty"`
	// Price is the price threshold specified for the Guaranteed Stop Loss Order. The associated Trade
	// will be closed at this price. Either price or distance may be specified, but not both.
	Price PriceValue `json:"price,omitempty"`
	// Distance specifies the distance (in price units) from the Account's current price to use as
	// the Guaranteed Stop Loss Order price. If the Trade is long the Order's price will be the bid
	// price minus the distance. If the Trade is short the Order's price will be the ask price plus
	// the distance. Either price or distance may be specified, but not both.
	Distance DecimalNumber `json:"distance,omitempty"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component. Default is DEFAULT.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// ClientExtensions are the client extensions to add to the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
}

func (r *GuaranteedStopLossOrderRequest) body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

func NewGuaranteedStopLossOrderRequest(tradeID TradeID, price PriceValue) *GuaranteedStopLossOrderRequest {
	return &GuaranteedStopLossOrderRequest{
		Type:             OrderTypeGuaranteedStopLoss,
		TradeID:          tradeID,
		Price:            price,
		TimeInForce:      TimeInForceGTC,
		TriggerCondition: OrderTriggerConditionDefault,
	}
}

func (r *GuaranteedStopLossOrderRequest) SetClientTradeID(clientID ClientID) *GuaranteedStopLossOrderRequest {
	r.ClientTradeID = clientID
	return r
}

func (r *GuaranteedStopLossOrderRequest) SetGTD(date DateTime) *GuaranteedStopLossOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = date
	return r
}

func (r *GuaranteedStopLossOrderRequest) SetGFD() *GuaranteedStopLossOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

func (r *GuaranteedStopLossOrderRequest) SetTriggerCondition(triggerCondition OrderTriggerCondition) *GuaranteedStopLossOrderRequest {
	r.TriggerCondition = triggerCondition
	return r
}

func (r *GuaranteedStopLossOrderRequest) SetClientExtensions(clientExtensions *ClientExtensions) *GuaranteedStopLossOrderRequest {
	r.ClientExtensions = clientExtensions
	return r
}

// TrailingStopLossOrderRequest is used to create a Trailing Stop Loss Order.
type TrailingStopLossOrderRequest struct {
	// Type is the type of the Order to Create. Must be set to "TRAILING_STOP_LOSS" when creating a
	// Trailing Stop Loss Order.
	Type OrderType `json:"type"`
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed when the price threshold is breached.
	ClientTradeID ClientID `json:"clientTradeID,omitempty"`
	// Distance is the price distance (in price units) specified for the Trailing Stop Loss Order.
	Distance DecimalNumber `json:"distance"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component. Default is DEFAULT.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// ClientExtensions are the client extensions to add to the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
}

func (r *TrailingStopLossOrderRequest) body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

func NewTrailingStopLossOrderRequest(tradeID TradeID, distance DecimalNumber) *TrailingStopLossOrderRequest {
	return &TrailingStopLossOrderRequest{
		Type:             OrderTypeTrailingStopLoss,
		TradeID:          tradeID,
		Distance:         distance,
		TimeInForce:      TimeInForceGTC,
		TriggerCondition: OrderTriggerConditionDefault,
	}
}

func (r *TrailingStopLossOrderRequest) SetClientTradeID(clientID ClientID) *TrailingStopLossOrderRequest {
	r.ClientTradeID = clientID
	return r
}

func (r *TrailingStopLossOrderRequest) SetGTD(date DateTime) *TrailingStopLossOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = date
	return r
}

func (r *TrailingStopLossOrderRequest) SetGFD() *TrailingStopLossOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

func (r *TrailingStopLossOrderRequest) SetTriggerCondition(triggerCondition OrderTriggerCondition) *TrailingStopLossOrderRequest {
	r.TriggerCondition = triggerCondition
	return r
}

func (r *TrailingStopLossOrderRequest) SetClientExtensions(clientExtensions *ClientExtensions) *TrailingStopLossOrderRequest {
	r.ClientExtensions = clientExtensions
	return r
}

// Order-related Definitions

// OrderID is the unique identifier for an Order within an Account.
type OrderID = string

// OrderType represents the type of an Order.
type OrderType string

const (
	OrderTypeMarket             OrderType = "MARKET"
	OrderTypeLimit              OrderType = "LIMIT"
	OrderTypeStop               OrderType = "STOP"
	OrderTypeMarketIfTouched    OrderType = "MARKET_IF_TOUCHED"
	OrderTypeFixedPrice         OrderType = "FIXED_PRICE"
	OrderTypeTakeProfit         OrderType = "TAKE_PROFIT"
	OrderTypeStopLoss           OrderType = "STOP_LOSS"
	OrderTypeGuaranteedStopLoss OrderType = "GUARANTEED_STOP_LOSS"
	OrderTypeTrailingStopLoss   OrderType = "TRAILING_STOP_LOSS"
)

// CancellableOrderType represents the type of Orders that can be cancelled.
// Market and FixedPrice orders cannot be cancelled as they are filled immediately.
type CancellableOrderType string

const (
	CancellableOrderTypeLimit              CancellableOrderType = "LIMIT"
	CancellableOrderTypeStop               CancellableOrderType = "STOP"
	CancellableOrderTypeMarketIfTouched    CancellableOrderType = "MARKET_IF_TOUCHED"
	CancellableOrderTypeTakeProfit         CancellableOrderType = "TAKE_PROFIT"
	CancellableOrderTypeStopLoss           CancellableOrderType = "STOP_LOSS"
	CancellableOrderTypeGuaranteedStopLoss CancellableOrderType = "GUARANTEED_STOP_LOSS"
	CancellableOrderTypeTrailingStopLoss   CancellableOrderType = "TRAILING_STOP_LOSS"
)

// OrderState represents the current state of an Order.
type OrderState string

const (
	OrderStatePending   OrderState = "PENDING"
	OrderStateFilled    OrderState = "FILLED"
	OrderStateTriggered OrderState = "TRIGGERED"
	OrderStateCancelled OrderState = "CANCELLED"
)

type OrderIdentifier struct {
	OrderID       OrderID  `json:"orderID"`
	ClientOrderID ClientID `json:"clientOrderID"`
}

type OrderSpecifier = string

// TimeInForce specifies how long an Order should remain pending before being automatically
// cancelled by the execution system.
type TimeInForce string

const (
	// TimeInForceFOK (Fill or Kill) means the Order must be immediately filled in its entirety
	// or be cancelled. Only valid for Market Orders.
	TimeInForceFOK TimeInForce = "FOK"
	// TimeInForceIOC (Immediate or Cancel) means the Order will be filled as much as possible
	// and any remaining units will be cancelled. Only valid for Market Orders.
	TimeInForceIOC TimeInForce = "IOC"
	// TimeInForceGTC (Good Till Cancelled) means the Order will remain pending until it is filled
	// or explicitly cancelled. Default for most pending orders.
	TimeInForceGTC TimeInForce = "GTC"
	// TimeInForceGTD (Good Till Date) means the Order will remain pending until the gtdTime or
	// until it is filled or cancelled.
	TimeInForceGTD TimeInForce = "GTD"
	// TimeInForceGFD (Good For Day) means the Order will remain pending until the end of the
	// trading day or until it is filled or cancelled.
	TimeInForceGFD TimeInForce = "GFD"
)

// OrderPositionFill specifies how Positions in the Account are modified when an Order is filled.
type OrderPositionFill string

const (
	// OrderPositionFillDefault means Positions are opened or closed using standard Account settings.
	OrderPositionFillDefault OrderPositionFill = "DEFAULT"
	// OrderPositionFillOpenOnly means the Order can only open new Positions.
	OrderPositionFillOpenOnly OrderPositionFill = "OPEN_ONLY"
	// OrderPositionFillReduceFirst means the Order will reduce an existing Position before opening
	// a new Position.
	OrderPositionFillReduceFirst OrderPositionFill = "REDUCE_FIRST"
	// OrderPositionFillReduceOnly means the Order can only reduce an existing Position.
	OrderPositionFillReduceOnly OrderPositionFill = "REDUCE_ONLY"
)

// OrderTriggerCondition specifies which price component should be used when determining if an Order
// should be triggered and filled.
type OrderTriggerCondition string

const (
	// OrderTriggerConditionDefault uses the ask price for long (buy) orders and bid price for
	// short (sell) orders.
	OrderTriggerConditionDefault OrderTriggerCondition = "DEFAULT"
	// OrderTriggerConditionInverse uses the bid price for long (buy) orders and ask price for
	// short (sell) orders.
	OrderTriggerConditionInverse OrderTriggerCondition = "INVERSE"
	// OrderTriggerConditionBid triggers the order when the bid price crosses the threshold.
	OrderTriggerConditionBid OrderTriggerCondition = "BID"
	// OrderTriggerConditionAsk triggers the order when the ask price crosses the threshold.
	OrderTriggerConditionAsk OrderTriggerCondition = "ASK"
	// OrderTriggerConditionMid triggers the order when the mid price crosses the threshold.
	OrderTriggerConditionMid OrderTriggerCondition = "MID"
)

//
// Endpoints https://developer.oanda.com/rest-live-v20/order-ep/
//

type OrderCreateResponse struct {
	OrderCreateTransaction        Transaction            `json:"orderCreateTransaction"`
	OrderFillTransaction          OrderFillTransaction   `json:"orderFillTransaction"`
	OrderCancelTransaction        OrderCancelTransaction `json:"orderCancelTransaction"`
	OrderReissueTransaction       Transaction            `json:"orderReissueTransaction"`
	OrderReissueRejectTransaction Transaction            `json:"orderReissueRejectTransaction"`
	RelatedTransactionIDs         []TransactionID        `json:"relatedTransactionIDs"`
	LastTransactionID             TransactionID          `json:"lastTransactionID"`
}

func orderRequestWrapper(req OrderRequest) (*bytes.Buffer, error) {
	request := struct {
		Order OrderRequest `json:"order"`
	}{Order: req}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(body), nil
}

func (c *Client) OrderCreate(ctx context.Context, accountID AccountID, req OrderRequest) (*OrderCreateResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/orders", accountID)
	body, err := orderRequestWrapper(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}
	resp, err := c.sendPostRequest(ctx, path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var orderCreateResponse OrderCreateResponse
	if err := decodeResponse(resp, &orderCreateResponse); err != nil {
		return nil, err
	}
	return &orderCreateResponse, nil
}

// OrderListRequest contains the parameters for retrieving a list of Orders for an Account.
// Use NewOrderListRequest to create a new request and the builder methods to configure options.
type OrderListRequest struct {
	AccountID  AccountID
	IDs        []OrderID
	State      *OrderState
	Instrument *InstrumentName
	Count      *int
	BeforeID   *OrderID
}

// NewOrderListRequest creates a new OrderListRequest for the specified account.
// Use the builder methods (AddIDs, SetState, SetInstrument, SetCount, SetBeforeID)
// to configure optional filtering parameters.
func NewOrderListRequest(accountID AccountID) *OrderListRequest {
	return &OrderListRequest{
		AccountID: accountID,
		IDs:       make([]OrderID, 0),
	}
}

// AddIDs adds Order IDs to filter the results. Only Orders with matching IDs will be returned.
func (r *OrderListRequest) AddIDs(ids ...OrderID) *OrderListRequest {
	r.IDs = append(r.IDs, ids...)
	return r
}

// SetState filters Orders by their state (PENDING, FILLED, TRIGGERED, CANCELLED).
func (r *OrderListRequest) SetState(state OrderState) *OrderListRequest {
	r.State = &state
	return r
}

// SetInstrument filters Orders by the specified instrument.
func (r *OrderListRequest) SetInstrument(instrument InstrumentName) *OrderListRequest {
	r.Instrument = &instrument
	return r
}

// SetCount sets the maximum number of Orders to return. Must be between 1 and 500.
func (r *OrderListRequest) SetCount(count int) *OrderListRequest {
	r.Count = &count
	return r
}

// SetBeforeID filters to return only Orders with an ID less than the specified ID.
// Used for pagination to retrieve older Orders.
func (r *OrderListRequest) SetBeforeID(beforeID OrderID) *OrderListRequest {
	r.BeforeID = &beforeID
	return r
}

func (r *OrderListRequest) validate() error {
	if r.Count != nil {
		if *r.Count <= 0 {
			return errors.New("count must be greater than zero")
		}
		if *r.Count > 500 {
			return errors.New("count must be less than or equal to 500")
		}
	}
	return nil
}

func (r *OrderListRequest) values() (url.Values, error) {
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

// OrderListResponse contains the response from the OrderList and OrderListPending endpoints.
type OrderListResponse struct {
	Orders            []Order       `json:"orders"`
	LastTransactionID TransactionID `json:"lastTransactionID"`
}

func unmarshalOrder(rawOrder json.RawMessage) (Order, error) {
	var typeOnly struct {
		Type OrderType `json:"type"`
	}
	if err := json.Unmarshal(rawOrder, &typeOnly); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order type: %w", err)
	}

	var order Order
	switch typeOnly.Type {
	case OrderTypeMarket:
		var marketOrder MarketOrder
		if err := json.Unmarshal(rawOrder, &marketOrder); err != nil {
			return nil, fmt.Errorf("failed to unmarshal market order: %w", err)
		}
		order = marketOrder
	case OrderTypeFixedPrice:
		var fixedPriceOrder FixedPriceOrder
		if err := json.Unmarshal(rawOrder, &fixedPriceOrder); err != nil {
			return nil, fmt.Errorf("failed to unmarshal fixed price order: %w", err)
		}
		order = fixedPriceOrder
	case OrderTypeLimit:
		var limitOrder LimitOrder
		if err := json.Unmarshal(rawOrder, &limitOrder); err != nil {
			return nil, fmt.Errorf("failed to unmarshal limit order: %w", err)
		}
		order = limitOrder
	case OrderTypeStop:
		var stopOrder StopOrder
		if err := json.Unmarshal(rawOrder, &stopOrder); err != nil {
			return nil, fmt.Errorf("failed to unmarshal stop order: %w", err)
		}
		order = stopOrder
	case OrderTypeMarketIfTouched:
		var marketIfTouchedOrder MarketIfTouchedOrder
		if err := json.Unmarshal(rawOrder, &marketIfTouchedOrder); err != nil {
			return nil, fmt.Errorf("failed to unmarshal market if touched order: %w", err)
		}
		order = marketIfTouchedOrder
	case OrderTypeTakeProfit:
		var takeProfitOrder TakeProfitOrder
		if err := json.Unmarshal(rawOrder, &takeProfitOrder); err != nil {
			return nil, fmt.Errorf("failed to unmarshal take profit order: %w", err)
		}
		order = takeProfitOrder
	case OrderTypeStopLoss:
		var stopLossOrder StopLossOrder
		if err := json.Unmarshal(rawOrder, &stopLossOrder); err != nil {
			return nil, fmt.Errorf("failed to unmarshal stop loss order: %w", err)
		}
		order = stopLossOrder
	case OrderTypeGuaranteedStopLoss:
		var guaranteedStopLossOrder GuaranteedStopLossOrder
		if err := json.Unmarshal(rawOrder, &guaranteedStopLossOrder); err != nil {
			return nil, fmt.Errorf("failed to unmarshal guaranteed stop loss order: %w", err)
		}
		order = guaranteedStopLossOrder
	case OrderTypeTrailingStopLoss:
		var trailingStopLossOrder TrailingStopLossOrder
		if err := json.Unmarshal(rawOrder, &trailingStopLossOrder); err != nil {
			return nil, fmt.Errorf("failed to unmarshal trailing stop loss order: %w", err)
		}
		order = trailingStopLossOrder
	}
	return order, nil
}

func (r *OrderListResponse) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Orders            []json.RawMessage `json:"orders"`
		LastTransactionID TransactionID     `json:"lastTransactionID"`
	}
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	r.Orders = make([]Order, 0, len(raw.Orders))

	for _, rawOrder := range raw.Orders {
		order, err := unmarshalOrder(rawOrder)
		if err != nil {
			return err
		}
		r.Orders = append(r.Orders, order)
	}
	raw.LastTransactionID = r.LastTransactionID
	return nil
}

// OrderList retrieves a list of Orders for an Account based on the specified request parameters.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/orders
//
// The request can be configured to filter by Order IDs, state, instrument, and to paginate
// results using count and beforeID parameters.
//
// Parameters:
//   - ctx: Context for the request.
//   - req: The OrderListRequest containing the account ID and optional filter parameters.
//     Use NewOrderListRequest to create and configure the request.
//
// Returns:
//   - []Order: A slice of Orders matching the request criteria.
//   - TransactionID: The ID of the most recent transaction created for the account.
//   - error: An error if the request fails, validation fails, or response cannot be decoded.
//
// Example:
//
//	req := oanda.NewOrderListRequest(accountID).
//	    SetState(oanda.OrderStatePending).
//	    SetInstrument("EUR_USD").
//	    SetCount(50)
//	orders, lastTxID, err := client.OrderList(ctx, req)
//
// Reference: https://developer.oanda.com/rest-live-v20/order-ep/#collapse_endpoint_2
func (c *Client) OrderList(ctx context.Context, req *OrderListRequest) (*OrderListResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/orders", req.AccountID)
	v, err := req.values()
	if err != nil {
		return nil, err
	}
	resp, err := c.sendGetRequest(ctx, path, v)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var orderListResp OrderListResponse
	if err := decodeResponse(resp, &orderListResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &orderListResp, nil
}

// OrderListPending retrieves all pending Orders in an Account.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/pendingOrders
//
// This is a convenience method that returns only Orders in the PENDING state.
// For more flexible filtering options, use OrderList with SetState(OrderStatePending).
//
// Parameters:
//   - ctx: Context for the request.
//   - accountID: The Account identifier.
//
// Returns:
//   - []Order: A slice of all pending Orders for the account.
//   - TransactionID: The ID of the most recent transaction created for the account.
//   - error: An error if the request fails or response cannot be decoded.
//
// Reference: https://developer.oanda.com/rest-live-v20/order-ep/#collapse_endpoint_3
func (c *Client) OrderListPending(ctx context.Context, accountID AccountID) (*OrderListResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/pendingOrders", accountID)
	resp, err := c.sendGetRequest(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var orderListResp OrderListResponse
	if err := decodeResponse(resp, &orderListResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &orderListResp, nil
}

type OrderDetailsResponse struct {
	Order             Order         `json:"order"`
	LastTransactionID TransactionID `json:"lastTransactionID"`
}

func (r *OrderDetailsResponse) UnmarshalJSON(b []byte) error {
	var raw struct {
		Order             json.RawMessage `json:"order"`
		LastTransactionID TransactionID   `json:"lastTransactionID"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	order, err := unmarshalOrder(raw.Order)
	if err != nil {
		return err
	}
	r.Order = order
	r.LastTransactionID = raw.LastTransactionID
	return nil
}

// OrderDetails retrieves the details of a single Order in an Account.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/orders/{orderSpecifier}
//
// Parameters:
//   - ctx: Context for the request.
//   - accountID: The Account identifier.
//   - specifier: The Order specifier, which can be either the Order's OANDA-assigned OrderID
//     (e.g., "1234") or the client-provided ClientID prefixed with "@" (e.g., "@my_order_id").
//
// Returns:
//   - *Order: The full details of the specified Order.
//   - TransactionID: The ID of the most recent transaction created for the account.
//   - error: An error if the request fails or response cannot be decoded.
//
// Reference: https://developer.oanda.com/rest-live-v20/order-ep/#collapse_endpoint_4
func (c *Client) OrderDetails(ctx context.Context, accountID AccountID, specifier OrderSpecifier) (*OrderDetailsResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/orders/%v", accountID, specifier)
	return doGet[OrderDetailsResponse](c, ctx, path, nil)
}

type OrderReplaceResponse struct {
	OrderCancelTransaction          OrderCancelTransaction `json:"orderCancelTransaction"`
	OrderCreateTransaction          Transaction            `json:"orderCreateTransaction"`
	OrderFillTransaction            OrderFillTransaction   `json:"orderFillTransaction"`
	OrderReissueTransaction         Transaction            `json:"orderReissueTransaction"`
	OrderReissueRejectTransaction   Transaction            `json:"orderReissueRejectTransaction"`
	ReplacingOrderCancelTransaction OrderCancelTransaction `json:"replacingOrderCancelTransaction"`
	RelatedTransactionIDs           []TransactionID        `json:"relatedTransactionIDs"`
	LastTransactionID               TransactionID          `json:"lastTransactionID"`
}

func (c *Client) OrderReplace(
	ctx context.Context, accountID AccountID, specifier OrderSpecifier, req OrderRequest,
) (*OrderReplaceResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/orders/%v", accountID, specifier)
	return doPut[OrderReplaceResponse](c, ctx, path, req)
}

type OrderCancelResponse struct {
	OrderCancelTransaction OrderCancelTransaction `json:"orderCancelTransaction"`
	RelatedTransactionIDs  []TransactionID        `json:"relatedTransactionIDs"`
	LastTransactionID      TransactionID          `json:"lastTransactionID"`
}

func (c *Client) OrderCancel(
	ctx context.Context, accountID AccountID, specifier OrderSpecifier,
) (*OrderCancelResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/orders/%v/cancel", accountID, specifier)
	return doPut[OrderCancelResponse](c, ctx, path, nil)
}

type OrderUpdateClientExtensionsRequest struct {
	ClientExtensions      ClientExtensions `json:"clientExtensions,omitempty"`
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions,omitempty"`
}

func (r *OrderUpdateClientExtensionsRequest) body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

type OrderUpdateClientExtensionsResponse struct {
}

func (c *Client) OrderUpdateClientExtensions(
	ctx context.Context,
	accountID AccountID,
	specifier OrderSpecifier,
	req OrderUpdateClientExtensionsRequest,
) (*OrderUpdateClientExtensionsResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/orders/%v/clientExtensions", accountID, specifier)
	return doPut[OrderUpdateClientExtensionsResponse](c, ctx, path, &req)
}
