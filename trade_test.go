package oanda

import "testing"

func TestTradeService(t *testing.T) {
	runEndpointTests(t, []endpointTest{
		{
			name:     `list`,
			method:   `GET`,
			path:     testAccountPath + "/trades",
			query:    `instrument=USD_JPY`,
			response: `{"trades":[{"id":"42","instrument":"USD_JPY","price":"150.000","state":"OPEN","currentUnits":"10000","openTime":"2025-01-01T00:00:00Z"}],"lastTransactionID":"50"}`,
			call: func(c *Client) (any, error) {
				return c.Trade.List(t.Context(), NewTradeListRequest().SetInstrument("USD_JPY"))
			},
		},
		{
			name:     `list open`,
			method:   `GET`,
			path:     testAccountPath + "/openTrades",
			response: `{"trades":[{"id":"42","instrument":"USD_JPY","price":"150.000","state":"OPEN","currentUnits":"10000","openTime":"2025-01-01T00:00:00Z"}],"lastTransactionID":"50"}`,
			call:     func(c *Client) (any, error) { return c.Trade.ListOpen(t.Context()) },
		},
		{
			name:     `details`,
			method:   `GET`,
			path:     testAccountPath + "/trades/42",
			response: `{"trade":{"id":"42","instrument":"USD_JPY","price":"150.000","state":"OPEN","currentUnits":"10000","openTime":"2025-01-01T00:00:00Z"},"lastTransactionID":"50"}`,
			call:     func(c *Client) (any, error) { return c.Trade.Details(t.Context(), "42") },
		},
		{
			name:     `client extensions`,
			method:   `PUT`,
			path:     testAccountPath + "/trades/42/clientExtensions",
			body:     `{"clientExtensions":{"id":"test-id","tag":"test-tag"}}`,
			response: `{"tradeClientExtensionsModifyTransaction":{"id":"51","type":"TRADE_CLIENT_EXTENSIONS_MODIFY","tradeID":"42","tradeClientExtensionsModify":{"id":"test-id","tag":"test-tag"}},"lastTransactionID":"51"}`,
			call: func(c *Client) (any, error) {
				return c.Trade.UpdateClientExtensions(t.Context(), "42", TradeUpdateClientExtensionsRequest{ClientExtensions: NewClientExtensions().SetID("test-id").SetTag("test-tag")})
			},
		},
		{
			name:     `dependent orders`,
			method:   `PUT`,
			path:     testAccountPath + "/trades/42/orders",
			body:     `{"takeProfit":{"price":"170.00","timeInForce":"GTC"},"stopLoss":{"distance":"10.00","timeInForce":"GTC"}}`,
			response: `{"takeProfitOrderTransaction":{"id":"51","type":"TAKE_PROFIT_ORDER","tradeID":"42","price":"170.00"},"stopLossOrderTransaction":{"id":"52","type":"STOP_LOSS_ORDER","tradeID":"42","distance":"10.00"},"lastTransactionID":"52"}`,
			call: func(c *Client) (any, error) {
				return c.Trade.UpdateOrders(t.Context(), "42", &TradeUpdateOrdersRequest{TakeProfit: NewTakeProfitDetails("170.00"), StopLoss: NewStopLossDetails().SetDistance("10.00")})
			},
		},
		{
			name:     `close`,
			method:   `PUT`,
			path:     testAccountPath + "/trades/42/close",
			body:     `{"units":"ALL"}`,
			response: `{"orderCreateTransaction":{"id":"51","type":"MARKET_ORDER","tradeClose":{"tradeID":"42","units":"ALL"}},"orderFillTransaction":{"id":"52","type":"ORDER_FILL","tradesClosed":[{"tradeID":"42","units":"10000"}]},"lastTransactionID":"52"}`,
			call:     func(c *Client) (any, error) { return c.Trade.Close(t.Context(), "42", NewTradeCloseALLRequest()) },
		},
	})
}
