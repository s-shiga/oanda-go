# oanda-go

[![Go Reference](https://pkg.go.dev/badge/github.com/s-shiga/oanda-go.svg)](https://pkg.go.dev/github.com/s-shiga/oanda-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go client library for the [OANDA v20 REST and Streaming API](https://developer.oanda.com/rest-live-v20/introduction/).

## Features

- Full coverage of the OANDA v20 REST API (accounts, orders, trades, positions, pricing, instruments, transactions)
- Real-time streaming for prices and transactions
- All 8 order types (Market, Limit, Stop, MarketIfTouched, TakeProfit, StopLoss, GuaranteedStopLoss, TrailingStopLoss)
- Builder pattern for constructing requests
- Context support for cancellation and timeouts
- Live and demo/practice environment support

## Installation

```sh
go get github.com/s-shiga/oanda-go
```

Requires Go 1.24.2 or later.

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	oanda "github.com/s-shiga/oanda-go"
)

func main() {
	client := oanda.NewDemoClient(
		"YOUR_API_KEY",
		oanda.WithAccountID("101-001-1234567-001"),
	)

	ctx := context.Background()

	// Get account summary
	summary, err := client.Account.Summary(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Balance:", summary.Account.Balance)

	// List open trades
	trades, err := client.Trade.ListOpen(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range trades.Trades {
		fmt.Printf("Trade %s: %s %s units\n", t.ID, t.Instrument, t.CurrentUnits)
	}
}
```

## Usage

### Client Initialization

```go
// Live environment
client := oanda.NewClient("YOUR_API_KEY", oanda.WithAccountID("your-account-id"))

// Demo/practice environment
client := oanda.NewDemoClient("YOUR_API_KEY", oanda.WithAccountID("your-account-id"))
```

#### Options

| Option | Description |
|--------|-------------|
| `WithAccountID(id)` | Set the default account ID for account-scoped calls |
| `WithHTTPClient(client)` | Replace the default HTTP client |
| `WithBaseURL(url)` | Override the default API base URL |
| `WithUserAgent(ua)` | Override the default User-Agent header |

### Orders

```go
// Place a market order
req := oanda.NewMarketOrderRequest("EUR_USD", "10000")
resp, err := client.Order.Create(ctx, req)

// Place a limit order
req := oanda.NewLimitOrderRequest("EUR_USD", "10000", "1.2500")
resp, err := client.Order.Create(ctx, req)

// List pending orders
orders, err := client.Order.ListPending(ctx)

// Cancel an order
resp, err := client.Order.Cancel(ctx, oanda.OrderSpecifier("123"))
```

### Trades

```go
// List open trades
trades, err := client.Trade.ListOpen(ctx)

// Get trade details
trade, err := client.Trade.Details(ctx, "123")

// Update dependent orders on a trade
req := oanda.NewTradeUpdateOrdersRequest().
	WithTakeProfit(oanda.NewTakeProfitDetails("1.3000")).
	WithStopLoss(oanda.NewStopLossDetails("1.2000"))
resp, err := client.Trade.UpdateOrders(ctx, "123", req)
```

### Positions

```go
// List open positions
positions, err := client.Position.ListOpen(ctx)

// Close a position
req := oanda.NewPositionCloseRequest().WithLongUnits("ALL")
resp, err := client.Position.Close(ctx, "EUR_USD", req)
```

### Pricing and Candlesticks

```go
// Get current prices
req := oanda.NewPriceInformationRequest("EUR_USD", "USD_JPY")
prices, err := client.Price.Information(ctx, req)

// Get candlestick data
req := oanda.NewPriceCandlesticksRequest("EUR_USD").
	WithGranularity(oanda.H1).
	WithCount(100)
candles, err := client.Price.Candlesticks(ctx, req)
```

### Instruments

```go
// List available instruments
instruments, err := client.Instrument.List(ctx)

// Get candlesticks for an instrument
req := oanda.NewCandlesticksRequest("EUR_USD").
	WithGranularity(oanda.D).
	WithCount(30)
candles, err := client.Instrument.Candlesticks(ctx, req)
```

### Transactions

```go
// List transaction IDs
req := oanda.NewTransactionListRequest(from, to)
txns, err := client.Transaction.List(ctx, req)

// Get transaction details
txn, err := client.Transaction.Details(ctx, "6356")
```

### Streaming

```go
// Stream prices
streamClient := oanda.NewDemoStreamClient(
	"YOUR_API_KEY",
	oanda.WithAccountID("your-account-id"),
)

ch := make(chan oanda.PriceStreamItem)
done := make(chan struct{})

go func() {
	err := streamClient.Price(ctx, oanda.NewPriceStreamRequest("EUR_USD"), ch, done)
	if err != nil {
		log.Fatal(err)
	}
}()

for item := range ch {
	switch v := item.(type) {
	case oanda.ClientPrice:
		fmt.Printf("Bid: %s Ask: %s\n", v.Bids[0].Price, v.Asks[0].Price)
	case oanda.PricingHeartbeat:
		fmt.Println("Heartbeat:", v.Time)
	}
}
```

```go
// Stream transactions
ch := make(chan oanda.TransactionStreamItem)
done := make(chan struct{})

go func() {
	err := streamClient.Transaction(ctx, ch, done)
	if err != nil {
		log.Fatal(err)
	}
}()
```

## API Coverage

| Service | Endpoints |
|---------|-----------|
| Account | List, Details, Summary, Configure, Changes |
| Order | Create, List, ListPending, Details, Cancel, UpdateClientExtensions |
| Trade | List, ListOpen, Details, UpdateClientExtensions, UpdateOrders |
| Position | List, ListOpen, ListByInstrument, Close |
| Pricing | Information, Candlesticks, LatestCandlesticks, Stream |
| Instrument | List, Candlesticks |
| Transaction | List, Details, GetByIDRange, GetBySinceID, Stream |

## Testing

Tests run against the OANDA demo environment. Set the following environment variables:

```sh
export OANDA_API_KEY_DEMO="your-demo-api-key"
export OANDA_ACCOUNT_ID_DEMO="your-demo-account-id"
go test ./...
```

## Disclaimer

This library is not affiliated with, endorsed by, or sponsored by OANDA Corporation. Use of this software is at your own risk. The authors and contributors are not responsible for any financial losses incurred through the use of this library.

## License

[MIT](LICENSE)
