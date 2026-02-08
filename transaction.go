package oanda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ---------------------------------------------------------------------
// Definitions https://developer.oanda.com/rest-live-v20/transaction-df/
// ---------------------------------------------------------------------

// Transactions

// Transaction represents the base specification for all Transactions.
type Transaction struct {
	// ID is the Transaction's Identifier.
	ID TransactionID `json:"id"`
	// Time is the date/time when the Transaction was created.
	Time DateTime `json:"time"`
	// UserID is the ID of the user that initiated the creation of the Transaction.
	UserID int `json:"userID"`
	// AccountID is the ID of the Account the Transaction was created for.
	AccountID AccountID `json:"accountID"`
	// BatchID is the ID of the "batch" that the Transaction belongs to. Transactions in the same
	// batch are applied to the Account simultaneously.
	BatchID TransactionID `json:"batchID"`
	// RequestID is the Request ID of the request which generated the transaction.
	RequestID RequestID `json:"requestID"`
	// Type is the Type of the Transaction.
	Type TransactionType `json:"type"`
}

func (t Transaction) GetType() string {
	return string(t.Type)
}

func (t Transaction) GetID() TransactionID {
	return t.ID
}

func (t Transaction) GetTime() DateTime {
	return t.Time
}

// CreateTransaction represents a Transaction that creates an Account.
type CreateTransaction struct {
	Transaction
	// DivisionID is the ID of the Division that the Account is in.
	DivisionID int `json:"divisionID"`
	// SiteID is the ID of the Site that the Account was created at.
	SiteID int `json:"siteID"`
	// AccountUserID is the ID of the user that the Account was created for.
	AccountUserID int `json:"accountUserID"`
	// AccountNumber is the number of the Account within the site/division/user.
	AccountNumber int `json:"accountNumber"`
	// HomeCurrency is the home currency of the Account.
	HomeCurrency Currency `json:"homeCurrency"`
}

// CloseTransaction represents a Transaction that closes an Account.
type CloseTransaction struct {
	Transaction
}

// ReopenTransaction represents a Transaction that reopens a closed Account.
type ReopenTransaction struct {
	Transaction
}

// ClientConfigureTransaction represents a Transaction that modifies an Account's client-provided configuration.
type ClientConfigureTransaction struct {
	Transaction
	// Alias is the client-provided alias for the Account.
	Alias string `json:"alias"`
	// MarginRate is the margin rate override for the Account.
	MarginRate DecimalNumber `json:"marginRate"`
}

// ClientConfigureRejectTransaction represents a Transaction that rejects the configuration of an Account's
// client-provided settings.
type ClientConfigureRejectTransaction struct {
	Transaction
	// Alias is the client-provided alias for the Account.
	Alias string `json:"alias"`
	// MarginRate is the margin rate override for the Account.
	MarginRate DecimalNumber `json:"marginRate"`
	// RejectReason is the reason that the Reject Transaction was created.
	RejectReason TransactionRejectReason `json:"rejectReason"`
}

// TransferFundsTransaction represents a Transaction that transfers funds between accounts.
type TransferFundsTransaction struct {
	Transaction
	// Amount is the amount to deposit/withdraw from the Account in the Account's home currency.
	// A positive value indicates a deposit, a negative value indicates a withdrawal.
	Amount AccountUnits `json:"amount"`
	// FundingReason is the reason that an Account is being funded.
	FundingReason FundingReason `json:"fundingReason"`
	// Comment is an optional comment that may be attached to a fund transfer for audit purposes.
	Comment string `json:"comment"`
	// AccountBalance is the Account's balance after funds are transferred.
	AccountBalance AccountUnits `json:"accountBalance"`
}

// TransferFundsRejectTransaction represents a Transaction that rejects the transfer of funds.
type TransferFundsRejectTransaction struct {
	Transaction
	// Amount is the amount to deposit/withdraw from the Account in the Account's home currency.
	Amount AccountUnits `json:"amount"`
	// FundingReason is the reason that an Account is being funded.
	FundingReason FundingReason `json:"fundingReason"`
	// Comment is an optional comment that may be attached to a fund transfer for audit purposes.
	Comment string `json:"comment"`
	// RejectReason is the reason that the Reject Transaction was created.
	RejectReason TransactionRejectReason `json:"rejectReason"`
}

// MarketOrderTransaction represents a Transaction that creates a Market Order.
type MarketOrderTransaction struct {
	Transaction
	// Instrument is the Market Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Market Order.
	Units DecimalNumber `json:"units"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// PriceBound is the worst price that the client is willing to have the Market Order filled at.
	PriceBound PriceValue `json:"priceBound"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TradeClose details the Trade requested to be closed.
	TradeClose MarketOrderTradeClose `json:"tradeClose"`
	// LongPositionCloseout details the long Position to closeout.
	LongPositionCloseout MarketOrderPositionCloseout `json:"longPositionCloseout"`
	// ShortPositionCloseout details the short Position to closeout.
	ShortPositionCloseout MarketOrderPositionCloseout `json:"shortPositionCloseout"`
	// MarginCloseout details the Margin Closeout that this Market Order was created for.
	MarginCloseout MarketOrderMarginCloseout `json:"marginCloseout"`
	// DelayedTradeClose details the delayed Trade close that this Market Order was created for.
	DelayedTradeClose MarketOrderDelayedTradeClose `json:"delayedTradeClose"`
	// Reason is the reason that the Market Order was created.
	Reason MarketOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
	// TakeProfitOnFill specifies the Take Profit Order details.
	TakeProfitOnFill TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the Stop Loss Order details.
	StopLossOnFill StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the Guaranteed Stop Loss Order details.
	GuaranteedStopLossOnFill GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the Trailing Stop Loss Order details.
	TrailingStopLossOnFill TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions for the Trade.
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions"`
}

// MarketOrderRejectTransaction represents a Transaction that rejects the creation of a Market Order.
type MarketOrderRejectTransaction struct {
	Transaction
	// Instrument is the Market Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Market Order.
	Units DecimalNumber `json:"units"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// PriceBound is the worst price that the client is willing to have the Market Order filled at.
	PriceBound PriceValue `json:"priceBound"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TradeClose details the Trade requested to be closed.
	TradeClose MarketOrderTradeClose `json:"tradeClose"`
	// LongPositionCloseout details the long Position to closeout.
	LongPositionCloseout MarketOrderPositionCloseout `json:"longPositionCloseout"`
	// ShortPositionCloseout details the short Position to closeout.
	ShortPositionCloseout MarketOrderPositionCloseout `json:"shortPositionCloseout"`
	// MarginCloseout details the Margin Closeout that this Market Order was created for.
	MarginCloseout MarketOrderMarginCloseout `json:"marginCloseout"`
	// DelayedTradeClose details the delayed Trade close that this Market Order was created for.
	DelayedTradeClose MarketOrderDelayedTradeClose `json:"delayedTradeClose"`
	// Reason is the reason that the Market Order was created.
	Reason MarketOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
	// TakeProfitOnFill specifies the Take Profit Order details.
	TakeProfitOnFill TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the Stop Loss Order details.
	StopLossOnFill StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the Guaranteed Stop Loss Order details.
	GuaranteedStopLossOnFill GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the Trailing Stop Loss Order details.
	TrailingStopLossOnFill TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions for the Trade.
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions"`
	// RejectReason is the reason that the Reject Transaction was created.
	RejectReason TransactionRejectReason `json:"rejectReason"`
}

// FixedPriceOrderTransaction represents a Transaction that creates a Fixed Price Order.
type FixedPriceOrderTransaction struct {
	Transaction
	// Instrument is the Fixed Price Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Fixed Price Order.
	Units DecimalNumber `json:"units"`
	// Price is the price specified for the Fixed Price Order.
	Price PriceValue `json:"price"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TradeState is the state that the trade resulting from the Fixed Price Order should be set to.
	TradeState string `json:"tradeState"`
	// Reason is the reason that the Fixed Price Order was created.
	Reason FixedPriceOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
	// TakeProfitOnFill specifies the Take Profit Order details.
	TakeProfitOnFill TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the Stop Loss Order details.
	StopLossOnFill StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the Guaranteed Stop Loss Order details.
	GuaranteedStopLossOnFill GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the Trailing Stop Loss Order details.
	TrailingStopLossOnFill TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions for the Trade.
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions"`
}

// LimitOrderTransaction represents a Transaction that creates a Limit Order.
type LimitOrderTransaction struct {
	Transaction
	// Instrument is the Limit Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Limit Order.
	Units DecimalNumber `json:"units"`
	// Price is the price threshold specified for the Limit Order.
	Price PriceValue `json:"price"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used for triggering.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Reason is the reason that the Limit Order was created.
	Reason LimitOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
	// TakeProfitOnFill specifies the Take Profit Order details.
	TakeProfitOnFill TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the Stop Loss Order details.
	StopLossOnFill StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the Guaranteed Stop Loss Order details.
	GuaranteedStopLossOnFill GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the Trailing Stop Loss Order details.
	TrailingStopLossOnFill TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions for the Trade.
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions"`
	// ReplacesOrderID is the ID of the Order that this Order replaces.
	ReplacesOrderID OrderID `json:"replacesOrderID"`
	// CancellingTransactionID is the ID of the Transaction that cancels the replaced Order.
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
}

// LimitOrderRejectTransaction represents a Transaction that rejects the creation of a Limit Order.
type LimitOrderRejectTransaction struct {
	Transaction
	// Instrument is the Limit Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Limit Order.
	Units DecimalNumber `json:"units"`
	// Price is the price threshold specified for the Limit Order.
	Price PriceValue `json:"price"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used for triggering.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Reason is the reason that the Limit Order was created.
	Reason LimitOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
	// TakeProfitOnFill specifies the Take Profit Order details.
	TakeProfitOnFill TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the Stop Loss Order details.
	StopLossOnFill StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the Guaranteed Stop Loss Order details.
	GuaranteedStopLossOnFill GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the Trailing Stop Loss Order details.
	TrailingStopLossOnFill TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions for the Trade.
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions"`
	// IntendedReplacesOrderID is the ID of the Order that this Order was intended to replace.
	IntendedReplacesOrderID OrderID `json:"intendedReplacesOrderID"`
	// RejectReason is the reason that the Reject Transaction was created.
	RejectReason TransactionRejectReason `json:"rejectReason"`
}

// StopOrderTransaction represents a Transaction that creates a Stop Order.
type StopOrderTransaction struct {
	Transaction
	// Instrument is the Stop Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Stop Order.
	Units DecimalNumber `json:"units"`
	// Price is the price threshold specified for the Stop Order.
	Price PriceValue `json:"price"`
	// PriceBound is the worst market price that may be used to fill this Stop Order.
	PriceBound PriceValue `json:"priceBound"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used for triggering.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Reason is the reason that the Stop Order was created.
	Reason StopOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
	// TakeProfitOnFill specifies the Take Profit Order details.
	TakeProfitOnFill TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the Stop Loss Order details.
	StopLossOnFill StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the Guaranteed Stop Loss Order details.
	GuaranteedStopLossOnFill GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the Trailing Stop Loss Order details.
	TrailingStopLossOnFill TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions for the Trade.
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions"`
	// ReplacesOrderID is the ID of the Order that this Order replaces.
	ReplacesOrderID OrderID `json:"replacesOrderID"`
	// CancellingTransactionID is the ID of the Transaction that cancels the replaced Order.
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
}

// StopOrderRejectTransaction represents a Transaction that rejects the creation of a Stop Order.
type StopOrderRejectTransaction struct {
	Transaction
	// Instrument is the Stop Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Stop Order.
	Units DecimalNumber `json:"units"`
	// Price is the price threshold specified for the Stop Order.
	Price PriceValue `json:"price"`
	// PriceBound is the worst market price that may be used to fill this Stop Order.
	PriceBound PriceValue `json:"priceBound"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used for triggering.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Reason is the reason that the Stop Order was created.
	Reason StopOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
	// TakeProfitOnFill specifies the Take Profit Order details.
	TakeProfitOnFill TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the Stop Loss Order details.
	StopLossOnFill StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the Guaranteed Stop Loss Order details.
	GuaranteedStopLossOnFill GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the Trailing Stop Loss Order details.
	TrailingStopLossOnFill TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions for the Trade.
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions"`
	// IntendedReplacesOrderID is the ID of the Order that this Order was intended to replace.
	IntendedReplacesOrderID OrderID `json:"intendedReplacesOrderID"`
	// RejectReason is the reason that the Reject Transaction was created.
	RejectReason TransactionRejectReason `json:"rejectReason"`
}

// MarketIfTouchedOrderTransaction represents a Transaction that creates a Market If Touched Order.
type MarketIfTouchedOrderTransaction struct {
	Transaction
	// Instrument is the Market If Touched Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Market If Touched Order.
	Units DecimalNumber `json:"units"`
	// Price is the price threshold specified for the Market If Touched Order.
	Price PriceValue `json:"price"`
	// PriceBound is the worst market price that may be used to fill this Order.
	PriceBound PriceValue `json:"priceBound"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime DateTime `json:"gtdTime"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used for triggering.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Reason is the reason that the Market If Touched Order was created.
	Reason MarketIfTouchedOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
	// TakeProfitOnFill specifies the Take Profit Order details.
	TakeProfitOnFill TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the Stop Loss Order details.
	StopLossOnFill StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the Guaranteed Stop Loss Order details.
	GuaranteedStopLossOnFill GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the Trailing Stop Loss Order details.
	TrailingStopLossOnFill TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions for the Trade.
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions"`
	// ReplacesOrderID is the ID of the Order that this Order replaces.
	ReplacesOrderID OrderID `json:"replacesOrderID"`
	// CancellingTransactionID is the ID of the Transaction that cancels the replaced Order.
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
}

// MarketIfTouchedOrderRejectTransaction represents a Transaction that rejects the creation of a Market If Touched Order.
type MarketIfTouchedOrderRejectTransaction struct {
	Transaction
	// Instrument is the Market If Touched Order's Instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the quantity requested to be filled by the Market If Touched Order.
	Units DecimalNumber `json:"units"`
	// Price is the price threshold specified for the Market If Touched Order.
	Price PriceValue `json:"price"`
	// PriceBound is the worst market price that may be used to fill this Order.
	PriceBound PriceValue `json:"priceBound"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// PositionFill specifies how Positions in the Account are modified when the Order is filled.
	PositionFill OrderPositionFill `json:"positionFill"`
	// TriggerCondition specifies which price component should be used for triggering.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Reason is the reason that the Market If Touched Order was created.
	Reason MarketIfTouchedOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
	// TakeProfitOnFill specifies the Take Profit Order details.
	TakeProfitOnFill TakeProfitDetails `json:"takeProfitOnFill"`
	// StopLossOnFill specifies the Stop Loss Order details.
	StopLossOnFill StopLossDetails `json:"stopLossOnFill"`
	// GuaranteedStopLossOnFill specifies the Guaranteed Stop Loss Order details.
	GuaranteedStopLossOnFill GuaranteedStopLossDetails `json:"guaranteedStopLossOnFill"`
	// TrailingStopLossOnFill specifies the Trailing Stop Loss Order details.
	TrailingStopLossOnFill TrailingStopLossDetails `json:"trailingStopLossOnFill"`
	// TradeClientExtensions are the client extensions for the Trade.
	TradeClientExtensions ClientExtensions `json:"tradeClientExtensions"`
	// IntendedReplacesOrderID is the ID of the Order that this Order was intended to replace.
	IntendedReplacesOrderID OrderID `json:"intendedReplacesOrderID"`
	// RejectReason is the reason that the Reject Transaction was created.
	RejectReason TransactionRejectReason `json:"rejectReason"`
}

// TakeProfitOrderTransaction represents a Transaction that creates a Take Profit Order.
type TakeProfitOrderTransaction struct {
	Transaction
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed.
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// Price is the price threshold specified for the Take Profit Order.
	Price PriceValue `json:"price"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used for triggering.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Reason is the reason that the Take Profit Order was created.
	Reason TakeProfitOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// OrderFillTransactionID is the ID of the OrderFill Transaction that caused this Order to be created.
	OrderFillTransactionID *TransactionID `json:"orderFillTransactionID,omitempty"`
	// ReplacesOrderID is the ID of the Order that this Order replaces.
	ReplacesOrderID *OrderID `json:"replacesOrderID,omitempty"`
	// CancellingTransactionID is the ID of the Transaction that cancels the replaced Order.
	CancellingTransactionID *TransactionID `json:"cancellingTransactionID,omitempty"`
}

