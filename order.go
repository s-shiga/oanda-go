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

// Definitions https://developer.oanda.com/rest-live-v20/order-df/

// Orders

// Order is the base specification of an Order as referred to by clients.
type Order struct {
	// ID is the Order's identifier, unique within the Order's Account.
	ID OrderID `json:"id"`
	// CreateTime is the time when the Order was created.
	CreateTime DateTime `json:"createTime"`
	// State is the current state of the Order.
	State OrderState `json:"state"`
	// ClientExtensions are the client extensions for the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions `json:"clientExtensions"`
}

// MarketOrder is an order that is filled immediately upon creation using the current market price.
type MarketOrder struct {
	// ID is the Order's identifier, unique within the Order's Account.
	ID OrderID `json:"id"`
	// CreateTime is the time when the Order was created.
	CreateTime DateTime `json:"createTime"`
	// State is the current state of the Order.
	State OrderState `json:"state"`
	// ClientExtensions are the client extensions for the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions `json:"clientExtensions"`
	// Type is the type of the Order. Always set to "MARKET" for Market Orders.
	Type OrderType `json:"type"`
	// Instrument is the Market Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Market Order. A positive number of units
	// results in a long Order, and a negative number of units results in a short Order.
	Units DecimalNumber `json:"units"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// PriceBound is the worst price that the client is willing to have the Market Order filled at.
	PriceBound PriceValue `json:"priceBound"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TradeClose details the Trade ID and the client Trade ID of the Trade to be closed when the
	// Order is filled.
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
	// TakeProfitOnFill specifies the details of a Take Profit Order to be created on behalf of a
	// client. This may happen when an Order is filled that opens a Trade requiring a Take Profit.
	TakeProfitOnFill TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the details of a Stop Loss Order to be created on behalf of a client.
	// This may happen when an Order is filled that opens a Trade requiring a Stop Loss.
	StopLossOnFill StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the details of a Guaranteed Stop Loss Order to be created
	// on behalf of a client. This may happen when an Order is filled that opens a Trade requiring
	// a Guaranteed Stop Loss.
	GuaranteedStopLossOnFill GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the details of a Trailing Stop Loss Order to be created on
	// behalf of a client. This may happen when an Order is filled that opens a Trade requiring a
	// Trailing Stop Loss.
	TrailingStopLossOnFill TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions to add to the Trade created when the Order is filled.
	// Do not set, modify, or delete tradeClientExtensions if your account is associated with MT4.
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions"`
	// FillingTransactionID is the ID of the Transaction that filled this Order (only provided when
	// the Order's state is FILLED).
	FillingTransactionID TransactionID `json:"fillingTransactionID"`
	// FilledTime is the date/time when the Order was filled (only provided when the Order's state is FILLED).
	FilledTime DateTime `json:"filledTime"`
	// TradeOpenedID is the ID of the Trade opened when the Order was filled (only provided when the
	// Order's state is FILLED and a Trade was opened as a result of the fill).
	TradeOpenedID TradeID `json:"tradeOpenedID"`
	// TradeReducedID is the ID of the Trade reduced when the Order was filled (only provided when the
	// Order's state is FILLED and a Trade was reduced as a result of the fill).
	TradeReducedID TradeID `json:"tradeReducedID"`
	// TradeClosedIDs are the IDs of the Trades closed when the Order was filled (only provided when
	// the Order's state is FILLED and one or more Trades were closed as a result of the fill).
	TradeClosedIDs []TradeID `json:"tradeClosedIDs"`
	// CancellingTransactionID is the ID of the Transaction that cancelled the Order (only provided
	// when the Order's state is CANCELLED).
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
	// CancelledTime is the date/time when the Order was cancelled (only provided when the Order's
	// state is CANCELLED).
	CancelledTime DateTime `json:"cancelledTime"`
}

