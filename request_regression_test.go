package oanda

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestDependentOrderCancellation(t *testing.T) {
	cases := []struct {
		name string
		req  TradeUpdateOrdersRequest
		want string
	}{
		{"unchanged", TradeUpdateOrdersRequest{}, `{}`},
		{"take profit", TradeUpdateOrdersRequest{CancelTakeProfit: true}, `{"takeProfit":null}`},
		{"stop loss", TradeUpdateOrdersRequest{CancelStopLoss: true}, `{"stopLoss":null}`},
		{"trailing stop loss", TradeUpdateOrdersRequest{CancelTrailingStopLoss: true}, `{"trailingStopLoss":null}`},
		{"guaranteed stop loss", TradeUpdateOrdersRequest{CancelGuaranteedStopLoss: true}, `{"guaranteedStopLoss":null}`},
		{"cancel all", TradeUpdateOrdersRequest{CancelTakeProfit: true, CancelStopLoss: true, CancelTrailingStopLoss: true, CancelGuaranteedStopLoss: true}, `{"takeProfit":null,"stopLoss":null,"trailingStopLoss":null,"guaranteedStopLoss":null}`},
		{"update and cancel different orders", TradeUpdateOrdersRequest{TakeProfit: NewTakeProfitDetails("1.3"), CancelStopLoss: true}, `{"takeProfit":{"price":"1.3","timeInForce":"GTC"},"stopLoss":null}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := json.Marshal(tc.req)
			if err != nil {
				t.Fatal(err)
			}
			var restored TradeUpdateOrdersRequest
			if err := json.Unmarshal(encoded, &restored); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(restored, tc.req) {
				t.Fatalf("round trip = %#v, want %#v", restored, tc.req)
			}
			fake := &fakeHTTPClient{}
			if _, err := newFakeClient(fake).Trade.UpdateOrders(t.Context(), "42", &restored); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decodeTestJSON(t, fake.bodies[0]), decodeTestJSON(t, tc.want)) {
				t.Errorf("body = %s, want %s", fake.bodies[0], tc.want)
			}
			b, err := json.Marshal(tc.req)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decodeTestJSON(t, string(b)), decodeTestJSON(t, tc.want)) {
				t.Errorf("MarshalJSON = %s, want %s", b, tc.want)
			}
		})
	}
}

func TestDependentOrderDecodeReusedRequest(t *testing.T) {
	var req TradeUpdateOrdersRequest
	for _, input := range []string{
		`{"takeProfit":null,"stopLoss":null,"trailingStopLoss":null,"guaranteedStopLoss":null}`,
		`{"takeProfit":{"price":"1.3","timeInForce":"GTC"},"stopLoss":{"price":"1.2","timeInForce":"GTC"},"trailingStopLoss":{"distance":"0.1","timeInForce":"GTC"},"guaranteedStopLoss":{"price":"1.1","timeInForce":"GTC"}}`,
		`{}`,
		`{"stopLoss": null }`,
		`{}`,
	} {
		if err := json.Unmarshal([]byte(input), &req); err != nil {
			t.Fatal(err)
		}
		encoded, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(decodeTestJSON(t, string(encoded)), decodeTestJSON(t, input)) {
			t.Fatalf("reused request = %s, want %s", encoded, input)
		}
	}
}

func TestDependentOrderCancellationConflicts(t *testing.T) {
	cases := []struct {
		field string
		req   TradeUpdateOrdersRequest
	}{
		{"takeProfit", TradeUpdateOrdersRequest{CancelTakeProfit: true, TakeProfit: NewTakeProfitDetails("1.3")}},
		{"stopLoss", TradeUpdateOrdersRequest{CancelStopLoss: true, StopLoss: NewStopLossDetails().SetPrice("1.2")}},
		{"trailingStopLoss", TradeUpdateOrdersRequest{CancelTrailingStopLoss: true, TrailingStopLoss: NewTrailingStopLossDetails("0.1")}},
		{"guaranteedStopLoss", TradeUpdateOrdersRequest{CancelGuaranteedStopLoss: true, GuaranteedStopLoss: NewGuaranteedStopLossDetails().SetPrice("1.2")}},
	}
	for _, tc := range cases {
		t.Run(tc.field, func(t *testing.T) {
			fake := &fakeHTTPClient{}
			_, err := newFakeClient(fake).Trade.UpdateOrders(t.Context(), "42", &tc.req)
			if err == nil || !strings.Contains(err.Error(), tc.field) {
				t.Fatalf("error = %v, want conflict for %s", err, tc.field)
			}
			if fake.calls != 0 {
				t.Errorf("sent %d HTTP requests for conflicting cancellation", fake.calls)
			}
		})
	}
}

func TestRequestTimePrecision(t *testing.T) {
	from := mustTime(t, "2026-09-13T10:00:00.100000001Z")
	to := mustTime(t, "2026-09-13T10:00:00.600000001Z")
	cases := []struct {
		name string
		call func(*Client) error
	}{
		{"transactions", func(c *Client) error {
			_, err := c.Transaction.List(t.Context(), NewTransactionListRequest().SetFrom(from).SetTo(to))
			return err
		}},
		{"instrument candles", func(c *Client) error {
			_, err := c.Instrument.Candlesticks(t.Context(), NewCandlesticksRequest("EUR_USD", S5).SetFrom(from).SetTo(to))
			return err
		}},
		{"account candles", func(c *Client) error {
			_, err := c.Price.Candlesticks(t.Context(), NewPriceCandlesticksRequest("EUR_USD", S5).SetFrom(from).SetTo(to))
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHTTPClient{}
			if err := tc.call(newFakeClient(fake)); err != nil {
				t.Fatal(err)
			}
			query := fake.requests[0].URL.Query()
			if query.Get("from") != "2026-09-13T10:00:00.100000001Z" || query.Get("to") != "2026-09-13T10:00:00.600000001Z" {
				t.Errorf("lost fractional seconds: %v", query)
			}
		})
	}
}

func TestOrderRequestStructLiterals(t *testing.T) {
	price := PriceValue("1.2")
	distance := DecimalNumber("0.01")
	cases := []struct {
		name string
		req  OrderRequest
		want string
	}{
		{"market", &MarketOrderRequest{Instrument: "EUR_USD", Units: "100"},
			`{"order":{"type":"MARKET","instrument":"EUR_USD","units":"100"}}`},
		{"limit", &LimitOrderRequest{Instrument: "EUR_USD", Units: "100", Price: "1.1"},
			`{"order":{"type":"LIMIT","instrument":"EUR_USD","units":"100","price":"1.1"}}`},
		{"stop", &StopOrderRequest{Instrument: "EUR_USD", Units: "100", Price: "1.1"},
			`{"order":{"type":"STOP","instrument":"EUR_USD","units":"100","price":"1.1"}}`},
		{"market if touched", &MarketIfTouchedOrderRequest{Instrument: "EUR_USD", Units: "100", Price: "1.1"},
			`{"order":{"type":"MARKET_IF_TOUCHED","instrument":"EUR_USD","units":"100","price":"1.1"}}`},
		{"take profit", &TakeProfitOrderRequest{TradeID: "42", Price: "1.2"},
			`{"order":{"type":"TAKE_PROFIT","tradeID":"42","price":"1.2"}}`},
		{"stop loss", &StopLossOrderRequest{TradeID: "42", Distance: &distance},
			`{"order":{"type":"STOP_LOSS","tradeID":"42","distance":"0.01"}}`},
		{"guaranteed stop loss", &GuaranteedStopLossOrderRequest{TradeID: "42", Price: &price},
			`{"order":{"type":"GUARANTEED_STOP_LOSS","tradeID":"42","price":"1.2"}}`},
		{"trailing stop loss", &TrailingStopLossOrderRequest{TradeID: "42", Distance: "0.01"},
			`{"order":{"type":"TRAILING_STOP_LOSS","tradeID":"42","distance":"0.01"}}`},
		{"dependent order details", &MarketOrderRequest{
			Instrument:             "EUR_USD",
			Units:                  "100",
			TakeProfitOnFill:       &TakeProfitDetails{Price: "1.2"},
			StopLossOnFill:         &StopLossDetails{Distance: &distance},
			TrailingStopLossOnFill: &TrailingStopLossDetails{Distance: "0.01"},
		}, `{"order":{"type":"MARKET","instrument":"EUR_USD","units":"100","takeProfitOnFill":{"price":"1.2"},"stopLossOnFill":{"distance":"0.01"},"trailingStopLossOnFill":{"distance":"0.01"}}}`},
		{"explicit type kept", &MarketOrderRequest{Type: OrderTypeMarket, Instrument: "EUR_USD", Units: "100", TimeInForce: TimeInForceIOC},
			`{"order":{"type":"MARKET","instrument":"EUR_USD","units":"100","timeInForce":"IOC"}}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHTTPClient{responses: []*http.Response{jsonResponse(http.StatusCreated, `{"lastTransactionID":"1"}`)}}
			if _, err := newFakeClient(fake).Order.Create(t.Context(), tc.req); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decodeTestJSON(t, fake.bodies[0]), decodeTestJSON(t, tc.want)) {
				t.Errorf("body = %s, want %s", fake.bodies[0], tc.want)
			}
		})
	}
	t.Run("request not modified", func(t *testing.T) {
		req := &MarketOrderRequest{Instrument: "EUR_USD", Units: "100"}
		fake := &fakeHTTPClient{responses: []*http.Response{jsonResponse(http.StatusCreated, `{"lastTransactionID":"1"}`)}}
		if _, err := newFakeClient(fake).Order.Create(t.Context(), req); err != nil {
			t.Fatal(err)
		}
		if req.Type != "" {
			t.Errorf("Type = %q, want the caller's request left unchanged", req.Type)
		}
	})
}

