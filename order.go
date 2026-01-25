package oanda

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
	PriceBound PriceValue `json:"priceBound"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	// Default is DEFAULT.
	PositionFill OrderPositionFill `json:"positionFill"`
	// ClientExtensions are the client extensions to add to the Order. Do not set, modify, or delete
	// clientExtensions if your account is associated with MT4.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
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
	ClientExtensions ClientExtensions `json:"clientExtensions"`
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
	PriceBound PriceValue `json:"priceBound"`
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
	ClientExtensions ClientExtensions `json:"clientExtensions"`
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
	ClientExtensions ClientExtensions `json:"clientExtensions"`
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
	ClientExtensions ClientExtensions `json:"clientExtensions"`
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

// Order-related Definitions

// OrderID is the unique identifier for an Order within an Account.
type OrderID uint

// OrderType represents the type of an Order.
type OrderType string

const (
	OrderTypeMarket                OrderType = "MARKET"
	OrderTypeLimit                 OrderType = "LIMIT"
	OrderTypeStop                  OrderType = "STOP"
	OrderTypeMarketIfTouched       OrderType = "MARKET_IF_TOUCHED"
	OrderTypeFixedPrice            OrderType = "FIXED_PRICE"
	OrderTypeTakeProfit            OrderType = "TAKE_PROFIT"
	OrderTypeStopLoss              OrderType = "STOP_LOSS"
	OrderTypeGuaranteedStopLoss    OrderType = "GUARANTEED_STOP_LOSS"
	OrderTypeTrailingStopLoss      OrderType = "TRAILING_STOP_LOSS"
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

// TakeProfitDetails specifies the details of a Take Profit Order to be created on behalf of a
// client. This may happen when an Order is filled that opens a Trade requiring a Take Profit.
type TakeProfitDetails struct {
	// Price is the price threshold specified for the Take Profit Order. The associated Trade will
	// be closed by a market price that is equal to or better than this threshold.
	Price PriceValue `json:"price"`
	// TimeInForce specifies how long the Take Profit Order should remain pending before being
	// automatically cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Take Profit Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// ClientExtensions are the client extensions to add to the Take Profit Order when created.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
}

// StopLossDetails specifies the details of a Stop Loss Order to be created on behalf of a client.
// This may happen when an Order is filled that opens a Trade requiring a Stop Loss.
type StopLossDetails struct {
	// Price is the price threshold specified for the Stop Loss Order. The associated Trade will be
	// closed by a market price that is equal to or worse than this threshold.
	Price PriceValue `json:"price"`
	// Distance specifies the distance (in price units) from the Trade's open price to use as the
	// Stop Loss Order price. If the Trade is long, the Stop Loss price will be the open price minus
	// the distance. If the Trade is short, the Stop Loss price will be the open price plus the distance.
	Distance DecimalNumber `json:"distance"`
	// TimeInForce specifies how long the Stop Loss Order should remain pending before being
	// automatically cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Stop Loss Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// ClientExtensions are the client extensions to add to the Stop Loss Order when created.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
	// Guaranteed is deprecated. Flag indicating that the Stop Loss Order is guaranteed. The default
	// value depends on the GuaranteedStopLossOrderMode of the account.
	Guaranteed bool `json:"guaranteed"`
}

// GuaranteedStopLossDetails specifies the details of a Guaranteed Stop Loss Order to be created on
// behalf of a client. This may happen when an Order is filled that opens a Trade requiring a
// Guaranteed Stop Loss.
type GuaranteedStopLossDetails struct {
	// Price is the price threshold specified for the Guaranteed Stop Loss Order. The associated Trade
	// will be closed at this price.
	Price PriceValue `json:"price"`
	// Distance specifies the distance (in price units) from the Trade's open price to use as the
	// Guaranteed Stop Loss Order price. If the Trade is long, the order price will be the open price
	// minus the distance. If the Trade is short, the order price will be the open price plus the distance.
	Distance DecimalNumber `json:"distance"`
	// TimeInForce specifies how long the Guaranteed Stop Loss Order should remain pending before
	// being automatically cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Guaranteed Stop Loss Order will be cancelled if its
	// timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// ClientExtensions are the client extensions to add to the Guaranteed Stop Loss Order when created.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
}

// TrailingStopLossDetails specifies the details of a Trailing Stop Loss Order to be created on
// behalf of a client. This may happen when an Order is filled that opens a Trade requiring a
// Trailing Stop Loss.
type TrailingStopLossDetails struct {
	// Distance is the price distance (in price units) specified for the Trailing Stop Loss Order.
	Distance DecimalNumber `json:"distance"`
	// TimeInForce specifies how long the Trailing Stop Loss Order should remain pending before
	// being automatically cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Trailing Stop Loss Order will be cancelled if its
	// timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// ClientExtensions are the client extensions to add to the Trailing Stop Loss Order when created.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
}

// MarketOrderTradeClose specifies the extensions to a Market Order that has been created specifically
// to close a Trade.
type MarketOrderTradeClose struct {
	// TradeID is the ID of the Trade requested to be closed.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade requested to be closed.
	ClientTradeID ClientID `json:"clientTradeID"`
	// Units indicates the number of units of the Trade to close. If not specified, all units of the
	// Trade will be closed.
	Units DecimalNumber `json:"units"`
}

// MarketOrderPositionCloseout specifies the extensions to a Market Order when it has been created to
// closeout a specific Position.
type MarketOrderPositionCloseout struct {
	// Instrument is the instrument of the Position being closed out.
	Instrument InstrumentName `json:"instrument"`
	// Units indicates the number of units of the Position being closed. If not specified, all units
	// of the Position will be closed.
	Units DecimalNumber `json:"units"`
}

// MarketOrderMarginCloseout details the reason that a Market Order was created as part of a Margin Closeout.
type MarketOrderMarginCloseout struct {
	// Reason is the reason the Market Order was created to perform a margin closeout.
	Reason string `json:"reason"`
}

// MarketOrderDelayedTradeClose details the reason that a Market Order was created for a delayed Trade close.
type MarketOrderDelayedTradeClose struct {
	// TradeID is the ID of the Trade being closed.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade being closed.
	ClientTradeID ClientID `json:"clientTradeID"`
	// SourceTransactionID is the Transaction ID of the DelayedTradeClosure transaction to which this
	// Delayed Trade Close belongs to.
	SourceTransactionID TransactionID `json:"sourceTransactionID"`
}
