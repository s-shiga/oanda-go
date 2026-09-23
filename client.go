package oanda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

const (
	// Version is the version of the oanda-go library.
	Version = "0.1.1"
	// FXTradeURL is the base URL for the OANDA fxTrade REST API (live).
	FXTradeURL = "https://api-fxtrade.oanda.com"
	// FXTradePracticeURL is the base URL for the OANDA fxTrade REST API (practice/demo).
	FXTradePracticeURL = "https://api-fxpractice.oanda.com"
	// FXTradeStreamingURL is the base URL for the OANDA fxTrade Streaming API (live).
	FXTradeStreamingURL = "https://stream-fxtrade.oanda.com"
	// FXTradeStreamingPracticeURL is the base URL for the OANDA fxTrade Streaming API (practice/demo).
	FXTradeStreamingPracticeURL = "https://stream-fxpractice.oanda.com"
)

// HTTPClient is an interface for executing HTTP requests.
// It is satisfied by *http.Client and can be replaced for testing.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func defaultUserAgent() string {
	return fmt.Sprintf(
		"oanda-go/%s (%s; %s/%s)",
		Version,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}

// defaultStreamStallTimeout is three missed heartbeats: OANDA sends one every
// 5 seconds on both the price and transaction streams.
const defaultStreamStallTimeout = 15 * time.Second

type clientConfig struct {
	baseURL            string
	apiKey             string
	userAgent          string
	accountID          AccountID
	httpClient         HTTPClient
	streamStallTimeout time.Duration
}

// Client is the OANDA v20 REST API client. Create one with [NewClient] (live)
// or [NewDemoClient] (practice). Each field exposes a service that maps to an
// OANDA API endpoint group.
type Client struct {
	clientConfig
	Account     *accountService
	Instrument  *instrumentService
	Order       *orderService
	Trade       *tradeService
	Position    *positionService
	Transaction *transactionService
	Price       *priceService
}

// Option configures a [Client] or [StreamClient]. Pass options to
// [NewClient], [NewDemoClient], [NewStreamClient], or [NewDemoStreamClient].
type Option func(*clientConfig)

// WithBaseURL overrides the default OANDA API base URL. Any path in baseURL,
// such as a proxy prefix, is kept in front of each endpoint path.
func WithBaseURL(baseURL string) Option {
	return func(c *clientConfig) {
		c.baseURL = baseURL
	}
}

// WithUserAgent overrides the default User-Agent header sent with every request.
func WithUserAgent(userAgent string) Option {
	return func(c *clientConfig) {
		c.userAgent = userAgent
	}
}

// WithAccountID sets the default account ID used by account-scoped API calls.
func WithAccountID(id AccountID) Option {
	return func(c *clientConfig) {
		c.accountID = id
	}
}

// WithStreamStallTimeout sets how long a [StreamClient] waits for data (price
// updates, transactions, or heartbeats) before it treats the connection as
// dead, closes it, and returns [ErrStreamStalled]. The default is 15 seconds,
// three missed heartbeats. A timeout of zero or less disables the check. The
// connection is closed by cancelling the request's context, so a custom
// [HTTPClient] must honour it, as *http.Client does. It has no effect on a
// [Client].
func WithStreamStallTimeout(timeout time.Duration) Option {
	return func(c *clientConfig) {
		c.streamStallTimeout = timeout
	}
}

// WithHTTPClient replaces the default HTTP client used for API requests.
// Any implementation of [HTTPClient] (including *http.Client) is accepted,
// which allows injecting a fake transport in tests.
func WithHTTPClient(client HTTPClient) Option {
	return func(c *clientConfig) {
		c.httpClient = client
	}
}

func defaultConfig(baseURL, apiKey string) clientConfig {
	return clientConfig{
		baseURL:            baseURL,
		apiKey:             apiKey,
		userAgent:          defaultUserAgent(),
		accountID:          "",
		httpClient:         http.DefaultClient,
		streamStallTimeout: defaultStreamStallTimeout,
	}
}

func buildClient(baseURL, apiKey string) *Client {
	client := &Client{
		clientConfig: defaultConfig(baseURL, apiKey),
	}
	client.Account = newAccountService(client)
	client.Instrument = newInstrumentService(client)
	client.Order = newOrderService(client)
	client.Trade = newTradeService(client)
	client.Position = newPositionService(client)
	client.Transaction = newTransactionService(client)
	client.Price = newPriceService(client)
	return client
}

// NewClient creates a new OANDA v20 REST API client for the live environment.
func NewClient(apiKey string, opts ...Option) *Client {
	client := buildClient(FXTradeURL, apiKey)
	for _, opt := range opts {
		opt(&client.clientConfig)
	}
	return client
}

// NewDemoClient creates a new OANDA v20 REST API client for the practice/demo environment.
func NewDemoClient(apiKey string, opts ...Option) *Client {
	client := buildClient(FXTradePracticeURL, apiKey)
	for _, opt := range opts {
		opt(&client.clientConfig)
	}
	return client
}

