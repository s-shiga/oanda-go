package oanda

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestStreamClientTransaction(t *testing.T) {
	cases := []struct {
		name      string
		status    int
		body      string
		wantTypes []TransactionType
		wantErr   error
	}{
		{
			name:   "transaction and heartbeat",
			status: http.StatusOK,
			body: "{\"id\":\"42\",\"type\":\"ORDER_FILL\",\"instrument\":\"USD_JPY\",\"units\":\"10000\",\"time\":\"2025-01-01T00:00:00Z\"}\n" +
				"{\"type\":\"HEARTBEAT\",\"lastTransactionID\":\"42\",\"time\":\"2025-01-01T00:00:05Z\"}\n",
			wantTypes: []TransactionType{TransactionTypeOrderFill, "HEARTBEAT"},
			wantErr:   ErrStreamEnded,
		},
		{
			name:      "unknown transaction type",
			status:    http.StatusOK,
			body:      "{\"id\":\"42\",\"type\":\"SOMETHING_NEW\",\"time\":\"2025-01-01T00:00:00Z\",\"newField\":\"x\"}\n",
			wantTypes: []TransactionType{"SOMETHING_NEW"},
			wantErr:   ErrStreamEnded,
		},
		{name: "empty stream", status: http.StatusOK, wantErr: ErrStreamEnded},
		{name: "unauthorized", status: http.StatusUnauthorized, body: `{"errorMessage":"invalid token"}`},
		{name: "malformed JSON", status: http.StatusOK, body: `{"type":`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHTTPClient{responses: []*http.Response{jsonResponse(tc.status, tc.body)}}
			client := newFakeStreamClient(fake)
			ctx, cancel := context.WithTimeout(t.Context(), time.Second)
			defer cancel()
			ch := make(chan TransactionStreamItem, 4)
			err := client.Transaction(ctx, ch, make(chan struct{}))
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err == nil {
				t.Fatal("expected an error")
			}
			if tc.status == http.StatusUnauthorized {
				var unauthorized Unauthorized
				if !errors.As(err, &unauthorized) {
					t.Fatalf("error = %T, want Unauthorized", err)
				}
			}
			if len(ch) != len(tc.wantTypes) {
				t.Fatalf("got %d items, want %d", len(ch), len(tc.wantTypes))
			}
			for _, want := range tc.wantTypes {
				item := <-ch
				if item.GetType() != want || item.GetID() != "42" || item.GetTime().IsZero() {
					t.Errorf("unexpected stream item: %#v", item)
				}
				if fill, ok := item.(OrderFillTransaction); ok && (fill.Instrument != "USD_JPY" || fill.Units != "10000") {
					t.Errorf("unexpected fill: %#v", fill)
				}
				if unknown, ok := item.(UnknownTransaction); ok && !strings.Contains(string(unknown.Raw), `"newField":"x"`) {
					t.Errorf("unknown transaction lost its raw JSON: %s", unknown.Raw)
				}
			}
			if len(fake.requests) != 1 {
				t.Fatalf("got %d requests, want 1", len(fake.requests))
			}
			req := fake.requests[0]
			if req.Method != http.MethodGet || req.URL.Path != testAccountPath+"/transactions/stream" || req.URL.RawQuery != "" {
				t.Errorf("unexpected streaming request: %s %s", req.Method, req.URL)
			}
			if req.Header.Get("Authorization") != "Bearer test-key" {
				t.Error("missing stream authorization")
			}
		})
	}
}

const testHeartbeat = `{"type":"HEARTBEAT","time":"2025-01-01T00:00:05Z"}` + "\n"

// stallingHTTPClient answers a stream request with body and then sends
// nothing more, like a connection that died without being closed. As with
// net/http, cancelling the request's context ends the pending read.
type stallingHTTPClient struct{ body string }

func (s stallingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	pr, pw := io.Pipe()
	go func() {
		if s.body != "" {
			if _, err := io.WriteString(pw, s.body); err != nil {
				return
			}
		}
		<-req.Context().Done()
		_ = pw.CloseWithError(req.Context().Err())
	}()
	return &http.Response{StatusCode: http.StatusOK, Body: pr, Header: http.Header{}}, nil
}

// hangingHTTPClient never answers, like a server that accepts the connection
// but never sends response headers.
type hangingHTTPClient struct{}

func (hangingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	<-req.Context().Done()
	return nil, req.Context().Err()
}

func newStallTestClient(httpClient HTTPClient, timeout time.Duration) *StreamClient {
	return NewStreamClient("test-key", WithAccountID("101-001-1234567-001"), WithHTTPClient(httpClient), WithStreamStallTimeout(timeout))
}