// TakeProfitOrderRejectTransaction represents a Transaction that rejects the creation of a Take Profit Order.
type TakeProfitOrderRejectTransaction struct {
	Transaction
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed.
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// Price is the price threshold specified for the Take Profit Order.
	Price PriceValue `json:"price"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used for triggering.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Reason is the reason that the Take Profit Order was created.
	Reason TakeProfitOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// OrderFillTransactionID is the ID of the OrderFill Transaction that caused this Order to be created.
	OrderFillTransactionID TransactionID `json:"orderFillTransactionID"`
	// IntendedReplacesOrderID is the ID of the Order that this Order was intended to replace.
	IntendedReplacesOrderID OrderID `json:"intendedReplacesOrderID"`
	// RejectReason is the reason that the Reject Transaction was created.
	RejectReason TransactionRejectReason `json:"rejectReason"`
}

// StopLossOrderTransaction represents a Transaction that creates a Stop Loss Order.
type StopLossOrderTransaction struct {
	Transaction
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed.
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// Price is the price threshold specified for the Stop Loss Order.
	Price *PriceValue `json:"price,omitempty"`
	// Distance specifies the distance from the current price to use as the Stop Loss Order price.
	Distance *DecimalNumber `json:"distance,omitempty"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used for triggering.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Guaranteed is deprecated. Indicates if the Stop Loss Order is guaranteed.
	Guaranteed bool `json:"guaranteed"`
	// GuaranteedExecutionPremium is the fee charged if the Stop Loss Order is guaranteed.
	GuaranteedExecutionPremium DecimalNumber `json:"guaranteedExecutionPremium"`
	// Reason is the reason that the Stop Loss Order was created.
	Reason StopLossOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
	// OrderFillTransactionID is the ID of the OrderFill Transaction that caused this Order to be created.
	OrderFillTransactionID TransactionID `json:"orderFillTransactionID"`
	// ReplacesOrderID is the ID of the Order that this Order replaces.
	ReplacesOrderID OrderID `json:"replacesOrderID"`
	// CancellingTransactionID is the ID of the Transaction that cancels the replaced Order.
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
}

// StopLossOrderRejectTransaction represents a Transaction that rejects the creation of a Stop Loss Order.
type StopLossOrderRejectTransaction struct {
	Transaction
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed.
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// Price is the price threshold specified for the Stop Loss Order.
	Price *PriceValue `json:"price,omitempty"`
	// Distance specifies the distance from the current price to use as the Stop Loss Order price.
	Distance *DecimalNumber `json:"distance,omitempty"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used for triggering.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Guaranteed is deprecated. Indicates if the Stop Loss Order is guaranteed.
	Guaranteed bool `json:"guaranteed"`
	// Reason is the reason that the Stop Loss Order was created.
	Reason StopLossOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// OrderFillTransactionID is the ID of the OrderFill Transaction that caused this Order to be created.
	OrderFillTransactionID TransactionID `json:"orderFillTransactionID"`
	// IntendedReplacesOrderID is the ID of the Order that this Order was intended to replace.
	IntendedReplacesOrderID OrderID `json:"intendedReplacesOrderID"`
	// RejectReason is the reason that the Reject Transaction was created.
	RejectReason TransactionRejectReason `json:"rejectReason"`
}

// GuaranteedStopLossOrderTransaction represents a Transaction that creates a Guaranteed Stop Loss Order.
type GuaranteedStopLossOrderTransaction struct {
	Transaction
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed.
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// Price is the price threshold specified for the Guaranteed Stop Loss Order.
	Price *PriceValue `json:"price,omitempty"`
	// Distance specifies the distance from the current price to use as the Order price.
	Distance *DecimalNumber `json:"distance,omitempty"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used for triggering.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// GuaranteedExecutionPremium is the fee charged for the Guaranteed Stop Loss Order.
	GuaranteedExecutionPremium DecimalNumber `json:"guaranteedExecutionPremium"`
	// Reason is the reason that the Guaranteed Stop Loss Order was created.
	Reason GuaranteedStopLossOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// OrderFillTransactionID is the ID of the OrderFill Transaction that caused this Order to be created.
	OrderFillTransactionID TransactionID `json:"orderFillTransactionID"`
	// ReplacesOrderID is the ID of the Order that this Order replaces.
	ReplacesOrderID OrderID `json:"replacesOrderID"`
	// CancellingTransactionID is the ID of the Transaction that cancels the replaced Order.
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
}

// GuaranteedStopLossOrderRejectTransaction represents a Transaction that rejects the creation of a Guaranteed Stop Loss Order.
type GuaranteedStopLossOrderRejectTransaction struct {
	Transaction
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed.
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// Price is the price threshold specified for the Guaranteed Stop Loss Order.
	Price *PriceValue `json:"price,omitempty"`
	// Distance specifies the distance from the current price to use as the Order price.
	Distance *DecimalNumber `json:"distance,omitempty"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used for triggering.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Reason is the reason that the Guaranteed Stop Loss Order was created.
	Reason GuaranteedStopLossOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// OrderFillTransactionID is the ID of the OrderFill Transaction that caused this Order to be created.
	OrderFillTransactionID TransactionID `json:"orderFillTransactionID"`
	// IntendedReplacesOrderID is the ID of the Order that this Order was intended to replace.
	IntendedReplacesOrderID OrderID `json:"intendedReplacesOrderID"`
	// RejectReason is the reason that the Reject Transaction was created.
	RejectReason TransactionRejectReason `json:"rejectReason"`
}

// TrailingStopLossOrderTransaction represents a Transaction that creates a Trailing Stop Loss Order.
type TrailingStopLossOrderTransaction struct {
	Transaction
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed.
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// Distance is the price distance specified for the Trailing Stop Loss Order.
	Distance DecimalNumber `json:"distance"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used for triggering.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Reason is the reason that the Trailing Stop Loss Order was created.
	Reason TrailingStopLossOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// OrderFillTransactionID is the ID of the OrderFill Transaction that caused this Order to be created.
	OrderFillTransactionID TransactionID `json:"orderFillTransactionID"`
	// ReplacesOrderID is the ID of the Order that this Order replaces.
	ReplacesOrderID OrderID `json:"replacesOrderID"`
	// CancellingTransactionID is the ID of the Transaction that cancels the replaced Order.
	CancellingTransactionID TransactionID `json:"cancellingTransactionID"`
}

// TrailingStopLossOrderRejectTransaction represents a Transaction that rejects the creation of a Trailing Stop Loss Order.
type TrailingStopLossOrderRejectTransaction struct {
	Transaction
	// TradeID is the ID of the Trade to close when the price threshold is breached.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade to be closed.
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// Distance is the price distance specified for the Trailing Stop Loss Order.
	Distance DecimalNumber `json:"distance"`
	// TimeInForce specifies how long the Order should remain pending.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// TriggerCondition specifies which price component should be used for triggering.
	TriggerCondition OrderTriggerCondition `json:"triggerCondition"`
	// Reason is the reason that the Trailing Stop Loss Order was created.
	Reason TrailingStopLossOrderReason `json:"reason"`
	// ClientExtensions are the client extensions for the Order.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// OrderFillTransactionID is the ID of the OrderFill Transaction that caused this Order to be created.
	OrderFillTransactionID TransactionID `json:"orderFillTransactionID"`
	// IntendedReplacesOrderID is the ID of the Order that this Order was intended to replace.
	IntendedReplacesOrderID OrderID `json:"intendedReplacesOrderID"`
	// RejectReason is the reason that the Reject Transaction was created.
	RejectReason TransactionRejectReason `json:"rejectReason"`
}

// OrderFillTransaction represents a Transaction that fills an Order.
type OrderFillTransaction struct {
	Transaction
	// OrderID is the ID of the Order filled.
	OrderID OrderID `json:"orderID"`
	// ClientOrderID is the client Order ID of the Order filled.
	ClientOrderID *ClientID `json:"clientOrderID,omitempty"`
	// Instrument is the name of the filled Order's instrument.
	Instrument InstrumentName `json:"instrument"`
	// Units is the number of units filled by the Order.
	Units DecimalNumber `json:"units"`
	// GainQuoteHomeConversionFactor is the conversion factor for gains.
	GainQuoteHomeConversionFactor DecimalNumber `json:"gainQuoteHomeConversionFactor"`
	// LossQuoteHomeConversionFactor is the conversion factor for losses.
	LossQuoteHomeConversionFactor DecimalNumber `json:"lossQuoteHomeConversionFactor"`
	// HomeConversionFactors is the HomeConversionFactors in effect at the time of the fill.
	HomeConversionFactors HomeConversionFactors `json:"homeConversionFactors"`
	// Price is the average price that the units were filled at (deprecated).
	Price PriceValue `json:"price"`
	// FullVWAP is the price in effect for the account at the time of the Order fill.
	FullVWAP PriceValue `json:"fullVWAP"`
	// FullPrice is the PriceValue containing information about the price.
	FullPrice ClientPrice `json:"fullPrice"`
	// Reason is the reason that the Order was filled.
	Reason OrderFillReason `json:"reason"`
	// PL is the profit or loss incurred when the Order was filled.
	PL AccountUnits `json:"pl"`
	// QuotePL is the profit or loss incurred in the quote currency.
	QuotePL DecimalNumber `json:"quotePL"`
	// Financing is the financing paid or collected when the Order was filled.
	Financing AccountUnits `json:"financing"`
	// BaseFinancing is the financing paid or collected in the base currency.
	BaseFinancing DecimalNumber `json:"baseFinancing"`
	// QuoteFinancing is the financing paid or collected in the quote currency.
	QuoteFinancing *DecimalNumber `json:"quoteFinancing,omitempty"`
	// Commission is the commission charged in the Account's home currency.
	Commission AccountUnits `json:"commission"`
	// GuaranteedExecutionFee is the total guaranteed execution fee charged.
	GuaranteedExecutionFee AccountUnits `json:"guaranteedExecutionFee"`
	// QuoteGuaranteedExecutionFee is the guaranteed execution fee in quote currency.
	QuoteGuaranteedExecutionFee DecimalNumber `json:"quoteGuaranteedExecutionFee"`
	// AccountBalance is the Account's balance after the Order was filled.
	AccountBalance AccountUnits `json:"accountBalance"`
	// TradeOpened is the Trade that was opened when the Order was filled.
	TradeOpened *TradeOpen `json:"tradeOpened,omitempty"`
	// TradesClosed are the Trades that were closed when the Order was filled.
	TradesClosed []TradeReduce `json:"tradesClosed,omitempty"`
	// TradeReduced is the Trade that was reduced when the Order was filled.
	TradeReduced *TradeReduce `json:"tradeReduced,omitempty"`
	// HalfSpreadCost is the half spread cost for the Order.
	HalfSpreadCost AccountUnits `json:"halfSpreadCost"`
}

// OrderCancelTransaction represents a Transaction that cancels an Order.
type OrderCancelTransaction struct {
	Transaction
	// OrderID is the ID of the Order cancelled.
	OrderID OrderID `json:"orderID"`
	// ClientOrderID is the client Order ID of the Order cancelled.
	ClientOrderID *ClientID `json:"clientOrderID,omitempty"`
	// Reason is the reason that the Order was cancelled.
	Reason OrderCancelReason `json:"reason"`
	// ReplacedByOrderID is the ID of the Order that replaced this Order.
	ReplacedByOrderID OrderID `json:"replacedByOrderID"`
}

// OrderCancelRejectTransaction represents a Transaction that rejects the cancellation of an Order.
type OrderCancelRejectTransaction struct {
	Transaction
	// OrderID is the ID of the Order intended to be cancelled.
	OrderID OrderID `json:"orderID"`
	// ClientOrderID is the client Order ID of the Order intended to be cancelled.
	ClientOrderID *ClientID `json:"clientOrderID,omitempty"`
	// RejectReason is the reason that the Reject Transaction was created.
	RejectReason TransactionRejectReason `json:"rejectReason"`
}

// OrderClientExtensionsModifyTransaction represents a Transaction that modifies an Order's client extensions.
type OrderClientExtensionsModifyTransaction struct {
	Transaction
	// OrderID is the ID of the Order whose client extensions are to be modified.
	OrderID OrderID `json:"orderID"`
	// ClientOrderID is the original client Order ID of the Order.
	ClientOrderID *ClientID `json:"clientOrderID,omitempty"`
	// ClientExtensionsModify contains the new client extensions for the Order.
	ClientExtensionsModify *ClientExtensions `json:"clientExtensionsModify,omitempty"`
	// TradeClientExtensionsModify contains the new client extensions for the Trade.
	TradeClientExtensionsModify *ClientExtensions `json:"tradeClientExtensionsModify,omitempty"`
}

// OrderClientExtensionsModifyRejectTransaction represents a Transaction that rejects the modification of an Order's client extensions.
type OrderClientExtensionsModifyRejectTransaction struct {
	Transaction
	// OrderID is the ID of the Order whose client extensions are to be modified.
	OrderID OrderID `json:"orderID"`
	// ClientOrderID is the original client Order ID of the Order.
	ClientOrderID *ClientID `json:"clientOrderID,omitempty"`
	// ClientExtensionsModify contains the new client extensions for the Order.
	ClientExtensionsModify *ClientExtensions `json:"clientExtensionsModify,omitempty"`
	// TradeClientExtensionsModify contains the new client extensions for the Trade.
	TradeClientExtensionsModify *ClientExtensions `json:"tradeClientExtensionsModify,omitempty"`
	// RejectReason is the reason that the Reject Transaction was created.
	RejectReason TransactionRejectReason `json:"rejectReason"`
}

// TradeClientExtensionsModifyTransaction represents a Transaction that modifies a Trade's client extensions.
type TradeClientExtensionsModifyTransaction struct {
	Transaction
	// TradeID is the ID of the Trade whose client extensions are to be modified.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the original client Trade ID of the Trade.
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// TradeClientExtensionsModify contains the new client extensions for the Trade.
	TradeClientExtensionsModify ClientExtensions `json:"tradeClientExtensionsModify"`
}