// FixedPriceOrder is an order that is filled immediately upon creation using a fixed price.
type FixedPriceOrder struct {
	// ID is the Order's identifier, unique within the Order's Account.
	ID OrderID `json:"id"`
	// CreateTime is the time when the Order was created.
	CreateTime DateTime `json:"createTime"`
	// State is the current state of the Order.
	State OrderState `json:"state"`
	// ClientExtensions are the client extensions for the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions `json:"clientExtensions"`
	// Type is the type of the Order. Always set to "FIXED_PRICE" for Fixed Price Orders.
	Type OrderType `json:"type"`
	// Instrument is the Fixed Price Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Fixed Price Order. A positive number of
	// units results in a long Order, and a negative number of units results in a short Order.
	Units DecimalNumber `json:"units"`
	// Price is the price specified for the Fixed Price Order. This price is the exact price that
	// the Fixed Price Order will be filled at.
	Price PriceValue `json:"price"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TradeState is the state that the trade resulting from the Fixed Price Order should be set to.
	TradeState string `json:"tradeState"`
	// TakeProfitOnFill specifies the details of a Take Profit Order to be created on behalf of a
	// client. This may happen when an Order is filled that opens a Trade requiring a Take Profit.
	TakeProfitOnFill TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the details of a Stop Loss Order to be created on behalf of a client.
	// This may happen when an Order is filled that opens a Trade requiring a Stop Loss.
	StopLossOnFill StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the details of a Guaranteed Stop Loss Order to be created
	// on behalf of a client. This may happen when an Order is filled that opens a Trade requiring
	// a Guaranteed Stop Loss.
	GuaranteedStopLossOnFill GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the details of a Trailing Stop Loss Order to be created on
	// behalf of a client. This may happen when an Order is filled that opens a Trade requiring a
	// Trailing Stop Loss.
	TrailingStopLossOnFill TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions to add to the Trade created when the Order is filled.
	// Do not set, modify, or delete tradeClientExtensions if your account is associated with MT4.
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions"`
	// FillingTransactionID is the ID of the Transaction that filled this Order (only provided when
	// the Order's state is FILLED).
	FillingTransactionID TransactionID `json:"fillingTransactionID"`
	// FilledTime is the date/time when the Order was filled (only provided when the Order's state is FILLED).
	FilledTime DateTime `json:"filledTime"`
	// TradeOpenedID is the ID of the Trade opened when the Order was filled (only provided when the
	// Order's state is FILLED and a Trade was opened as a result of the fill).
	TradeOpenedID TradeID `json:"tradeOpenedID"`
	// TradeReducedID is the ID of the Trade reduced when the Order was filled (only provided when the
	// Order's state is FILLED and a Trade was reduced as a result of the fill).
	TradeReducedID TradeID `json:"tradeReducedID"`
	// TradeClosedIDs are the IDs of the Trades closed when the Order was filled (only provided when
	// the Order's state is FILLED and one or more Trades were closed as a result of the fill).
	TradeClosedIDs []TradeID `json:"tradeClosedIDs"`
	// CancellingTransactionID is the ID of the Transaction that cancelled the Order (only provided
	// when the Order's state is CANCELLED).
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
	// CancelledTime is the date/time when the Order was cancelled (only provided when the Order's
	// state is CANCELLED).
	CancelledTime DateTime `json:"cancelledTime"`
}

// LimitOrder is an order that is created with a price threshold, and will only be filled by a price
// that is equal to or better than the threshold.
type LimitOrder struct {
	// ID is the Order's identifier, unique within the Order's Account.
	ID OrderID `json:"id"`
	// CreateTime is the time when the Order was created.
	CreateTime DateTime `json:"createTime"`
	// State is the current state of the Order.
	State OrderState `json:"state"`
	// ClientExtensions are the client extensions for the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions `json:"clientExtensions"`
	// Type is the type of the Order. Always set to "LIMIT" for Limit Orders.
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
	// cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// TakeProfitOnFill specifies the details of a Take Profit Order to be created on behalf of a
	// client. This may happen when an Order is filled that opens a Trade requiring a Take Profit.
	TakeProfitOnFill TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the details of a Stop Loss Order to be created on behalf of a client.
	// This may happen when an Order is filled that opens a Trade requiring a Stop Loss.
	StopLossOnFill StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the details of a Guaranteed Stop Loss Order to be created
	// on behalf of a client. This may happen when an Order is filled that opens a Trade requiring
	// a Guaranteed Stop Loss.
	GuaranteedStopLossOnFill GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the details of a Trailing Stop Loss Order to be created on
	// behalf of a client. This may happen when an Order is filled that opens a Trade requiring a
	// Trailing Stop Loss.
	TrailingStopLossOnFill TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions to add to the Trade created when the Order is filled.
	// Do not set, modify, or delete tradeClientExtensions if your account is associated with MT4.
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions"`
	// FillingTransactionID is the ID of the Transaction that filled this Order (only provided when
	// the Order's state is FILLED).
	FillingTransactionID TransactionID `json:"fillingTransactionID"`
	// FilledTime is the date/time when the Order was filled (only provided when the Order's state is FILLED).
	FilledTime DateTime `json:"filledTime"`
	// TradeOpenedID is the ID of the Trade opened when the Order was filled (only provided when the
	// Order's state is FILLED and a Trade was opened as a result of the fill).
	TradeOpenedID TradeID `json:"tradeOpenedID"`
	// TradeReducedID is the ID of the Trade reduced when the Order was filled (only provided when the
	// Order's state is FILLED and a Trade was reduced as a result of the fill).
	TradeReducedID TradeID `json:"tradeReducedID"`
	// TradeClosedIDs are the IDs of the Trades closed when the Order was filled (only provided when
	// the Order's state is FILLED and one or more Trades were closed as a result of the fill).
	TradeClosedIDs []TradeID `json:"tradeClosedIDs"`
	// CancellingTransactionID is the ID of the Transaction that cancelled the Order (only provided
	// when the Order's state is CANCELLED).
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
	// CancelledTime is the date/time when the Order was cancelled (only provided when the Order's
	// state is CANCELLED).
	CancelledTime DateTime `json:"cancelledTime"`
	// ReplacesOrderID is the ID of the Order that was replaced by this Order (only provided if this
	// Order was created as part of a cancel/replace).
	ReplacesOrderID OrderID `json:"replacesOrderID"`
	// ReplacedByOrderID is the ID of the Order that replaced this Order (only provided if this Order
	// was cancelled as part of a cancel/replace).
	ReplacedByOrderID OrderID `json:"replacedByOrderID"`
}

