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
	Version                     = "0.1.0"
	FXTradeURL                  = "https://api-fxtrade.oanda.com"
	FXTradePracticeURL          = "https://api-fxpractice.oanda.com"
	FXTradeStreamingURL         = "https://stream-fxtrade.oanda.com"
	FXTradeStreamingPracticeURL = "https://stream-fxpractice.oanda.com"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	baseURL          string
	baseStreamingURL string
	apiKey           string
	userAgent        string
	accountID        AccountID
	httpClient       HTTPClient
	Account          *AccountService
	Instrument       *InstrumentService
	Order            *OrderService
	Trade            *TradeService
	Position         *PositionService
	Transaction      *TransactionService
	Price            *PriceService
}

type Option func(*Client)

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithBaseStreamingURL(baseStreamingURL string) Option {
	return func(c *Client) {
		c.baseStreamingURL = baseStreamingURL
	}
}

func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

func WithAccountID(id AccountID) Option {
	return func(c *Client) {
		c.accountID = id
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

func buildClient(baseURL, baseStreamingURL string, apiKey string) *Client {
	goVersion := runtime.Version()
	osArch := runtime.GOOS + "/" + runtime.GOARCH
	client := &Client{
		baseURL:          baseURL,
		baseStreamingURL: baseStreamingURL,
		apiKey:           apiKey,
		userAgent:        fmt.Sprintf("oanda-go/%s (%s; %s)", Version, goVersion, osArch),
		accountID:        "",
		httpClient:       http.DefaultClient,
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

func NewClient(apiKey string, opts ...Option) *Client {
	client := buildClient(FXTradeURL, FXTradeStreamingURL, apiKey)
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func NewDemoClient(apiKey string, opts ...Option) *Client {
	client := buildClient(FXTradePracticeURL, FXTradeStreamingPracticeURL, apiKey)
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *Client) getURL(path string, query url.Values) (string, error) {
	u, err := url.Parse(c.baseURL)
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

type Request interface {
	body() (*bytes.Buffer, error)
}

func (c *Client) sendGetRequest(ctx context.Context, path string, values url.Values) (*http.Response, error) {
	u, err := c.getURL(path, values)
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
	u, err := c.getURL(path, nil)
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
	u, err := c.getURL(path, nil)
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
	u, err := c.getURL(path, nil)
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

func doPatch[R any](c *Client, ctx context.Context, path string, req Request) (*R, error) {
	var body io.Reader
	var err error
	if req != nil {
		body, err = req.body()
		if err != nil {
			return nil, err
		}
	}
	httpResp, err := c.sendPatchRequest(ctx, path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to send PATCH request: %w", err)
	}
	var resp R
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
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
		panic(fmt.Errorf("failed to decode response body: %w", err))
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
