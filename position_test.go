package oanda

import "testing"

func TestPositionService(t *testing.T) {
	runEndpointTests(t, []endpointTest{
		{
			name:     `list`,
			method:   `GET`,
			path:     testAccountPath + "/positions",
			response: `{"positions":[{"instrument":"USD_JPY","long":{"units":"10000","averagePrice":"150.000","tradeIDs":["42"]},"short":{"units":"0"},"pl":"1.50"}],"lastTransactionID":"50"}`,
			// The response type currently marshals this field as lastTransactionId.
			want: `{"positions":[{"instrument":"USD_JPY","long":{"units":"10000","averagePrice":"150.000","tradeIDs":["42"]},"short":{"units":"0"},"pl":"1.50"}],"lastTransactionId":"50"}`,
			call: func(c *Client) (any, error) { return c.Position.List(t.Context()) },
		},
		{
			name:     `list open`,
			method:   `GET`,
			path:     testAccountPath + "/openPositions",
			response: `{"positions":[{"instrument":"USD_JPY","long":{"units":"10000","averagePrice":"150.000","tradeIDs":["42"]},"short":{"units":"0"},"pl":"1.50"}],"lastTransactionID":"50"}`,
			// The response type currently marshals this field as lastTransactionId.
			want: `{"positions":[{"instrument":"USD_JPY","long":{"units":"10000","averagePrice":"150.000","tradeIDs":["42"]},"short":{"units":"0"},"pl":"1.50"}],"lastTransactionId":"50"}`,
			call: func(c *Client) (any, error) { return c.Position.ListOpen(t.Context()) },
		},
		{
			name:     `by instrument`,
			method:   `GET`,
			path:     testAccountPath + "/positions/USD_JPY",
			response: `{"position":{"instrument":"USD_JPY","long":{"units":"10000","averagePrice":"150.000","tradeIDs":["42"]},"short":{"units":"0"},"pl":"1.50"},"lastTransactionID":"50"}`,
			call:     func(c *Client) (any, error) { return c.Position.ListByInstrument(t.Context(), "USD_JPY") },
		},
		{
			name:     `close`,
			method:   `PUT`,
			path:     testAccountPath + "/positions/USD_JPY/close",
			body:     `{"longUnits":"ALL"}`,
			response: `{"longOrderCreateTransaction":{"id":"51","type":"MARKET_ORDER","units":"-10000"},"longOrderFillTransaction":{"id":"52","type":"ORDER_FILL","units":"-10000"},"relatedTransactionIDs":["51","52"],"lastTransactionID":"52"}`,
			call: func(c *Client) (any, error) {
				return c.Position.Close(t.Context(), "USD_JPY", NewPositionCloseRequest().SetLongAll())
			},
		},
	})
}
