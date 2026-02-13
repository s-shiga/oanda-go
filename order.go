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
// Definitions https://developer.oanda.com/rest-live-v20/order-df/
// ---------------------------------------------------------------

// Orders

// Order is the interface implemented by all order types returned by the OANDA v20 API.
type Order interface {
	GetID() OrderID
	GetCreateTime() DateTime
	GetState() OrderState
	GetClientExtensions() *ClientExtensions
	GetType() OrderType
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

func unmarshalOrders(src []json.RawMessage) ([]Order, error) {
	dest := make([]Order, 0, len(src))
	for _, rawOrder := range src {
		order, err := unmarshalOrder(rawOrder)
		if err != nil {
			return nil, err
		}
		dest = append(dest, order)
	}
	return dest, nil
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

// TradeClosingDetails contains the Trade ID and client Trade ID of the Trade to close.
type TradeClosingDetails struct {
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed when the price threshold is breached.
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
}

// OrdersOnFill contains details of dependent Orders to create when a Trade is opened by filling this Order.
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

// PositionClosingDetails contains details of Trade/Position closeouts associated with a MarketOrder.
type PositionClosingDetails struct {
	// TradeClose is details of the Trade requested to be closed, only provided when the
	// MarketOrder is being used to explicitly close a Trade.
	TradeClose *MarketOrderTradeClose `json:"tradeClose,omitempty"`
	// LongPositionCloseout details the long Position to closeout when the Order is filled and
	// whether the long Position should be fully closed or only partially closed.
	LongPositionCloseout *MarketOrderPositionCloseout `json:"longPositionCloseout,omitempty"`
	// ShortPositionCloseout details the short Position to closeout when the Order is filled and
	// whether the short Position should be fully closed or only partially closed.
	ShortPositionCloseout *MarketOrderPositionCloseout `json:"shortPositionCloseout,omitempty"`
	// MarginCloseout details the Margin Closeout that this Market Order was created for.
	MarginCloseout *MarketOrderMarginCloseout `json:"marginCloseout,omitempty"`
	// DelayedTradeClose details the delayed Trade close that this Market Order was created for.
	DelayedTradeClose *MarketOrderDelayedTradeClose `json:"delayedTradeClose,omitempty"`
}

// FillingDetails contains the Transaction ID and time when an Order was filled.
type FillingDetails struct {
	// FillingTransactionID is the ID of the Transaction that filled this Order (only provided when
	// the Order's state is FILLED).
	FillingTransactionID *TransactionID `json:"fillingTransactionID,omitempty"`
	// FilledTime is the date/time when the Order was filled (only provided when the Order's state is FILLED).
	FilledTime *DateTime `json:"filledTime,omitempty"`
}

// CancellingDetails contains the Transaction ID and time when an Order was cancelled.
type CancellingDetails struct {
	// CancellingTransactionID is the ID of the Transaction that cancelled the Order (only provided
	// when the Order's state is CANCELLED).
	CancellingTransactionID *TransactionID `json:"cancellingTransactionID,omitempty"`
	// CancelledTime is the date/time when the Order was cancelled (only provided when the Order's
	// state is CANCELLED).
	CancelledTime *DateTime `json:"cancelledTime,omitempty"`
}

// RelatedTradeIDs contains the IDs of Trades opened, reduced, or closed when an Order was filled.
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

// ReplaceDetails contains IDs linking Orders involved in a cancel/replace operation.
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
	PriceBound *PriceValue `json:"priceBound,omitempty"`
	PositionClosingDetails
	OrdersOnFill
	FillingDetails
	RelatedTradeIDs
	CancellingDetails
}

func (o MarketOrder) GetID() OrderID {
	return o.ID
}

func (o MarketOrder) GetCreateTime() DateTime {
	return o.CreateTime
}

func (o MarketOrder) GetState() OrderState {
	return o.State
}

func (o MarketOrder) GetClientExtensions() *ClientExtensions {
	return o.ClientExtensions
}

func (o MarketOrder) GetType() OrderType {
	return o.Type
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

func (o FixedPriceOrder) GetID() OrderID {
	return o.ID
}

func (o FixedPriceOrder) GetCreateTime() DateTime {
	return o.CreateTime
}

func (o FixedPriceOrder) GetState() OrderState {
	return o.State
}

func (o FixedPriceOrder) GetClientExtensions() *ClientExtensions {
	return o.ClientExtensions
}

func (o FixedPriceOrder) GetType() OrderType {
	return o.Type
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

func (o LimitOrder) GetID() OrderID {
	return o.ID
}

func (o LimitOrder) GetCreateTime() DateTime {
	return o.CreateTime
}

func (o LimitOrder) GetState() OrderState {
	return o.State
}

func (o LimitOrder) GetClientExtensions() *ClientExtensions {
	return o.ClientExtensions
}

func (o LimitOrder) GetType() OrderType {
	return o.Type
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
	PriceBound *PriceValue `json:"priceBound,omitempty"`
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

func (o StopOrder) GetID() OrderID {
	return o.ID
}

func (o StopOrder) GetCreateTime() DateTime {
	return o.CreateTime
}

func (o StopOrder) GetState() OrderState {
	return o.State
}

func (o StopOrder) GetClientExtensions() *ClientExtensions {
	return o.ClientExtensions
}

func (o StopOrder) GetType() OrderType {
	return o.Type
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
	PriceBound *PriceValue `json:"priceBound,omitempty"`
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

func (o MarketIfTouchedOrder) GetID() OrderID {
	return o.ID
}

func (o MarketIfTouchedOrder) GetCreateTime() DateTime {
	return o.CreateTime
}

func (o MarketIfTouchedOrder) GetState() OrderState {
	return o.State
}

func (o MarketIfTouchedOrder) GetClientExtensions() *ClientExtensions {
	return o.ClientExtensions
}

func (o MarketIfTouchedOrder) GetType() OrderType {
	return o.Type
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

func (o TakeProfitOrder) GetID() OrderID {
	return o.ID
}

func (o TakeProfitOrder) GetCreateTime() DateTime {
	return o.CreateTime
}

func (o TakeProfitOrder) GetState() OrderState {
	return o.State
}

func (o TakeProfitOrder) GetClientExtensions() *ClientExtensions {
	return o.ClientExtensions
}

func (o TakeProfitOrder) GetType() OrderType {
	return o.Type
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
	Distance *DecimalNumber `json:"distance,omitempty"`
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

func (o StopLossOrder) GetID() OrderID {
	return o.ID
}

func (o StopLossOrder) GetCreateTime() DateTime {
	return o.CreateTime
}

func (o StopLossOrder) GetState() OrderState {
	return o.State
}

func (o StopLossOrder) GetClientExtensions() *ClientExtensions {
	return o.ClientExtensions
}

func (o StopLossOrder) GetType() OrderType {
	return o.Type
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
	Distance *DecimalNumber `json:"distance"`
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

func (o GuaranteedStopLossOrder) GetID() OrderID {
	return o.ID
}

func (o GuaranteedStopLossOrder) GetCreateTime() DateTime {
	return o.CreateTime
}

func (o GuaranteedStopLossOrder) GetState() OrderState {
	return o.State
}

func (o GuaranteedStopLossOrder) GetClientExtensions() *ClientExtensions {
	return o.ClientExtensions
}

func (o GuaranteedStopLossOrder) GetType() OrderType {
	return o.Type
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

func (o TrailingStopLossOrder) GetID() OrderID {
	return o.ID
}

func (o TrailingStopLossOrder) GetCreateTime() DateTime {
	return o.CreateTime
}

func (o TrailingStopLossOrder) GetState() OrderState {
	return o.State
}

func (o TrailingStopLossOrder) GetClientExtensions() *ClientExtensions {
	return o.ClientExtensions
}

func (o TrailingStopLossOrder) GetType() OrderType {
	return o.Type
}

// Order Requests

// OrderRequest is the interface implemented by all order request types (e.g. MarketOrderRequest).
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
	PriceBound *PriceValue `json:"priceBound,omitempty"`
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
	return orderRequestWrapper(r)
}

// NewMarketOrderRequest creates a new MarketOrderRequest with default TimeInForce FOK and PositionFill DEFAULT.
func NewMarketOrderRequest(instrument InstrumentName, units DecimalNumber) *MarketOrderRequest {
	return &MarketOrderRequest{
		Type:         OrderTypeMarket,
		Instrument:   instrument,
		Units:        units,
		TimeInForce:  TimeInForceFOK,
		PositionFill: OrderPositionFillDefault,
	}
}

// SetIOC sets the TimeInForce to IOC (Immediate or Cancel).
func (r *MarketOrderRequest) SetIOC() *MarketOrderRequest {
	r.TimeInForce = TimeInForceIOC
	return r
}

// SetPriceBound sets the worst price that the client is willing to have the Market Order filled at.
func (r *MarketOrderRequest) SetPriceBound(priceBound PriceValue) *MarketOrderRequest {
	r.PriceBound = &priceBound
	return r
}

// SetPositionFill sets how Positions in the Account are modified when the Order is filled.
func (r *MarketOrderRequest) SetPositionFill(positionFill OrderPositionFill) *MarketOrderRequest {
	r.PositionFill = positionFill
	return r
}

// SetClientExtensions sets the client extensions for the Order.
func (r *MarketOrderRequest) SetClientExtensions(clientExtensions *ClientExtensions) *MarketOrderRequest {
	r.ClientExtensions = clientExtensions
	return r
}

// SetTakeProfitOnFill sets the Take Profit Order details for when the Order is filled.
func (r *MarketOrderRequest) SetTakeProfitOnFill(details *TakeProfitDetails) *MarketOrderRequest {
	r.TakeProfitOnFill = details
	return r
}

// SetStopLossOnFill sets the Stop Loss Order details for when the Order is filled.
func (r *MarketOrderRequest) SetStopLossOnFill(details *StopLossDetails) *MarketOrderRequest {
	r.StopLossOnFill = details
	return r
}

// SetGuaranteedStopLossOnFill sets the Guaranteed Stop Loss Order details for when the Order is filled.
func (r *MarketOrderRequest) SetGuaranteedStopLossOnFill(details *GuaranteedStopLossDetails) *MarketOrderRequest {
	r.GuaranteedStopLossOnFill = details
	return r
}

// SetTrailingStopLossOnFill sets the Trailing Stop Loss Order details for when the Order is filled.
func (r *MarketOrderRequest) SetTrailingStopLossOnFill(details *TrailingStopLossDetails) *MarketOrderRequest {
	r.TrailingStopLossOnFill = details
	return r
}

// SetTradeClientExtensions sets the client extensions for the Trade created when the Order is filled.
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
	return orderRequestWrapper(r)
}

// NewLimitOrderRequest creates a new LimitOrderRequest with default TimeInForce GTC, PositionFill DEFAULT, and TriggerCondition DEFAULT.
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

// SetGTD sets the TimeInForce to GTD (Good Till Date) with the specified expiry time.
func (r *LimitOrderRequest) SetGTD(date DateTime) *LimitOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = &date
	return r
}

// SetGFD sets the TimeInForce to GFD (Good For Day).
func (r *LimitOrderRequest) SetGFD() *LimitOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

// SetPositionFill sets how Positions in the Account are modified when the Order is filled.
func (r *LimitOrderRequest) SetPositionFill(positionFill OrderPositionFill) *LimitOrderRequest {
	r.PositionFill = positionFill
	return r
}

// SetTriggerCondition sets which price component is used to determine if the Order should be triggered.
func (r *LimitOrderRequest) SetTriggerCondition(triggerCondition OrderTriggerCondition) *LimitOrderRequest {
	r.TriggerCondition = triggerCondition
	return r
}

// SetClientExtensions sets the client extensions for the Order.
func (r *LimitOrderRequest) SetClientExtensions(clientExtensions *ClientExtensions) *LimitOrderRequest {
	r.ClientExtensions = clientExtensions
	return r
}

// SetTakeProfitOnFill sets the Take Profit Order details for when the Order is filled.
func (r *LimitOrderRequest) SetTakeProfitOnFill(details *TakeProfitDetails) *LimitOrderRequest {
	r.TakeProfitOnFill = details
	return r
}

// SetStopLossOnFill sets the Stop Loss Order details for when the Order is filled.
func (r *LimitOrderRequest) SetStopLossOnFill(details *StopLossDetails) *LimitOrderRequest {
	r.StopLossOnFill = details
	return r
}

// SetGuaranteedStopLossOnFill sets the Guaranteed Stop Loss Order details for when the Order is filled.
func (r *LimitOrderRequest) SetGuaranteedStopLossOnFill(details *GuaranteedStopLossDetails) *LimitOrderRequest {
	r.GuaranteedStopLossOnFill = details
	return r
}

// SetTrailingStopLossOnFill sets the Trailing Stop Loss Order details for when the Order is filled.
func (r *LimitOrderRequest) SetTrailingStopLossOnFill(details *TrailingStopLossDetails) *LimitOrderRequest {
	r.TrailingStopLossOnFill = details
	return r
}

// SetTradeClientExtensions sets the client extensions for the Trade created when the Order is filled.
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
	PriceBound *PriceValue `json:"priceBound,omitempty"`
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

func (r *StopOrderRequest) body() (*bytes.Buffer, error) {
	return orderRequestWrapper(r)
}

// NewStopOrderRequest creates a new StopOrderRequest with default TimeInForce GTC, PositionFill DEFAULT, and TriggerCondition DEFAULT.
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

// SetPriceBound sets the worst market price that may be used to fill this Stop Order.
func (r *StopOrderRequest) SetPriceBound(priceBound PriceValue) *StopOrderRequest {
	r.PriceBound = &priceBound
	return r
}

// SetGTD sets the TimeInForce to GTD (Good Till Date) with the specified expiry time.
func (r *StopOrderRequest) SetGTD(date DateTime) *StopOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = &date
	return r
}

// SetGFD sets the TimeInForce to GFD (Good For Day).
func (r *StopOrderRequest) SetGFD() *StopOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

// SetPositionFill sets how Positions in the Account are modified when the Order is filled.
func (r *StopOrderRequest) SetPositionFill(positionFill OrderPositionFill) *StopOrderRequest {
	r.PositionFill = positionFill
	return r
}

// SetTriggerCondition sets which price component is used to determine if the Order should be triggered.
func (r *StopOrderRequest) SetTriggerCondition(triggerCondition OrderTriggerCondition) *StopOrderRequest {
	r.TriggerCondition = triggerCondition
	return r
}

// SetClientExtensions sets the client extensions for the Order.
func (r *StopOrderRequest) SetClientExtensions(clientExtensions *ClientExtensions) *StopOrderRequest {
	r.ClientExtensions = clientExtensions
	return r
}

// SetTakeProfitOnFill sets the Take Profit Order details for when the Order is filled.
func (r *StopOrderRequest) SetTakeProfitOnFill(details *TakeProfitDetails) *StopOrderRequest {
	r.TakeProfitOnFill = details
	return r
}

// SetStopLossOnFill sets the Stop Loss Order details for when the Order is filled.
func (r *StopOrderRequest) SetStopLossOnFill(details *StopLossDetails) *StopOrderRequest {
	r.StopLossOnFill = details
	return r
}

// SetGuaranteedStopLossOnFill sets the Guaranteed Stop Loss Order details for when the Order is filled.
func (r *StopOrderRequest) SetGuaranteedStopLossOnFill(details *GuaranteedStopLossDetails) *StopOrderRequest {
	r.GuaranteedStopLossOnFill = details
	return r
}

// SetTrailingStopLossOnFill sets the Trailing Stop Loss Order details for when the Order is filled.
func (r *StopOrderRequest) SetTrailingStopLossOnFill(details *TrailingStopLossDetails) *StopOrderRequest {
	r.TrailingStopLossOnFill = details
	return r
}

// SetTradeClientExtensions sets the client extensions for the Trade created when the Order is filled.
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
	PriceBound *PriceValue `json:"priceBound,omitempty"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
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

func (r *MarketIfTouchedOrderRequest) body() (*bytes.Buffer, error) {
	return orderRequestWrapper(r)
}

// NewMarketIfTouchedOrderRequest creates a new MarketIfTouchedOrderRequest with default TimeInForce GTC, PositionFill DEFAULT, and TriggerCondition DEFAULT.
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

// SetPriceBound sets the worst market price that may be used to fill this Order.
func (r *MarketIfTouchedOrderRequest) SetPriceBound(priceBound PriceValue) *MarketIfTouchedOrderRequest {
	r.PriceBound = &priceBound
	return r
}

// SetGTD sets the TimeInForce to GTD (Good Till Date) with the specified expiry time.
func (r *MarketIfTouchedOrderRequest) SetGTD(date DateTime) *MarketIfTouchedOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = &date
	return r
}

// SetGFD sets the TimeInForce to GFD (Good For Day).
func (r *MarketIfTouchedOrderRequest) SetGFD() *MarketIfTouchedOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

// SetOpenOnly sets the PositionFill to OPEN_ONLY so the Order can only open new Positions.
func (r *MarketIfTouchedOrderRequest) SetOpenOnly() *MarketIfTouchedOrderRequest {
	r.PositionFill = OrderPositionFillOpenOnly
	return r
}

// SetReduceFirst sets the PositionFill to REDUCE_FIRST so existing Positions are reduced before opening new ones.
func (r *MarketIfTouchedOrderRequest) SetReduceFirst() *MarketIfTouchedOrderRequest {
	r.PositionFill = OrderPositionFillReduceFirst
	return r
}

// SetReduceOnly sets the PositionFill to REDUCE_ONLY so the Order can only reduce existing Positions.
func (r *MarketIfTouchedOrderRequest) SetReduceOnly() *MarketIfTouchedOrderRequest {
	r.PositionFill = OrderPositionFillReduceOnly
	return r
}

// SetTriggerCondition sets which price component is used to determine if the Order should be triggered.
func (r *MarketIfTouchedOrderRequest) SetTriggerCondition(triggerCondition OrderTriggerCondition) *MarketIfTouchedOrderRequest {
	r.TriggerCondition = triggerCondition
	return r
}

// SetClientExtensions sets the client extensions for the Order.
func (r *MarketIfTouchedOrderRequest) SetClientExtensions(clientExtensions *ClientExtensions) *MarketIfTouchedOrderRequest {
	r.ClientExtensions = clientExtensions
	return r
}

// SetTakeProfitOnFill sets the Take Profit Order details for when the Order is filled.
func (r *MarketIfTouchedOrderRequest) SetTakeProfitOnFill(details *TakeProfitDetails) *MarketIfTouchedOrderRequest {
	r.TakeProfitOnFill = details
	return r
}

// SetStopLossOnFill sets the Stop Loss Order details for when the Order is filled.
func (r *MarketIfTouchedOrderRequest) SetStopLossOnFill(details *StopLossDetails) *MarketIfTouchedOrderRequest {
	r.StopLossOnFill = details
	return r
}

// SetGuaranteedStopLossOnFill sets the Guaranteed Stop Loss Order details for when the Order is filled.
func (r *MarketIfTouchedOrderRequest) SetGuaranteedStopLossOnFill(details *GuaranteedStopLossDetails) *MarketIfTouchedOrderRequest {
	r.GuaranteedStopLossOnFill = details
	return r
}

// SetTrailingStopLossOnFill sets the Trailing Stop Loss Order details for when the Order is filled.
func (r *MarketIfTouchedOrderRequest) SetTrailingStopLossOnFill(details *TrailingStopLossDetails) *MarketIfTouchedOrderRequest {
	r.TrailingStopLossOnFill = details
	return r
}

// SetTradeClientExtensions sets the client extensions for the Trade created when the Order is filled.
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
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// Price is the price threshold specified for the Take Profit Order. The associated Trade will be
	// closed by a market price that is equal to or better than this threshold.
	Price PriceValue `json:"price"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
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
	return orderRequestWrapper(r)
}

// NewTakeProfitOrderRequest creates a new TakeProfitOrderRequest with default TimeInForce GTC and TriggerCondition DEFAULT.
func NewTakeProfitOrderRequest(tradeID TradeID, price PriceValue) *TakeProfitOrderRequest {
	return &TakeProfitOrderRequest{
		Type:             OrderTypeTakeProfit,
		TradeID:          tradeID,
		Price:            price,
		TimeInForce:      TimeInForceGTC,
		TriggerCondition: OrderTriggerConditionDefault,
	}
}

// SetClientTradeID sets the client Trade ID of the Trade to be closed.
func (r *TakeProfitOrderRequest) SetClientTradeID(clientID ClientID) *TakeProfitOrderRequest {
	r.ClientTradeID = &clientID
	return r
}

// SetGTD sets the TimeInForce to GTD (Good Till Date) with the specified expiry time.
func (r *TakeProfitOrderRequest) SetGTD(date DateTime) *TakeProfitOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = &date
	return r
}

// SetGFD sets the TimeInForce to GFD (Good For Day).
func (r *TakeProfitOrderRequest) SetGFD() *TakeProfitOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

// SetTriggerCondition sets which price component is used to determine if the Order should be triggered.
func (r *TakeProfitOrderRequest) SetTriggerCondition(triggerCondition OrderTriggerCondition) *TakeProfitOrderRequest {
	r.TriggerCondition = triggerCondition
	return r
}

// SetClientExtensions sets the client extensions for the Order.
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
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// Price is the price threshold specified for the Stop Loss Order. The associated Trade will be
	// closed by a market price that is equal to or worse than this threshold. Either price or distance
	// may be specified, but not both.
	Price *PriceValue `json:"price,omitempty"`
	// Distance specifies the distance (in price units) from the Account's current price to use as
	// the Stop Loss Order price. If the Trade is long the Order's price will be the bid price minus
	// the distance. If the Trade is short the Order's price will be the ask price plus the distance.
	// Either price or distance may be specified, but not both.
	Distance *DecimalNumber `json:"distance,omitempty"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
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
	if r.Price == nil && r.Distance == nil {
		return nil, errors.New("price or distance must be set")
	}
	if r.Price != nil && r.Distance != nil {
		return nil, errors.New("price and distance cannot be set at the same time")
	}
	return orderRequestWrapper(r)
}

// NewStopLossOrderRequest creates a new StopLossOrderRequest with default TimeInForce GTC and TriggerCondition DEFAULT.
func NewStopLossOrderRequest(tradeID TradeID) *StopLossOrderRequest {
	return &StopLossOrderRequest{
		Type:             OrderTypeStopLoss,
		TradeID:          tradeID,
		Price:            nil,
		Distance:         nil,
		TimeInForce:      TimeInForceGTC,
		GtdTime:          nil,
		TriggerCondition: OrderTriggerConditionDefault,
		ClientExtensions: nil,
	}
}

// SetClientTradeID sets the client Trade ID of the Trade to be closed.
func (r *StopLossOrderRequest) SetClientTradeID(clientID ClientID) *StopLossOrderRequest {
	r.ClientTradeID = &clientID
	return r
}

// SetPrice sets the price for StopLossOrder.
func (r *StopLossOrderRequest) SetPrice(price PriceValue) *StopLossOrderRequest {
	r.Price = &price
	return r
}

// SetDistance sets the distance from the Account's current price to use as the StopLossOrder price.
func (r *StopLossOrderRequest) SetDistance(distance DecimalNumber) *StopLossOrderRequest {
	r.Distance = &distance
	return r
}

// SetGTD sets the TimeInForce to GTD (Good Till Date) with the specified expiry time.
func (r *StopLossOrderRequest) SetGTD(date DateTime) *StopLossOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = &date
	return r
}

// SetGFD sets the TimeInForce to GFD (Good For Day).
func (r *StopLossOrderRequest) SetGFD() *StopLossOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

// SetTriggerCondition sets which price component is used to determine if the Order should be triggered.
func (r *StopLossOrderRequest) SetTriggerCondition(triggerCondition OrderTriggerCondition) *StopLossOrderRequest {
	r.TriggerCondition = triggerCondition
	return r
}

// SetClientExtensions sets the client extensions for the Order.
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
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// Price is the price threshold specified for the Guaranteed Stop Loss Order. The associated Trade
	// will be closed at this price. Either price or distance may be specified, but not both.
	Price *PriceValue `json:"price,omitempty"`
	// Distance specifies the distance (in price units) from the Account's current price to use as
	// the Guaranteed Stop Loss Order price. If the Trade is long the Order's price will be the bid
	// price minus the distance. If the Trade is short the Order's price will be the ask price plus
	// the distance. Either price or distance may be specified, but not both.
	Distance *DecimalNumber `json:"distance,omitempty"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
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
	return orderRequestWrapper(r)
}

// NewGuaranteedStopLossOrderRequest creates a new GuaranteedStopLossOrderRequest with default TimeInForce GTC and TriggerCondition DEFAULT.
func NewGuaranteedStopLossOrderRequest(tradeID TradeID, price PriceValue) *GuaranteedStopLossOrderRequest {
	return &GuaranteedStopLossOrderRequest{
		Type:             OrderTypeGuaranteedStopLoss,
		TradeID:          tradeID,
		TimeInForce:      TimeInForceGTC,
		TriggerCondition: OrderTriggerConditionDefault,
	}
}

func (r *GuaranteedStopLossOrderRequest) SetPrice(price PriceValue) *GuaranteedStopLossOrderRequest {
	r.Price = &price
	return r
}

func (r *GuaranteedStopLossOrderRequest) SetDistance(distance DecimalNumber) *GuaranteedStopLossOrderRequest {
	r.Distance = &distance
	return r
}

// SetClientTradeID sets the client Trade ID of the Trade to be closed.
func (r *GuaranteedStopLossOrderRequest) SetClientTradeID(clientID ClientID) *GuaranteedStopLossOrderRequest {
	r.ClientTradeID = &clientID
	return r
}

// SetGTD sets the TimeInForce to GTD (Good Till Date) with the specified expiry time.
func (r *GuaranteedStopLossOrderRequest) SetGTD(date DateTime) *GuaranteedStopLossOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = &date
	return r
}

// SetGFD sets the TimeInForce to GFD (Good For Day).
func (r *GuaranteedStopLossOrderRequest) SetGFD() *GuaranteedStopLossOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

// SetTriggerCondition sets which price component is used to determine if the Order should be triggered.
func (r *GuaranteedStopLossOrderRequest) SetTriggerCondition(triggerCondition OrderTriggerCondition) *GuaranteedStopLossOrderRequest {
	r.TriggerCondition = triggerCondition
	return r
}

// SetClientExtensions sets the client extensions for the Order.
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
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// Distance is the price distance (in price units) specified for the Trailing Stop Loss Order.
	Distance DecimalNumber `json:"distance"`
	// TimeInForce specifies how long the Order should remain pending before being automatically
	// cancelled by the execution system. Valid options are GTC, GFD, and GTD. Default is GTC.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
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
	return orderRequestWrapper(r)
}

