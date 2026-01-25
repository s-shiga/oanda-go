package oanda

// Definitions https://developer.oanda.com/rest-live-v20/transaction-df/

// TransactionID is the unique identifier of a Transaction.
type TransactionID string

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
	ID ClientID `json:"id"`
	// Tag is a client-provided tag that can be associated with an Order or Trade.
	Tag ClientTag `json:"tag"`
	// Comment is a client-provided comment that can be associated with an Order or Trade.
	Comment ClientComment `json:"comment"`
}

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

// RequestID is the unique identifier for a client request.
type RequestID string
