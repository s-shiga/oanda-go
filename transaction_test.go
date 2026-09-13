package oanda

import "testing"

func TestTransactionService(t *testing.T) {
	runEndpointTests(t, []endpointTest{
		{
			name:     `list`,
			method:   `GET`,
			path:     testAccountPath + "/transactions",
			query:    `pageSize=50`,
			response: `{"from":"2025-01-01T00:00:00Z","to":"2025-01-02T00:00:00Z","pageSize":50,"count":1,"lastTransactionID":"42","pages":["https://example.invalid/transactions/idrange?from=42&to=42"]}`,
			call: func(c *Client) (any, error) {
				return c.Transaction.List(t.Context(), NewTransactionListRequest().SetPageSize(50))
			},
		},
		{
			name:     `details`,
			method:   `GET`,
			path:     testAccountPath + "/transactions/42",
			response: `{"transaction":{"id":"42","type":"ORDER_FILL","time":"2025-01-01T00:00:00.123456789Z","instrument":"USD_JPY","units":"10000","price":"150.000"},"lastTransactionID":"42"}`,
			call:     func(c *Client) (any, error) { return c.Transaction.Details(t.Context(), "42") },
		},
		{
			name:     `ID range`,
			method:   `GET`,
			path:     testAccountPath + "/transactions/idrange",
			query:    `from=40&to=50`,
			response: `{"transactions":[{"id":"42","type":"ORDER_FILL","time":"2025-01-01T00:00:00.123456789Z","instrument":"USD_JPY","units":"10000","price":"150.000"}],"lastTransactionID":"50"}`,
			call: func(c *Client) (any, error) {
				return c.Transaction.GetByIDRange(t.Context(), NewTransactionGetByIDRangeRequest("40", "50"))
			},
		},
		{
			name:     `since ID`,
			method:   `GET`,
			path:     testAccountPath + "/transactions/sinceid",
			query:    `id=40`,
			response: `{"transactions":[{"id":"42","type":"ORDER_FILL","time":"2025-01-01T00:00:00.123456789Z","instrument":"USD_JPY","units":"10000","price":"150.000"}],"lastTransactionID":"50"}`,
			call: func(c *Client) (any, error) {
				return c.Transaction.GetBySinceID(t.Context(), NewTransactionGetBySinceIDRequest("40"))
			},
		},
	})
}
