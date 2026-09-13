package oanda

import "testing"

func TestAccountService(t *testing.T) {
	runEndpointTests(t, []endpointTest{
		{
			name:     `list`,
			method:   `GET`,
			path:     "/v3/accounts",
			response: `{"accounts":[{"id":"101-001-1234567-001","tags":["demo"]}]}`,
			call:     func(c *Client) (any, error) { return c.Account.List(t.Context()) },
		},
		{
			name:     `details`,
			method:   `GET`,
			path:     testAccountPath + "",
			response: `{"account":{"id":"101-001-1234567-001","balance":"1000.00","orders":[{"id":"42","type":"LIMIT","instrument":"USD_JPY","price":"100.00","state":"PENDING"}]},"lastTransactionID":"50"}`,
			call:     func(c *Client) (any, error) { return c.Account.Details(t.Context()) },
		},
		{
			name:     `summary`,
			method:   `GET`,
			path:     testAccountPath + "/summary",
			response: `{"account":{"id":"101-001-1234567-001","balance":"1000.00","openTradeCount":1},"lastTransactionID":"50"}`,
			call:     func(c *Client) (any, error) { return c.Account.Summary(t.Context()) },
		},
		{
			name:     `configure`,
			method:   `PATCH`,
			path:     testAccountPath + "/configuration",
			body:     `{"alias":"TestAlias"}`,
			response: `{"clientConfigureTransaction":{"id":"51","type":"CLIENT_CONFIGURE","alias":"TestAlias"},"lastTransactionID":"51"}`,
			call: func(c *Client) (any, error) {
				return c.Account.Configure(t.Context(), NewAccountConfigureRequest().SetAlias("TestAlias"))
			},
		},
		{
			name:     `changes`,
			method:   `GET`,
			path:     testAccountPath + "/changes",
			query:    `sinceTransactionID=40`,
			response: `{"changes":{"ordersCreated":[{"id":"42","type":"LIMIT","price":"100.00"}],"transactions":[{"id":"42","type":"LIMIT_ORDER","price":"100.00"}]},"state":{"NAV":"1001.00"},"lastTransactionID":"50"}`,
			call:     func(c *Client) (any, error) { return c.Account.Changes(t.Context(), "40") },
		},
	})
}