// StopOrder is an order that is created with a price threshold, and will only be filled by a price
// that is equal to or worse than the threshold.
type StopOrder struct {
	// ID is the Order's identifier, unique within the Order's Account.
	ID OrderID `json:"id"`
	// CreateTime is the time when the Order was created.
	CreateTime DateTime `json:"createTime"`
	// State is the current state of the Order.
	State OrderState `json:"state"`
	// ClientExtensions are the client extensions for the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions `json:"clientExtensions"`
	// Type is the type of the Order. Always set to "STOP" for Stop Orders.
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
	PriceBound PriceValue `json:"priceBound"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// TakeProfitOnFill specifies the details of a Take Profit Order to be created on behalf of a
	// client. This may happen when an Order is filled that opens a Trade requiring a Take Profit.
	TakeProfitOnFill TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the details of a Stop Loss Order to be created on behalf of a client.
	// This may happen when an Order is filled that opens a Trade requiring a Stop Loss.
	StopLossOnFill StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the details of a Guaranteed Stop Loss Order to be created
	// on behalf of a client. This may happen when an Order is filled that opens a Trade requiring
	// a Guaranteed Stop Loss.
	GuaranteedStopLossOnFill GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the details of a Trailing Stop Loss Order to be created on
	// behalf of a client. This may happen when an Order is filled that opens a Trade requiring a
	// Trailing Stop Loss.
	TrailingStopLossOnFill TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions to add to the Trade created when the Order is filled.
	// Do not set, modify, or delete tradeClientExtensions if your account is associated with MT4.
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions"`
	// FillingTransactionID is the ID of the Transaction that filled this Order (only provided when
	// the Order's state is FILLED).
	FillingTransactionID TransactionID `json:"fillingTransactionID"`
	// FilledTime is the date/time when the Order was filled (only provided when the Order's state is FILLED).
	FilledTime DateTime `json:"filledTime"`
	// TradeOpenedID is the ID of the Trade opened when the Order was filled (only provided when the
	// Order's state is FILLED and a Trade was opened as a result of the fill).
	TradeOpenedID TradeID `json:"tradeOpenedID"`
	// TradeReducedID is the ID of the Trade reduced when the Order was filled (only provided when the
	// Order's state is FILLED and a Trade was reduced as a result of the fill).
	TradeReducedID TradeID `json:"tradeReducedID"`
	// TradeClosedIDs are the IDs of the Trades closed when the Order was filled (only provided when
	// the Order's state is FILLED and one or more Trades were closed as a result of the fill).
	TradeClosedIDs []TradeID `json:"tradeClosedIDs"`
	// CancellingTransactionID is the ID of the Transaction that cancelled the Order (only provided
	// when the Order's state is CANCELLED).
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
	// CancelledTime is the date/time when the Order was cancelled (only provided when the Order's
	// state is CANCELLED).
	CancelledTime DateTime `json:"cancelledTime"`
	// ReplacesOrderID is the ID of the Order that was replaced by this Order (only provided if this
	// Order was created as part of a cancel/replace).
	ReplacesOrderID OrderID `json:"replacesOrderID"`
	// ReplacedByOrderID is the ID of the Order that replaced this Order (only provided if this Order
	// was cancelled as part of a cancel/replace).
	ReplacedByOrderID OrderID `json:"replacedByOrderID"`
}