// accountPath returns the escaped path of an endpoint under the Account set
// with [WithAccountID], built from the given segments.
func (c *clientConfig) accountPath(segments ...string) (string, error) {
	if c.accountID == "" {
		return "", ErrNoAccountID
	}
	return endpointPath(append([]string{"v3", "accounts", c.accountID}, segments...)...)
}

// endpointPath joins segments into an absolute path, escaping each one so an
// ID containing "/", "?" or "#" stays a single segment. An empty segment is an
// error: it would silently address a different endpoint.
func endpointPath(segments ...string) (string, error) {
	var b strings.Builder
	empty := false
	for _, segment := range segments {
		empty = empty || segment == ""
		b.WriteByte('/')
		b.WriteString(url.PathEscape(segment))
	}
	if empty {
		return "", fmt.Errorf("request path %s has an empty segment", b.String())
	}
	return b.String(), nil
}

// joinURL appends the escaped path to baseURL, keeping any path baseURL has.
func joinURL(baseURL string, path string, query url.Values) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	rawPath := strings.TrimSuffix(u.EscapedPath(), "/") + path
	if u.Path, err = url.PathUnescape(rawPath); err != nil {
		return "", err
	}
	u.RawPath = rawPath
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}
	return u.String(), nil
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("User-Agent", c.userAgent)
	req.Header.Add("Authorization", "Bearer "+c.apiKey)
}

// Request is implemented by types that can serialize themselves into an HTTP request body.
type Request interface {
	body() (*bytes.Buffer, error)
}

func (c *Client) sendRequest(ctx context.Context, method, path string, query url.Values, body io.Reader) (*http.Response, error) {
	u, err := joinURL(c.baseURL, path, query)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.setHeaders(req)
	return c.httpClient.Do(req)
}

func (c *Client) sendGetRequest(ctx context.Context, path string, values url.Values) (*http.Response, error) {
	return c.sendRequest(ctx, http.MethodGet, path, values, nil)
}

func (c *Client) sendPostRequest(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	return c.sendRequest(ctx, http.MethodPost, path, nil, body)
}

func (c *Client) sendPutRequest(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	return c.sendRequest(ctx, http.MethodPut, path, nil, body)
}

func (c *Client) sendPatchRequest(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	return c.sendRequest(ctx, http.MethodPatch, path, nil, body)
}

func doGet[R any](c *Client, ctx context.Context, path string, query url.Values) (*R, error) {
	httpResp, err := c.sendGetRequest(ctx, path, query)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}
	defer closeBody(httpResp)
	if httpResp.StatusCode != http.StatusOK {
		return nil, decodeErrorResponse(httpResp)
	}
	return decodeJSON[R](httpResp)
}

// StreamClient is the OANDA v20 Streaming API client. Create one with
// [NewStreamClient] (live) or [NewDemoStreamClient] (practice).
type StreamClient struct {
	clientConfig
}

func buildStreamClient(baseURL string, apiKey string) *StreamClient {
	client := &StreamClient{
		clientConfig: defaultConfig(baseURL, apiKey),
	}
	return client
}

// NewStreamClient creates a new OANDA v20 Streaming API client for the live environment.
func NewStreamClient(apiKey string, opts ...Option) *StreamClient {
	client := buildStreamClient(FXTradeStreamingURL, apiKey)
	for _, opt := range opts {
		opt(&client.clientConfig)
	}
	return client
}

// NewDemoStreamClient creates a new OANDA v20 Streaming API client for the practice/demo environment.
func NewDemoStreamClient(apiKey string, opts ...Option) *StreamClient {
	client := buildStreamClient(FXTradeStreamingPracticeURL, apiKey)
	for _, opt := range opts {
		opt(&client.clientConfig)
	}
	return client
}

func (c *StreamClient) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Add("User-Agent", c.userAgent)
	req.Header.Add("Authorization", "Bearer "+c.apiKey)
}

// closeBody closes a response body. By then the body has been decoded or its
// error returned, so a Close error is not actionable and is ignored rather
// than written to the application's logs.
func closeBody(resp *http.Response) {
	_ = resp.Body.Close()
}

// decodeJSON decodes an HTTP response body into a value of type R.
// It does not close the body; callers are responsible for that.
func decodeJSON[R any](resp *http.Response) (*R, error) {
	var v R
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &v, nil
}

// decodeTypedError decodes an HTTP error response body into an
// endpoint-specific error type E and wraps it in the error matching the
// response status code. A body that is not the documented JSON (such as an
// empty body or an HTML page from a proxy) is handled as in
// decodeErrorResponse, so the status code is never lost. It does not close
// the body; callers are responsible for that.
func decodeTypedError[E error](resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return wrapHTTPError(resp.StatusCode, fmt.Errorf("failed to read error response body: %w", err))
	}
	var e E
	if err := json.Unmarshal(body, &e); err != nil {
		return errorFromBody(resp.StatusCode, body)
	}
	return wrapHTTPError(resp.StatusCode, e)
}