// TradeClientExtensionsModifyRejectTransaction represents a Transaction that rejects the modification of a Trade's client extensions.
type TradeClientExtensionsModifyRejectTransaction struct {
	Transaction
	// TradeID is the ID of the Trade whose client extensions are to be modified.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the original client Trade ID of the Trade.
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// TradeClientExtensionsModify contains the new client extensions for the Trade.
	TradeClientExtensionsModify ClientExtensions `json:"tradeClientExtensionsModify"`
	// RejectReason is the reason that the Reject Transaction was created.
	RejectReason TransactionRejectReason `json:"rejectReason"`
} // MarginCallEnterTransaction represents a Transaction that indicates an Account has entered a margin call state.
type MarginCallEnterTransaction struct {
	Transaction
}

// MarginCallExtendTransaction represents a Transaction that indicates a margin call state has been extended.
type MarginCallExtendTransaction struct {
	Transaction
	// ExtensionNumber is the number of the extension within the current margin call.
	ExtensionNumber int `json:"extensionNumber"`
}

// MarginCallExitTransaction represents a Transaction that indicates an Account has exited a margin call state.
type MarginCallExitTransaction struct {
	Transaction
}

// DelayedTradeClosureTransaction represents a Transaction that indicates a delayed trade closure.
type DelayedTradeClosureTransaction struct {
	Transaction
	// Reason is the reason for the delayed trade closure.
	Reason MarketOrderReason `json:"reason"`
	// TradeIDs are the IDs of the Trades that will be closed.
	TradeIDs []TradeID `json:"tradeIDs"`
}

// DailyFinancingTransaction represents a Transaction that accounts for daily financing charges.
type DailyFinancingTransaction struct {
	Transaction
	// Financing is the amount of financing paid or collected for the Account.
	Financing AccountUnits `json:"financing"`
	// AccountBalance is the Account's balance after daily financing.
	AccountBalance AccountUnits `json:"accountBalance"`
	// AccountFinancingMode describes how financing charges are calculated.
	AccountFinancingMode AccountFinancingMode `json:"accountFinancingMode"`
	// PositionFinancings contains the financing paid/collected for each Position.
	PositionFinancings []PositionFinancing `json:"positionFinancings"`
}

// DividendAdjustmentTransaction represents a Transaction that accounts for dividend adjustments.
type DividendAdjustmentTransaction struct {
	Transaction
	// Instrument is the name of the Instrument for the dividend adjustment.
	Instrument InstrumentName `json:"instrument"`
	// DividendAdjustment is the total dividend adjustment for the Account.
	DividendAdjustment AccountUnits `json:"dividendAdjustment"`
	// QuoteDividendAdjustment is the dividend adjustment in the quote currency.
	QuoteDividendAdjustment DecimalNumber `json:"quoteDividendAdjustment"`
	// HomeConversionFactors is the HomeConversionFactors in effect at the time.
	HomeConversionFactors HomeConversionFactors `json:"homeConversionFactors"`
	// AccountBalance is the Account's balance after the dividend adjustment.
	AccountBalance AccountUnits `json:"accountBalance"`
	// OpenTradeDividendAdjustments are the adjustments for each open Trade.
	OpenTradeDividendAdjustments []OpenTradeDividendAdjustment `json:"openTradeDividendAdjustments"`
}

// ResetResettablePLTransaction represents a Transaction that resets the Account's resettable PL counters.
type ResetResettablePLTransaction struct {
	Transaction
}

// Transaction-related Definitions

// TransactionID is the unique identifier of a Transaction.
type TransactionID = string

// TransactionType represents the type of a Transaction.
type TransactionType string

const (
	// TransactionTypeCreate represents the creation of an Account.
	TransactionTypeCreate TransactionType = "CREATE"
	// TransactionTypeClose represents the closing of an Account.
	TransactionTypeClose TransactionType = "CLOSE"
	// TransactionTypeReopen represents the re-opening of a closed Account.
	TransactionTypeReopen TransactionType = "REOPEN"
	// TransactionTypeClientConfigure represents the configuration of an Account by a client.
	TransactionTypeClientConfigure TransactionType = "CLIENT_CONFIGURE"
	// TransactionTypeClientConfigureReject represents the rejection of configuration of an Account by a client.
	TransactionTypeClientConfigureReject TransactionType = "CLIENT_CONFIGURE_REJECT"
	// TransactionTypeTransferFunds represents a transfer of funds in/out of an Account.
	TransactionTypeTransferFunds TransactionType = "TRANSFER_FUNDS"
	// TransactionTypeTransferFundsReject represents the rejection of a transfer of funds in/out of an Account.
	TransactionTypeTransferFundsReject TransactionType = "TRANSFER_FUNDS_REJECT"
	// TransactionTypeMarketOrder represents the creation of a Market Order in an Account.
	TransactionTypeMarketOrder TransactionType = "MARKET_ORDER"
	// TransactionTypeMarketOrderReject represents the rejection of the creation of a Market Order.
	TransactionTypeMarketOrderReject TransactionType = "MARKET_ORDER_REJECT"
	// TransactionTypeFixedPriceOrder represents the creation of a Fixed Price Order in an Account.
	TransactionTypeFixedPriceOrder TransactionType = "FIXED_PRICE_ORDER"
	// TransactionTypeLimitOrder represents the creation of a Limit Order in an Account.
	TransactionTypeLimitOrder TransactionType = "LIMIT_ORDER"
	// TransactionTypeLimitOrderReject represents the rejection of the creation of a Limit Order.
	TransactionTypeLimitOrderReject TransactionType = "LIMIT_ORDER_REJECT"
	// TransactionTypeStopOrder represents the creation of a Stop Order in an Account.
	TransactionTypeStopOrder TransactionType = "STOP_ORDER"
	// TransactionTypeStopOrderReject represents the rejection of the creation of a Stop Order.
	TransactionTypeStopOrderReject TransactionType = "STOP_ORDER_REJECT"
	// TransactionTypeMarketIfTouchedOrder represents the creation of a Market If Touched Order in an Account.
	TransactionTypeMarketIfTouchedOrder TransactionType = "MARKET_IF_TOUCHED_ORDER"
	// TransactionTypeMarketIfTouchedOrderReject represents the rejection of the creation of a Market If Touched Order.
	TransactionTypeMarketIfTouchedOrderReject TransactionType = "MARKET_IF_TOUCHED_ORDER_REJECT"
	// TransactionTypeTakeProfitOrder represents the creation of a Take Profit Order in an Account.
	TransactionTypeTakeProfitOrder TransactionType = "TAKE_PROFIT_ORDER"
	// TransactionTypeTakeProfitOrderReject represents the rejection of the creation of a Take Profit Order.
	TransactionTypeTakeProfitOrderReject TransactionType = "TAKE_PROFIT_ORDER_REJECT"
	// TransactionTypeStopLossOrder represents the creation of a Stop Loss Order in an Account.
	TransactionTypeStopLossOrder TransactionType = "STOP_LOSS_ORDER"
	// TransactionTypeStopLossOrderReject represents the rejection of the creation of a Stop Loss Order.
	TransactionTypeStopLossOrderReject TransactionType = "STOP_LOSS_ORDER_REJECT"
	// TransactionTypeGuaranteedStopLossOrder represents the creation of a Guaranteed Stop Loss Order in an Account.
	TransactionTypeGuaranteedStopLossOrder TransactionType = "GUARANTEED_STOP_LOSS_ORDER"
	// TransactionTypeGuaranteedStopLossOrderReject represents the rejection of the creation of a Guaranteed Stop Loss Order.
	TransactionTypeGuaranteedStopLossOrderReject TransactionType = "GUARANTEED_STOP_LOSS_ORDER_REJECT"
	// TransactionTypeTrailingStopLossOrder represents the creation of a Trailing Stop Loss Order in an Account.
	TransactionTypeTrailingStopLossOrder TransactionType = "TRAILING_STOP_LOSS_ORDER"
	// TransactionTypeTrailingStopLossOrderReject represents the rejection of the creation of a Trailing Stop Loss Order.
	TransactionTypeTrailingStopLossOrderReject TransactionType = "TRAILING_STOP_LOSS_ORDER_REJECT"
	// TransactionTypeOrderFill represents the filling of an Order in an Account.
	TransactionTypeOrderFill TransactionType = "ORDER_FILL"
	// TransactionTypeOrderCancel represents the cancellation of an Order in an Account.
	TransactionTypeOrderCancel TransactionType = "ORDER_CANCEL"
	// TransactionTypeOrderCancelReject represents the rejection of the cancellation of an Order.
	TransactionTypeOrderCancelReject TransactionType = "ORDER_CANCEL_REJECT"
	// TransactionTypeOrderClientExtensionsModify represents the modification of an Order's client extensions.
	TransactionTypeOrderClientExtensionsModify TransactionType = "ORDER_CLIENT_EXTENSIONS_MODIFY"
	// TransactionTypeOrderClientExtensionsModifyReject represents the rejection of the modification of an Order's client extensions.
	TransactionTypeOrderClientExtensionsModifyReject TransactionType = "ORDER_CLIENT_EXTENSIONS_MODIFY_REJECT"
	// TransactionTypeTradeClientExtensionsModify represents the modification of a Trade's client extensions.
	TransactionTypeTradeClientExtensionsModify TransactionType = "TRADE_CLIENT_EXTENSIONS_MODIFY"
	// TransactionTypeTradeClientExtensionsModifyReject represents the rejection of the modification of a Trade's client extensions.
	TransactionTypeTradeClientExtensionsModifyReject TransactionType = "TRADE_CLIENT_EXTENSIONS_MODIFY_REJECT"
	// TransactionTypeMarginCallEnter represents the entering of a margin call state for an Account.
	TransactionTypeMarginCallEnter TransactionType = "MARGIN_CALL_ENTER"
	// TransactionTypeMarginCallExtend represents the extension of an Account's margin call state.
	TransactionTypeMarginCallExtend TransactionType = "MARGIN_CALL_EXTEND"
	// TransactionTypeMarginCallExit represents the exit from a margin call state for an Account.
	TransactionTypeMarginCallExit TransactionType = "MARGIN_CALL_EXIT"
	// TransactionTypeDelayedTradeClosure represents the delayed closure of a Trade.
	TransactionTypeDelayedTradeClosure TransactionType = "DELAYED_TRADE_CLOSURE"
	// TransactionTypeDailyFinancing represents the daily financing of an Account's open positions.
	TransactionTypeDailyFinancing TransactionType = "DAILY_FINANCING"
	// TransactionTypeDividendAdjustment represents a dividend adjustment for an Account.
	TransactionTypeDividendAdjustment TransactionType = "DIVIDEND_ADJUSTMENT"
	// TransactionTypeResetResettablePL represents the resetting of the Account's resettable PL counters.
	TransactionTypeResetResettablePL TransactionType = "RESET_RESETTABLE_PL"
)

// FundingReason represents the reason that an Account is being funded.
type FundingReason string

const (
	// FundingReasonClientFunding indicates the client has initiated the funding.
	FundingReasonClientFunding FundingReason = "CLIENT_FUNDING"
	// FundingReasonAccountTransfer indicates funds are being transferred between two accounts.
	FundingReasonAccountTransfer FundingReason = "ACCOUNT_TRANSFER"
	// FundingReasonDivisionMigration indicates funds are being transferred as part of a division migration.
	FundingReasonDivisionMigration FundingReason = "DIVISION_MIGRATION"
	// FundingReasonSiteMigration indicates funds are being transferred as part of a site migration.
	FundingReasonSiteMigration FundingReason = "SITE_MIGRATION"
	// FundingReasonAdjustment indicates funds are being transferred as a general adjustment.
	FundingReasonAdjustment FundingReason = "ADJUSTMENT"
)

// MarketOrderReason represents the reason that a Market Order was created.
type MarketOrderReason string

const (
	// MarketOrderReasonClientOrder indicates the Market Order was created at the request of a client.
	MarketOrderReasonClientOrder MarketOrderReason = "CLIENT_ORDER"
	// MarketOrderReasonTradeClose indicates the Market Order was created to close a Trade.
	MarketOrderReasonTradeClose MarketOrderReason = "TRADE_CLOSE"
	// MarketOrderReasonPositionCloseout indicates the Market Order was created to close a Position.
	MarketOrderReasonPositionCloseout MarketOrderReason = "POSITION_CLOSEOUT"
	// MarketOrderReasonMarginCloseout indicates the Market Order was created due to margin closeout.
	MarketOrderReasonMarginCloseout MarketOrderReason = "MARGIN_CLOSEOUT"
	// MarketOrderReasonDelayedTradeClose indicates the Market Order was created for a delayed Trade close.
	MarketOrderReasonDelayedTradeClose MarketOrderReason = "DELAYED_TRADE_CLOSE"
)

// FixedPriceOrderReason represents the reason that a Fixed Price Order was created.
type FixedPriceOrderReason string

const (
	// FixedPriceOrderReasonPlatformAccountMigration indicates the order was created as part of platform account migration.
	FixedPriceOrderReasonPlatformAccountMigration FixedPriceOrderReason = "PLATFORM_ACCOUNT_MIGRATION"
	// FixedPriceOrderReasonTradeCloseDivisionAccountMigration indicates the order was created to close a Trade due to division account migration.
	FixedPriceOrderReasonTradeCloseDivisionAccountMigration FixedPriceOrderReason = "TRADE_CLOSE_DIVISION_ACCOUNT_MIGRATION"
	// FixedPriceOrderReasonTradeCloseSiteAccountMigration indicates the order was created to close a Trade due to site account migration.
	FixedPriceOrderReasonTradeCloseSiteAccountMigration FixedPriceOrderReason = "TRADE_CLOSE_SITE_ACCOUNT_MIGRATION"
	// FixedPriceOrderReasonTradeCloseAdministrativeAction indicates the order was created to close a Trade due to administrative action.
	FixedPriceOrderReasonTradeCloseAdministrativeAction FixedPriceOrderReason = "TRADE_CLOSE_ADMINISTRATIVE_ACTION"
)

// LimitOrderReason represents the reason that a Limit Order was created.
type LimitOrderReason string

const (
	// LimitOrderReasonClientOrder indicates the Limit Order was created at the request of a client.
	LimitOrderReasonClientOrder LimitOrderReason = "CLIENT_ORDER"
	// LimitOrderReasonReplacement indicates the Limit Order was created to replace an existing Order.
	LimitOrderReasonReplacement LimitOrderReason = "REPLACEMENT"
)

// StopOrderReason represents the reason that a Stop Order was created.
type StopOrderReason string

const (
	// StopOrderReasonClientOrder indicates the Stop Order was created at the request of a client.
	StopOrderReasonClientOrder StopOrderReason = "CLIENT_ORDER"
	// StopOrderReasonReplacement indicates the Stop Order was created to replace an existing Order.
	StopOrderReasonReplacement StopOrderReason = "REPLACEMENT"
)

// MarketIfTouchedOrderReason represents the reason that a Market If Touched Order was created.
type MarketIfTouchedOrderReason string

const (
	// MarketIfTouchedOrderReasonClientOrder indicates the MIT Order was created at the request of a client.
	MarketIfTouchedOrderReasonClientOrder MarketIfTouchedOrderReason = "CLIENT_ORDER"
	// MarketIfTouchedOrderReasonReplacement indicates the MIT Order was created to replace an existing Order.
	MarketIfTouchedOrderReasonReplacement MarketIfTouchedOrderReason = "REPLACEMENT"
)

