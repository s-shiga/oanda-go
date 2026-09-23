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

The snippets below assume a configured `client` and a `context.Context` named
`ctx`. Handle returned errors as shown in Quick Start.

### Client Initialization

```go
// Live environment
client := oanda.NewClient("YOUR_API_KEY", oanda.WithAccountID("your-account-id"))

// Demo/practice environment
client = oanda.NewDemoClient("YOUR_API_KEY", oanda.WithAccountID("your-account-id"))
```

#### Options

| Option | Description |
|--------|-------------|
| `WithAccountID(id)` | Set the default account ID for account-scoped calls |
| `WithHTTPClient(client)` | Replace the default HTTP client |
| `WithBaseURL(url)` | Override the default API base URL |
| `WithUserAgent(ua)` | Override the default User-Agent header |
| `WithStreamStallTimeout(d)` | How long a stream may go without data, even a heartbeat, before it is closed as stalled (default 15s; 0 disables) |

### Orders

```go
// Place a market order
marketReq := oanda.NewMarketOrderRequest("EUR_USD", "10000")
marketResp, err := client.Order.Create(ctx, marketReq)

// Place a limit order
limitReq := oanda.NewLimitOrderRequest("EUR_USD", "10000", "1.2500")
limitResp, err := client.Order.Create(ctx, limitReq)

// List pending orders
orders, err := client.Order.ListPending(ctx)

// Replace a pending order
replaceReq := oanda.NewLimitOrderRequest("EUR_USD", "10000", "1.2600")
replaceResp, err := client.Order.Replace(ctx, oanda.OrderSpecifier("123"), replaceReq)

// Cancel an order
cancelResp, err := client.Order.Cancel(ctx, oanda.OrderSpecifier("123"))
```

### Trades

```go
// List open trades
trades, err := client.Trade.ListOpen(ctx)

// Get trade details
trade, err := client.Trade.Details(ctx, "123")

// Close a trade (fully or partially)
resp, err := client.Trade.Close(ctx, "123", oanda.NewTradeCloseALLRequest())

// Update dependent orders on a trade
req := &oanda.TradeUpdateOrdersRequest{
	TakeProfit: oanda.NewTakeProfitDetails("1.3000"),
	StopLoss:   oanda.NewStopLossDetails().SetPrice("1.2000"),
}
updateResp, err := client.Trade.UpdateOrders(ctx, "123", req)

// Cancel the take-profit order, leaving other dependent orders unchanged
cancelReq := &oanda.TradeUpdateOrdersRequest{CancelTakeProfit: true}
cancelResp, err := client.Trade.UpdateOrders(ctx, "123", cancelReq)
```

### Positions

```go
// List open positions
positions, err := client.Position.ListOpen(ctx)

// Close a position
req := oanda.NewPositionCloseRequest().SetLongAll()
resp, err := client.Position.Close(ctx, "EUR_USD", req)
```

### Pricing and Candlesticks

```go
// Get current prices
priceReq := oanda.NewPriceInformationRequest().AddInstruments("EUR_USD", "USD_JPY")
prices, err := client.Price.Information(ctx, priceReq)

// Get candlestick data
req := oanda.NewPriceCandlesticksRequest("EUR_USD", oanda.H1).SetCount(100)
candles, err := client.Price.Candlesticks(ctx, req)
```

### Instruments

```go
// List available instruments
instruments, err := client.Instrument.List(ctx)

// Get candlesticks for an instrument
req := oanda.NewCandlesticksRequest("EUR_USD", oanda.D).SetCount(30)
candles, err := client.Instrument.Candlesticks(ctx, req)
```

### Transactions

```go
// List transaction IDs
req := oanda.NewTransactionListRequest().SetFrom(from).SetTo(to)
txns, err := client.Transaction.List(ctx, req)

// Get transaction details
txn, err := client.Transaction.Details(ctx, "6356")
```

### Streaming

OANDA sends a heartbeat every 5 seconds and routinely drops idle or slow
connections. A stream method returns `ErrStreamEnded` when the server closes
the connection and `ErrStreamStalled` when no data has arrived within the
stall timeout; reconnect after either.

```go
// Stream prices
streamClient := oanda.NewDemoStreamClient(
	"YOUR_API_KEY",
	oanda.WithAccountID("your-account-id"),
)

ch := make(chan oanda.PriceStreamItem)
done := make(chan struct{})

go func() {
	defer close(ch) // The producer owns channel closure; the client does not close it.
	for {
		err := streamClient.Price(ctx, oanda.NewPriceStreamRequest("EUR_USD"), ch, done)
		if errors.Is(err, oanda.ErrStreamEnded) || errors.Is(err, oanda.ErrStreamStalled) {
			time.Sleep(time.Second)
			continue // reconnect
		}
		if err != nil && ctx.Err() == nil {
			log.Println("price stream:", err)
		}
		return
	}
}()

for item := range ch {
	switch v := item.(type) {
	case *oanda.ClientPrice:
		if len(v.Bids) > 0 && len(v.Asks) > 0 {
			fmt.Printf("Bid: %s Ask: %s\n", v.Bids[0].Price, v.Asks[0].Price)
		}
	case *oanda.PricingHeartbeat:
		fmt.Println("Heartbeat:", v.Time)
	}
}
```

```go
// Stream transactions
ch := make(chan oanda.TransactionStreamItem)
done := make(chan struct{})

go func() {
	defer close(ch)
	for {
		err := streamClient.Transaction(ctx, ch, done)
		if errors.Is(err, oanda.ErrStreamEnded) || errors.Is(err, oanda.ErrStreamStalled) {
			time.Sleep(time.Second)
			continue // reconnect
		}
		if err != nil && ctx.Err() == nil {
			log.Println("transaction stream:", err)
		}
		return
	}
}()

for item := range ch {
	fmt.Printf("%s: %s\n", item.GetID(), item.GetType())
}
```

## API Coverage

| Service | Endpoints |
|---------|-----------|
| Account | List, Details, Summary, Configure, Changes |
| Order | Create, List, ListPending, Details, Replace, Cancel, UpdateClientExtensions |
| Trade | List, ListOpen, Details, Close, UpdateClientExtensions, UpdateOrders |
| Position | List, ListOpen, ListByInstrument, Close |
| Pricing | Information, Candlesticks, LatestCandlesticks, Stream |
| Instrument | List, Candlesticks |
| Transaction | List, Details, GetByIDRange, GetBySinceID, Stream |

## Testing

All tests run offline using an injected fake HTTP client. No API credentials or
OANDA connection are required, and the test suite disables the default network
transport to catch accidental API calls. Tests check request construction,
response decoding, errors, and streaming with fixed responses.

```sh
go test -race ./...
go vet ./...
```

## Disclaimer

This library is not affiliated with, endorsed by, or sponsored by OANDA Corporation. Use of this software is at your own risk. The authors and contributors are not responsible for any financial losses incurred through the use of this library.

## License

[MIT](LICENSE)
