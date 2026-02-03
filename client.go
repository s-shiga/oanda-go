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
	"os"
	"runtime"
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
	apiKey, ok := os.LookupEnv("OANDA_API_KEY_DEMO")
	if !ok {
		return nil, errors.New("OANDA_API_KEY_DEMO not set")
	}
	return &Client{
		URL:          fxTradePracticeURL,
		StreamingURL: fxTradeStreamingPracticeURL,
		APIKey:       apiKey,
		HTTPClient:   &http.Client{},
	}, nil
}

func (c *Client) getURL(path string, query url.Values) (string, error) {
	u, err := url.Parse(c.URL)
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
	goVersion := runtime.Version()
	osArch := runtime.GOOS + "/" + runtime.GOARCH
	req.Header.Add("User-Agent", fmt.Sprintf("github.com/S-Shiga/oanda-go (%s; %s)", goVersion, osArch))
	req.Header.Add("Authorization", "Bearer "+c.APIKey)
}

type Request interface {
	Body() (*bytes.Buffer, error)
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
	return c.HTTPClient.Do(req)
}

func doGet[R any](c *Client, ctx context.Context, path string, query url.Values) (*R, error) {
	httpResp, err := c.sendGetRequest(ctx, path, query)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}
	var resp R
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
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
	return c.HTTPClient.Do(req)
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
	return c.HTTPClient.Do(req)
}

func doPut[R any](c *Client, ctx context.Context, path string, req Request) (*R, error) {
	var body *bytes.Buffer
	var err error
	if req != nil {
		body, err = req.Body()
		if err != nil {
			return nil, fmt.Errorf("failed to marshall request body: %w", err)
		}
	}
	httpResp, err := c.sendPutRequest(ctx, path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to send PUT request: %w", err)
	}
	var resp *R
	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return resp, nil
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
	return c.HTTPClient.Do(req)
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
