package oanda

import (
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

func (c *Client) sendGetRequest(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.URL+path, nil)
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

func decodeError(resp *http.Response) error {
	errResp := struct {
		Message string `json:"errorMessage"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		panic(fmt.Errorf("failed to decode response body: %w", err))
	}
	switch resp.StatusCode {
	case http.StatusBadRequest:
		return BadRequest{Code: resp.StatusCode, Err: errors.New(errResp.Message)}
	case http.StatusUnauthorized:
		return Unauthorized{Code: resp.StatusCode, Err: errors.New(errResp.Message)}
	case http.StatusForbidden:
		return Forbidden{Code: resp.StatusCode, Err: errors.New(errResp.Message)}
	case http.StatusNotFound:
		return NotFoundError{Code: resp.StatusCode, Err: errors.New(errResp.Message)}
	case http.StatusMethodNotAllowed:
		return MethodNotAllowed{Code: resp.StatusCode, Err: errors.New(errResp.Message)}
	default:
		panic(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}
}
