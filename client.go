package oanda

import (
	"errors"
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
	}, nil
}