// NewTrailingStopLossOrderRequest creates a new TrailingStopLossOrderRequest with default TimeInForce GTC and TriggerCondition DEFAULT.
func NewTrailingStopLossOrderRequest(tradeID TradeID, distance DecimalNumber) *TrailingStopLossOrderRequest {
	return &TrailingStopLossOrderRequest{
		Type:             OrderTypeTrailingStopLoss,
		TradeID:          tradeID,
		Distance:         distance,
		TimeInForce:      TimeInForceGTC,
		TriggerCondition: OrderTriggerConditionDefault,
	}
}

// SetClientTradeID sets the client Trade ID of the Trade to be closed.
func (r *TrailingStopLossOrderRequest) SetClientTradeID(clientID ClientID) *TrailingStopLossOrderRequest {
	r.ClientTradeID = &clientID
	return r
}

// SetGTD sets the TimeInForce to GTD (Good Till Date) with the specified expiry time.
func (r *TrailingStopLossOrderRequest) SetGTD(date DateTime) *TrailingStopLossOrderRequest {
	r.TimeInForce = TimeInForceGTD
	r.GtdTime = &date
	return r
}

// SetGFD sets the TimeInForce to GFD (Good For Day).
func (r *TrailingStopLossOrderRequest) SetGFD() *TrailingStopLossOrderRequest {
	r.TimeInForce = TimeInForceGFD
	return r
}

