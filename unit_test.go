package oanda

// Unit tests that run without network access or credentials. They inject a
// fake HTTPClient to assert on the exact requests the library builds and to
// exercise decode and streaming paths against canned responses.

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// fakeHTTPClient records every request it receives and plays back canned
// responses in order. Responses beyond the configured list default to 200 {}.
type fakeHTTPClient struct {
	requests  []*http.Request
	bodies    []string
	responses []*http.Response
	calls     int
}

func (f *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	i := f.calls
	f.calls++
	f.requests = append(f.requests, req)
	var body string
	if req.Body != nil {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		body = string(b)
	}
	f.bodies = append(f.bodies, body)
	if i < len(f.responses) {
		return f.responses[i], nil
	}
	return jsonResponse(http.StatusOK, `{}`), nil
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func newFakeClient(fake *fakeHTTPClient) *Client {
	return NewClient("test-key", WithAccountID("101-001-1234567-001"), WithHTTPClient(fake))
}

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}

// --- Request construction ---

func TestTransactionListRequestURL(t *testing.T) {
	fake := &fakeHTTPClient{}
	client := newFakeClient(fake)
	req := NewTransactionListRequest().SetPageSize(50)
	if _, err := client.Transaction.List(t.Context(), req); err != nil {
		t.Fatalf("List: %v", err)
	}
	got := fake.requests[0].URL
	if got.Path != "/v3/accounts/101-001-1234567-001/transactions" {
		t.Errorf("path = %q", got.Path)
	}
	if got.Query().Get("pageSize") != "50" {
		t.Errorf("query = %q, want pageSize=50", got.RawQuery)
	}
	if auth := fake.requests[0].Header.Get("Authorization"); auth != "Bearer test-key" {
		t.Errorf("Authorization = %q", auth)
	}
}

func TestPriceStreamRequestSnapshotParam(t *testing.T) {
	values, err := NewPriceStreamRequest("EUR_USD").DisableSnapshot().values()
	if err != nil {
		t.Fatal(err)
	}
	if got := values.Get("snapshot"); got != "false" {
		t.Errorf("snapshot = %q, want %q (raw values: %v)", got, "false", values)
	}
}

func TestPriceInformationRequestSinceParam(t *testing.T) {
	since := mustTime(t, "2024-05-01T12:30:45.123456789Z")
	values, err := NewPriceInformationRequest().
		AddInstruments("EUR_USD").
		SetSince(DateTime{since}).
		values()
	if err != nil {
		t.Fatal(err)
	}
	if got := values.Get("since"); got != "2024-05-01T12:30:45.123456789Z" {
		t.Errorf("since = %q, want RFC3339Nano", got)
	}
}

func TestPriceInformationRequestZeroSince(t *testing.T) {
	values, err := NewPriceInformationRequest().
		AddInstruments("EUR_USD").
		SetSince(DateTime{}). // no panic, param omitted
		values()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := values["since"]; ok {
		t.Errorf("zero-value since should be omitted, got %q", values.Get("since"))
	}
}

func TestTradeUpdateClientExtensionsBody(t *testing.T) {
	id := ClientID("my-id")
	req := TradeUpdateClientExtensionsRequest{
		ClientExtensions: &ClientExtensions{ID: &id},
	}
	body, err := req.body()
	if err != nil {
		t.Fatal(err)
	}
	want := `{"clientExtensions":{"id":"my-id"}}`
	if got := body.String(); got != want {
		t.Errorf("body = %s, want %s", got, want)
	}
}

func TestMarketOrderSetTradeClientExtensions(t *testing.T) {
	ext := &ClientExtensions{}
	req := NewMarketOrderRequest("EUR_USD", "100").SetTradeClientExtensions(ext)
	if req.TradeClientExtensions != ext {
		t.Error("TradeClientExtensions not set")
	}
	if req.ClientExtensions != nil {
		t.Error("ClientExtensions must remain unset")
	}
}