// TakeProfitOrderReason represents the reason that a Take Profit Order was created.
type TakeProfitOrderReason string

const (
	// TakeProfitOrderReasonClientOrder indicates the Take Profit Order was created at the request of a client.
	TakeProfitOrderReasonClientOrder TakeProfitOrderReason = "CLIENT_ORDER"
	// TakeProfitOrderReasonReplacement indicates the Take Profit Order was created to replace an existing Order.
	TakeProfitOrderReasonReplacement TakeProfitOrderReason = "REPLACEMENT"
	// TakeProfitOrderReasonOnFill indicates the Take Profit Order was created when an Order was filled.
	TakeProfitOrderReasonOnFill TakeProfitOrderReason = "ON_FILL"
)

// StopLossOrderReason represents the reason that a Stop Loss Order was created.
type StopLossOrderReason string

const (
	// StopLossOrderReasonClientOrder indicates the Stop Loss Order was created at the request of a client.
	StopLossOrderReasonClientOrder StopLossOrderReason = "CLIENT_ORDER"
	// StopLossOrderReasonReplacement indicates the Stop Loss Order was created to replace an existing Order.
	StopLossOrderReasonReplacement StopLossOrderReason = "REPLACEMENT"
	// StopLossOrderReasonOnFill indicates the Stop Loss Order was created when an Order was filled.
	StopLossOrderReasonOnFill StopLossOrderReason = "ON_FILL"
)

// GuaranteedStopLossOrderReason represents the reason that a Guaranteed Stop Loss Order was created.
type GuaranteedStopLossOrderReason string

const (
	// GuaranteedStopLossOrderReasonClientOrder indicates the Guaranteed Stop Loss Order was created at the request of a client.
	GuaranteedStopLossOrderReasonClientOrder GuaranteedStopLossOrderReason = "CLIENT_ORDER"
	// GuaranteedStopLossOrderReasonReplacement indicates the Guaranteed Stop Loss Order was created to replace an existing Order.
	GuaranteedStopLossOrderReasonReplacement GuaranteedStopLossOrderReason = "REPLACEMENT"
	// GuaranteedStopLossOrderReasonOnFill indicates the Guaranteed Stop Loss Order was created when an Order was filled.
	GuaranteedStopLossOrderReasonOnFill GuaranteedStopLossOrderReason = "ON_FILL"
)

// TrailingStopLossOrderReason represents the reason that a Trailing Stop Loss Order was created.
type TrailingStopLossOrderReason string

const (
	// TrailingStopLossOrderReasonClientOrder indicates the Trailing Stop Loss Order was created at the request of a client.
	TrailingStopLossOrderReasonClientOrder TrailingStopLossOrderReason = "CLIENT_ORDER"
	// TrailingStopLossOrderReasonReplacement indicates the Trailing Stop Loss Order was created to replace an existing Order.
	TrailingStopLossOrderReasonReplacement TrailingStopLossOrderReason = "REPLACEMENT"
	// TrailingStopLossOrderReasonOnFill indicates the Trailing Stop Loss Order was created when an Order was filled.
	TrailingStopLossOrderReasonOnFill TrailingStopLossOrderReason = "ON_FILL"
)

// OrderFillReason represents the reason that an Order was filled.
type OrderFillReason string

const (
	// OrderFillReasonLimitOrder indicates a Limit Order was filled.
	OrderFillReasonLimitOrder OrderFillReason = "LIMIT_ORDER"
	// OrderFillReasonStopOrder indicates a Stop Order was filled.
	OrderFillReasonStopOrder OrderFillReason = "STOP_ORDER"
	// OrderFillReasonMarketIfTouchedOrder indicates a Market If Touched Order was filled.
	OrderFillReasonMarketIfTouchedOrder OrderFillReason = "MARKET_IF_TOUCHED_ORDER"
	// OrderFillReasonTakeProfitOrder indicates a Take Profit Order was filled.
	OrderFillReasonTakeProfitOrder OrderFillReason = "TAKE_PROFIT_ORDER"
	// OrderFillReasonStopLossOrder indicates a Stop Loss Order was filled.
	OrderFillReasonStopLossOrder OrderFillReason = "STOP_LOSS_ORDER"
	// OrderFillReasonGuaranteedStopLossOrder indicates a Guaranteed Stop Loss Order was filled.
	OrderFillReasonGuaranteedStopLossOrder OrderFillReason = "GUARANTEED_STOP_LOSS_ORDER"
	// OrderFillReasonTrailingStopLossOrder indicates a Trailing Stop Loss Order was filled.
	OrderFillReasonTrailingStopLossOrder OrderFillReason = "TRAILING_STOP_LOSS_ORDER"
	// OrderFillReasonMarketOrder indicates a Market Order was filled.
	OrderFillReasonMarketOrder OrderFillReason = "MARKET_ORDER"
	// OrderFillReasonMarketOrderTradeClose indicates a Market Order to close a Trade was filled.
	OrderFillReasonMarketOrderTradeClose OrderFillReason = "MARKET_ORDER_TRADE_CLOSE"
	// OrderFillReasonMarketOrderPositionCloseout indicates a Market Order to closeout a Position was filled.
	OrderFillReasonMarketOrderPositionCloseout OrderFillReason = "MARKET_ORDER_POSITION_CLOSEOUT"
	// OrderFillReasonMarketOrderMarginCloseout indicates a Market Order for margin closeout was filled.
	OrderFillReasonMarketOrderMarginCloseout OrderFillReason = "MARKET_ORDER_MARGIN_CLOSEOUT"
	// OrderFillReasonMarketOrderDelayedTradeClose indicates a Market Order for delayed Trade close was filled.
	OrderFillReasonMarketOrderDelayedTradeClose OrderFillReason = "MARKET_ORDER_DELAYED_TRADE_CLOSE"
	// OrderFillReasonFixedPriceOrder indicates a Fixed Price Order was filled.
	OrderFillReasonFixedPriceOrder OrderFillReason = "FIXED_PRICE_ORDER"
	// OrderFillReasonFixedPriceOrderPlatformAccountMigration indicates a Fixed Price Order for platform account migration was filled.
	OrderFillReasonFixedPriceOrderPlatformAccountMigration OrderFillReason = "FIXED_PRICE_ORDER_PLATFORM_ACCOUNT_MIGRATION"
	// OrderFillReasonFixedPriceOrderDivisionAccountMigration indicates a Fixed Price Order for division account migration was filled.
	OrderFillReasonFixedPriceOrderDivisionAccountMigration OrderFillReason = "FIXED_PRICE_ORDER_DIVISION_ACCOUNT_MIGRATION"
	// OrderFillReasonFixedPriceOrderAdministrativeAction indicates a Fixed Price Order for administrative action was filled.
	OrderFillReasonFixedPriceOrderAdministrativeAction OrderFillReason = "FIXED_PRICE_ORDER_ADMINISTRATIVE_ACTION"
)

// OrderCancelReason represents the reason that an Order was cancelled.
type OrderCancelReason string

const (
	// OrderCancelReasonInternalServerError indicates an internal server error cancelled the Order.
	OrderCancelReasonInternalServerError OrderCancelReason = "INTERNAL_SERVER_ERROR"
	// OrderCancelReasonAccountLocked indicates the Account is locked.
	OrderCancelReasonAccountLocked OrderCancelReason = "ACCOUNT_LOCKED"
	// OrderCancelReasonAccountNewPositionsLocked indicates new positions are locked for the Account.
	OrderCancelReasonAccountNewPositionsLocked OrderCancelReason = "ACCOUNT_NEW_POSITIONS_LOCKED"
	// OrderCancelReasonAccountOrderCreationLocked indicates order creation is locked for the Account.
	OrderCancelReasonAccountOrderCreationLocked OrderCancelReason = "ACCOUNT_ORDER_CREATION_LOCKED"
	// OrderCancelReasonAccountOrderFillLocked indicates order fill is locked for the Account.
	OrderCancelReasonAccountOrderFillLocked OrderCancelReason = "ACCOUNT_ORDER_FILL_LOCKED"
	// OrderCancelReasonClientRequest indicates the Order was cancelled due to a client request.
	OrderCancelReasonClientRequest OrderCancelReason = "CLIENT_REQUEST"
	// OrderCancelReasonMigration indicates the Order was cancelled due to migration.
	OrderCancelReasonMigration OrderCancelReason = "MIGRATION"
	// OrderCancelReasonMarketHalted indicates the Order was cancelled because the market was halted.
	OrderCancelReasonMarketHalted OrderCancelReason = "MARKET_HALTED"
	// OrderCancelReasonLinkedTradeClosed indicates the Order was cancelled because the linked Trade was closed.
	OrderCancelReasonLinkedTradeClosed OrderCancelReason = "LINKED_TRADE_CLOSED"
	// OrderCancelReasonTimeInForceExpired indicates the Order's time in force expired.
	OrderCancelReasonTimeInForceExpired OrderCancelReason = "TIME_IN_FORCE_EXPIRED"
	// OrderCancelReasonInsufficientMargin indicates there was insufficient margin.
	OrderCancelReasonInsufficientMargin OrderCancelReason = "INSUFFICIENT_MARGIN"
	// OrderCancelReasonFifoViolation indicates a FIFO violation cancelled the Order.
	OrderCancelReasonFifoViolation OrderCancelReason = "FIFO_VIOLATION"
	// OrderCancelReasonBoundsViolation indicates a price bounds violation cancelled the Order.
	OrderCancelReasonBoundsViolation OrderCancelReason = "BOUNDS_VIOLATION"
	// OrderCancelReasonClientRequestReplaced indicates the Order was cancelled and replaced by client request.
	OrderCancelReasonClientRequestReplaced OrderCancelReason = "CLIENT_REQUEST_REPLACED"
	// OrderCancelReasonInsufficientLiquidity indicates there was insufficient liquidity.
	OrderCancelReasonInsufficientLiquidity OrderCancelReason = "INSUFFICIENT_LIQUIDITY"
	// OrderCancelReasonTakeProfitOnFillGtdTimestampInPast indicates the Take Profit on fill GTD timestamp is in the past.
	OrderCancelReasonTakeProfitOnFillGtdTimestampInPast OrderCancelReason = "TAKE_PROFIT_ON_FILL_GTD_TIMESTAMP_IN_PAST"
	// OrderCancelReasonTakeProfitOnFillLoss indicates the Take Profit on fill would result in a loss.
	OrderCancelReasonTakeProfitOnFillLoss OrderCancelReason = "TAKE_PROFIT_ON_FILL_LOSS"
	// OrderCancelReasonLosingTakeProfit indicates losing Take Profit.
	OrderCancelReasonLosingTakeProfit OrderCancelReason = "LOSING_TAKE_PROFIT"
	// OrderCancelReasonStopLossOnFillGtdTimestampInPast indicates the Stop Loss on fill GTD timestamp is in the past.
	OrderCancelReasonStopLossOnFillGtdTimestampInPast OrderCancelReason = "STOP_LOSS_ON_FILL_GTD_TIMESTAMP_IN_PAST"
	// OrderCancelReasonStopLossOnFillLoss indicates the Stop Loss on fill would result in a loss.
	OrderCancelReasonStopLossOnFillLoss OrderCancelReason = "STOP_LOSS_ON_FILL_LOSS"
	// OrderCancelReasonStopLossOnFillPriceDistanceMaximumExceeded indicates the Stop Loss price distance maximum was exceeded.
	OrderCancelReasonStopLossOnFillPriceDistanceMaximumExceeded OrderCancelReason = "STOP_LOSS_ON_FILL_PRICE_DISTANCE_MAXIMUM_EXCEEDED"
	// OrderCancelReasonStopLossOnFillRequired indicates a Stop Loss on fill is required.
	OrderCancelReasonStopLossOnFillRequired OrderCancelReason = "STOP_LOSS_ON_FILL_REQUIRED"
	// OrderCancelReasonStopLossOnFillGuaranteedRequired indicates a guaranteed Stop Loss on fill is required.
	OrderCancelReasonStopLossOnFillGuaranteedRequired OrderCancelReason = "STOP_LOSS_ON_FILL_GUARANTEED_REQUIRED"
	// OrderCancelReasonStopLossOnFillGuaranteedNotAllowed indicates a guaranteed Stop Loss on fill is not allowed.
	OrderCancelReasonStopLossOnFillGuaranteedNotAllowed OrderCancelReason = "STOP_LOSS_ON_FILL_GUARANTEED_NOT_ALLOWED"
	// OrderCancelReasonStopLossOnFillGuaranteedMinimumDistanceNotMet indicates the guaranteed Stop Loss minimum distance was not met.
	OrderCancelReasonStopLossOnFillGuaranteedMinimumDistanceNotMet OrderCancelReason = "STOP_LOSS_ON_FILL_GUARANTEED_MINIMUM_DISTANCE_NOT_MET"
	// OrderCancelReasonStopLossOnFillGuaranteedLevelRestrictionExceeded indicates the guaranteed Stop Loss level restriction was exceeded.
	OrderCancelReasonStopLossOnFillGuaranteedLevelRestrictionExceeded OrderCancelReason = "STOP_LOSS_ON_FILL_GUARANTEED_LEVEL_RESTRICTION_EXCEEDED"
	// OrderCancelReasonStopLossOnFillGuaranteedHedgingNotAllowed indicates guaranteed Stop Loss on fill hedging is not allowed.
	OrderCancelReasonStopLossOnFillGuaranteedHedgingNotAllowed OrderCancelReason = "STOP_LOSS_ON_FILL_GUARANTEED_HEDGING_NOT_ALLOWED"
	// OrderCancelReasonStopLossOnFillTimeInForceInvalid indicates invalid time in force for Stop Loss on fill.
	OrderCancelReasonStopLossOnFillTimeInForceInvalid OrderCancelReason = "STOP_LOSS_ON_FILL_TIME_IN_FORCE_INVALID"
	// OrderCancelReasonStopLossOnFillTriggerConditionInvalid indicates invalid trigger condition for Stop Loss on fill.
	OrderCancelReasonStopLossOnFillTriggerConditionInvalid OrderCancelReason = "STOP_LOSS_ON_FILL_TRIGGER_CONDITION_INVALID"
	// OrderCancelReasonTakeProfitOnFillPriceDistanceMaximumExceeded indicates the Take Profit price distance maximum was exceeded.
	OrderCancelReasonTakeProfitOnFillPriceDistanceMaximumExceeded OrderCancelReason = "TAKE_PROFIT_ON_FILL_PRICE_DISTANCE_MAXIMUM_EXCEEDED"
	// OrderCancelReasonTrailingStopLossOnFillGtdTimestampInPast indicates the Trailing Stop Loss on fill GTD timestamp is in the past.
	OrderCancelReasonTrailingStopLossOnFillGtdTimestampInPast OrderCancelReason = "TRAILING_STOP_LOSS_ON_FILL_GTD_TIMESTAMP_IN_PAST"
	// OrderCancelReasonClientTradeIdAlreadyExists indicates the client Trade ID already exists.
	OrderCancelReasonClientTradeIdAlreadyExists OrderCancelReason = "CLIENT_TRADE_ID_ALREADY_EXISTS"
	// OrderCancelReasonPositionCloseoutFailed indicates the position closeout failed.
	OrderCancelReasonPositionCloseoutFailed OrderCancelReason = "POSITION_CLOSEOUT_FAILED"
	// OrderCancelReasonOpenTradesAllowedExceeded indicates the open Trades allowed limit was exceeded.
	OrderCancelReasonOpenTradesAllowedExceeded OrderCancelReason = "OPEN_TRADES_ALLOWED_EXCEEDED"
	// OrderCancelReasonPendingOrdersAllowedExceeded indicates the pending Orders allowed limit was exceeded.
	OrderCancelReasonPendingOrdersAllowedExceeded OrderCancelReason = "PENDING_ORDERS_ALLOWED_EXCEEDED"
	// OrderCancelReasonTakeProfitOnFillClientOrderIdAlreadyExists indicates the Take Profit on fill client Order ID already exists.
	OrderCancelReasonTakeProfitOnFillClientOrderIdAlreadyExists OrderCancelReason = "TAKE_PROFIT_ON_FILL_CLIENT_ORDER_ID_ALREADY_EXISTS"
	// OrderCancelReasonStopLossOnFillClientOrderIdAlreadyExists indicates the Stop Loss on fill client Order ID already exists.
	OrderCancelReasonStopLossOnFillClientOrderIdAlreadyExists OrderCancelReason = "STOP_LOSS_ON_FILL_CLIENT_ORDER_ID_ALREADY_EXISTS"
	// OrderCancelReasonTrailingStopLossOnFillClientOrderIdAlreadyExists indicates the Trailing Stop Loss on fill client Order ID already exists.
	OrderCancelReasonTrailingStopLossOnFillClientOrderIdAlreadyExists OrderCancelReason = "TRAILING_STOP_LOSS_ON_FILL_CLIENT_ORDER_ID_ALREADY_EXISTS"
	// OrderCancelReasonPositionSizeExceeded indicates the position size limit was exceeded.
	OrderCancelReasonPositionSizeExceeded OrderCancelReason = "POSITION_SIZE_EXCEEDED"
	// OrderCancelReasonHedgingGsloViolation indicates hedging with guaranteed Stop Loss Orders is not allowed.
	OrderCancelReasonHedgingGsloViolation OrderCancelReason = "HEDGING_GSLO_VIOLATION"
	// OrderCancelReasonAccountPositionValueLimitExceeded indicates the Account position value limit was exceeded.
	OrderCancelReasonAccountPositionValueLimitExceeded OrderCancelReason = "ACCOUNT_POSITION_VALUE_LIMIT_EXCEEDED"
	// OrderCancelReasonInstrumentBidReduceOnly indicates the Instrument is bid reduce only.
	OrderCancelReasonInstrumentBidReduceOnly OrderCancelReason = "INSTRUMENT_BID_REDUCE_ONLY"
	// OrderCancelReasonInstrumentAskReduceOnly indicates the Instrument is ask reduce only.
	OrderCancelReasonInstrumentAskReduceOnly OrderCancelReason = "INSTRUMENT_ASK_REDUCE_ONLY"
	// OrderCancelReasonInstrumentBidHalted indicates the Instrument bid is halted.
	OrderCancelReasonInstrumentBidHalted OrderCancelReason = "INSTRUMENT_BID_HALTED"
	// OrderCancelReasonInstrumentAskHalted indicates the Instrument ask is halted.
	OrderCancelReasonInstrumentAskHalted OrderCancelReason = "INSTRUMENT_ASK_HALTED"
	// OrderCancelReasonStopLossOnFillGuaranteedBidHalted indicates the guaranteed Stop Loss bid is halted.
	OrderCancelReasonStopLossOnFillGuaranteedBidHalted OrderCancelReason = "STOP_LOSS_ON_FILL_GUARANTEED_BID_HALTED"
	// OrderCancelReasonStopLossOnFillGuaranteedAskHalted indicates the guaranteed Stop Loss ask is halted.
	OrderCancelReasonStopLossOnFillGuaranteedAskHalted OrderCancelReason = "STOP_LOSS_ON_FILL_GUARANTEED_ASK_HALTED"
	// OrderCancelReasonGuaranteedStopLossOnFillBidHalted indicates the Guaranteed Stop Loss on fill bid is halted.
	OrderCancelReasonGuaranteedStopLossOnFillBidHalted OrderCancelReason = "GUARANTEED_STOP_LOSS_ON_FILL_BID_HALTED"
	// OrderCancelReasonGuaranteedStopLossOnFillAskHalted indicates the Guaranteed Stop Loss on fill ask is halted.
	OrderCancelReasonGuaranteedStopLossOnFillAskHalted OrderCancelReason = "GUARANTEED_STOP_LOSS_ON_FILL_ASK_HALTED"
)

