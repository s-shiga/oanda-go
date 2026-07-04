package oanda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"runtime"
)

const (
	// Version is the version of the oanda-go library.
	Version = "0.1.0"
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

type clientConfig struct {
	baseURL    string
	apiKey     string
	userAgent  string
	accountID  AccountID
	httpClient HTTPClient
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

// WithBaseURL overrides the default OANDA API base URL.
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

// WithHTTPClient replaces the default HTTP client used for API requests.
func WithHTTPClient(client *http.Client) Option {
	return func(c *clientConfig) {
		c.httpClient = client
	}
}

func defaultConfig(baseURL, apiKey string) clientConfig {
	return clientConfig{
		baseURL:    baseURL,
		apiKey:     apiKey,
		userAgent:  defaultUserAgent(),
		accountID:  "",
		httpClient: http.DefaultClient,
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

func joinURL(baseURL string, path string, query url.Values) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	u.Path = path
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

func closeBody(resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		slog.Error(err.Error())
	}
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
// response status code. It does not close the body; callers are
// responsible for that.
func decodeTypedError[E error](resp *http.Response) error {
	var e E
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
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
		return err
	}
}

func decodeErrorResponse(resp *http.Response) error {
	errResp := struct {
		Message string `json:"errorMessage"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return fmt.Errorf("failed to decode error response body: %w", err)
	}
	return wrapHTTPError(resp.StatusCode, errors.New(errResp.Message))
}

// streamLoop opens a streaming GET connection and decodes newline-delimited
// JSON objects until done is closed, the context is cancelled, or the server
// ends the stream. Each object is passed to parse; items it accepts are sent
// to ch.
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
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	c.setHeaders(httpReq)
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %w", err)
	}
	defer closeBody(httpResp)
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
		if err := dec.Decode(&raw); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
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