// SetTriggerCondition sets which price component is used to determine if the Order should be triggered.
func (r *TrailingStopLossOrderRequest) SetTriggerCondition(triggerCondition OrderTriggerCondition) *TrailingStopLossOrderRequest {
	r.TriggerCondition = triggerCondition
	return r
}

// SetClientExtensions sets the client extensions for the Order.
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
	// OrderTypeMarket represents a Market Order.
	OrderTypeMarket OrderType = "MARKET"
	// OrderTypeLimit represents a Limit Order.
	OrderTypeLimit OrderType = "LIMIT"
	// OrderTypeStop represents a Stop Order.
	OrderTypeStop OrderType = "STOP"
	// OrderTypeMarketIfTouched represents a Market If Touched Order.
	OrderTypeMarketIfTouched OrderType = "MARKET_IF_TOUCHED"
	// OrderTypeFixedPrice represents a Fixed Price Order.
	OrderTypeFixedPrice OrderType = "FIXED_PRICE"
	// OrderTypeTakeProfit represents a Take Profit Order.
	OrderTypeTakeProfit OrderType = "TAKE_PROFIT"
	// OrderTypeStopLoss represents a Stop Loss Order.
	OrderTypeStopLoss OrderType = "STOP_LOSS"
	// OrderTypeGuaranteedStopLoss represents a Guaranteed Stop Loss Order.
	OrderTypeGuaranteedStopLoss OrderType = "GUARANTEED_STOP_LOSS"
	// OrderTypeTrailingStopLoss represents a Trailing Stop Loss Order.
	OrderTypeTrailingStopLoss OrderType = "TRAILING_STOP_LOSS"
)