// MarketIfTouchedOrder is an order that is created with a price threshold, and will only be filled
// by a market price that touches or crosses the threshold.
type MarketIfTouchedOrder struct {
	// ID is the Order's identifier, unique within the Order's Account.
	ID OrderID `json:"id"`
	// CreateTime is the time when the Order was created.
	CreateTime DateTime `json:"createTime"`
	// State is the current state of the Order.
	State OrderState `json:"state"`
	// ClientExtensions are the client extensions for the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions `json:"clientExtensions"`
	// Type is the type of the Order. Always set to "MARKET_IF_TOUCHED" for Market If Touched Orders.
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
	PriceBound PriceValue `json:"priceBound"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// InitialMarketPrice is the Market price at the time when the MarketIfTouched Order was created.
	InitialMarketPrice PriceValue `json:"initialMarketPrice"`
	// TakeProfitOnFill specifies the details of a Take Profit Order to be created on behalf of a
	// client. This may happen when an Order is filled that opens a Trade requiring a Take Profit.
	TakeProfitOnFill TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the details of a Stop Loss Order to be created on behalf of a client.
	// This may happen when an Order is filled that opens a Trade requiring a Stop Loss.
	StopLossOnFill StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the details of a Guaranteed Stop Loss Order to be created
	// on behalf of a client. This may happen when an Order is filled that opens a Trade requiring
	// a Guaranteed Stop Loss.
	GuaranteedStopLossOnFill GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the details of a Trailing Stop Loss Order to be created on
	// behalf of a client. This may happen when an Order is filled that opens a Trade requiring a
	// Trailing Stop Loss.
	TrailingStopLossOnFill TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions to add to the Trade created when the Order is filled.
	// Do not set, modify, or delete tradeClientExtensions if your account is associated with MT4.
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions"`
	// FillingTransactionID is the ID of the Transaction that filled this Order (only provided when
	// the Order's state is FILLED).
	FillingTransactionID TransactionID `json:"fillingTransactionID"`
	// FilledTime is the date/time when the Order was filled (only provided when the Order's state is FILLED).
	FilledTime DateTime `json:"filledTime"`
	// TradeOpenedID is the ID of the Trade opened when the Order was filled (only provided when the
	// Order's state is FILLED and a Trade was opened as a result of the fill).
	TradeOpenedID TradeID `json:"tradeOpenedID"`
	// TradeReducedID is the ID of the Trade reduced when the Order was filled (only provided when the
	// Order's state is FILLED and a Trade was reduced as a result of the fill).
	TradeReducedID TradeID `json:"tradeReducedID"`
	// TradeClosedIDs are the IDs of the Trades closed when the Order was filled (only provided when
	// the Order's state is FILLED and one or more Trades were closed as a result of the fill).
	TradeClosedIDs []TradeID `json:"tradeClosedIDs"`
	// CancellingTransactionID is the ID of the Transaction that cancelled the Order (only provided
	// when the Order's state is CANCELLED).
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
	// CancelledTime is the date/time when the Order was cancelled (only provided when the Order's
	// state is CANCELLED).
	CancelledTime DateTime `json:"cancelledTime"`
	// ReplacesOrderID is the ID of the Order that was replaced by this Order (only provided if this
	// Order was created as part of a cancel/replace).
	ReplacesOrderID OrderID `json:"replacesOrderID"`
	// ReplacedByOrderID is the ID of the Order that replaced this Order (only provided if this Order
	// was cancelled as part of a cancel/replace).
	ReplacedByOrderID OrderID `json:"replacedByOrderID"`
}

// TakeProfitOrder is an order that is linked to an open Trade and created with a price threshold.
// The Order will be filled (closing the Trade) by the first price that is equal to or better than
// the threshold. A Take Profit Order cannot be used to open a new Position.
type TakeProfitOrder struct {
	// ID is the Order's identifier, unique within the Order's Account.
	ID OrderID `json:"id"`
	// CreateTime is the time when the Order was created.
	CreateTime DateTime `json:"createTime"`
	// State is the current state of the Order.
	State OrderState `json:"state"`
	// ClientExtensions are the client extensions for the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions `json:"clientExtensions"`
	// Type is the type of the Order. Always set to "TAKE_PROFIT" for Take Profit Orders.
	Type OrderType `json:"type"`
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed when the price threshold is breached.
	ClientTradeID ClientID `json:"clientTradeID"`
	// Price is the price threshold specified for the Take Profit Order. The associated Trade will be
	// closed by a market price that is equal to or better than this threshold.
	Price PriceValue `json:"price"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// FillingTransactionID is the ID of the Transaction that filled this Order (only provided when
	// the Order's state is FILLED).
	FillingTransactionID TransactionID `json:"fillingTransactionID"`
	// FilledTime is the date/time when the Order was filled (only provided when the Order's state is FILLED).
	FilledTime DateTime `json:"filledTime"`
	// TradeClosedID is the ID of the Trade that was closed when the Order was filled (only provided
	// when the Order's state is FILLED and a Trade was closed as a result of the fill).
	TradeClosedID TradeID `json:"tradeClosedID"`
	// CancellingTransactionID is the ID of the Transaction that cancelled the Order (only provided
	// when the Order's state is CANCELLED).
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
	// CancelledTime is the date/time when the Order was cancelled (only provided when the Order's
	// state is CANCELLED).
	CancelledTime DateTime `json:"cancelledTime"`
	// ReplacesOrderID is the ID of the Order that was replaced by this Order (only provided if this
	// Order was created as part of a cancel/replace).
	ReplacesOrderID OrderID `json:"replacesOrderID"`
	// ReplacedByOrderID is the ID of the Order that replaced this Order (only provided if this Order
	// was cancelled as part of a cancel/replace).
	ReplacedByOrderID OrderID `json:"replacedByOrderID"`
}

// StopLossOrder is an order that is linked to an open Trade and created with a price threshold.
// The Order will be filled (closing the Trade) by the first price that is equal to or worse than
// the threshold. A Stop Loss Order cannot be used to open a new Position.
type StopLossOrder struct {
	// ID is the Order's identifier, unique within the Order's Account.
	ID OrderID `json:"id"`
	// CreateTime is the time when the Order was created.
	CreateTime DateTime `json:"createTime"`
	// State is the current state of the Order.
	State OrderState `json:"state"`
	// ClientExtensions are the client extensions for the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions `json:"clientExtensions"`
	// Type is the type of the Order. Always set to "STOP_LOSS" for Stop Loss Orders.
	Type OrderType `json:"type"`
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed when the price threshold is breached.
	ClientTradeID ClientID `json:"clientTradeID"`
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
	GtdTime DateTime `json:"gtdTime"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Guaranteed is a deprecated field. Indicates if the Stop Loss Order is guaranteed. The default
	// value depends on the GuaranteedStopLossOrderMode of the account. If true, the Stop Loss Order
	// is guaranteed and the Order's guaranteed execution premium will be charged if the Order is filled.
	// If false, the Stop Loss Order is not guaranteed and the Order's guaranteed execution premium
	// will not be charged if the Order is filled.
	Guaranteed bool `json:"guaranteed"`
	// GuaranteedExecutionPremium is the premium that will be charged if the Stop Loss Order is
	// guaranteed and the Order is filled at the guaranteed price. The value is determined at Order
	// creation time. It is in price units and is charged for each unit of the Trade.
	GuaranteedExecutionPremium DecimalNumber `json:"guaranteedExecutionPremium"`
	// FillingTransactionID is the ID of the Transaction that filled this Order (only provided when
	// the Order's state is FILLED).
	FillingTransactionID TransactionID `json:"fillingTransactionID"`
	// FilledTime is the date/time when the Order was filled (only provided when the Order's state is FILLED).
	FilledTime DateTime `json:"filledTime"`
	// TradeClosedID is the ID of the Trade that was closed when the Order was filled (only provided
	// when the Order's state is FILLED and a Trade was closed as a result of the fill).
	TradeClosedID TradeID `json:"tradeClosedID"`
	// CancellingTransactionID is the ID of the Transaction that cancelled the Order (only provided
	// when the Order's state is CANCELLED).
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
	// CancelledTime is the date/time when the Order was cancelled (only provided when the Order's
	// state is CANCELLED).
	CancelledTime DateTime `json:"cancelledTime"`
	// ReplacesOrderID is the ID of the Order that was replaced by this Order (only provided if this
	// Order was created as part of a cancel/replace).
	ReplacesOrderID OrderID `json:"replacesOrderID"`
	// ReplacedByOrderID is the ID of the Order that replaced this Order (only provided if this Order
	// was cancelled as part of a cancel/replace).
	ReplacedByOrderID OrderID `json:"replacedByOrderID"`
}

