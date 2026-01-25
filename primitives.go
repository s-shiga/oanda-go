package oanda

import (
	"encoding/json"
	"time"
)

type DecimalNumber string

type AccountUnits string

type Currency string

type InstrumentName string

type DateTime time.Time

func (t *DateTime) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "0" {
		return nil
	}
	tm, err := time.Parse("2006-01-02T15:04:05.999999999Z", s)
	if err != nil {
		return err
	}
	*t = DateTime(tm)
	return nil
}