// CancellableOrderType represents the type of Orders that can be cancelled.
// Market and FixedPrice orders cannot be cancelled as they are filled immediately.
type CancellableOrderType string

const (
	// CancellableOrderTypeLimit represents a cancellable Limit Order.
	CancellableOrderTypeLimit CancellableOrderType = "LIMIT"
	// CancellableOrderTypeStop represents a cancellable Stop Order.
	CancellableOrderTypeStop CancellableOrderType = "STOP"
	// CancellableOrderTypeMarketIfTouched represents a cancellable Market If Touched Order.
	CancellableOrderTypeMarketIfTouched CancellableOrderType = "MARKET_IF_TOUCHED"
	// CancellableOrderTypeTakeProfit represents a cancellable Take Profit Order.
	CancellableOrderTypeTakeProfit CancellableOrderType = "TAKE_PROFIT"
	// CancellableOrderTypeStopLoss represents a cancellable Stop Loss Order.
	CancellableOrderTypeStopLoss CancellableOrderType = "STOP_LOSS"
	// CancellableOrderTypeGuaranteedStopLoss represents a cancellable Guaranteed Stop Loss Order.
	CancellableOrderTypeGuaranteedStopLoss CancellableOrderType = "GUARANTEED_STOP_LOSS"
	// CancellableOrderTypeTrailingStopLoss represents a cancellable Trailing Stop Loss Order.
	CancellableOrderTypeTrailingStopLoss CancellableOrderType = "TRAILING_STOP_LOSS"
)