func wrapHTTPError(statusCode int, err error) error {
	switch statusCode {
	case http.StatusBadRequest:
		return BadRequest{HTTPError{statusCode, "bad request", err}}
	case http.StatusUnauthorized:
		return Unauthorized{HTTPError{statusCode, "unauthorized", err}}
	case http.StatusForbidden:
		return Forbidden{HTTPError{statusCode, "forbidden", err}}
	case http.StatusNotFound:
		return NotFound{HTTPError{statusCode, "not found", err}}
	case http.StatusMethodNotAllowed:
		return MethodNotAllowed{HTTPError{statusCode, "method not allowed", err}}
	default:
		return HTTPError{statusCode, http.StatusText(statusCode), err}
	}
}

func decodeErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return wrapHTTPError(resp.StatusCode, fmt.Errorf("failed to read error response body: %w", err))
	}
	return errorFromBody(resp.StatusCode, body)
}

// errorFromBody builds the error for an HTTP error response, using OANDA's
// errorCode and errorMessage when the body has them and the raw text otherwise.
func errorFromBody(statusCode int, body []byte) error {
	errResp := struct {
		Code    string `json:"errorCode"`
		Message string `json:"errorMessage"`
	}{}
	if err := json.Unmarshal(body, &errResp); err != nil || errResp.Message == "" {
		// Non-JSON body (e.g. a gateway/HTML error page): keep the raw text.
		text := string(bytes.TrimSpace(body))
		if text == "" {
			text = "empty response body"
		}
		return wrapHTTPError(statusCode, errors.New(text))
	}
	if errResp.Code != "" {
		return wrapHTTPError(statusCode, fmt.Errorf("%s: %s", errResp.Code, errResp.Message))
	}
	return wrapHTTPError(statusCode, errors.New(errResp.Message))
}

// stallWatchdog aborts a stream's request when no data arrives for timeout,
// so a connection that died without closing surfaces as ErrStreamStalled
// instead of a read that blocks forever.
type stallWatchdog struct {
	timer   *time.Timer
	timeout time.Duration
	stalled atomic.Bool
}

// newStallWatchdog starts a countdown that calls abort after timeout unless it
// is stopped first. A timeout of zero or less disables the watchdog.
func newStallWatchdog(timeout time.Duration, abort context.CancelFunc) *stallWatchdog {
	w := &stallWatchdog{timeout: timeout}
	if timeout > 0 {
		w.timer = time.AfterFunc(timeout, func() {
			w.stalled.Store(true)
			abort()
		})
	}
	return w
}

// reset restarts the countdown from the full timeout.
func (w *stallWatchdog) reset() {
	if w.timer != nil {
		w.timer.Reset(w.timeout)
	}
}

// stop pauses the countdown until the next reset.
func (w *stallWatchdog) stop() {
	if w.timer != nil {
		w.timer.Stop()
	}
}

// streamLoop opens a streaming GET connection and decodes newline-delimited
// JSON objects until done is closed, the context is cancelled, the server
// ends the stream, or the stream stalls. Each object is passed to parse;
// items it accepts are sent to ch.
//
// ch is never closed by streamLoop; the caller detects the end of the stream
// by streamLoop returning. done is only checked between messages, so it
// cannot interrupt a read that is blocked waiting for data — cancel ctx to
// abort the connection reliably. When the server ends the stream, streamLoop
// returns ErrStreamEnded. When no data arrives for the stall timeout while
// connecting or waiting for the next message, it closes the connection and
// returns ErrStreamStalled; time spent waiting for the consumer to receive
// from ch does not count.
func streamLoop[T any](
	ctx context.Context,
	c *StreamClient,
	path string,
	values url.Values,
	ch chan<- T,
	done <-chan struct{},
	parse func(json.RawMessage) (T, bool, error),
) error {
	u, err := joinURL(c.baseURL, path, values)
	if err != nil {
		return err
	}
	// The watchdog cancels only the request; ctx still reports whether the
	// caller cancelled.
	reqCtx, abort := context.WithCancel(ctx)
	defer abort()
	watchdog := newStallWatchdog(c.streamStallTimeout, abort)
	defer watchdog.stop()
	httpReq, err := http.NewRequestWithContext(reqCtx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	c.setHeaders(httpReq)
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		if ctx.Err() == nil && watchdog.stalled.Load() {
			return ErrStreamStalled
		}
		return fmt.Errorf("failed to send GET request: %w", err)
	}
	defer closeBody(httpResp)
	if httpResp.StatusCode != http.StatusOK {
		return decodeErrorResponse(httpResp)
	}
	dec := json.NewDecoder(httpResp.Body)
	for {
		select {
		case <-done:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		var raw json.RawMessage
		watchdog.reset()
		err := dec.Decode(&raw)
		watchdog.stop()
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if watchdog.stalled.Load() {
				return ErrStreamStalled
			}
			if errors.Is(err, io.EOF) {
				return ErrStreamEnded
			}
			return fmt.Errorf("failed to decode JSON response: %w", err)
		}
		item, ok, err := parse(raw)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		select {
		case ch <- item:
		case <-done:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
