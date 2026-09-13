package oanda

import (
	"encoding/json"
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