// OrderState represents the current state of an Order.
type OrderState string

const (
	// OrderStatePending means the Order is currently pending execution.
	OrderStatePending OrderState = "PENDING"
	// OrderStateFilled means the Order has been filled.
	OrderStateFilled OrderState = "FILLED"
	// OrderStateTriggered means the Order has been triggered.
	OrderStateTriggered OrderState = "TRIGGERED"
	// OrderStateCancelled means the Order has been cancelled.
	OrderStateCancelled OrderState = "CANCELLED"
)

// OrderIdentifier contains both the OANDA-assigned and client-assigned identifiers for an Order.
type OrderIdentifier struct {
	OrderID       OrderID  `json:"orderID"`
	ClientOrderID ClientID `json:"clientOrderID"`
}

// OrderSpecifier is either an Order's OANDA-assigned OrderID or the client-provided ClientID prefixed with "@".
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

// -------------------------------------------------------------
// Endpoints https://developer.oanda.com/rest-live-v20/order-ep/
// -------------------------------------------------------------

type orderService struct {
	client *Client
}

func newOrderService(client *Client) *orderService {
	return &orderService{client}
}

// OrderCreateResponse is the successful response returned by [orderService.Create].
type OrderCreateResponse struct {
	OrderCreateTransaction        Transaction             `json:"orderCreateTransaction"`
	OrderFillTransaction          *OrderFillTransaction   `json:"orderFillTransaction,omitempty"`
	OrderCancelTransaction        *OrderCancelTransaction `json:"orderCancelTransaction,omitempty"`
	OrderReissueTransaction       Transaction             `json:"orderReissueTransaction,omitempty"`
	OrderReissueRejectTransaction Transaction             `json:"orderReissueRejectTransaction,omitempty"`
	RelatedTransactionIDs         []TransactionID         `json:"relatedTransactionIDs"`
	LastTransactionID             TransactionID           `json:"lastTransactionID"`
}

