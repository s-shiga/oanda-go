package oanda

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Definitions https://developer.oanda.com/rest-live-v20/instrument-df/

// Instrument declared in primitives.go

// CandlestickGranularity specifies the granularity of the candlestick.
type CandlestickGranularity string

const (
	// S5 represents 5 second candlesticks, minute alignment.
	S5 CandlestickGranularity = "S5"
	// S10 represents 10 second candlesticks, minute alignment.
	S10 CandlestickGranularity = "S10"
	// S15 represents 15 second candlesticks, minute alignment.
	S15 CandlestickGranularity = "S15"
	// S30 represents 30 second candlesticks, minute alignment.
	S30 CandlestickGranularity = "S30"
	// M1 represents 1 minute candlesticks, minute alignment.
	M1 CandlestickGranularity = "M1"
	// M2 represents 2 minute candlesticks, hour alignment.
	M2 CandlestickGranularity = "M2"
	// M4 represents 4 minute candlesticks, hour alignment.
	M4 CandlestickGranularity = "M4"
	// M5 represents 5 minute candlesticks, hour alignment.
	M5 CandlestickGranularity = "M5"
	// M10 represents 10 minute candlesticks, hour alignment.
	M10 CandlestickGranularity = "M10"
	// M15 represents 15 minute candlesticks, hour alignment.
	M15 CandlestickGranularity = "M15"
	// M30 represents 30 minute candlesticks, hour alignment.
	M30 CandlestickGranularity = "M30"
	// H1 represents 1 hour candlesticks, hour alignment.
	H1 CandlestickGranularity = "H1"
	// H2 represents 2 hour candlesticks, day alignment.
	H2 CandlestickGranularity = "H2"
	// H3 represents 3 hour candlesticks, day alignment.
	H3 CandlestickGranularity = "H3"
	// H4 represents 4 hour candlesticks, day alignment.
	H4 CandlestickGranularity = "H4"
	// H6 represents 6 hour candlesticks, day alignment.
	H6 CandlestickGranularity = "H6"
	// H8 represents 8 hour candlesticks, day alignment.
	H8 CandlestickGranularity = "H8"
	// H12 represents 12 hour candlesticks, day alignment.
	H12 CandlestickGranularity = "H12"
	// D represents 1 day candlesticks, day alignment.
	D CandlestickGranularity = "D"
	// W represents 1 week candlesticks, aligned to start of week.
	W CandlestickGranularity = "W"
	// M represents 1 month candlesticks, aligned to first day of the month.
	M CandlestickGranularity = "M"
)

// WeeklyAlignment specifies the day of the week used for granularity that has weekly alignment.
type WeeklyAlignment string

const (
	// WeeklyAlignmentMonday means weekly candlesticks are aligned to Mondays.
	WeeklyAlignmentMonday WeeklyAlignment = "Monday"
	// WeeklyAlignmentTuesday means weekly candlesticks are aligned to Tuesdays.
	WeeklyAlignmentTuesday WeeklyAlignment = "Tuesday"
	// WeeklyAlignmentWednesday means weekly candlesticks are aligned to Wednesdays.
	WeeklyAlignmentWednesday WeeklyAlignment = "Wednesday"
	// WeeklyAlignmentThursday means weekly candlesticks are aligned to Thursdays.
	WeeklyAlignmentThursday WeeklyAlignment = "Thursday"
	// WeeklyAlignmentFriday means weekly candlesticks are aligned to Fridays.
	WeeklyAlignmentFriday WeeklyAlignment = "Friday"
	// WeeklyAlignmentSaturday means weekly candlesticks are aligned to Saturdays.
	WeeklyAlignmentSaturday WeeklyAlignment = "Saturday"
	// WeeklyAlignmentSunday means weekly candlesticks are aligned to Sundays.
	WeeklyAlignmentSunday WeeklyAlignment = "Sunday"
)

// Candlestick represents a candlestick for an instrument.
type Candlestick struct {
	// Time is the start time of the candlestick.
	Time DateTime `json:"time"`
	// Bid contains the candlestick data based on bid prices. Only provided if bid-based candles
	// were requested.
	Bid CandlestickData `json:"bid"`
	// Ask contains the candlestick data based on ask prices. Only provided if ask-based candles
	// were requested.
	Ask CandlestickData `json:"ask"`
	// Mid contains the candlestick data based on midpoint prices. Only provided if midpoint-based
	// candles were requested.
	Mid CandlestickData `json:"mid"`
	// Volume is the number of prices created during the time-range represented by the candlestick.
	Volume int `json:"volume"`
	// Complete indicates whether or not the candlestick is complete. A complete candlestick is one
	// whose ending time is not in the future.
	Complete bool `json:"complete"`
}

// CandlestickData contains the price data (open, high, low, close) for a candlestick.
type CandlestickData struct {
	// O is the first (open) price in the time-range represented by the candlestick.
	O PriceValue `json:"o"`
	// H is the highest price in the time-range represented by the candlestick.
	H PriceValue `json:"h"`
	// L is the lowest price in the time-range represented by the candlestick.
	L PriceValue `json:"l"`
	// C is the last (closing) price in the time-range represented by the candlestick.
	C PriceValue `json:"c"`
}

// CandlestickResponse represents the response containing candlestick data for an instrument.
type CandlestickResponse struct {
	// Instrument is the instrument whose Prices are represented by the candlesticks.
	Instrument InstrumentName `json:"instrument"`
	// Granularity is the granularity of the candlesticks provided.
	Granularity CandlestickGranularity `json:"granularity"`
	// Candles is the list of candlesticks that satisfy the request.
	Candles []Candlestick `json:"candles"`
}