func TestNewGuaranteedStopLossOrderRequest(t *testing.T) {
	req := NewGuaranteedStopLossOrderRequest("42", "1.2345")
	if req.Price == nil || *req.Price != "1.2345" {
		t.Fatalf("Price = %v, want 1.2345", req.Price)
	}
	if _, err := req.body(); err != nil {
		t.Errorf("body: %v", err)
	}
	req.SetDistance("0.005")
	if req.Price != nil {
		t.Error("SetDistance should clear Price")
	}
	if body, err := req.body(); err != nil {
		t.Errorf("body with distance: %v", err)
	} else if !strings.Contains(body.String(), `"distance":"0.005"`) {
		t.Errorf("body with distance = %s", body)
	}
	req.SetPrice("1.2500")
	if req.Distance != nil {
		t.Error("SetPrice should clear Distance")
	}
	if _, err := req.body(); err != nil {
		t.Errorf("body with price: %v", err)
	}
	distance := DecimalNumber("0.005")
	req.Distance = &distance
	if _, err := req.body(); err == nil {
		t.Error("body should reject price and distance both set")
	}
	req.Price = nil
	req.Distance = nil
	if _, err := req.body(); err == nil {
		t.Error("body should reject neither price nor distance set")
	}
}

func TestAccountConfigurePartialBody(t *testing.T) {
	fake := &fakeHTTPClient{responses: []*http.Response{
		jsonResponse(http.StatusOK, `{"clientConfigureTransaction":{},"lastTransactionID":"1"}`),
	}}
	client := newFakeClient(fake)
	req := NewAccountConfigureRequest().SetAlias("my-alias")
	if _, err := client.Account.Configure(t.Context(), req); err != nil {
		t.Fatalf("Configure: %v", err)
	}
	if method := fake.requests[0].Method; method != http.MethodPatch {
		t.Errorf("method = %s, want PATCH", method)
	}
	want := `{"alias":"my-alias"}`
	if got := fake.bodies[0]; got != want {
		t.Errorf("body = %s, want %s (marginRate must be omitted)", got, want)
	}
}

// --- Response decoding ---

func TestPositionGuaranteedExecutionFeesDecode(t *testing.T) {
	var p Position
	if err := json.Unmarshal([]byte(`{"guaranteedExecutionFees":"1.5"}`), &p); err != nil {
		t.Fatal(err)
	}
	if p.GuaranteedExecutionFees == nil || *p.GuaranteedExecutionFees != "1.5" {
		t.Errorf("GuaranteedExecutionFees = %v, want 1.5", p.GuaranteedExecutionFees)
	}
}

func TestClientPriceDecode(t *testing.T) {
	raw := `{"type":"PRICE","instrument":"EUR_USD","time":"2024-05-01T12:30:45.000000000Z","tradeable":true,"status":"tradeable"}`
	var p ClientPrice
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatal(err)
	}
	if p.Instrument != "EUR_USD" {
		t.Errorf("Instrument = %q", p.Instrument)
	}
	if !p.Tradeable {
		t.Error("Tradeable = false, want true")
	}
}

