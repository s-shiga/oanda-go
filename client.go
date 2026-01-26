package oanda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

const (
	fxTradeURL                  = "https://api-fxtrade.oanda.com"
	fxTradePracticeURL          = "https://api-fxpractice.oanda.com"
	fxTradeStreamingURL         = "https://stream-fxtrade.oanda.com"
	fxTradeStreamingPracticeURL = "https://stream-fxpractice.oanda.com"
)

type Client struct {
	URL          string
	StreamingURL string
	APIKey       string
	HTTPClient   *http.Client
}

func NewClient() (*Client, error) {
	apiKey, ok := os.LookupEnv("OANDA_API_KEY")
	if !ok {
		return nil, errors.New("OANDA_API_KEY not set")
	}
	return &Client{
		URL:          fxTradeURL,
		StreamingURL: fxTradeStreamingURL,
		APIKey:       apiKey,
		HTTPClient:   &http.Client{},
	}, nil
}

func NewPracticeClient() (*Client, error) {
	apiKey, ok := os.LookupEnv("OANDA_API_KEY")
	if !ok {
		return nil, errors.New("OANDA_API_KEY not set")
	}
	return &Client{
		URL:          fxTradePracticeURL,
		StreamingURL: fxTradeStreamingPracticeURL,
		APIKey:       apiKey,
		HTTPClient:   &http.Client{},
	}, nil
}

func (c *Client) sendGetRequest(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.URL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+c.APIKey)
	return c.HTTPClient.Do(req)
}

func closeBody(resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		slog.Error(err.Error())
	}
}

func decodeResponse(resp *http.Response, v any) error {
	defer closeBody(resp)
	if resp.StatusCode != http.StatusOK {
		return decodeErrorResponse(resp)
	}
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return err
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