// OrderErrorResponse is the error response returned by order endpoints when a request is rejected.
type OrderErrorResponse struct {
	OrderRejectTransaction Transaction     `json:"orderRejectTransaction"`
	RelatedTransactionIDs  []TransactionID `json:"relatedTransactionIDs"`
	LastTransactionID      TransactionID   `json:"lastTransactionID"`
	ErrorCode              string          `json:"errorCode"`
	ErrorMessage           string          `json:"errorMessage"`
}

// Error implements the error interface.
func (e OrderErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode, e.ErrorMessage)
}

func unmarshalOrderErrorResponse(resp *http.Response) (*OrderErrorResponse, error) {
	var r OrderErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OrderErrorResponse: %w", err)
	}
	return &r, nil
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

// Create submits a new Order for the Account configured via WithAccountID.
//
// This corresponds to the OANDA API endpoint: POST /v3/accounts/{accountID}/orders
//
// Reference: https://developer.oanda.com/rest-live-v20/order-ep/#collapse_endpoint_1
func (s *orderService) Create(ctx context.Context, req OrderRequest) (*OrderCreateResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/orders", s.client.accountID)
	body, err := req.body()
	if err != nil {
		return nil, err
	}
	httpResp, err := s.client.sendPostRequest(ctx, path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %w", err)
	}
	defer closeBody(httpResp)
	switch httpResp.StatusCode {
	case http.StatusCreated:
		var resp OrderCreateResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &resp, nil
	case http.StatusBadRequest:
		r, err := unmarshalOrderErrorResponse(httpResp)
		if err != nil {
			return nil, err
		}
		return nil, BadRequest{HTTPError{httpResp.StatusCode, "bad request", r}}
	case http.StatusNotFound:
		r, err := unmarshalOrderErrorResponse(httpResp)
		if err != nil {
			return nil, err
		}
		return nil, NotFoundError{HTTPError{httpResp.StatusCode, "not found", r}}
	default:
		return nil, decodeErrorResponse(httpResp)
	}
}

