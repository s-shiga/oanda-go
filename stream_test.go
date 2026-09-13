package oanda

import (
	"context"
	"errors"
	"net/http"
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
				if item.GetType() != want || item.GetID() != "42" || item.GetTime().Time == nil {
					t.Errorf("unexpected stream item: %#v", item)
				}
				if fill, ok := item.(OrderFillTransaction); ok && (fill.Instrument != "USD_JPY" || fill.Units != "10000") {
					t.Errorf("unexpected fill: %#v", fill)
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
