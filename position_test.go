package oanda

import (
	"strings"
	"testing"
)

func TestPositionService(t *testing.T) {
	runEndpointTests(t, []endpointTest{
		{
			name:     `list`,
			method:   `GET`,
			path:     testAccountPath + "/positions",
			response: `{"positions":[{"instrument":"USD_JPY","long":{"units":"10000","averagePrice":"150.000","tradeIDs":["42"]},"short":{"units":"0"},"pl":"1.50"}],"lastTransactionID":"50"}`,
			want:     `{"positions":[{"instrument":"USD_JPY","long":{"units":"10000","averagePrice":"150.000","tradeIDs":["42"]},"short":{"units":"0"},"pl":"1.50"}],"lastTransactionID":"50"}`,
			call:     func(c *Client) (any, error) { return c.Position.List(t.Context()) },
		},
		{
			name:     `list open`,
			method:   `GET`,
			path:     testAccountPath + "/openPositions",
			response: `{"positions":[{"instrument":"USD_JPY","long":{"units":"10000","averagePrice":"150.000","tradeIDs":["42"]},"short":{"units":"0"},"pl":"1.50"}],"lastTransactionID":"50"}`,
			want:     `{"positions":[{"instrument":"USD_JPY","long":{"units":"10000","averagePrice":"150.000","tradeIDs":["42"]},"short":{"units":"0"},"pl":"1.50"}],"lastTransactionID":"50"}`,
			call:     func(c *Client) (any, error) { return c.Position.ListOpen(t.Context()) },
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
			body:     `{"longUnits":"ALL","shortUnits":"NONE"}`,
			response: `{"longOrderCreateTransaction":{"id":"51","type":"MARKET_ORDER","units":"-10000"},"longOrderFillTransaction":{"id":"52","type":"ORDER_FILL","units":"-10000"},"relatedTransactionIDs":["51","52"],"lastTransactionID":"52"}`,
			call: func(c *Client) (any, error) {
				return c.Position.Close(t.Context(), "USD_JPY", NewPositionCloseRequest().SetLongAll())
			},
		},
	})
}

func TestPositionCloseRequestLeavesUnselectedSideOpen(t *testing.T) {
	tests := []struct {
		name string
		req  *PositionCloseRequest
		want string
	}{
		{"long only", NewPositionCloseRequest().SetLongUnits("100"), `{"longUnits":"100","shortUnits":"NONE"}`},
		{"short only", NewPositionCloseRequest().SetShortAll(), `{"longUnits":"NONE","shortUnits":"ALL"}`},
		{"both sides", NewPositionCloseRequest().SetLongAll().SetShortUnits("50"), `{"longUnits":"ALL","shortUnits":"50"}`},
		{"direct struct", &PositionCloseRequest{}, `{"longUnits":"NONE","shortUnits":"NONE"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := tt.req.body()
			if err != nil {
				t.Fatal(err)
			}
			if got := body.String(); got != tt.want {
				t.Errorf("body = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestPositionCloseUnits(t *testing.T) {
	fake := &fakeHTTPClient{}
	req := NewPositionCloseRequest().SetLongUnits("0.5").SetShortUnits("1000")
	if _, err := newFakeClient(fake).Position.Close(t.Context(), "XAU_USD", req); err != nil {
		t.Fatal(err)
	}
	if want := `{"longUnits":"0.5","shortUnits":"1000"}`; fake.bodies[0] != want {
		t.Errorf("body = %s, want %s", fake.bodies[0], want)
	}

	for _, units := range []DecimalNumber{"0", "-10", "", "ten", "all"} {
		for side, req := range map[string]*PositionCloseRequest{
			"long":  NewPositionCloseRequest().SetLongUnits(units),
			"short": NewPositionCloseRequest().SetShortUnits(units),
		} {
			t.Run(side+" "+string(units), func(t *testing.T) {
				fake := &fakeHTTPClient{}
				_, err := newFakeClient(fake).Position.Close(t.Context(), "EUR_USD", req)
				if err == nil || !strings.Contains(err.Error(), side+" units") {
					t.Fatalf("err = %v, want a %s units error", err, side)
				}
				if fake.calls != 0 {
					t.Errorf("sent %d HTTP requests for invalid units", fake.calls)
				}
			})
		}
	}
}