// OrderListRequest contains the parameters for retrieving a list of Orders for an Account.
// Use NewOrderListRequest to create a new request and the builder methods to configure options.
type OrderListRequest struct {
	IDs        []OrderID
	State      *OrderState
	Instrument *InstrumentName
	Count      *int
	BeforeID   *OrderID
}

// NewOrderListRequest creates a new OrderListRequest for the specified account.
// Use the builder methods (AddIDs, SetState, SetInstrument, SetCount, SetBeforeID)
// to configure optional filtering parameters.
func NewOrderListRequest() *OrderListRequest {
	return &OrderListRequest{
		IDs: make([]OrderID, 0),
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

func (r *OrderListResponse) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Orders            []json.RawMessage `json:"orders"`
		LastTransactionID TransactionID     `json:"lastTransactionID"`
	}
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	orders, err := unmarshalOrders(raw.Orders)
	if err != nil {
		return err
	}
	r.Orders = orders
	r.LastTransactionID = raw.LastTransactionID
	return nil
}

// List retrieves a list of Orders for the Account configured via WithAccountID.
// Use [NewOrderListRequest] to create and configure filter parameters.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/orders
//
// Reference: https://developer.oanda.com/rest-live-v20/order-ep/#collapse_endpoint_2
func (s *orderService) List(ctx context.Context, req *OrderListRequest) (*OrderListResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/orders", s.client.accountID)
	v, err := req.values()
	if err != nil {
		return nil, err
	}
	resp, err := s.client.sendGetRequest(ctx, path, v)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var orderListResp OrderListResponse
	if err := decodeResponse(resp, &orderListResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &orderListResp, nil
}

// ListPending retrieves all pending Orders for the Account configured via WithAccountID.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/pendingOrders
//
// Reference: https://developer.oanda.com/rest-live-v20/order-ep/#collapse_endpoint_3
func (s *orderService) ListPending(ctx context.Context) (*OrderListResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/pendingOrders", s.client.accountID)
	resp, err := s.client.sendGetRequest(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var orderListResp OrderListResponse
	if err := decodeResponse(resp, &orderListResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &orderListResp, nil
}

// OrderDetailsResponse is the response returned by [orderService.Details].
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

// Details retrieves the details of a single Order for the Account configured via WithAccountID.
// The specifier can be the Order's OANDA-assigned OrderID or a client ID prefixed with "@".
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/orders/{orderSpecifier}
//
// Reference: https://developer.oanda.com/rest-live-v20/order-ep/#collapse_endpoint_4
func (s *orderService) Details(ctx context.Context, specifier OrderSpecifier) (*OrderDetailsResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/orders/%v", s.client.accountID, specifier)
	return doGet[OrderDetailsResponse](s.client, ctx, path, nil)
}

// OrderReplaceResponse is the successful response returned by [Client.OrderReplace].
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

// OrderReplace cancels and replaces an existing Order with a new one.
//
// This corresponds to the OANDA API endpoint: PUT /v3/accounts/{accountID}/orders/{orderSpecifier}
//
// Reference: https://developer.oanda.com/rest-live-v20/order-ep/#collapse_endpoint_5
func (c *Client) OrderReplace(ctx context.Context, specifier OrderSpecifier, req OrderRequest) (*OrderReplaceResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/orders/%v", c.accountID, specifier)
	body, err := req.body()
	if err != nil {
		return nil, err
	}
	httpResp, err := c.sendPutRequest(ctx, path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to send PUT request: %w", err)
	}
	defer closeBody(httpResp)
	switch httpResp.StatusCode {
	case http.StatusCreated:
		var resp OrderReplaceResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &resp, nil
	case http.StatusBadRequest:
		r, err := unmarshalOrderErrorResponse(httpResp)
		if err != nil {
			return nil, err
		}
		return nil, BadRequest{HTTPError{httpResp.StatusCode, "bad request", r}}
	case http.StatusNotFound:
		r, err := unmarshalOrderErrorResponse(httpResp)
		if err != nil {
			return nil, err
		}
		return nil, NotFoundError{HTTPError{httpResp.StatusCode, "not found", r}}
	default:
		return nil, decodeErrorResponse(httpResp)
	}
}

// OrderCancelResponse is the successful response returned by [orderService.Cancel].
type OrderCancelResponse struct {
	OrderCancelTransaction OrderCancelTransaction `json:"orderCancelTransaction"`
	RelatedTransactionIDs  []TransactionID        `json:"relatedTransactionIDs"`
	LastTransactionID      TransactionID          `json:"lastTransactionID"`
}

// Cancel cancels a pending Order for the Account configured via WithAccountID.
//
// This corresponds to the OANDA API endpoint: PUT /v3/accounts/{accountID}/orders/{orderSpecifier}/cancel
//
// Reference: https://developer.oanda.com/rest-live-v20/order-ep/#collapse_endpoint_6
func (s *orderService) Cancel(ctx context.Context, specifier OrderSpecifier) (*OrderCancelResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/orders/%v/cancel", s.client.accountID, specifier)
	httpResp, err := s.client.sendPutRequest(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send PUT request: %w", err)
	}
	defer closeBody(httpResp)
	switch httpResp.StatusCode {
	case http.StatusOK:
		var resp OrderCancelResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &resp, nil
	case http.StatusNotFound:
		resp, err := unmarshalOrderErrorResponse(httpResp)
		if err != nil {
			return nil, err
		}
		return nil, NotFoundError{HTTPError{httpResp.StatusCode, "not found", resp}}
	default:
		return nil, decodeErrorResponse(httpResp)
	}
}

// OrderUpdateClientExtensionsRequest is the request body for updating client extensions on an Order.
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

// OrderUpdateClientExtensionsResponse is the successful response returned by [orderService.UpdateClientExtensions].
type OrderUpdateClientExtensionsResponse struct {
	OrderClientExtensionsModifyTransaction OrderClientExtensionsModifyTransaction `json:"orderClientExtensionsModifyTransaction"`
	LastTransactionID                      TransactionID                          `json:"lastTransactionID"`
	RelatedTransactionIDs                  []TransactionID                        `json:"relatedTransactionIDs"`
}

// UpdateClientExtensions updates the client extensions for an Order.
//
// This corresponds to the OANDA API endpoint: PUT /v3/accounts/{accountID}/orders/{orderSpecifier}/clientExtensions
//
// Reference: https://developer.oanda.com/rest-live-v20/order-ep/#collapse_endpoint_7
func (s *orderService) UpdateClientExtensions(
	ctx context.Context,
	specifier OrderSpecifier,
	req OrderUpdateClientExtensionsRequest,
) (*OrderUpdateClientExtensionsResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/orders/%v/clientExtensions", s.client.accountID, specifier)
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
		var resp OrderUpdateClientExtensionsResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &resp, nil
	case http.StatusBadRequest:
		r, err := unmarshalOrderErrorResponse(httpResp)
		if err != nil {
			return nil, err
		}
		return nil, BadRequest{HTTPError{httpResp.StatusCode, "bad request", r}}
	case http.StatusNotFound:
		r, err := unmarshalOrderErrorResponse(httpResp)
		if err != nil {
			return nil, err
		}
		return nil, NotFoundError{HTTPError{httpResp.StatusCode, "not found", r}}
	default:
		return nil, decodeErrorResponse(httpResp)
	}
}