func TestDateTimeMarshalByValue(t *testing.T) {
	ts := mustTime(t, "2024-05-01T12:30:45.123456789Z")
	got, err := json.Marshal(struct {
		Time DateTime `json:"time"`
	}{Time: DateTime{ts}})
	if err != nil {
		t.Fatal(err)
	}
	want := `{"time":"2024-05-01T12:30:45.123456789Z"}`
	if string(got) != want {
		t.Errorf("marshal = %s, want %s", got, want)
	}
	got, err = json.Marshal(struct {
		Time DateTime `json:"time"`
	}{})
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"time":null}`; string(got) != want {
		t.Errorf("zero marshal = %s, want %s", got, want)
	}
}

func TestDateTimeUnmarshalZeroSentinel(t *testing.T) {
	var dt DateTime
	if err := json.Unmarshal([]byte(`"0"`), &dt); err != nil {
		t.Fatal(err)
	}
	if !dt.IsZero() {
		t.Errorf("Time = %v, want the zero time for sentinel \"0\"", dt.Time)
	}
}

func TestUnmarshalOrderUnknownType(t *testing.T) {
	raw := `{"id":"42","type":"SOMETHING_NEW","state":"PENDING","createTime":"2025-01-01T00:00:00Z","newField":"x"}`
	order, err := unmarshalOrder([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	unknown, ok := order.(UnknownOrder)
	if !ok {
		t.Fatalf("order = %T, want UnknownOrder", order)
	}
	if unknown.GetID() != "42" || unknown.GetType() != "SOMETHING_NEW" || unknown.GetState() != OrderStatePending || unknown.GetCreateTime().IsZero() {
		t.Errorf("common fields not decoded: %#v", unknown)
	}
	encoded, err := json.Marshal(unknown)
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != raw {
		t.Errorf("MarshalJSON = %s, want the original JSON %s", encoded, raw)
	}
	if _, err := unmarshalOrder([]byte(`{"id":"42"}`)); err == nil {
		t.Error("want error for an order without a type, got nil")
	}
}

func TestUnmarshalTransactionUnknownType(t *testing.T) {
	raw := `{"id":"42","type":"SOMETHING_NEW","time":"2025-01-01T00:00:00Z","accountID":"101-001-1234567-001","newField":"x"}`
	transaction, err := unmarshalTransaction([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	unknown, ok := transaction.(*UnknownTransaction)
	if !ok {
		t.Fatalf("transaction = %T, want *UnknownTransaction", transaction)
	}
	if unknown.GetID() != "42" || unknown.GetType() != "SOMETHING_NEW" || unknown.GetTime().IsZero() || unknown.AccountID != "101-001-1234567-001" {
		t.Errorf("common fields not decoded: %#v", unknown)
	}
	encoded, err := json.Marshal(unknown)
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != raw {
		t.Errorf("MarshalJSON = %s, want the original JSON %s", encoded, raw)
	}
	if _, err := unmarshalTransaction([]byte(`{"id":"42"}`)); err == nil {
		t.Error("want error for a transaction without a type, got nil")
	}
}

func TestBaseURLPathPrefix(t *testing.T) {
	for _, base := range []string{"https://proxy.example.com/oanda", "https://proxy.example.com/oanda/"} {
		t.Run(base, func(t *testing.T) {
			fake := &fakeHTTPClient{}
			c := NewClient("test-key", WithBaseURL(base), WithAccountID("101-001-1234567-001"), WithHTTPClient(fake))
			if _, err := c.Trade.List(t.Context(), NewTradeListRequest().SetInstrument("USD_JPY")); err != nil {
				t.Fatal(err)
			}
			want := "https://proxy.example.com/oanda" + testAccountPath + "/trades?instrument=USD_JPY"
			if got := fake.requests[0].URL.String(); got != want {
				t.Errorf("REST URL = %s, want %s", got, want)
			}

			streamFake := &fakeHTTPClient{responses: []*http.Response{jsonResponse(http.StatusOK, "")}}
			sc := NewStreamClient("test-key", WithBaseURL(base), WithAccountID("101-001-1234567-001"), WithHTTPClient(streamFake))
			if err := sc.Transaction(t.Context(), make(chan TransactionStreamItem), make(chan struct{})); !errors.Is(err, ErrStreamEnded) {
				t.Fatalf("err = %v, want ErrStreamEnded", err)
			}
			want = "https://proxy.example.com/oanda" + testAccountPath + "/transactions/stream"
			if got := streamFake.requests[0].URL.String(); got != want {
				t.Errorf("stream URL = %s, want %s", got, want)
			}
		})
	}
}

func TestNilListRequestsUseDefaults(t *testing.T) {
	cases := []struct {
		name string
		path string
		call func(*Client) error
	}{
		{"orders", testAccountPath + "/orders", func(c *Client) error {
			_, err := c.Order.List(t.Context(), nil)
			return err
		}},
		{"trades", testAccountPath + "/trades", func(c *Client) error {
			_, err := c.Trade.List(t.Context(), nil)
			return err
		}},
		{"transactions", testAccountPath + "/transactions", func(c *Client) error {
			_, err := c.Transaction.List(t.Context(), nil)
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHTTPClient{}
			if err := tc.call(newFakeClient(fake)); err != nil {
				t.Fatal(err)
			}
			if len(fake.requests) != 1 {
				t.Fatalf("got %d requests, want 1", len(fake.requests))
			}
			if req := fake.requests[0]; req.URL.Path != tc.path || req.URL.RawQuery != "" {
				t.Errorf("request = %s, want %s with no query", req.URL, tc.path)
			}
		})
	}
}

func TestNilRequestsRejected(t *testing.T) {
	var nilMarketOrder *MarketOrderRequest
	cases := []struct {
		name string
		call func(*Client) error
	}{
		{"order create", func(c *Client) error {
			_, err := c.Order.Create(t.Context(), nil)
			return err
		}},
		{"order create typed nil", func(c *Client) error {
			_, err := c.Order.Create(t.Context(), nilMarketOrder)
			return err
		}},
		{"order replace", func(c *Client) error {
			_, err := c.Order.Replace(t.Context(), "42", nil)
			return err
		}},
		{"order replace typed nil", func(c *Client) error {
			_, err := c.Order.Replace(t.Context(), "42", nilMarketOrder)
			return err
		}},
		{"trade dependent orders", func(c *Client) error {
			_, err := c.Trade.UpdateOrders(t.Context(), "42", nil)
			return err
		}},
		{"position close", func(c *Client) error {
			_, err := c.Position.Close(t.Context(), "USD_JPY", nil)
			return err
		}},
		{"transactions by ID range", func(c *Client) error {
			_, err := c.Transaction.GetByIDRange(t.Context(), nil)
			return err
		}},
		{"transactions since ID", func(c *Client) error {
			_, err := c.Transaction.GetBySinceID(t.Context(), nil)
			return err
		}},
		{"price information", func(c *Client) error {
			_, err := c.Price.Information(t.Context(), nil)
			return err
		}},
		{"latest candlesticks", func(c *Client) error {
			_, err := c.Price.LatestCandlesticks(t.Context(), nil)
			return err
		}},
		{"account candlesticks", func(c *Client) error {
			_, err := c.Price.Candlesticks(t.Context(), nil)
			return err
		}},
		{"instrument candlesticks", func(c *Client) error {
			_, err := c.Instrument.Candlesticks(t.Context(), nil)
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHTTPClient{}
			if err := tc.call(newFakeClient(fake)); !errors.Is(err, ErrNilRequest) {
				t.Fatalf("err = %v, want ErrNilRequest", err)
			}
			if fake.calls != 0 {
				t.Errorf("sent %d HTTP requests for a nil request", fake.calls)
			}
		})
	}
	t.Run("price stream", func(t *testing.T) {
		fake := &fakeHTTPClient{}
		err := newFakeStreamClient(fake).Price(t.Context(), nil, make(chan PriceStreamItem), make(chan struct{}))
		if !errors.Is(err, ErrNilRequest) {
			t.Fatalf("err = %v, want ErrNilRequest", err)
		}
		if fake.calls != 0 {
			t.Errorf("sent %d HTTP requests for a nil request", fake.calls)
		}
	})
}

func TestNoAccountID(t *testing.T) {
	fake := &fakeHTTPClient{}
	c := NewClient("test-key", WithHTTPClient(fake))
	for name, call := range map[string]func() error{
		"account details": func() error { _, err := c.Account.Details(t.Context()); return err },
		"account changes": func() error { _, err := c.Account.Changes(t.Context(), "40"); return err },
		"instrument list": func() error { _, err := c.Instrument.List(t.Context()); return err },
		"pending orders":  func() error { _, err := c.Order.ListPending(t.Context()); return err },
		"order create": func() error {
			_, err := c.Order.Create(t.Context(), NewMarketOrderRequest("USD_JPY", "100"))
			return err
		},
		"open trades":        func() error { _, err := c.Trade.ListOpen(t.Context()); return err },
		"positions":          func() error { _, err := c.Position.List(t.Context()); return err },
		"transaction detail": func() error { _, err := c.Transaction.Details(t.Context(), "42"); return err },
		"price information": func() error {
			_, err := c.Price.Information(t.Context(), NewPriceInformationRequest().AddInstruments("USD_JPY"))
			return err
		},
	} {
		if err := call(); !errors.Is(err, ErrNoAccountID) {
			t.Errorf("%s: err = %v, want ErrNoAccountID", name, err)
		}
	}
	if fake.calls != 0 {
		t.Errorf("sent %d HTTP requests without an account ID", fake.calls)
	}

	// Endpoints that are not scoped to an Account still work.
	if _, err := c.Account.List(t.Context()); err != nil {
		t.Errorf("account list: %v", err)
	}
	if _, err := c.Instrument.Candlesticks(t.Context(), NewCandlesticksRequest("USD_JPY", H1)); err != nil {
		t.Errorf("instrument candlesticks: %v", err)
	}

	stream := NewStreamClient("test-key", WithHTTPClient(fake))
	if err := stream.Transaction(t.Context(), make(chan TransactionStreamItem), make(chan struct{})); !errors.Is(err, ErrNoAccountID) {
		t.Errorf("transaction stream: err = %v, want ErrNoAccountID", err)
	}
	if err := stream.Price(t.Context(), NewPriceStreamRequest("USD_JPY"), make(chan PriceStreamItem), make(chan struct{})); !errors.Is(err, ErrNoAccountID) {
		t.Errorf("price stream: err = %v, want ErrNoAccountID", err)
	}
}

func TestPathSegmentsEscaped(t *testing.T) {
	fake := &fakeHTTPClient{responses: []*http.Response{jsonResponse(http.StatusOK, `{"order":{"id":"42","type":"LIMIT"}}`)}}
	c := newFakeClient(fake)
	if _, err := c.Order.Details(t.Context(), "@my/strategy?x=1#a b"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Trade.Close(t.Context(), "@my/trade", NewTradeCloseALLRequest()); err != nil {
		t.Fatal(err)
	}
	for i, want := range []string{
		testAccountPath + "/orders/@my%2Fstrategy%3Fx=1%23a%20b",
		testAccountPath + "/trades/@my%2Ftrade/close",
	} {
		if got := fake.requests[i].URL.EscapedPath(); got != want {
			t.Errorf("request %d path = %s, want %s", i, got, want)
		}
		if fake.requests[i].URL.RawQuery != "" || fake.requests[i].URL.Fragment != "" {
			t.Errorf("request %d = %s, want the ID kept in the path", i, fake.requests[i].URL)
		}
	}
}

func TestEmptyPathSegment(t *testing.T) {
	fake := &fakeHTTPClient{}
	c := newFakeClient(fake)
	for name, call := range map[string]func() error{
		"order details": func() error { _, err := c.Order.Details(t.Context(), ""); return err },
		"order cancel":  func() error { _, err := c.Order.Cancel(t.Context(), ""); return err },
		"trade details": func() error { _, err := c.Trade.Details(t.Context(), ""); return err },
		"position close": func() error {
			_, err := c.Position.Close(t.Context(), "", NewPositionCloseRequest().SetLongAll())
			return err
		},
		"transaction details": func() error { _, err := c.Transaction.Details(t.Context(), ""); return err },
		"instrument candlesticks": func() error {
			_, err := c.Instrument.Candlesticks(t.Context(), &CandlesticksRequest{Granularity: H1})
			return err
		},
	} {
		if err := call(); err == nil || !strings.Contains(err.Error(), "empty segment") {
			t.Errorf("%s: err = %v, want an empty segment error", name, err)
		}
	}
	if fake.calls != 0 {
		t.Errorf("sent %d HTTP requests with an empty path segment", fake.calls)
	}
}

// --- Error decoding ---

func TestDecodeErrorResponseWithErrorCode(t *testing.T) {
	resp := jsonResponse(http.StatusBadRequest, `{"errorCode":"INVALID_RANGE","errorMessage":"bad range"}`)
	err := decodeErrorResponse(resp)
	var badRequest BadRequest
	if !errors.As(err, &badRequest) {
		t.Fatalf("err = %T (%v), want BadRequest", err, err)
	}
	if !strings.Contains(err.Error(), "INVALID_RANGE: bad range") {
		t.Errorf("message %q should contain errorCode and errorMessage", err.Error())
	}
}

func TestDecodeErrorResponseNonJSON(t *testing.T) {
	resp := jsonResponse(http.StatusBadGateway, `<html>bad gateway</html>`)
	err := decodeErrorResponse(resp)
	var httpErr HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("err = %T (%v), want HTTPError", err, err)
	}
	if httpErr.StatusCode != http.StatusBadGateway {
		t.Errorf("StatusCode = %d, want 502", httpErr.StatusCode)
	}
	if !strings.Contains(err.Error(), "<html>bad gateway</html>") {
		t.Errorf("message %q should keep the raw body", err.Error())
	}
}

func TestWrapHTTPErrorKeepsUnmappedStatus(t *testing.T) {
	err := wrapHTTPError(http.StatusTooManyRequests, errors.New("rate limited"))
	var httpErr HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("err = %T, want HTTPError", err)
	}
	if httpErr.StatusCode != http.StatusTooManyRequests {
		t.Errorf("StatusCode = %d, want 429", httpErr.StatusCode)
	}
}

func TestTypedErrorNonJSONBody(t *testing.T) {
	cases := []struct {
		name   string
		status int
		call   func(*Client) error
	}{
		{"account configure", http.StatusForbidden, func(c *Client) error {
			_, err := c.Account.Configure(t.Context(), NewAccountConfigureRequest().SetAlias("alias"))
			return err
		}},
		{"order create", http.StatusBadRequest, func(c *Client) error {
			_, err := c.Order.Create(t.Context(), NewMarketOrderRequest("USD_JPY", "10000"))
			return err
		}},
		{"order replace", http.StatusBadRequest, func(c *Client) error {
			_, err := c.Order.Replace(t.Context(), "42", NewLimitOrderRequest("USD_JPY", "10000", "110.00"))
			return err
		}},
		{"order replace missing order", http.StatusNotFound, func(c *Client) error {
			_, err := c.Order.Replace(t.Context(), "42", NewLimitOrderRequest("USD_JPY", "10000", "110.00"))
			return err
		}},
		{"order cancel", http.StatusNotFound, func(c *Client) error {
			_, err := c.Order.Cancel(t.Context(), "42")
			return err
		}},
		{"order client extensions", http.StatusBadRequest, func(c *Client) error {
			_, err := c.Order.UpdateClientExtensions(t.Context(), "42", OrderUpdateClientExtensionsRequest{ClientExtensions: NewClientExtensions().SetID("order-id")})
			return err
		}},
		{"trade close", http.StatusBadRequest, func(c *Client) error {
			_, err := c.Trade.Close(t.Context(), "42", NewTradeCloseALLRequest())
			return err
		}},
		{"trade close missing trade", http.StatusNotFound, func(c *Client) error {
			_, err := c.Trade.Close(t.Context(), "42", NewTradeCloseALLRequest())
			return err
		}},
		{"trade client extensions", http.StatusNotFound, func(c *Client) error {
			_, err := c.Trade.UpdateClientExtensions(t.Context(), "42", TradeUpdateClientExtensionsRequest{ClientExtensions: NewClientExtensions().SetID("trade-id")})
			return err
		}},
		{"trade dependent orders", http.StatusBadRequest, func(c *Client) error {
			_, err := c.Trade.UpdateOrders(t.Context(), "42", &TradeUpdateOrdersRequest{CancelTakeProfit: true})
			return err
		}},
		{"position close", http.StatusBadRequest, func(c *Client) error {
			_, err := c.Position.Close(t.Context(), "USD_JPY", NewPositionCloseRequest().SetLongAll())
			return err
		}},
	}
	bodies := []struct {
		name string
		body string
		want string
	}{
		{"HTML page", "<html><body>Bad Gateway</body></html>\n", "<html><body>Bad Gateway</body></html>"},
		{"empty body", "", "empty response body"},
	}
	for _, tc := range cases {
		for _, body := range bodies {
			t.Run(tc.name+"/"+body.name, func(t *testing.T) {
				fake := &fakeHTTPClient{responses: []*http.Response{jsonResponse(tc.status, body.body)}}
				err := tc.call(newFakeClient(fake))
				if got := httpStatusOf(err); got != tc.status {
					t.Fatalf("error = %T (%v), want an HTTP %d error", err, err, tc.status)
				}
				if !strings.Contains(err.Error(), body.want) {
					t.Errorf("error = %q, want it to contain %q", err, body.want)
				}
			})
		}
	}
}

// httpStatusOf returns the status code of an error built by wrapHTTPError for
// the statuses that have typed endpoint errors, or 0 for any other error.
func httpStatusOf(err error) int {
	var badRequest BadRequest
	var forbidden Forbidden
	var notFound NotFound
	switch {
	case errors.As(err, &badRequest):
		return badRequest.StatusCode
	case errors.As(err, &forbidden):
		return forbidden.StatusCode
	case errors.As(err, &notFound):
		return notFound.StatusCode
	}
	return 0
}

// --- Streaming ---

func newFakeStreamClient(fake *fakeHTTPClient) *StreamClient {
	return NewStreamClient("test-key", WithAccountID("101-001-1234567-001"), WithHTTPClient(fake))
}

func TestPriceStreamAuthFailure(t *testing.T) {
	fake := &fakeHTTPClient{responses: []*http.Response{
		jsonResponse(http.StatusUnauthorized, `{"errorMessage":"Insufficient authorization to perform request."}`),
	}}
	sc := newFakeStreamClient(fake)
	ch := make(chan PriceStreamItem, 1)
	err := sc.Price(t.Context(), NewPriceStreamRequest("EUR_USD"), ch, make(chan struct{}))
	var unauthorized Unauthorized
	if !errors.As(err, &unauthorized) {
		t.Fatalf("err = %T (%v), want Unauthorized", err, err)
	}
}

func TestPriceStreamServerEnd(t *testing.T) {
	body := `{"type":"PRICE","instrument":"EUR_USD","time":"2024-05-01T12:30:45.000000000Z","tradeable":true}` + "\n" +
		`{"type":"HEARTBEAT","time":"2024-05-01T12:30:50.000000000Z"}` + "\n"
	fake := &fakeHTTPClient{responses: []*http.Response{jsonResponse(http.StatusOK, body)}}
	sc := newFakeStreamClient(fake)
	ch := make(chan PriceStreamItem, 4)
	err := sc.Price(t.Context(), NewPriceStreamRequest("EUR_USD"), ch, make(chan struct{}))
	if !errors.Is(err, ErrStreamEnded) {
		t.Fatalf("err = %v, want ErrStreamEnded", err)
	}
	price, ok := (<-ch).(ClientPrice)
	if !ok || price.Instrument != "EUR_USD" {
		t.Errorf("first item = %#v, want ClientPrice for EUR_USD", price)
	}
	if _, ok := (<-ch).(PricingHeartbeat); !ok {
		t.Error("second item should be a PricingHeartbeat")
	}
}

func TestPriceStreamSkipsUnknownTypes(t *testing.T) {
	body := `{"type":"SOMETHING_NEW"}` + "\n" +
		`{"type":"HEARTBEAT","time":"2024-05-01T12:30:50.000000000Z"}` + "\n"
	fake := &fakeHTTPClient{responses: []*http.Response{jsonResponse(http.StatusOK, body)}}
	sc := newFakeStreamClient(fake)
	ch := make(chan PriceStreamItem, 4)
	err := sc.Price(t.Context(), NewPriceStreamRequest("EUR_USD"), ch, make(chan struct{}))
	if !errors.Is(err, ErrStreamEnded) {
		t.Fatalf("err = %v, want ErrStreamEnded", err)
	}
	if got := len(ch); got != 1 {
		t.Errorf("received %d items, want 1 (unknown type skipped)", got)
	}
}

func TestPriceStreamDecimalLiquidity(t *testing.T) {
	body := `{"type":"PRICE","instrument":"XAU_USD","time":"2024-05-01T12:30:45.000000000Z","bids":[{"price":"2300.10","liquidity":"1.5"}],"asks":[{"price":"2300.50","liquidity":0.5}]}` + "\n" +
		`{"type":"HEARTBEAT","time":"2024-05-01T12:30:50.000000000Z"}` + "\n"
	fake := &fakeHTTPClient{responses: []*http.Response{jsonResponse(http.StatusOK, body)}}
	sc := newFakeStreamClient(fake)
	ch := make(chan PriceStreamItem, 4)
	err := sc.Price(t.Context(), NewPriceStreamRequest("XAU_USD"), ch, make(chan struct{}))
	if !errors.Is(err, ErrStreamEnded) {
		t.Fatalf("err = %v, want ErrStreamEnded", err)
	}
	price, ok := (<-ch).(ClientPrice)
	if !ok || len(price.Bids) != 1 || len(price.Asks) != 1 {
		t.Fatalf("first item = %#v, want ClientPrice with one bid and one ask", price)
	}
	if price.Bids[0].Liquidity != "1.5" || price.Asks[0].Liquidity != "0.5" {
		t.Errorf("liquidity = %q/%q, want 1.5/0.5", price.Bids[0].Liquidity, price.Asks[0].Liquidity)
	}
	if _, ok := (<-ch).(PricingHeartbeat); !ok {
		t.Error("second item should be a PricingHeartbeat")
	}
}

func TestPriceStreamDoneAlreadyClosed(t *testing.T) {
	pr, pw := io.Pipe()
	defer func() {
		if err := pw.Close(); err != nil {
			t.Error(err)
		}
	}()
	fake := &fakeHTTPClient{responses: []*http.Response{{StatusCode: http.StatusOK, Body: pr}}}
	sc := newFakeStreamClient(fake)
	ch := make(chan PriceStreamItem)
	done := make(chan struct{})
	close(done)
	if err := sc.Price(t.Context(), NewPriceStreamRequest("EUR_USD"), ch, done); err != nil {
		t.Fatalf("err = %v, want nil when done is closed", err)
	}
}

func TestPriceStreamContextCancel(t *testing.T) {
	pr, pw := io.Pipe()
	fake := &fakeHTTPClient{responses: []*http.Response{{StatusCode: http.StatusOK, Body: pr}}}
	sc := newFakeStreamClient(fake)
	ctx, cancel := context.WithCancel(t.Context())
	ch := make(chan PriceStreamItem)
	errCh := make(chan error, 1)
	go func() {
		errCh <- sc.Price(ctx, NewPriceStreamRequest("EUR_USD"), ch, make(chan struct{}))
	}()
	if _, err := pw.Write([]byte(`{"type":"HEARTBEAT","time":"2024-05-01T12:30:50.000000000Z"}` + "\n")); err != nil {
		t.Fatal(err)
	}
	<-ch
	cancel()
	// The real transport aborts the body read when the request context is
	// cancelled; the pipe stands in for that here.
	if err := pw.CloseWithError(context.Canceled); err != nil {
		t.Fatal(err)
	}
	if err := <-errCh; !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
}