// Endpoints https://developer.oanda.com/rest-live-v20/instrument-ep/

type CandlesticksRequest struct {
	Instrument        InstrumentName
	Price             PricingComponent
	Granularity       CandlestickGranularity
	Count             *int
	From              *time.Time
	To                *time.Time
	Smooth            bool
	IncludeFirst      bool
	DailyAlignment    *int
	AlignmentTimezone *string
	WeeklyAlignment   WeeklyAlignment
}

func NewCandlesticksRequest(instrument InstrumentName, granularity CandlestickGranularity) *CandlesticksRequest {
	return &CandlesticksRequest{
		Instrument:      instrument,
		Price:           "",
		Granularity:     granularity,
		Smooth:          false,
		IncludeFirst:    true,
		WeeklyAlignment: WeeklyAlignmentFriday,
	}
}

func (req *CandlesticksRequest) Mid() *CandlesticksRequest {
	if !strings.Contains(req.Price, "M") {
		req.Price += "M"
	}
	return req
}

func (req *CandlesticksRequest) Bid() *CandlesticksRequest {
	if !strings.Contains(req.Price, "B") {
		req.Price += "B"
	}
	return req
}

func (req *CandlesticksRequest) Ask() *CandlesticksRequest {
	if !strings.Contains(req.Price, "A") {
		req.Price += "A"
	}
	return req
}

func (req *CandlesticksRequest) SetCount(count int) *CandlesticksRequest {
	req.Count = &count
	return req
}

func (req *CandlesticksRequest) SetFrom(from time.Time) *CandlesticksRequest {
	req.From = &from
	return req
}

func (req *CandlesticksRequest) SetTo(to time.Time) *CandlesticksRequest {
	req.To = &to
	return req
}

func (req *CandlesticksRequest) SetSmooth() *CandlesticksRequest {
	req.Smooth = true
	return req
}

func (req *CandlesticksRequest) SetExcludeFirst() *CandlesticksRequest {
	req.IncludeFirst = false
	return req
}

func (req *CandlesticksRequest) SetDailyAlignment(dailyAlignment int) *CandlesticksRequest {
	req.DailyAlignment = &dailyAlignment
	return req
}

func (req *CandlesticksRequest) SetAlignmentTimezone(alignmentTimezone string) *CandlesticksRequest {
	req.AlignmentTimezone = &alignmentTimezone
	return req
}

func (req *CandlesticksRequest) SetWeeklyAlignment(weeklyAlignment WeeklyAlignment) *CandlesticksRequest {
	req.WeeklyAlignment = weeklyAlignment
	return req
}

func (req *CandlesticksRequest) validate() error {
	if req.Count != nil {
		if req.From != nil && req.To != nil {
			return errors.New("count cannot be set with both from and to")
		}
		if *req.Count < 0 {
			return errors.New("count cannot be negative")
		}
		if *req.Count > 5000 {
			return errors.New("max count is 5000")
		}
	}
	if req.AlignmentTimezone != nil {
		if _, err := time.LoadLocation(*req.AlignmentTimezone); err != nil {
			return fmt.Errorf("invalid timezone: %s", *req.AlignmentTimezone)
		}
	}
	if req.DailyAlignment != nil {
		if *req.DailyAlignment < 0 || *req.DailyAlignment > 23 {
			return errors.New("daily alignment must be between 0 and 23")
		}
	}
	return nil
}

// values validates parameters and returns url.Values
// Gives empty fields if set default
func (req *CandlesticksRequest) values() (url.Values, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if req.Price != "" {
		v.Set("price", req.Price)
	}
	if req.Granularity != S5 {
		v.Set("granularity", string(req.Granularity))
	}
	if req.Count != nil {
		v.Set("count", strconv.Itoa(*req.Count))
	}
	if req.From != nil {
		v.Set("from", req.From.Format(time.RFC3339))
	}
	if req.To != nil {
		v.Set("to", req.To.Format(time.RFC3339))
	}
	if req.Smooth {
		v.Set("smooth", "True")
	}
	if !req.IncludeFirst {
		v.Set("includeFirst", "False")
	}
	if req.DailyAlignment != nil {
		v.Set("dailyAlignment", strconv.Itoa(*req.DailyAlignment))
	}
	if req.AlignmentTimezone != nil {
		v.Set("alignmentTimezone", *req.AlignmentTimezone)
	}
	if req.WeeklyAlignment != WeeklyAlignmentFriday {
		v.Set("weeklyAlignment", string(req.WeeklyAlignment))
	}
	return v, nil
}

func (c *Client) Candlesticks(ctx context.Context, req *CandlesticksRequest) ([]Candlestick, error) {
	path := fmt.Sprintf("/v3/instruments/%s/candles", req.Instrument)
	v, err := req.values()
	if err != nil {
		return nil, err
	}
	resp, err := c.sendGetRequest(ctx, path, v)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	candlesticksResp := struct {
		Instrument  InstrumentName         `json:"instrument"`
		Granularity CandlestickGranularity `json:"granularity"`
		Candles     []Candlestick          `json:"candles"`
	}{}
	if err := decodeResponse(resp, &candlesticksResp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}
	if candlesticksResp.Instrument != req.Instrument {
		return nil, fmt.Errorf("expected instrument %q but got %q", req.Instrument, candlesticksResp.Instrument)
	}
	if candlesticksResp.Granularity != req.Granularity {
		return nil, fmt.Errorf("expected granularity %q but got %q", req.Granularity, candlesticksResp.Granularity)
	}
	return candlesticksResp.Candles, nil
}
