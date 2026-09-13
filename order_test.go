package oanda

import (
	"errors"
	"net/http"
	"testing"
)

func TestOrderCreateRejected(t *testing.T) {
	fake := &fakeHTTPClient{responses: []*http.Response{jsonResponse(http.StatusBadRequest,
		`{"orderRejectTransaction":{"id":"42","type":"MARKET_ORDER_REJECT","rejectReason":"UNITS_INVALID"},"errorCode":"UNITS_INVALID","errorMessage":"Invalid units","lastTransactionID":"42"}`)}}
	_, err := newFakeClient(fake).Order.Create(t.Context(), NewMarketOrderRequest("USD_JPY", "0"))
	var badRequest BadRequest
	var rejection OrderErrorResponse
	if !errors.As(err, &badRequest) || !errors.As(err, &rejection) {
		t.Fatalf("error = %T (%v), want BadRequest wrapping OrderErrorResponse", err, err)
	}
	txn, ok := rejection.OrderRejectTransaction.(*MarketOrderRejectTransaction)
	if !ok || txn.GetID() != "42" || txn.RejectReason != "UNITS_INVALID" {
		t.Errorf("unexpected rejection transaction: %#v", rejection.OrderRejectTransaction)
	}
	if rejection.ErrorCode != "UNITS_INVALID" || rejection.LastTransactionID != "42" {
		t.Errorf("unexpected rejection: %#v", rejection)
	}
}

func TestOrderService(t *testing.T) {
	runEndpointTests(t, []endpointTest{
		{
			name:     `create market`,
			method:   `POST`,
			path:     testAccountPath + "/orders",
			body:     `{"order":{"type":"MARKET","instrument":"USD_JPY","units":"10000","timeInForce":"FOK","positionFill":"DEFAULT"}}`,
			status:   201,
			response: `{"orderCreateTransaction":{"id":"41","type":"MARKET_ORDER","instrument":"USD_JPY","units":"10000"},"orderFillTransaction":{"id":"42","type":"ORDER_FILL","tradeOpened":{"tradeID":"42","units":"10000","price":"150.000"}},"lastTransactionID":"42"}`,
			call: func(c *Client) (any, error) {
				return c.Order.Create(t.Context(), NewMarketOrderRequest("USD_JPY", "10000"))
			},
		},
		{
			name:     `create limit`,
			method:   `POST`,
			path:     testAccountPath + "/orders",
			body:     `{"order":{"type":"LIMIT","instrument":"USD_JPY","units":"10000","price":"100.00","timeInForce":"GTC","positionFill":"DEFAULT","triggerCondition":"DEFAULT"}}`,
			status:   201,
			response: `{"orderCreateTransaction":{"id":"42","type":"LIMIT_ORDER","price":"100.00"},"lastTransactionID":"42"}`,
			call: func(c *Client) (any, error) {
				return c.Order.Create(t.Context(), NewLimitOrderRequest("USD_JPY", "10000", "100.00"))
			},
		},
		{
			name:     `create take profit`,
			method:   `POST`,
			path:     testAccountPath + "/orders",
			body:     `{"order":{"type":"TAKE_PROFIT","tradeID":"42","price":"170.00","timeInForce":"GTC","triggerCondition":"DEFAULT"}}`,
			status:   201,
			response: `{"orderCreateTransaction":{"id":"43","type":"TAKE_PROFIT_ORDER","tradeID":"42","price":"170.00"},"lastTransactionID":"43"}`,
			call: func(c *Client) (any, error) {
				return c.Order.Create(t.Context(), NewTakeProfitOrderRequest("42", "170.00"))
			},
		},
		{
			name:     `create stop loss`,
			method:   `POST`,
			path:     testAccountPath + "/orders",
			body:     `{"order":{"type":"STOP_LOSS","tradeID":"42","distance":"10.000","timeInForce":"GTC","triggerCondition":"DEFAULT"}}`,
			status:   201,
			response: `{"orderCreateTransaction":{"id":"44","type":"STOP_LOSS_ORDER","tradeID":"42","distance":"10.000"},"lastTransactionID":"44"}`,
			call: func(c *Client) (any, error) {
				return c.Order.Create(t.Context(), NewStopLossOrderRequest("42").SetDistance("10.000"))
			},
		},
		{
			name:     `list`,
			method:   `GET`,
			path:     testAccountPath + "/orders",
			query:    `instrument=USD_JPY`,
			response: `{"orders":[{"id":"42","type":"LIMIT","instrument":"USD_JPY","units":"10000","price":"100.00","state":"PENDING"}],"lastTransactionID":"50"}`,
			call: func(c *Client) (any, error) {
				return c.Order.List(t.Context(), NewOrderListRequest().SetInstrument("USD_JPY"))
			},
		},
		{
			name:     `list pending`,
			method:   `GET`,
			path:     testAccountPath + "/pendingOrders",
			response: `{"orders":[{"id":"42","type":"LIMIT","instrument":"USD_JPY","units":"10000","price":"100.00","state":"PENDING"}],"lastTransactionID":"50"}`,
			call:     func(c *Client) (any, error) { return c.Order.ListPending(t.Context()) },
		},
		{
			name:     `details`,
			method:   `GET`,
			path:     testAccountPath + "/orders/42",
			response: `{"order":{"id":"42","type":"LIMIT","instrument":"USD_JPY","units":"10000","price":"100.00","state":"PENDING"},"lastTransactionID":"50"}`,
			call:     func(c *Client) (any, error) { return c.Order.Details(t.Context(), "42") },
		},
		{
			name:     `replace`,
			method:   `PUT`,
			path:     testAccountPath + "/orders/42",
			body:     `{"order":{"type":"LIMIT","instrument":"USD_JPY","units":"10000","price":"110.00","timeInForce":"GTC","positionFill":"DEFAULT","triggerCondition":"DEFAULT"}}`,
			status:   201,
			response: `{"orderCancelTransaction":{"id":"51","type":"ORDER_CANCEL","orderID":"42"},"orderCreateTransaction":{"id":"52","type":"LIMIT_ORDER","price":"110.00"},"lastTransactionID":"52"}`,
			call: func(c *Client) (any, error) {
				return c.Order.Replace(t.Context(), "42", NewLimitOrderRequest("USD_JPY", "10000", "110.00"))
			},
		},
		{
			name:     `client extensions`,
			method:   `PUT`,
			path:     testAccountPath + "/orders/42/clientExtensions",
			body:     `{"clientExtensions":{"id":"order-id"},"tradeClientExtensions":{"id":"trade-id"}}`,
			response: `{"orderClientExtensionsModifyTransaction":{"id":"51","type":"ORDER_CLIENT_EXTENSIONS_MODIFY","orderID":"42","clientExtensionsModify":{"id":"order-id"},"tradeClientExtensionsModify":{"id":"trade-id"}},"lastTransactionID":"51"}`,
			call: func(c *Client) (any, error) {
				return c.Order.UpdateClientExtensions(t.Context(), "42", OrderUpdateClientExtensionsRequest{ClientExtensions: NewClientExtensions().SetID("order-id"), TradeClientExtensions: NewClientExtensions().SetID("trade-id")})
			},
		},
		{
			name:     `cancel`,
			method:   `PUT`,
			path:     testAccountPath + "/orders/42/cancel",
			response: `{"orderCancelTransaction":{"id":"51","type":"ORDER_CANCEL","orderID":"42","reason":"CLIENT_REQUEST"},"lastTransactionID":"51"}`,
			call:     func(c *Client) (any, error) { return c.Order.Cancel(t.Context(), "42") },
		},
	})
}