// stallTestContext bounds a test that expects a stall, so a broken watchdog
// fails with context.DeadlineExceeded instead of hanging the test run.
func stallTestContext(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func TestStreamStalled(t *testing.T) {
	t.Run("silent after a heartbeat", func(t *testing.T) {
		sc := newStallTestClient(stallingHTTPClient{body: testHeartbeat}, 50*time.Millisecond)
		ch := make(chan PriceStreamItem, 4)
		err := sc.Price(stallTestContext(t), NewPriceStreamRequest("EUR_USD"), ch, make(chan struct{}))
		if !errors.Is(err, ErrStreamStalled) {
			t.Fatalf("err = %v, want ErrStreamStalled", err)
		}
		if len(ch) != 1 {
			t.Errorf("got %d items before the stall, want the heartbeat", len(ch))
		}
	})
	t.Run("no response headers", func(t *testing.T) {
		sc := newStallTestClient(hangingHTTPClient{}, 50*time.Millisecond)
		err := sc.Transaction(stallTestContext(t), make(chan TransactionStreamItem), make(chan struct{}))
		if !errors.Is(err, ErrStreamStalled) {
			t.Fatalf("err = %v, want ErrStreamStalled", err)
		}
	})
	t.Run("caller cancels first", func(t *testing.T) {
		sc := newStallTestClient(stallingHTTPClient{body: testHeartbeat}, time.Minute)
		ctx, cancel := context.WithCancel(t.Context())
		ch := make(chan PriceStreamItem)
		errCh := make(chan error, 1)
		go func() {
			errCh <- sc.Price(ctx, NewPriceStreamRequest("EUR_USD"), ch, make(chan struct{}))
		}()
		<-ch
		cancel()
		if err := <-errCh; !errors.Is(err, context.Canceled) {
			t.Fatalf("err = %v, want context.Canceled", err)
		}
	})
}

func TestStreamNotStalledWhileDataFlows(t *testing.T) {
	// Heartbeats keep arriving for longer than the timeout, so each one must
	// restart the countdown.
	pr, pw := io.Pipe()
	go func() {
		for range 30 {
			time.Sleep(10 * time.Millisecond)
			if _, err := io.WriteString(pw, testHeartbeat); err != nil {
				return
			}
		}
		_ = pw.Close()
	}()
	fake := &fakeHTTPClient{responses: []*http.Response{{StatusCode: http.StatusOK, Body: pr}}}
	ch := make(chan PriceStreamItem, 30)
	err := newStallTestClient(fake, 200*time.Millisecond).Price(stallTestContext(t), NewPriceStreamRequest("EUR_USD"), ch, make(chan struct{}))
	if !errors.Is(err, ErrStreamEnded) {
		t.Fatalf("err = %v, want ErrStreamEnded", err)
	}
	if len(ch) != 30 {
		t.Errorf("got %d heartbeats, want 30", len(ch))
	}
}

func TestStreamNotStalledBySlowConsumer(t *testing.T) {
	// Waiting for the consumer is not a stalled connection.
	fake := &fakeHTTPClient{responses: []*http.Response{jsonResponse(http.StatusOK, testHeartbeat+testHeartbeat)}}
	ch := make(chan PriceStreamItem)
	errCh := make(chan error, 1)
	go func() {
		errCh <- newStallTestClient(fake, 50*time.Millisecond).Price(t.Context(), NewPriceStreamRequest("EUR_USD"), ch, make(chan struct{}))
	}()
	for range 2 {
		time.Sleep(150 * time.Millisecond)
		<-ch
	}
	if err := <-errCh; !errors.Is(err, ErrStreamEnded) {
		t.Fatalf("err = %v, want ErrStreamEnded", err)
	}
}

func TestStreamStallTimeoutOption(t *testing.T) {
	if got := NewStreamClient("test-key").streamStallTimeout; got != 15*time.Second {
		t.Errorf("default stall timeout = %v, want 15s", got)
	}
	if got := NewStreamClient("test-key", WithStreamStallTimeout(0)).streamStallTimeout; got != 0 {
		t.Errorf("stall timeout = %v, want 0", got)
	}
	aborted := make(chan struct{})
	w := newStallWatchdog(0, func() { close(aborted) })
	w.reset()
	defer w.stop()
	select {
	case <-aborted:
		t.Fatal("a disabled watchdog aborted the stream")
	case <-time.After(50 * time.Millisecond):
	}
}