// OpenTradeDividendAdjustment contains the dividend adjustment information for an open Trade.
type OpenTradeDividendAdjustment struct {
	// TradeID is the ID of the Trade for which the dividend adjustment is calculated.
	TradeID TradeID `json:"tradeID"`
	// DividendAdjustment is the dividend adjustment applied to the Trade.
	DividendAdjustment AccountUnits `json:"dividendAdjustment"`
	// QuoteDividendAdjustment is the dividend adjustment in the quote currency.
	QuoteDividendAdjustment DecimalNumber `json:"quoteDividendAdjustment"`
}

// ClientID is a client-provided identifier, used by clients to refer to their Orders or Trades
// with an identifier that they have provided.
type ClientID string

// ClientTag is a client-provided tag that can be associated with an Order or Trade. Tags are
// typically used for organization and filtering.
type ClientTag string

// ClientComment is a client-provided comment that can be associated with an Order or Trade.
type ClientComment string

// ClientExtensions represents the client-configurable portions of an Order, Trade, or client-related
// elements in the OANDA platform. Do not set, modify, or delete clientExtensions if your account
// is associated with MT4.
type ClientExtensions struct {
	// ID is a client-provided identifier, used by clients to refer to their Orders or Trades
	// with an identifier that they have provided.
	ID ClientID `json:"id,omitempty"`
	// Tag is a client-provided tag that can be associated with an Order or Trade.
	Tag ClientTag `json:"tag,omitempty"`
	// Comment is a client-provided comment that can be associated with an Order or Trade.
	Comment ClientComment `json:"comment,omitempty"`
}

func NewClientExtensions(id ClientID, tag ClientTag, comment ClientComment) *ClientExtensions {
	return &ClientExtensions{
		ID:      id,
		Tag:     tag,
		Comment: comment,
	}
}

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
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// ClientExtensions are the client extensions to add to the Take Profit Order when created.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
}

func NewTakeProfitDetails(price PriceValue) *TakeProfitDetails {
	return &TakeProfitDetails{
		Price:       price,
		TimeInForce: TimeInForceGTC,
	}
}

func (d *TakeProfitDetails) SetGTD(date DateTime) *TakeProfitDetails {
	d.TimeInForce = TimeInForceGTD
	d.GtdTime = &date
	return d
}

func (d *TakeProfitDetails) SetGFD() *TakeProfitDetails {
	d.TimeInForce = TimeInForceGFD
	return d
}

func (d *TakeProfitDetails) SetClientExtensions(clientExtensions *ClientExtensions) *TakeProfitDetails {
	d.ClientExtensions = clientExtensions
	return d
}

// StopLossDetails specifies the details of a Stop Loss Order to be created on behalf of a client.
// This may happen when an Order is filled that opens a Trade requiring a Stop Loss.
type StopLossDetails struct {
	// Price is the price threshold specified for the Stop Loss Order. The associated Trade will be
	// closed by a market price that is equal to or worse than this threshold.
	Price *PriceValue `json:"price,omitempty"`
	// Distance specifies the distance (in price units) from the Trade's open price to use as the
	// Stop Loss Order price. If the Trade is long, the Stop Loss price will be the open price minus
	// the distance. If the Trade is short, the Stop Loss price will be the open price plus the distance.
	Distance *DecimalNumber `json:"distance,omitempty"`
	// TimeInForce specifies how long the Stop Loss Order should remain pending before being
	// automatically cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Stop Loss Order will be cancelled if its timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// ClientExtensions are the client extensions to add to the Stop Loss Order when created.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
	// Guaranteed is deprecated. Flag indicating that the Stop Loss Order is guaranteed. The default
	// value depends on the GuaranteedStopLossOrderMode of the account.
	Guaranteed bool `json:"guaranteed"`
}

func NewStopLossDetails() *StopLossDetails {
	return &StopLossDetails{
		TimeInForce: TimeInForceGTC,
	}
}

func (d *StopLossDetails) SetPrice(price PriceValue) *StopLossDetails {
	d.Price = &price
	return d
}

func (d *StopLossDetails) SetDistance(distance DecimalNumber) *StopLossDetails {
	d.Distance = &distance
	return d
}

func (d *StopLossDetails) SetGTD(date DateTime) *StopLossDetails {
	d.TimeInForce = TimeInForceGTD
	d.GtdTime = &date
	return d
}

func (d *StopLossDetails) SetGFD() *StopLossDetails {
	d.TimeInForce = TimeInForceGFD
	return d
}

func (d *StopLossDetails) SetClientExtensions(clientExtensions *ClientExtensions) *StopLossDetails {
	d.ClientExtensions = clientExtensions
	return d
}

// GuaranteedStopLossDetails specifies the details of a Guaranteed Stop Loss Order to be created on
// behalf of a client. This may happen when an Order is filled that opens a Trade requiring a
// Guaranteed Stop Loss.
type GuaranteedStopLossDetails struct {
	// Price is the price threshold specified for the Guaranteed Stop Loss Order. The associated Trade
	// will be closed at this price.
	Price *PriceValue `json:"price,omitempty"`
	// Distance specifies the distance (in price units) from the Trade's open price to use as the
	// Guaranteed Stop Loss Order price. If the Trade is long, the order price will be the open price
	// minus the distance. If the Trade is short, the order price will be the open price plus the distance.
	Distance *DecimalNumber `json:"distance,omitempty"`
	// TimeInForce specifies how long the Guaranteed Stop Loss Order should remain pending before
	// being automatically cancelled by the execution system.
	TimeInForce TimeInForce `json:"timeInForce"`
	// GtdTime is the date/time when the Guaranteed Stop Loss Order will be cancelled if its
	// timeInForce is "GTD".
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// ClientExtensions are the client extensions to add to the Guaranteed Stop Loss Order when created.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
}

func NewGuaranteedStopLossDetails() *GuaranteedStopLossDetails {
	return &GuaranteedStopLossDetails{
		TimeInForce: TimeInForceGTC,
	}
}

func (d *GuaranteedStopLossDetails) SetPrice(price PriceValue) *GuaranteedStopLossDetails {
	d.Price = &price
	return d
}

func (d *GuaranteedStopLossDetails) SetDistance(distance DecimalNumber) *GuaranteedStopLossDetails {
	d.Distance = &distance
	return d
}

func (d *GuaranteedStopLossDetails) SetGTD(date DateTime) *GuaranteedStopLossDetails {
	d.TimeInForce = TimeInForceGTD
	d.GtdTime = &date
	return d
}

func (d *GuaranteedStopLossDetails) SetGFD() *GuaranteedStopLossDetails {
	d.TimeInForce = TimeInForceGFD
	return d
}

func (d *GuaranteedStopLossDetails) SetClientExtensions(clientExtensions *ClientExtensions) *GuaranteedStopLossDetails {
	d.ClientExtensions = clientExtensions
	return d
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
	GtdTime *DateTime `json:"gtdTime,omitempty"`
	// ClientExtensions are the client extensions to add to the Trailing Stop Loss Order when created.
	ClientExtensions *ClientExtensions `json:"clientExtensions,omitempty"`
}

func NewTrailingStopLossDetails(distance DecimalNumber) *TrailingStopLossDetails {
	return &TrailingStopLossDetails{
		Distance:    distance,
		TimeInForce: TimeInForceGTC,
	}
}

func (d *TrailingStopLossDetails) SetGTD(date DateTime) *TrailingStopLossDetails {
	d.TimeInForce = TimeInForceGTD
	d.GtdTime = &date
	return d
}

func (d *TrailingStopLossDetails) SetGFD() *TrailingStopLossDetails {
	d.TimeInForce = TimeInForceGFD
	return d
}

func (d *TrailingStopLossDetails) SetClientExtensions(clientExtensions *ClientExtensions) *TrailingStopLossDetails {
	d.ClientExtensions = clientExtensions
	return d
}

// TradeOpen contains the details of a Trade opened by an OrderFillTransaction.
type TradeOpen struct {
	// TradeID is the ID of the Trade that was opened.
	TradeID TradeID `json:"tradeID"`
	// Units is the number of units opened by the Trade.
	Units DecimalNumber `json:"units"`
	// Price is the average price that the units were opened at.
	Price PriceValue `json:"price"`
	// GuaranteedExecutionFee is the fee charged for the Trade.
	GuaranteedExecutionFee AccountUnits `json:"guaranteedExecutionFee"`
	// QuoteGuaranteedExecutionFee is the fee in quote currency.
	QuoteGuaranteedExecutionFee DecimalNumber `json:"quoteGuaranteedExecutionFee"`
	// ClientExtensions are the client extensions for the Trade.
	ClientExtensions ClientExtensions `json:"clientExtensions"`
	// HalfSpreadCost is the half spread cost for the Trade.
	HalfSpreadCost AccountUnits `json:"halfSpreadCost"`
	// InitialMarginRequired is the margin required at the time the Trade was created.
	InitialMarginRequired AccountUnits `json:"initialMarginRequired"`
}

// TradeReduce contains the details of a Trade reduced by an OrderFillTransaction.
type TradeReduce struct {
	// TradeID is the ID of the Trade that was reduced.
	TradeID TradeID `json:"tradeID"`
	// Units is the number of units that the Trade was reduced by.
	Units DecimalNumber `json:"units"`
	// Price is the average price that the units were closed at.
	Price PriceValue `json:"price"`
	// RealizedPL is the PL realized when reducing the Trade.
	RealizedPL AccountUnits `json:"realizedPL"`
	// Financing is the financing paid/collected when reducing the Trade.
	Financing AccountUnits `json:"financing"`
	// BaseFinancing is the financing in the base currency.
	BaseFinancing DecimalNumber `json:"baseFinancing"`
	// QuoteFinancing is the financing in the quote currency.
	QuoteFinancing DecimalNumber `json:"quoteFinancing"`
	// FinancingRate is the financing rate in effect for the Instrument.
	FinancingRate DecimalNumber `json:"financingRate"`
	// GuaranteedExecutionFee is the fee charged for closing the Trade.
	GuaranteedExecutionFee AccountUnits `json:"guaranteedExecutionFee"`
	// QuoteGuaranteedExecutionFee is the fee in quote currency.
	QuoteGuaranteedExecutionFee DecimalNumber `json:"quoteGuaranteedExecutionFee"`
	// HalfSpreadCost is the half spread cost for the Trade.
	HalfSpreadCost AccountUnits `json:"halfSpreadCost"`
}

