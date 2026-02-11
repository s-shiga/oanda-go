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
	Account     *AccountService
	Instrument  *InstrumentService
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
	if query != nil && len(query) > 0 {
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

func (c *Client) sendGetRequest(ctx context.Context, path string, values url.Values) (*http.Response, error) {
	u, err := joinURL(c.baseURL, path, values)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.setHeaders(req)
	return c.httpClient.Do(req)
}

func doGet[R any](c *Client, ctx context.Context, path string, query url.Values) (*R, error) {
	httpResp, err := c.sendGetRequest(ctx, path, query)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}
	var resp R
	if httpResp.StatusCode != http.StatusOK {
		return nil, decodeErrorResponse(httpResp)
	}
	if err := decodeResponse(httpResp, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &resp, nil
}

func (c *Client) sendPostRequest(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	u, err := joinURL(c.baseURL, path, nil)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", u, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.setHeaders(req)
	return c.httpClient.Do(req)
}

func (c *Client) sendPutRequest(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	u, err := joinURL(c.baseURL, path, nil)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", u, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.setHeaders(req)
	return c.httpClient.Do(req)
}

func (c *Client) sendPatchRequest(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	u, err := joinURL(c.baseURL, path, nil)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "PATCH", u, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.setHeaders(req)
	return c.httpClient.Do(req)
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

func decodeResponse(resp *http.Response, v any) error {
	defer closeBody(resp)
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusMethodNotAllowed:
		return decodeErrorResponse(resp)
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func decodeErrorResponse(resp *http.Response) error {
	errResp := struct {
		Message string `json:"errorMessage"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return fmt.Errorf("failed to decode error response body: %w", err)
	}
	err := errors.New(errResp.Message)
	switch resp.StatusCode {
	case http.StatusBadRequest:
		return BadRequest{HTTPError{resp.StatusCode, "bad request", err}}
	case http.StatusUnauthorized:
		return Unauthorized{HTTPError{resp.StatusCode, "unauthorized", err}}
	case http.StatusForbidden:
		return Forbidden{HTTPError{resp.StatusCode, "forbidden", err}}
	case http.StatusNotFound:
		return NotFoundError{HTTPError{resp.StatusCode, "not found", err}}
	case http.StatusMethodNotAllowed:
		return MethodNotAllowed{HTTPError{resp.StatusCode, "method not allowed", err}}
	default:
		return err
	}
}