// GuaranteedStopLossOrder is an order that is linked to an open Trade and created with a price threshold
// which is guaranteed against slippage that may occur as the market crosses the price set for that order.
// The Order will be filled (closing the Trade) by the first price that is equal to or worse than the threshold.
// The price level specified for the Guaranteed Stop Loss Order must be at least the configured minimum
// distance (in price units) away from the entry price for the traded instrument. A Guaranteed Stop Loss
// Order cannot be used to open new Positions.
type GuaranteedStopLossOrder struct {
	// ID is the Order's identifier, unique within the Order's Account.
	ID OrderID `json:"id"`
	// CreateTime is the time when the Order was created.
	CreateTime DateTime `json:"createTime"`
	// State is the current state of the Order.
	State OrderState `json:"state"`
	// ClientExtensions are the client extensions for the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions `json:"clientExtensions"`
	// Type is the type of the Order. Always set to "GUARANTEED_STOP_LOSS" for Guaranteed Stop Loss Orders.
	Type OrderType `json:"type"`
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed when the price threshold is breached.
	ClientTradeID ClientID `json:"clientTradeID"`
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
	GtdTime DateTime `json:"gtdTime"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component. This field
	// is restricted to "DEFAULT" for long Trades (BID price used), and "INVERSE" for short Trades
	// (ASK price used).
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// GuaranteedExecutionPremium is the premium that will be charged if the Guaranteed Stop Loss Order
	// is filled at the guaranteed price. It is in price units and is charged for each unit of the Trade.
	GuaranteedExecutionPremium DecimalNumber `json:"guaranteedExecutionPremium"`
	// FillingTransactionID is the ID of the Transaction that filled this Order (only provided when
	// the Order's state is FILLED).
	FillingTransactionID TransactionID `json:"fillingTransactionID"`
	// FilledTime is the date/time when the Order was filled (only provided when the Order's state is FILLED).
	FilledTime DateTime `json:"filledTime"`
	// TradeClosedID is the ID of the Trade that was closed when the Order was filled (only provided
	// when the Order's state is FILLED and a Trade was closed as a result of the fill).
	TradeClosedID TradeID `json:"tradeClosedID"`
	// CancellingTransactionID is the ID of the Transaction that cancelled the Order (only provided
	// when the Order's state is CANCELLED).
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
	// CancelledTime is the date/time when the Order was cancelled (only provided when the Order's
	// state is CANCELLED).
	CancelledTime DateTime `json:"cancelledTime"`
	// ReplacesOrderID is the ID of the Order that was replaced by this Order (only provided if this
	// Order was created as part of a cancel/replace).
	ReplacesOrderID OrderID `json:"replacesOrderID"`
	// ReplacedByOrderID is the ID of the Order that replaced this Order (only provided if this Order
	// was cancelled as part of a cancel/replace).
	ReplacedByOrderID OrderID `json:"replacedByOrderID"`
}

// TrailingStopLossOrder is an order that is linked to an open Trade and created with a price distance.
// The price distance is used to calculate a trailing stop value for the order that is in the losing
// direction from the market price at the time of the order's creation. The trailing stop value will
// follow the market price as it moves in the winning direction, and the order will be filled (closing
// the Trade) by the first price that is equal to or worse than the trailing stop value. A Trailing
// Stop Loss Order cannot be used to open new Positions.
type TrailingStopLossOrder struct {
	// ID is the Order's identifier, unique within the Order's Account.
	ID OrderID `json:"id"`
	// CreateTime is the time when the Order was created.
	CreateTime DateTime `json:"createTime"`
	// State is the current state of the Order.
	State OrderState `json:"state"`
	// ClientExtensions are the client extensions for the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions `json:"clientExtensions"`
	// Type is the type of the Order. Always set to "TRAILING_STOP_LOSS" for Trailing Stop Loss Orders.
	Type OrderType `json:"type"`
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed when the price threshold is breached.
	ClientTradeID ClientID `json:"clientTradeID"`
	// Distance is the price distance (in price units) specified for the Trailing Stop Loss Order.
	Distance DecimalNumber `json:"distance"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
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
	// FillingTransactionID is the ID of the Transaction that filled this Order (only provided when
	// the Order's state is FILLED).
	FillingTransactionID TransactionID `json:"fillingTransactionID"`
	// FilledTime is the date/time when the Order was filled (only provided when the Order's state is FILLED).
	FilledTime DateTime `json:"filledTime"`
	// TradeClosedID is the ID of the Trade that was closed when the Order was filled (only provided
	// when the Order's state is FILLED and a Trade was closed as a result of the fill).
	TradeClosedID TradeID `json:"tradeClosedID"`
	// CancellingTransactionID is the ID of the Transaction that cancelled the Order (only provided
	// when the Order's state is CANCELLED).
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
	// CancelledTime is the date/time when the Order was cancelled (only provided when the Order's
	// state is CANCELLED).
	CancelledTime DateTime `json:"cancelledTime"`
	// ReplacesOrderID is the ID of the Order that was replaced by this Order (only provided if this
	// Order was created as part of a cancel/replace).
	ReplacesOrderID OrderID `json:"replacesOrderID"`
	// ReplacedByOrderID is the ID of the Order that replaced this Order (only provided if this Order
	// was cancelled as part of a cancel/replace).
	ReplacedByOrderID OrderID `json:"replacedByOrderID"`
}

