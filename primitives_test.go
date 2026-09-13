package oanda

import (
	"encoding/json"
	"testing"
)

func TestDateTimeUnsetRoundTrip(t *testing.T) {
	for _, input := range []string{`null`, `"0"`, " \n null \t"} {
		t.Run(input, func(t *testing.T) {
			// Decoding an unset value must also clear an existing timestamp.
			ts := mustTime(t, "2026-09-13T10:00:00.123456789Z")
			dt := DateTime{&ts}
			if err := json.Unmarshal([]byte(input), &dt); err != nil {
				t.Fatal(err)
			}
			if dt.Time != nil {
				t.Fatalf("Time = %v, want nil", dt.Time)
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
			if restored.Time != nil {
				t.Fatalf("round trip Time = %v, want nil", restored.Time)
			}
		})
	}
}

func TestDateTimeTimestampRoundTrip(t *testing.T) {
	ts := mustTime(t, "2026-09-13T10:00:00.123456789+09:00")
	encoded, err := json.Marshal(DateTime{&ts})
	if err != nil {
		t.Fatal(err)
	}
	var restored DateTime
	if err := json.Unmarshal(encoded, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.Time == nil || !restored.Equal(ts) {
		t.Fatalf("round trip Time = %v, want %v", restored.Time, ts)
	}
}

func TestDateTimeInvalidJSONPreservesValue(t *testing.T) {
	for _, input := range []string{`""`, `"invalid"`, `0`, `false`, `{}`} {
		t.Run(input, func(t *testing.T) {
			ts := mustTime(t, "2026-09-13T10:00:00Z")
			dt := DateTime{&ts}
			if err := json.Unmarshal([]byte(input), &dt); err == nil {
				t.Fatal("expected invalid timestamp to be rejected")
			}
			if dt.Time == nil || !dt.Equal(ts) {
				t.Fatalf("invalid input changed Time to %v", dt.Time)
			}
		})
	}
}