// MarketOrderTradeClose specifies the extensions to a Market Order that has been created specifically
// to close a Trade.
type MarketOrderTradeClose struct {
	// TradeID is the ID of the Trade requested to be closed.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade requested to be closed.
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// Units indicates the number of units of the Trade to close. If not specified, all units of the
	// Trade will be closed.
	Units DecimalNumber `json:"units"`
}

// MarketOrderMarginCloseout details the reason that a Market Order was created as part of a Margin Closeout.
type MarketOrderMarginCloseout struct {
	// Reason is the reason the Market Order was created to perform a margin closeout.
	Reason string `json:"reason"`
}

// MarketOrderMarginCloseoutReason represents the reason that a Market Order was created for margin closeout.
type MarketOrderMarginCloseoutReason string

const (
	// MarketOrderMarginCloseoutReasonMarginCheckViolation indicates the Trade was closed due to margin check violation.
	MarketOrderMarginCloseoutReasonMarginCheckViolation MarketOrderMarginCloseoutReason = "MARGIN_CHECK_VIOLATION"
	// MarketOrderMarginCloseoutReasonRegulatoryMarginCallViolation indicates the Trade was closed due to regulatory margin call violation.
	MarketOrderMarginCloseoutReasonRegulatoryMarginCallViolation MarketOrderMarginCloseoutReason = "REGULATORY_MARGIN_CALL_VIOLATION"
	// MarketOrderMarginCloseoutReasonRegulatoryMarginCheckViolation indicates the Trade was closed due to regulatory margin check violation.
	MarketOrderMarginCloseoutReasonRegulatoryMarginCheckViolation MarketOrderMarginCloseoutReason = "REGULATORY_MARGIN_CHECK_VIOLATION"
)

