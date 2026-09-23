package oanda

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestDateTimeUnsetRoundTrip(t *testing.T) {
	for _, input := range []string{`null`, `"0"`, " \n null \t"} {
		t.Run(input, func(t *testing.T) {
			// Decoding an unset value must also clear an existing timestamp.
			ts := mustTime(t, "2026-09-13T10:00:00.123456789Z")
			dt := DateTime{ts}
			if err := json.Unmarshal([]byte(input), &dt); err != nil {
				t.Fatal(err)
			}
			if !dt.IsZero() {
				t.Fatalf("Time = %v, want the zero time", dt.Time)
			}
			encoded, err := json.Marshal(dt)
			if err != nil {
				t.Fatal(err)
			}
			if string(encoded) != "null" {
				t.Fatalf("encoded = %s, want null", encoded)
			}
			var restored DateTime
			if err := json.Unmarshal(encoded, &restored); err != nil {
				t.Fatal(err)
			}
			if !restored.IsZero() {
				t.Fatalf("round trip Time = %v, want the zero time", restored.Time)
			}
		})
	}
}

func TestDateTimeTimestampRoundTrip(t *testing.T) {
	ts := mustTime(t, "2026-09-13T10:00:00.123456789+09:00")
	encoded, err := json.Marshal(DateTime{ts})
	if err != nil {
		t.Fatal(err)
	}
	var restored DateTime
	if err := json.Unmarshal(encoded, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.IsZero() || !restored.Equal(ts) {
		t.Fatalf("round trip Time = %v, want %v", restored.Time, ts)
	}
}

func TestDateTimeInvalidJSONPreservesValue(t *testing.T) {
	for _, input := range []string{`""`, `"invalid"`, `0`, `false`, `{}`} {
		t.Run(input, func(t *testing.T) {
			ts := mustTime(t, "2026-09-13T10:00:00Z")
			dt := DateTime{ts}
			if err := json.Unmarshal([]byte(input), &dt); err == nil {
				t.Fatal("expected invalid timestamp to be rejected")
			}
			if dt.IsZero() || !dt.Equal(ts) {
				t.Fatalf("invalid input changed Time to %v", dt.Time)
			}
		})
	}
}

func TestDateTimeUnsetIsUsable(t *testing.T) {
	for _, input := range []string{`{"type":"HEARTBEAT"}`, `{"type":"HEARTBEAT","time":null}`, `{"type":"HEARTBEAT","time":"0"}`} {
		t.Run(input, func(t *testing.T) {
			var heartbeat PricingHeartbeat
			if err := json.Unmarshal([]byte(input), &heartbeat); err != nil {
				t.Fatal(err)
			}
			dt := heartbeat.GetTime()
			if !dt.IsZero() {
				t.Fatalf("time = %v, want the zero time", dt)
			}
			// An unset time must behave like the zero time.Time, not panic.
			if s := fmt.Sprint(dt, heartbeat); strings.Contains(s, "PANIC") {
				t.Errorf("printing an unset DateTime panicked: %s", s)
			}
			if dt.Format(time.RFC3339) != "0001-01-01T00:00:00Z" || dt.Unix() != (time.Time{}).Unix() || !dt.Before(time.Now()) {
				t.Errorf("unset DateTime does not behave like the zero time: %v", dt)
			}
		})
	}
}

func TestIsPositiveDecimal(t *testing.T) {
	for _, s := range []string{"1", "10", "0.125", "0.5", "1.0", "007"} {
		if !isPositiveDecimal(s) {
			t.Errorf("%q should be a positive decimal", s)
		}
	}
	for _, s := range []string{"", "0", "0.0", "00", "-1", "+2", ".5", "5.", "1e3", "1E3", "0x10", "0b11", "0o7",
		"0x1p3", "1p3", "1_000", "1,000", "Inf", "NaN", "1/2", " 1", "1 ", "１"} {
		if isPositiveDecimal(s) {
			t.Errorf("%q should not be a positive decimal", s)
		}
	}
}