func TestTradeRequestStructLiterals(t *testing.T) {
	fake := &fakeHTTPClient{}
	c := newFakeClient(fake)
	if _, err := c.Trade.Close(t.Context(), "42", TradeCloseRequest{}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Trade.UpdateClientExtensions(t.Context(), "42", TradeUpdateClientExtensionsRequest{}); err != nil {
		t.Fatal(err)
	}
	for i, body := range fake.bodies {
		if body != "{}" {
			t.Errorf("request %d body = %s, want {} so the API defaults apply", i, body)
		}
	}
}

func TestCandlesticksRequestDefaults(t *testing.T) {
	from := mustTime(t, "2026-09-13T10:00:00Z")
	cases := []struct {
		name  string
		path  string
		query string
		call  func(*Client) error
	}{
		{"instrument struct literal", "/v3/instruments/EUR_USD/candles", "", func(c *Client) error {
			_, err := c.Instrument.Candlesticks(t.Context(), &CandlesticksRequest{Instrument: "EUR_USD"})
			return err
		}},
		{"account struct literal", testAccountPath + "/instruments/EUR_USD/candles", "", func(c *Client) error {
			_, err := c.Price.Candlesticks(t.Context(), &PriceCandlesticksRequest{CandlesticksRequest: CandlesticksRequest{Instrument: "EUR_USD"}})
			return err
		}},
		{"constructor without granularity", "/v3/instruments/EUR_USD/candles", "", func(c *Client) error {
			_, err := c.Instrument.Candlesticks(t.Context(), NewCandlesticksRequest("EUR_USD", ""))
			return err
		}},
		{"non-default values", "/v3/instruments/EUR_USD/candles", "granularity=H1&from=2026-09-13T10:00:00Z&includeFirst=False&weeklyAlignment=Monday", func(c *Client) error {
			_, err := c.Instrument.Candlesticks(t.Context(), &CandlesticksRequest{Instrument: "EUR_USD", Granularity: H1, From: &from, ExcludeFirst: true, WeeklyAlignment: WeeklyAlignmentMonday})
			return err
		}},
		{"exclude first setter", "/v3/instruments/EUR_USD/candles", "granularity=H1&from=2026-09-13T10:00:00Z&includeFirst=False", func(c *Client) error {
			_, err := c.Instrument.Candlesticks(t.Context(), NewCandlesticksRequest("EUR_USD", H1).SetFrom(from).SetExcludeFirst())
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHTTPClient{}
			if err := tc.call(newFakeClient(fake)); err != nil {
				t.Fatal(err)
			}
			req := fake.requests[0]
			wantQuery, err := url.ParseQuery(tc.query)
			if err != nil {
				t.Fatal(err)
			}
			if req.URL.Path != tc.path || !reflect.DeepEqual(req.URL.Query(), wantQuery) {
				t.Errorf("request = %s, want %s?%s", req.URL, tc.path, tc.query)
			}
		})
	}
}