// Order Requests

type OrderRequest interface {
	Body() (*bytes.Buffer, error)
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
	PositionFill OrderPositionFill `json:"positionFill,omitempty"`
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

func (r *MarketOrderRequest) Body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

func NewMarketOrderRequest(instrument InstrumentName, units DecimalNumber) *MarketOrderRequest {
	return &MarketOrderRequest{
		Type:        OrderTypeMarket,
		Instrument:  instrument,
		Units:       units,
		TimeInForce: TimeInForceFOK,
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
	GtdTime DateTime `json:"gtdTime,omitempty"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	// Default is DEFAULT.
	PositionFill OrderPositionFill `json:"positionFill,omitempty"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component. Default is DEFAULT.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition,omitempty"`
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

func (r *LimitOrderRequest) Body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

func NewLimitOrderRequest(instrument InstrumentName, units DecimalNumber, price PriceValue) *LimitOrderRequest {
	return &LimitOrderRequest{
		Type:        OrderTypeLimit,
		Instrument:  instrument,
		Units:       units,
		Price:       price,
		TimeInForce: TimeInForceGTC,
	}
}

func (r *LimitOrderRequest) SetGTD(date DateTime) *LimitOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = date
	return r
}

func (r *LimitOrderRequest) SetGFD() *LimitOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

func (r *LimitOrderRequest) SetOpenOnly() *LimitOrderRequest {
	r.PositionFill = OrderPositionFillOpenOnly
	return r
}

func (r *LimitOrderRequest) SetReduceFirst() *LimitOrderRequest {
	r.PositionFill = OrderPositionFillReduceFirst
	return r
}

func (r *LimitOrderRequest) SetReduceOnly() *LimitOrderRequest {
	r.PositionFill = OrderPositionFillReduceOnly
	return r
}

func (r *LimitOrderRequest) SetInverse() *LimitOrderRequest {
	r.TriggerCondition = OrderTriggerConditionInverse
	return r
}

func (r *LimitOrderRequest) SetBid() *LimitOrderRequest {
	r.TriggerCondition = OrderTriggerConditionBid
	return r
}

func (r *LimitOrderRequest) SetAsk() *LimitOrderRequest {
	r.TriggerCondition = OrderTriggerConditionAsk
	return r
}

func (r *LimitOrderRequest) SetMid() *LimitOrderRequest {
	r.TriggerCondition = OrderTriggerConditionMid
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
	GtdTime DateTime `json:"gtdTime"`
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
	ClientExtensions *ClientExtensions `json:"clientExtensions"`
	// TakeProfitOnFill specifies the details of a Take Profit Order to be created on behalf of a
	// client. This may happen when an Order is filled that opens a Trade requiring a Take Profit.
	TakeProfitOnFill *TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the details of a Stop Loss Order to be created on behalf of a client.
	// This may happen when an Order is filled that opens a Trade requiring a Stop Loss.
	StopLossOnFill *StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the details of a Guaranteed Stop Loss Order to be created
	// on behalf of a client. This may happen when an Order is filled that opens a Trade requiring
	// a Guaranteed Stop Loss.
	GuaranteedStopLossOnFill *GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the details of a Trailing Stop Loss Order to be created on
	// behalf of a client. This may happen when an Order is filled that opens a Trade requiring a
	// Trailing Stop Loss.
	TrailingStopLossOnFill *TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions to add to the Trade created when the Order is filled.
	// Do not set, modify, or delete tradeClientExtensions if your account is associated with MT4.
	TradeClientExtensions *ClientExtensions `json:"tradeClientExtensions"`
}

func (r *StopOrderRequest) Body() (*bytes.Buffer, error) {
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

func (r *StopOrderRequest) SetOpenOnly() *StopOrderRequest {
	r.PositionFill = OrderPositionFillOpenOnly
	return r
}

func (r *StopOrderRequest) SetReduceFirst() *StopOrderRequest {
	r.PositionFill = OrderPositionFillReduceFirst
	return r
}

func (r *StopOrderRequest) SetReduceOnly() *StopOrderRequest {
	r.PositionFill = OrderPositionFillReduceOnly
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
	PriceBound PriceValue `json:"priceBound"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
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
	ClientExtensions *ClientExtensions `json:"clientExtensions"`
	// TakeProfitOnFill specifies the details of a Take Profit Order to be created on behalf of a
	// client. This may happen when an Order is filled that opens a Trade requiring a Take Profit.
	TakeProfitOnFill *TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the details of a Stop Loss Order to be created on behalf of a client.
	// This may happen when an Order is filled that opens a Trade requiring a Stop Loss.
	StopLossOnFill *StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the details of a Guaranteed Stop Loss Order to be created
	// on behalf of a client. This may happen when an Order is filled that opens a Trade requiring
	// a Guaranteed Stop Loss.
	GuaranteedStopLossOnFill *GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the details of a Trailing Stop Loss Order to be created on
	// behalf of a client. This may happen when an Order is filled that opens a Trade requiring a
	// Trailing Stop Loss.
	TrailingStopLossOnFill *TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions to add to the Trade created when the Order is filled.
	// Do not set, modify, or delete tradeClientExtensions if your account is associated with MT4.
	TradeClientExtensions *ClientExtensions `json:"tradeClientExtensions"`
}

func (r *MarketIfTouchedOrderRequest) Body() (*bytes.Buffer, error) {
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
	ClientTradeID ClientID `json:"clientTradeID"`
	// Price is the price threshold specified for the Take Profit Order. The associated Trade will be
	// closed by a market price that is equal to or better than this threshold.
	Price PriceValue `json:"price"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component. Default is DEFAULT.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// ClientExtensions are the client extensions to add to the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions *ClientExtensions `json:"clientExtensions"`
}

func (r *TakeProfitOrderRequest) Body() (*bytes.Buffer, error) {
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
	ClientTradeID ClientID `json:"clientTradeID"`
	// Price is the price threshold specified for the Stop Loss Order. The associated Trade will be
	// closed by a market price that is equal to or worse than this threshold. Either price or distance
	// may be specified, but not both.
	Price PriceValue `json:"price"`
	// Distance specifies the distance (in price units) from the Account's current price to use as
	// the Stop Loss Order price. If the Trade is long the Order's price will be the bid price minus
	// the distance. If the Trade is short the Order's price will be the ask price plus the distance.
	// Either price or distance may be specified, but not both.
	Distance DecimalNumber `json:"distance"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component. Default is DEFAULT.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Guaranteed is a deprecated field. Indicates if the Stop Loss Order is guaranteed. The default
	// value depends on the GuaranteedStopLossOrderMode of the account. If true, the Stop Loss Order
	// is guaranteed and the Order's guaranteed execution premium will be charged if the Order is filled.
	// If false, the Stop Loss Order is not guaranteed and the Order's guaranteed execution premium
	// will not be charged if the Order is filled.
	Guaranteed bool `json:"guaranteed"`
	// ClientExtensions are the client extensions to add to the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
}

func (r StopLossOrderRequest) Body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

// GuaranteedStopLossOrderRequest is used to create a Guaranteed Stop Loss Order.
type GuaranteedStopLossOrderRequest struct {
	// Type is the type of the Order to Create. Must be set to "GUARANTEED_STOP_LOSS" when creating
	// a Guaranteed Stop Loss Order.
	Type OrderType `json:"type"`
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed when the price threshold is breached.
	ClientTradeID ClientID `json:"clientTradeID"`
	// Price is the price threshold specified for the Guaranteed Stop Loss Order. The associated Trade
	// will be closed at this price. Either price or distance may be specified, but not both.
	Price PriceValue `json:"price"`
	// Distance specifies the distance (in price units) from the Account's current price to use as
	// the Guaranteed Stop Loss Order price. If the Trade is long the Order's price will be the bid
	// price minus the distance. If the Trade is short the Order's price will be the ask price plus
	// the distance. Either price or distance may be specified, but not both.
	Distance DecimalNumber `json:"distance"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component. Default is DEFAULT.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// ClientExtensions are the client extensions to add to the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
}

func (r GuaranteedStopLossOrderRequest) Body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

// TrailingStopLossOrderRequest is used to create a Trailing Stop Loss Order.
type TrailingStopLossOrderRequest struct {
	// Type is the type of the Order to Create. Must be set to "TRAILING_STOP_LOSS" when creating a
	// Trailing Stop Loss Order.
	Type OrderType `json:"type"`
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed when the price threshold is breached.
	ClientTradeID ClientID `json:"clientTradeID"`
	// Distance is the price distance (in price units) specified for the Trailing Stop Loss Order.
	Distance DecimalNumber `json:"distance"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// TriggerCondition specifies which price component should be used when determining if an Order
	// should be triggered and filled. This allows Orders to be triggered based on the bid, ask, mid,
	// default (ask for buy, bid for sell) or inverse (bid for buy, ask for sell) price depending on
	// the desired behaviour. Orders are always filled using their default price component. Default is DEFAULT.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// ClientExtensions are the client extensions to add to the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
}

func (r TrailingStopLossOrderRequest) Body() (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
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

// Endpoints https://developer.oanda.com/rest-live-v20/order-ep/

func (c *Client) OrderCreate(ctx context.Context, req *OrderRequest) {}

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
func (req *OrderListRequest) AddIDs(ids ...OrderID) *OrderListRequest {
	req.IDs = append(req.IDs, ids...)
	return req
}

// SetState filters Orders by their state (PENDING, FILLED, TRIGGERED, CANCELLED).
func (req *OrderListRequest) SetState(state OrderState) *OrderListRequest {
	req.State = &state
	return req
}

// SetInstrument filters Orders by the specified instrument.
func (req *OrderListRequest) SetInstrument(instrument InstrumentName) *OrderListRequest {
	req.Instrument = &instrument
	return req
}

// SetCount sets the maximum number of Orders to return. Must be between 1 and 500.
func (req *OrderListRequest) SetCount(count int) *OrderListRequest {
	req.Count = &count
	return req
}

// SetBeforeID filters to return only Orders with an ID less than the specified ID.
// Used for pagination to retrieve older Orders.
func (req *OrderListRequest) SetBeforeID(beforeID OrderID) *OrderListRequest {
	req.BeforeID = &beforeID
	return req
}

func (req *OrderListRequest) validate() error {
	if req.Count != nil {
		if *req.Count <= 0 {
			return errors.New("count must be greater than zero")
		}
		if *req.Count > 500 {
			return errors.New("count must be less than or equal to 500")
		}
	}
	return nil
}

func (req *OrderListRequest) values() (url.Values, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if len(req.IDs) > 0 {
		v.Set("ids", strings.Join(req.IDs, ","))
	}
	if req.State != nil {
		v.Set("state", string(*req.State))
	}
	if req.Instrument != nil {
		v.Set("instrument", *req.Instrument)
	}
	if req.Count != nil {
		v.Set("count", strconv.Itoa(*req.Count))
	}
	if req.BeforeID != nil {
		v.Set("beforeID", *req.BeforeID)
	}
	return v, nil
}

// OrderListResponse contains the response from the OrderList and OrderListPending endpoints.
type OrderListResponse struct {
	Orders            []Order       `json:"orders"`
	LastTransactionID TransactionID `json:"lastTransactionID"`
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
func (c *Client) OrderList(ctx context.Context, req *OrderListRequest) ([]Order, TransactionID, error) {
	path := fmt.Sprintf("/v3/accounts/%v/orders", req.AccountID)
	v, err := req.values()
	if err != nil {
		return nil, "", err
	}
	resp, err := c.sendGetRequest(ctx, path, v)
	if err != nil {
		return nil, "", fmt.Errorf("failed to send request: %w", err)
	}
	var orderListResp OrderListResponse
	if err := decodeResponse(resp, &orderListResp); err != nil {
		return nil, "", fmt.Errorf("failed to decode response: %w", err)
	}
	return orderListResp.Orders, orderListResp.LastTransactionID, nil
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
func (c *Client) OrderListPending(ctx context.Context, accountID AccountID) ([]Order, TransactionID, error) {
	path := fmt.Sprintf("/v3/accounts/%v/pendingOrders", accountID)
	resp, err := c.sendGetRequest(ctx, path, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to send request: %w", err)
	}
	var orderListResp OrderListResponse
	if err := decodeResponse(resp, &orderListResp); err != nil {
		return nil, "", fmt.Errorf("failed to decode response: %w", err)
	}
	return orderListResp.Orders, orderListResp.LastTransactionID, nil
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
func (c *Client) OrderDetails(ctx context.Context, accountID AccountID, specifier OrderSpecifier) (*Order, TransactionID, error) {
	path := fmt.Sprintf("/v3/accounts/%v/orders/%v", accountID, specifier)
	resp, err := c.sendGetRequest(ctx, path, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to send request: %w", err)
	}
	orderDetailsResponse := struct {
		Order             Order         `json:"order"`
		LastTransactionID TransactionID `json:"lastTransactionID"`
	}{}
	if err := decodeResponse(resp, &orderDetailsResponse); err != nil {
		return nil, "", fmt.Errorf("failed to decode response: %w", err)
	}
	return &orderDetailsResponse.Order, orderDetailsResponse.LastTransactionID, nil
}