// MarketOrderDelayedTradeClose details the reason that a Market Order was created for a delayed Trade close.
type MarketOrderDelayedTradeClose struct {
	// TradeID is the ID of the Trade being closed.
	TradeID TradeID `json:"tradeID"`
	// ClientTradeID is the client ID of the Trade being closed.
	ClientTradeID *ClientID `json:"clientTradeID,omitempty"`
	// SourceTransactionID is the Transaction ID of the DelayedTradeClosure transaction to which this
	// Delayed Trade Close belongs to.
	SourceTransactionID TransactionID `json:"sourceTransactionID"`
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

// LiquidityRegenerationSchedule indicates how liquidity that is used when filling an Order
// is regenerated following the fill.
type LiquidityRegenerationSchedule struct {
	// Steps are the steps in the liquidity regeneration schedule.
	Steps []LiquidityRegenerationScheduleStep `json:"steps"`
}

// LiquidityRegenerationScheduleStep indicates the amount of bid/ask liquidity that is used
// by the Account at a certain time.
type LiquidityRegenerationScheduleStep struct {
	// Timestamp is the timestamp of the schedule step.
	Timestamp DateTime `json:"timestamp"`
	// BidLiquidityUsed is the amount of bid liquidity used at this step.
	BidLiquidityUsed DecimalNumber `json:"bidLiquidityUsed"`
	// AskLiquidityUsed is the amount of ask liquidity used at this step.
	AskLiquidityUsed DecimalNumber `json:"askLiquidityUsed"`
}

// OpenTradeFinancing contains the financing information for an open Trade.
type OpenTradeFinancing struct {
	// TradeID is the ID of the Trade that financing is being paid/collected for.
	TradeID TradeID `json:"tradeID"`
	// Financing is the amount of financing paid/collected for the Trade.
	Financing AccountUnits `json:"financing"`
	// BaseFinancing is the financing in the base currency.
	BaseFinancing DecimalNumber `json:"baseFinancing"`
	// QuoteFinancing is the financing in the quote currency.
	QuoteFinancing DecimalNumber `json:"quoteFinancing"`
	// FinancingRate is the financing rate in effect for the Instrument.
	FinancingRate DecimalNumber `json:"financingRate"`
}

// PositionFinancing contains the financing information for a Position.
type PositionFinancing struct {
	// Instrument is the Instrument of the Position.
	Instrument InstrumentName `json:"instrument"`
	// Financing is the amount of financing paid/collected for the Position.
	Financing AccountUnits `json:"financing"`
	// BaseFinancing is the financing in the base currency.
	BaseFinancing DecimalNumber `json:"baseFinancing"`
	// QuoteFinancing is the financing in the quote currency.
	QuoteFinancing DecimalNumber `json:"quoteFinancing"`
	// HomeConversionFactors is the HomeConversionFactors in effect at the time.
	HomeConversionFactors HomeConversionFactors `json:"homeConversionFactors"`
	// OpenTradeFinancings are the financing paid/collected for each open Trade.
	OpenTradeFinancings []OpenTradeFinancing `json:"openTradeFinancings"`
	// AccountFinancingMode describes how financing charges are calculated.
	AccountFinancingMode AccountFinancingMode `json:"accountFinancingMode"`
}

// RequestID is the unique identifier for a client request.
type RequestID string

// ClientRequestID is a client-provided request identifier.
type ClientRequestID = string

// TransactionRejectReason represents the reason that a Transaction was rejected.
type TransactionRejectReason string

const (
	// TransactionRejectReasonInternalServerError indicates an internal server error.
	TransactionRejectReasonInternalServerError TransactionRejectReason = "INTERNAL_SERVER_ERROR"
	// TransactionRejectReasonInstrumentPriceUnknown indicates the price for the Instrument is unknown.
	TransactionRejectReasonInstrumentPriceUnknown TransactionRejectReason = "INSTRUMENT_PRICE_UNKNOWN"
	// TransactionRejectReasonAccountNotActive indicates the Account is not active.
	TransactionRejectReasonAccountNotActive TransactionRejectReason = "ACCOUNT_NOT_ACTIVE"
	// TransactionRejectReasonAccountLocked indicates the Account is locked.
	TransactionRejectReasonAccountLocked TransactionRejectReason = "ACCOUNT_LOCKED"
	// TransactionRejectReasonAccountOrderCreationLocked indicates order creation is locked for the Account.
	TransactionRejectReasonAccountOrderCreationLocked TransactionRejectReason = "ACCOUNT_ORDER_CREATION_LOCKED"
	// TransactionRejectReasonAccountConfigurationLocked indicates configuration is locked for the Account.
	TransactionRejectReasonAccountConfigurationLocked TransactionRejectReason = "ACCOUNT_CONFIGURATION_LOCKED"
	// TransactionRejectReasonAccountDepositLocked indicates deposits are locked for the Account.
	TransactionRejectReasonAccountDepositLocked TransactionRejectReason = "ACCOUNT_DEPOSIT_LOCKED"
	// TransactionRejectReasonAccountWithdrawalLocked indicates withdrawals are locked for the Account.
	TransactionRejectReasonAccountWithdrawalLocked TransactionRejectReason = "ACCOUNT_WITHDRAWAL_LOCKED"
	// TransactionRejectReasonAccountOrderCancelLocked indicates order cancellation is locked for the Account.
	TransactionRejectReasonAccountOrderCancelLocked TransactionRejectReason = "ACCOUNT_ORDER_CANCEL_LOCKED"
	// TransactionRejectReasonInstrumentNotTradeable indicates the Instrument is not tradeable.
	TransactionRejectReasonInstrumentNotTradeable TransactionRejectReason = "INSTRUMENT_NOT_TRADEABLE"
	// TransactionRejectReasonPendingOrdersAllowedExceeded indicates too many pending Orders.
	TransactionRejectReasonPendingOrdersAllowedExceeded TransactionRejectReason = "PENDING_ORDERS_ALLOWED_EXCEEDED"
	// TransactionRejectReasonOrderIdUnspecified indicates the Order ID was not specified.
	TransactionRejectReasonOrderIdUnspecified TransactionRejectReason = "ORDER_ID_UNSPECIFIED"
	// TransactionRejectReasonOrderDoesntExist indicates the Order does not exist.
	TransactionRejectReasonOrderDoesntExist TransactionRejectReason = "ORDER_DOESNT_EXIST"
	// TransactionRejectReasonOrderIdentifierInconsistency indicates an Order identifier inconsistency.
	TransactionRejectReasonOrderIdentifierInconsistency TransactionRejectReason = "ORDER_IDENTIFIER_INCONSISTENCY"
	// TransactionRejectReasonTradeIdUnspecified indicates the Trade ID was not specified.
	TransactionRejectReasonTradeIdUnspecified TransactionRejectReason = "TRADE_ID_UNSPECIFIED"
	// TransactionRejectReasonTradeDoesntExist indicates the Trade does not exist.
	TransactionRejectReasonTradeDoesntExist TransactionRejectReason = "TRADE_DOESNT_EXIST"
	// TransactionRejectReasonTradeIdentifierInconsistency indicates a Trade identifier inconsistency.
	TransactionRejectReasonTradeIdentifierInconsistency TransactionRejectReason = "TRADE_IDENTIFIER_INCONSISTENCY"
	// TransactionRejectReasonInsufficientMargin indicates insufficient margin.
	TransactionRejectReasonInsufficientMargin TransactionRejectReason = "INSUFFICIENT_MARGIN"
	// TransactionRejectReasonInstrumentMissing indicates the Instrument is missing.
	TransactionRejectReasonInstrumentMissing TransactionRejectReason = "INSTRUMENT_MISSING"
	// TransactionRejectReasonInstrumentUnknown indicates the Instrument is unknown.
	TransactionRejectReasonInstrumentUnknown TransactionRejectReason = "INSTRUMENT_UNKNOWN"
	// TransactionRejectReasonUnitsMissing indicates the units are missing.
	TransactionRejectReasonUnitsMissing TransactionRejectReason = "UNITS_MISSING"
	// TransactionRejectReasonUnitsInvalid indicates the units are invalid.
	TransactionRejectReasonUnitsInvalid TransactionRejectReason = "UNITS_INVALID"
	// TransactionRejectReasonUnitsPrecisionExceeded indicates the units precision was exceeded.
	TransactionRejectReasonUnitsPrecisionExceeded TransactionRejectReason = "UNITS_PRECISION_EXCEEDED"
	// TransactionRejectReasonUnitsLimitExceeded indicates the units limit was exceeded.
	TransactionRejectReasonUnitsLimitExceeded TransactionRejectReason = "UNITS_LIMIT_EXCEEDED"
	// TransactionRejectReasonUnitsMinimumNotMet indicates the minimum units were not met.
	TransactionRejectReasonUnitsMinimumNotMet TransactionRejectReason = "UNITS_MINIMUM_NOT_MET"
	// TransactionRejectReasonPriceInvalid indicates the price is invalid.
	TransactionRejectReasonPriceInvalid TransactionRejectReason = "PRICE_INVALID"
	// TransactionRejectReasonPricePrecisionExceeded indicates the price precision was exceeded.
	TransactionRejectReasonPricePrecisionExceeded TransactionRejectReason = "PRICE_PRECISION_EXCEEDED"
	// TransactionRejectReasonPriceDistanceMissing indicates the price distance is missing.
	TransactionRejectReasonPriceDistanceMissing TransactionRejectReason = "PRICE_DISTANCE_MISSING"
	// TransactionRejectReasonPriceDistanceInvalid indicates the price distance is invalid.
	TransactionRejectReasonPriceDistanceInvalid TransactionRejectReason = "PRICE_DISTANCE_INVALID"
	// TransactionRejectReasonPriceDistancePrecisionExceeded indicates the price distance precision was exceeded.
	TransactionRejectReasonPriceDistancePrecisionExceeded TransactionRejectReason = "PRICE_DISTANCE_PRECISION_EXCEEDED"
	// TransactionRejectReasonPriceDistanceMaximumExceeded indicates the price distance maximum was exceeded.
	TransactionRejectReasonPriceDistanceMaximumExceeded TransactionRejectReason = "PRICE_DISTANCE_MAXIMUM_EXCEEDED"
	// TransactionRejectReasonPriceDistanceMinimumNotMet indicates the minimum price distance was not met.
	TransactionRejectReasonPriceDistanceMinimumNotMet TransactionRejectReason = "PRICE_DISTANCE_MINIMUM_NOT_MET"
	// TransactionRejectReasonTimeInForceMissing indicates the time in force is missing.
	TransactionRejectReasonTimeInForceMissing TransactionRejectReason = "TIME_IN_FORCE_MISSING"
	// TransactionRejectReasonTimeInForceInvalid indicates the time in force is invalid.
	TransactionRejectReasonTimeInForceInvalid TransactionRejectReason = "TIME_IN_FORCE_INVALID"
	// TransactionRejectReasonTimeInForceGtdTimestampMissing indicates the GTD timestamp is missing.
	TransactionRejectReasonTimeInForceGtdTimestampMissing TransactionRejectReason = "TIME_IN_FORCE_GTD_TIMESTAMP_MISSING"
	// TransactionRejectReasonTimeInForceGtdTimestampInPast indicates the GTD timestamp is in the past.
	TransactionRejectReasonTimeInForceGtdTimestampInPast TransactionRejectReason = "TIME_IN_FORCE_GTD_TIMESTAMP_IN_PAST"
	// TransactionRejectReasonPriceBoundInvalid indicates the price bound is invalid.
	TransactionRejectReasonPriceBoundInvalid TransactionRejectReason = "PRICE_BOUND_INVALID"
	// TransactionRejectReasonPriceBoundPrecisionExceeded indicates the price bound precision was exceeded.
	TransactionRejectReasonPriceBoundPrecisionExceeded TransactionRejectReason = "PRICE_BOUND_PRECISION_EXCEEDED"
	// TransactionRejectReasonOrdersOnFillDuplicateClientOrderIds indicates duplicate client Order IDs on fill.
	TransactionRejectReasonOrdersOnFillDuplicateClientOrderIds TransactionRejectReason = "ORDERS_ON_FILL_DUPLICATE_CLIENT_ORDER_IDS"
	// TransactionRejectReasonTradeOnFillClientExtensionsNotSupported indicates Trade on fill client extensions are not supported.
	TransactionRejectReasonTradeOnFillClientExtensionsNotSupported TransactionRejectReason = "TRADE_ON_FILL_CLIENT_EXTENSIONS_NOT_SUPPORTED"
	// TransactionRejectReasonClientOrderIdInvalid indicates the client Order ID is invalid.
	TransactionRejectReasonClientOrderIdInvalid TransactionRejectReason = "CLIENT_ORDER_ID_INVALID"
	// TransactionRejectReasonClientOrderIdAlreadyExists indicates the client Order ID already exists.
	TransactionRejectReasonClientOrderIdAlreadyExists TransactionRejectReason = "CLIENT_ORDER_ID_ALREADY_EXISTS"
	// TransactionRejectReasonClientOrderTagInvalid indicates the client Order tag is invalid.
	TransactionRejectReasonClientOrderTagInvalid TransactionRejectReason = "CLIENT_ORDER_TAG_INVALID"
	// TransactionRejectReasonClientOrderCommentInvalid indicates the client Order comment is invalid.
	TransactionRejectReasonClientOrderCommentInvalid TransactionRejectReason = "CLIENT_ORDER_COMMENT_INVALID"
	// TransactionRejectReasonClientTradeIdInvalid indicates the client Trade ID is invalid.
	TransactionRejectReasonClientTradeIdInvalid TransactionRejectReason = "CLIENT_TRADE_ID_INVALID"
	// TransactionRejectReasonClientTradeIdAlreadyExists indicates the client Trade ID already exists.
	TransactionRejectReasonClientTradeIdAlreadyExists TransactionRejectReason = "CLIENT_TRADE_ID_ALREADY_EXISTS"
	// TransactionRejectReasonClientTradeTagInvalid indicates the client Trade tag is invalid.
	TransactionRejectReasonClientTradeTagInvalid TransactionRejectReason = "CLIENT_TRADE_TAG_INVALID"
	// TransactionRejectReasonClientTradeCommentInvalid indicates the client Trade comment is invalid.
	TransactionRejectReasonClientTradeCommentInvalid TransactionRejectReason = "CLIENT_TRADE_COMMENT_INVALID"
	// TransactionRejectReasonOrderFillPositionActionMissing indicates the Order fill position action is missing.
	TransactionRejectReasonOrderFillPositionActionMissing TransactionRejectReason = "ORDER_FILL_POSITION_ACTION_MISSING"
	// TransactionRejectReasonOrderFillPositionActionInvalid indicates the Order fill position action is invalid.
	TransactionRejectReasonOrderFillPositionActionInvalid TransactionRejectReason = "ORDER_FILL_POSITION_ACTION_INVALID"
	// TransactionRejectReasonTriggerConditionMissing indicates the trigger condition is missing.
	TransactionRejectReasonTriggerConditionMissing TransactionRejectReason = "TRIGGER_CONDITION_MISSING"
	// TransactionRejectReasonTriggerConditionInvalid indicates the trigger condition is invalid.
	TransactionRejectReasonTriggerConditionInvalid TransactionRejectReason = "TRIGGER_CONDITION_INVALID"
	// TransactionRejectReasonTakeProfitOrderAlreadyExists indicates a Take Profit Order already exists.
	TransactionRejectReasonTakeProfitOrderAlreadyExists TransactionRejectReason = "TAKE_PROFIT_ORDER_ALREADY_EXISTS"
	// TransactionRejectReasonStopLossOrderAlreadyExists indicates a Stop Loss Order already exists.
	TransactionRejectReasonStopLossOrderAlreadyExists TransactionRejectReason = "STOP_LOSS_ORDER_ALREADY_EXISTS"
	// TransactionRejectReasonGuaranteedStopLossOrderAlreadyExists indicates a Guaranteed Stop Loss Order already exists.
	TransactionRejectReasonGuaranteedStopLossOrderAlreadyExists TransactionRejectReason = "GUARANTEED_STOP_LOSS_ORDER_ALREADY_EXISTS"
	// TransactionRejectReasonTrailingStopLossOrderAlreadyExists indicates a Trailing Stop Loss Order already exists.
	TransactionRejectReasonTrailingStopLossOrderAlreadyExists TransactionRejectReason = "TRAILING_STOP_LOSS_ORDER_ALREADY_EXISTS"
	// TransactionRejectReasonCloseTradeTypeMissing indicates the close Trade type is missing.
	TransactionRejectReasonCloseTradeTypeMissing TransactionRejectReason = "CLOSE_TRADE_TYPE_MISSING"
	// TransactionRejectReasonCloseTradePartialUnitsMissing indicates the close Trade partial units are missing.
	TransactionRejectReasonCloseTradePartialUnitsMissing TransactionRejectReason = "CLOSE_TRADE_PARTIAL_UNITS_MISSING"
	// TransactionRejectReasonCloseTradeUnitsExceedTradeSize indicates the close Trade units exceed Trade size.
	TransactionRejectReasonCloseTradeUnitsExceedTradeSize TransactionRejectReason = "CLOSE_TRADE_UNITS_EXCEED_TRADE_SIZE"
	// TransactionRejectReasonCloseoutPositionDoesntExist indicates the closeout position does not exist.
	TransactionRejectReasonCloseoutPositionDoesntExist TransactionRejectReason = "CLOSEOUT_POSITION_DOESNT_EXIST"
	// TransactionRejectReasonCloseoutPositionIncompleteSpecification indicates incomplete closeout position specification.
	TransactionRejectReasonCloseoutPositionIncompleteSpecification TransactionRejectReason = "CLOSEOUT_POSITION_INCOMPLETE_SPECIFICATION"
	// TransactionRejectReasonCloseoutPositionUnitsExceedPositionSize indicates the closeout position units exceed position size.
	TransactionRejectReasonCloseoutPositionUnitsExceedPositionSize TransactionRejectReason = "CLOSEOUT_POSITION_UNITS_EXCEED_POSITION_SIZE"
	// TransactionRejectReasonCloseoutPositionReject indicates the closeout position was rejected.
	TransactionRejectReasonCloseoutPositionReject TransactionRejectReason = "CLOSEOUT_POSITION_REJECT"
	// TransactionRejectReasonCloseoutPositionPartialUnitsMissing indicates the closeout position partial units are missing.
	TransactionRejectReasonCloseoutPositionPartialUnitsMissing TransactionRejectReason = "CLOSEOUT_POSITION_PARTIAL_UNITS_MISSING"
	// TransactionRejectReasonMarkupGroupIdInvalid indicates the markup group ID is invalid.
	TransactionRejectReasonMarkupGroupIdInvalid TransactionRejectReason = "MARKUP_GROUP_ID_INVALID"
	// TransactionRejectReasonPositionAggregationModeInvalid indicates the position aggregation mode is invalid.
	TransactionRejectReasonPositionAggregationModeInvalid TransactionRejectReason = "POSITION_AGGREGATION_MODE_INVALID"
	// TransactionRejectReasonAdminConfigureDataMissing indicates the admin configure data is missing.
	TransactionRejectReasonAdminConfigureDataMissing TransactionRejectReason = "ADMIN_CONFIGURE_DATA_MISSING"
	// TransactionRejectReasonMarginRateInvalid indicates the margin rate is invalid.
	TransactionRejectReasonMarginRateInvalid TransactionRejectReason = "MARGIN_RATE_INVALID"
	// TransactionRejectReasonMarginRateWouldTriggerCloseout indicates the margin rate would trigger closeout.
	TransactionRejectReasonMarginRateWouldTriggerCloseout TransactionRejectReason = "MARGIN_RATE_WOULD_TRIGGER_CLOSEOUT"
	// TransactionRejectReasonAliasInvalid indicates the alias is invalid.
	TransactionRejectReasonAliasInvalid TransactionRejectReason = "ALIAS_INVALID"
	// TransactionRejectReasonClientConfigureDataMissing indicates the client configure data is missing.
	TransactionRejectReasonClientConfigureDataMissing TransactionRejectReason = "CLIENT_CONFIGURE_DATA_MISSING"
	// TransactionRejectReasonMarginRateWouldTriggerMarginCall indicates the margin rate would trigger margin call.
	TransactionRejectReasonMarginRateWouldTriggerMarginCall TransactionRejectReason = "MARGIN_RATE_WOULD_TRIGGER_MARGIN_CALL"
	// TransactionRejectReasonAmountInvalid indicates the amount is invalid.
	TransactionRejectReasonAmountInvalid TransactionRejectReason = "AMOUNT_INVALID"
	// TransactionRejectReasonInsufficientFunds indicates insufficient funds.
	TransactionRejectReasonInsufficientFunds TransactionRejectReason = "INSUFFICIENT_FUNDS"
	// TransactionRejectReasonAmountMissing indicates the amount is missing.
	TransactionRejectReasonAmountMissing TransactionRejectReason = "AMOUNT_MISSING"
	// TransactionRejectReasonFundingReasonMissing indicates the funding reason is missing.
	TransactionRejectReasonFundingReasonMissing TransactionRejectReason = "FUNDING_REASON_MISSING"
	// TransactionRejectReasonClientExtensionsDataMissing indicates the client extensions data is missing.
	TransactionRejectReasonClientExtensionsDataMissing TransactionRejectReason = "CLIENT_EXTENSIONS_DATA_MISSING"
	// TransactionRejectReasonReplacingOrderInvalid indicates the replacing Order is invalid.
	TransactionRejectReasonReplacingOrderInvalid TransactionRejectReason = "REPLACING_ORDER_INVALID"
	// TransactionRejectReasonReplacingTradeIdInvalid indicates the replacing Trade ID is invalid.
	TransactionRejectReasonReplacingTradeIdInvalid TransactionRejectReason = "REPLACING_TRADE_ID_INVALID"
	// TransactionRejectReasonOrderCannotBeReplaced indicates the Order cannot be replaced.
	TransactionRejectReasonOrderCannotBeReplaced TransactionRejectReason = "ORDER_CANNOT_BE_REPLACED"
	// TransactionRejectReasonOrderCannotBeCancelled indicates the Order cannot be cancelled.
	TransactionRejectReasonOrderCannotBeCancelled TransactionRejectReason = "ORDER_CANNOT_BE_CANCELLED"
)

// TransactionFilter represents the types of Transactions that can be filtered on.
type TransactionFilter string

const (
	// TransactionFilterOrder filters for Order-related Transactions.
	TransactionFilterOrder TransactionFilter = "ORDER"
	// TransactionFilterFunding filters for Funding-related Transactions.
	TransactionFilterFunding TransactionFilter = "FUNDING"
	// TransactionFilterAdmin filters for Administrative Transactions.
	TransactionFilterAdmin TransactionFilter = "ADMIN"
	// TransactionFilterCreate filters for Account Create Transaction.
	TransactionFilterCreate TransactionFilter = "CREATE"
	// TransactionFilterClose filters for Account Close Transaction.
	TransactionFilterClose TransactionFilter = "CLOSE"
	// TransactionFilterReopen filters for Account Reopen Transaction.
	TransactionFilterReopen TransactionFilter = "REOPEN"
	// TransactionFilterClientConfigure filters for Client Configure Transaction.
	TransactionFilterClientConfigure TransactionFilter = "CLIENT_CONFIGURE"
	// TransactionFilterClientConfigureReject filters for Client Configure Reject Transaction.
	TransactionFilterClientConfigureReject TransactionFilter = "CLIENT_CONFIGURE_REJECT"
	// TransactionFilterTransferFunds filters for Transfer Funds Transaction.
	TransactionFilterTransferFunds TransactionFilter = "TRANSFER_FUNDS"
	// TransactionFilterTransferFundsReject filters for Transfer Funds Reject Transaction.
	TransactionFilterTransferFundsReject TransactionFilter = "TRANSFER_FUNDS_REJECT"
	// TransactionFilterMarketOrder filters for Market Order Transaction.
	TransactionFilterMarketOrder TransactionFilter = "MARKET_ORDER"
	// TransactionFilterMarketOrderReject filters for Market Order Reject Transaction.
	TransactionFilterMarketOrderReject TransactionFilter = "MARKET_ORDER_REJECT"
	// TransactionFilterLimitOrder filters for Limit Order Transaction.
	TransactionFilterLimitOrder TransactionFilter = "LIMIT_ORDER"
	// TransactionFilterLimitOrderReject filters for Limit Order Reject Transaction.
	TransactionFilterLimitOrderReject TransactionFilter = "LIMIT_ORDER_REJECT"
	// TransactionFilterStopOrder filters for Stop Order Transaction.
	TransactionFilterStopOrder TransactionFilter = "STOP_ORDER"
	// TransactionFilterStopOrderReject filters for Stop Order Reject Transaction.
	TransactionFilterStopOrderReject TransactionFilter = "STOP_ORDER_REJECT"
	// TransactionFilterMarketIfTouchedOrder filters for Market If Touched Order Transaction.
	TransactionFilterMarketIfTouchedOrder TransactionFilter = "MARKET_IF_TOUCHED_ORDER"
	// TransactionFilterMarketIfTouchedOrderReject filters for Market If Touched Order Reject Transaction.
	TransactionFilterMarketIfTouchedOrderReject TransactionFilter = "MARKET_IF_TOUCHED_ORDER_REJECT"
	// TransactionFilterTakeProfitOrder filters for Take Profit Order Transaction.
	TransactionFilterTakeProfitOrder TransactionFilter = "TAKE_PROFIT_ORDER"
	// TransactionFilterTakeProfitOrderReject filters for Take Profit Order Reject Transaction.
	TransactionFilterTakeProfitOrderReject TransactionFilter = "TAKE_PROFIT_ORDER_REJECT"
	// TransactionFilterStopLossOrder filters for Stop Loss Order Transaction.
	TransactionFilterStopLossOrder TransactionFilter = "STOP_LOSS_ORDER"
	// TransactionFilterStopLossOrderReject filters for Stop Loss Order Reject Transaction.
	TransactionFilterStopLossOrderReject TransactionFilter = "STOP_LOSS_ORDER_REJECT"
	// TransactionFilterGuaranteedStopLossOrder filters for Guaranteed Stop Loss Order Transaction.
	TransactionFilterGuaranteedStopLossOrder TransactionFilter = "GUARANTEED_STOP_LOSS_ORDER"
	// TransactionFilterGuaranteedStopLossOrderReject filters for Guaranteed Stop Loss Order Reject Transaction.
	TransactionFilterGuaranteedStopLossOrderReject TransactionFilter = "GUARANTEED_STOP_LOSS_ORDER_REJECT"
	// TransactionFilterTrailingStopLossOrder filters for Trailing Stop Loss Order Transaction.
	TransactionFilterTrailingStopLossOrder TransactionFilter = "TRAILING_STOP_LOSS_ORDER"
	// TransactionFilterTrailingStopLossOrderReject filters for Trailing Stop Loss Order Reject Transaction.
	TransactionFilterTrailingStopLossOrderReject TransactionFilter = "TRAILING_STOP_LOSS_ORDER_REJECT"
	// TransactionFilterOneCancelsAllOrder filters for One Cancels All Order Transaction.
	TransactionFilterOneCancelsAllOrder TransactionFilter = "ONE_CANCELS_ALL_ORDER"
	// TransactionFilterOneCancelsAllOrderReject filters for One Cancels All Order Reject Transaction.
	TransactionFilterOneCancelsAllOrderReject TransactionFilter = "ONE_CANCELS_ALL_ORDER_REJECT"
	// TransactionFilterOneCancelsAllOrderTriggered filters for One Cancels All Order Triggered Transaction.
	TransactionFilterOneCancelsAllOrderTriggered TransactionFilter = "ONE_CANCELS_ALL_ORDER_TRIGGERED"
	// TransactionFilterOrderFill filters for Order Fill Transaction.
	TransactionFilterOrderFill TransactionFilter = "ORDER_FILL"
	// TransactionFilterOrderCancel filters for Order Cancel Transaction.
	TransactionFilterOrderCancel TransactionFilter = "ORDER_CANCEL"
	// TransactionFilterOrderCancelReject filters for Order Cancel Reject Transaction.
	TransactionFilterOrderCancelReject TransactionFilter = "ORDER_CANCEL_REJECT"
	// TransactionFilterOrderClientExtensionsModify filters for Order Client Extensions Modify Transaction.
	TransactionFilterOrderClientExtensionsModify TransactionFilter = "ORDER_CLIENT_EXTENSIONS_MODIFY"
	// TransactionFilterOrderClientExtensionsModifyReject filters for Order Client Extensions Modify Reject Transaction.
	TransactionFilterOrderClientExtensionsModifyReject TransactionFilter = "ORDER_CLIENT_EXTENSIONS_MODIFY_REJECT"
	// TransactionFilterTradeClientExtensionsModify filters for Trade Client Extensions Modify Transaction.
	TransactionFilterTradeClientExtensionsModify TransactionFilter = "TRADE_CLIENT_EXTENSIONS_MODIFY"
	// TransactionFilterTradeClientExtensionsModifyReject filters for Trade Client Extensions Modify Reject Transaction.
	TransactionFilterTradeClientExtensionsModifyReject TransactionFilter = "TRADE_CLIENT_EXTENSIONS_MODIFY_REJECT"
	// TransactionFilterMarginCallEnter filters for Margin Call Enter Transaction.
	TransactionFilterMarginCallEnter TransactionFilter = "MARGIN_CALL_ENTER"
	// TransactionFilterMarginCallExtend filters for Margin Call Extend Transaction.
	TransactionFilterMarginCallExtend TransactionFilter = "MARGIN_CALL_EXTEND"
	// TransactionFilterMarginCallExit filters for Margin Call Exit Transaction.
	TransactionFilterMarginCallExit TransactionFilter = "MARGIN_CALL_EXIT"
	// TransactionFilterDelayedTradeClosure filters for Delayed Trade Closure Transaction.
	TransactionFilterDelayedTradeClosure TransactionFilter = "DELAYED_TRADE_CLOSURE"
	// TransactionFilterDailyFinancing filters for Daily Financing Transaction.
	TransactionFilterDailyFinancing TransactionFilter = "DAILY_FINANCING"
	// TransactionFilterResetResettablePL filters for Reset Resettable PL Transaction.
	TransactionFilterResetResettablePL TransactionFilter = "RESET_RESETTABLE_PL"
)

// TransactionHeartbeat represents a heartbeat message sent for a Transaction stream.
type TransactionHeartbeat struct {
	// Type is the string "HEARTBEAT".
	Type string `json:"type"`
	// LastTransactionID is the ID of the most recent Transaction created for the Account.
	LastTransactionID TransactionID `json:"lastTransactionID"`
	// Time is the date/time when the TransactionHeartbeat was created.
	Time DateTime `json:"time"`
}

func (t TransactionHeartbeat) GetType() string {
	return t.Type
}

func (t TransactionHeartbeat) GetID() string {
	return t.LastTransactionID
}

func (t TransactionHeartbeat) GetTime() DateTime {
	return t.Time
}

// -------------------------------------------------------------------
// Endpoints https://developer.oanda.com/rest-live-v20/transaction-ep/
// -------------------------------------------------------------------

type transactionService struct {
	client *Client
}

func newTransactionService(client *Client) *transactionService {
	return &transactionService{client}
}

type TransactionListRequest struct {
	From     *time.Time
	To       *time.Time
	PageSize *int
	Type     []TransactionType
}

func NewTransactionListRequest() *TransactionListRequest {
	return &TransactionListRequest{}
}

func (req *TransactionListRequest) SetFrom(from time.Time) *TransactionListRequest {
	req.From = &from
	return req
}

func (req *TransactionListRequest) SetTo(to time.Time) *TransactionListRequest {
	req.To = &to
	return req
}

func (req *TransactionListRequest) SetPageSize(pageSize int) *TransactionListRequest {
	req.PageSize = &pageSize
	return req
}

func (req *TransactionListRequest) AddType(transactionType TransactionType) *TransactionListRequest {
	req.Type = append(req.Type, transactionType)
	return req
}

func (req *TransactionListRequest) validate() error {
	if req.PageSize != nil {
		if *req.PageSize < 1 {
			return errors.New("page size must be greater than zero")
		}
		if *req.PageSize > 1000 {
			return errors.New("page size must be equal or less than 1000")
		}
	}
	return nil
}

func (req *TransactionListRequest) values() (url.Values, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if req.From != nil {
		v.Set("from", req.From.Format(time.RFC3339))
	}
	if req.To != nil {
		v.Set("to", req.To.Format(time.RFC3339))
	}
	if req.PageSize != nil {
		v.Set("page_size", strconv.Itoa(*req.PageSize))
	}
	if len(req.Type) > 0 {
		var s []string
		for _, t := range req.Type {
			s = append(s, string(t))
		}
		v.Set("type", strings.Join(s, ","))
	}
	return v, nil
}

type TransactionListResponse struct {
	From              DateTime          `json:"from"`
	To                DateTime          `json:"to"`
	PageSize          int               `json:"pageSize"`
	Type              []TransactionType `json:"type"`
	Count             int               `json:"count"`
	Pages             []string          `json:"pages"`
	LastTransactionID TransactionID     `json:"lastTransactionID"`
}

func (s *transactionService) List(ctx context.Context, req *TransactionListRequest) (*TransactionListResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/transactions", s.client.accountID)
	v, err := req.values()
	if err != nil {
		return nil, err
	}
	resp, err := s.client.sendGetRequest(ctx, path, v)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	var transactionListResp TransactionListResponse
	if err := decodeResponse(resp, &transactionListResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &transactionListResp, nil
}

type TransactionDetailsResponse struct {
	Transaction       Transaction   `json:"transaction"`
	LastTransactionID TransactionID `json:"lastTransactionID"`
}

func (s *transactionService) Details(ctx context.Context, transactionID TransactionID) (*TransactionDetailsResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%v/transactions/%v", s.client.accountID, transactionID)
	return doGet[TransactionDetailsResponse](s.client, ctx, path, nil)
}

type TransactionGetByIDRangeRequest struct {
	From TransactionID
	To   TransactionID
	Type []TransactionFilter
}

func NewTransactionGetByIDRangeRequest(from, to TransactionID) *TransactionGetByIDRangeRequest {
	return &TransactionGetByIDRangeRequest{
		From: from,
		To:   to,
		Type: make([]TransactionFilter, 0),
	}
}

func (r *TransactionGetByIDRangeRequest) SetFilters(filters ...TransactionFilter) *TransactionGetByIDRangeRequest {
	r.Type = append(r.Type, filters...)
	return r
}

func (r *TransactionGetByIDRangeRequest) values() (url.Values, error) {
	values := url.Values{}
	values.Set("from", r.From)
	values.Set("to", r.To)
	if len(r.Type) > 0 {
		var s []string
		for _, t := range r.Type {
			s = append(s, string(t))
		}
		values.Set("type", strings.Join(s, ","))
	}
	return values, nil
}

type TransactionsResponse struct {
	Transactions      []Transaction `json:"transactions"`
	LastTransactionID TransactionID `json:"lastTransactionID"`
}

func (s *transactionService) GetByIDRange(ctx context.Context, req *TransactionGetByIDRangeRequest) (*TransactionsResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%s/transactions/idrange", s.client.accountID)
	v, err := req.values()
	if err != nil {
		return nil, err
	}
	return doGet[TransactionsResponse](s.client, ctx, path, v)
}

type TransactionGetBySinceIDRequest struct {
	ID   TransactionID
	Type []TransactionFilter
}

func NewTransactionGetBySinceIDRequest(id TransactionID) *TransactionGetBySinceIDRequest {
	return &TransactionGetBySinceIDRequest{
		ID:   id,
		Type: make([]TransactionFilter, 0),
	}
}

func (r *TransactionGetBySinceIDRequest) SetFilters(filters ...TransactionFilter) *TransactionGetBySinceIDRequest {
	r.Type = append(r.Type, filters...)
	return r
}

func (r *TransactionGetBySinceIDRequest) values() (url.Values, error) {
	values := url.Values{}
	values.Set("id", string(r.ID))
	if len(r.Type) > 0 {
		var s []string
		for _, t := range r.Type {
			s = append(s, string(t))
		}
		values.Set("type", strings.Join(s, ","))
	}
	return values, nil
}

func (s *transactionService) GetBySinceID(ctx context.Context, req *TransactionGetBySinceIDRequest) (*TransactionsResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%s/transactions/sinceid", s.client.accountID)
	v, err := req.values()
	if err != nil {
		return nil, err
	}
	return doGet[TransactionsResponse](s.client, ctx, path, v)
}

type transactionStreamService struct {
	client *StreamClient
}

func newTransactionStreamService(client *StreamClient) *transactionStreamService {
	return &transactionStreamService{client}
}

type TransactionStreamItem interface {
	GetType() string
	GetID() TransactionID
	GetTime() DateTime
}

func (s *transactionStreamService) Stream(ctx context.Context, ch chan<- TransactionStreamItem, done <-chan struct{}) error {
	path := fmt.Sprintf("/v3/accounts/%s/transactions/stream", s.client.accountID)
	u, err := joinURL(s.client.baseURL, path, nil)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return err
	}
	s.client.setHeaders(httpReq)
	httpResp, err := s.client.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %w", err)
	}
	defer closeBody(httpResp)
	dec := json.NewDecoder(httpResp.Body)
	for {
		select {
		case <-done:
			slog.Info("transaction stream closed")
			return nil
		case <-ctx.Done():
			return ctx.Err()
		default:
			var typeOnly struct {
				Type string `json:"type"`
			}
			if err := dec.Decode(&typeOnly); err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("failed to decode JSON response: %w", err)
			}
			switch typeOnly.Type {
			case "CREATE":
				if err := decodeItem[CreateTransaction](dec, ch); err != nil {
					return err
				}
			case "CLOSE":
				if err := decodeItem[CloseTransaction](dec, ch); err != nil {
					return err
				}
			case "REOPEN":
				if err := decodeItem[ReopenTransaction](dec, ch); err != nil {
					return err
				}
			case "CLIENT_CONFIGURE":
				if err := decodeItem[ClientConfigureTransaction](dec, ch); err != nil {
					return err
				}
			case "CLIENT_CONFIGURE_REJECT":
				if err := decodeItem[ClientConfigureRejectTransaction](dec, ch); err != nil {
					return err
				}
			case "TRANSFER_FUNDS":
				if err := decodeItem[TransferFundsTransaction](dec, ch); err != nil {
					return err
				}
			case "TRANSFER_FUNDS_REJECT":
				if err := decodeItem[TransferFundsRejectTransaction](dec, ch); err != nil {
					return err
				}
			case "MARKET_ORDER":
				if err := decodeItem[MarketOrderTransaction](dec, ch); err != nil {
					return err
				}
			case "MARKET_ORDER_REJECT":
				if err := decodeItem[MarketOrderRejectTransaction](dec, ch); err != nil {
					return err
				}
			case "FIXED_PRICE_ORDER":
				if err := decodeItem[FixedPriceOrderTransaction](dec, ch); err != nil {
					return err
				}
			case "LIMIT_ORDER":
				if err := decodeItem[LimitOrderTransaction](dec, ch); err != nil {
					return err
				}
			case "LIMIT_ORDER_REJECT":
				if err := decodeItem[LimitOrderRejectTransaction](dec, ch); err != nil {
					return err
				}
			case "STOP_ORDER":
				if err := decodeItem[StopOrderTransaction](dec, ch); err != nil {
					return err
				}
			case "STOP_ORDER_REJECT":
				if err := decodeItem[StopOrderRejectTransaction](dec, ch); err != nil {
					return err
				}
			case "MARKET_IF_TOUCHED_ORDER":
				if err := decodeItem[MarketIfTouchedOrderTransaction](dec, ch); err != nil {
					return err
				}
			case "MARKET_IF_TOUCHED_ORDER_REJECT":
				if err := decodeItem[MarketIfTouchedOrderRejectTransaction](dec, ch); err != nil {
					return err
				}
			case "TAKE_PROFIT_ORDER":
				if err := decodeItem[TakeProfitOrderTransaction](dec, ch); err != nil {
					return err
				}
			case "TAKE_PROFIT_ORDER_REJECT":
				if err := decodeItem[TakeProfitOrderRejectTransaction](dec, ch); err != nil {
					return err
				}
			case "STOP_LOSS_ORDER":
				if err := decodeItem[StopLossOrderTransaction](dec, ch); err != nil {
					return err
				}
			case "STOP_LOSS_ORDER_REJECT":
				if err := decodeItem[StopLossOrderRejectTransaction](dec, ch); err != nil {
					return err
				}
			case "GUARANTEED_STOP_LOSS_ORDER":
				if err := decodeItem[GuaranteedStopLossOrderTransaction](dec, ch); err != nil {
					return err
				}
			case "GUARANTEED_STOP_LOSS_ORDER_REJECT":
				if err := decodeItem[GuaranteedStopLossOrderRejectTransaction](dec, ch); err != nil {
					return err
				}
			case "TRAILING_STOP_LOSS_ORDER":
				if err := decodeItem[TrailingStopLossOrderTransaction](dec, ch); err != nil {
					return err
				}
			case "TRAILING_STOP_LOSS_ORDER_REJECT":
				if err := decodeItem[TrailingStopLossOrderRejectTransaction](dec, ch); err != nil {
					return err
				}
			case "ORDER_FILL":
				if err := decodeItem[OrderFillTransaction](dec, ch); err != nil {
					return err
				}
			case "ORDER_CANCEL":
				if err := decodeItem[OrderCancelTransaction](dec, ch); err != nil {
					return err
				}
			case "ORDER_CANCEL_REJECT":
				if err := decodeItem[OrderCancelRejectTransaction](dec, ch); err != nil {
					return err
				}
			case "ORDER_CLIENT_EXTENSIONS_MODIFY":
				if err := decodeItem[OrderClientExtensionsModifyTransaction](dec, ch); err != nil {
					return err
				}
			case "ORDER_CLIENT_EXTENSIONS_MODIFY_REJECT":
				if err := decodeItem[OrderClientExtensionsModifyRejectTransaction](dec, ch); err != nil {
					return err
				}
			case "TRADE_CLIENT_EXTENSIONS_MODIFY":
				if err := decodeItem[TradeClientExtensionsModifyTransaction](dec, ch); err != nil {
					return err
				}
			case "TRADE_CLIENT_EXTENSIONS_MODIFY_REJECT":
				if err := decodeItem[TradeClientExtensionsModifyRejectTransaction](dec, ch); err != nil {
					return err
				}
			case "MARGIN_CALL_ENTER":
				if err := decodeItem[MarginCallEnterTransaction](dec, ch); err != nil {
					return err
				}
			case "MARGIN_CALL_EXTEND":
				if err := decodeItem[MarginCallExtendTransaction](dec, ch); err != nil {
					return err
				}
			case "MARGIN_CALL_EXIT":
				if err := decodeItem[MarginCallExitTransaction](dec, ch); err != nil {
					return err
				}
			case "DELAYED_TRADE_CLOSURE":
				if err := decodeItem[DelayedTradeClosureTransaction](dec, ch); err != nil {
					return err
				}
			case "DAILY_FINANCING":
				if err := decodeItem[DailyFinancingTransaction](dec, ch); err != nil {
					return err
				}
			case "DIVIDEND_ADJUSTMENT":
				if err := decodeItem[DividendAdjustmentTransaction](dec, ch); err != nil {
					return err
				}
			case "RESET_RESETTABLE_PL":
				if err := decodeItem[ResetResettablePLTransaction](dec, ch); err != nil {
					return err
				}
			case "HEARTBEAT":
				if err := decodeItem[TransactionHeartbeat](dec, ch); err != nil {
					return err
				}
			}

		}
	}
}

func decodeItem[R TransactionStreamItem](dec *json.Decoder, ch chan<- TransactionStreamItem) error {
	var t R
	if err := dec.Decode(&t); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}
	ch <- t
	return nil
}
