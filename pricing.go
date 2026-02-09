package oanda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// -----------------------------------------------------------------
// Definitions https://developer.oanda.com/rest-live-v20/pricing-df/
// -----------------------------------------------------------------

// ClientPrice represents the price available for an Account at a given time.
type ClientPrice struct {
	Type string   `json:"type"`
	Time DateTime `json:"time"`
	// Bids are the bid prices available.
	Bids []PriceBucket `json:"bids"`
	// Asks are the ask prices available.
	Asks []PriceBucket `json:"asks"`
	// CloseoutBid is the closeout bid price.
	CloseoutBid PriceValue `json:"closeoutBid"`
	// CloseoutAsk is the closeout ask price.
	CloseoutAsk PriceValue `json:"closeoutAsk"`
}

func (p ClientPrice) GetType() string {
	return p.Type
}

func (p ClientPrice) GetTime() DateTime {
	return p.Time
}

// PriceStatus represents the status of the Price.
type PriceStatus string

const (
	// PriceStatusTradeable indicates the Instrument's price is tradeable.
	PriceStatusTradeable PriceStatus = "tradeable"
	// PriceStatusNonTradeable indicates the Instrument's price is not tradeable.
	PriceStatusNonTradeable PriceStatus = "non-tradeable"
	// PriceStatusInvalid indicates a price for this Instrument has not been set.
	PriceStatusInvalid PriceStatus = "invalid"
)

// QuoteHomeConversionFactors represents the factors used to convert quantities
// of a Price's Instrument's quote currency into a quantity of the Account's home currency.
type QuoteHomeConversionFactors struct {
	// PositiveUnits is the factor used to convert a positive amount of the Price's
	// Instrument's quote currency into a positive amount of the Account's home currency.
	PositiveUnits DecimalNumber `json:"positiveUnits"`
	// NegativeUnits is the factor used to convert a negative amount of the Price's
	// Instrument's quote currency into a negative amount of the Account's home currency.
	NegativeUnits DecimalNumber `json:"negativeUnits"`
}

// HomeConversions represents the factors to use to convert quantities of a given
// currency into the Account's home currency.
type HomeConversions struct {
	// Currency is the currency to be converted into the home currency.
	Currency Currency `json:"currency"`
	// AccountGain is the factor used to convert any gains for an Account in the
	// specified currency into the Account's home currency.
	AccountGain DecimalNumber `json:"accountGain"`
	// AccountLoss is the factor used to convert any losses for an Account in the
	// specified currency into the Account's home currency.
	AccountLoss DecimalNumber `json:"accountLoss"`
	// PositionValue is the factor used to convert a Position or Trade Value in the
	// specified currency into the Account's home currency.
	PositionValue DecimalNumber `json:"positionValue"`
}

// PricingHeartbeat represents a heartbeat message sent for a Pricing stream.
type PricingHeartbeat struct {
	// Type is the string "HEARTBEAT".
	Type string `json:"type"`
	// Time is the date/time when the PricingHeartbeat was created.
	Time DateTime `json:"time"`
}

func (p PricingHeartbeat) GetType() string {
	return p.Type
}

func (p PricingHeartbeat) GetTime() DateTime {
	return p.Time
}

// CandleSpecification is a string containing the following, all delimited by ":"
// characters: 1) InstrumentName 2) CandlestickGranularity 3) PricingComponent
// e.g. "EUR_USD:S10:BM"
type CandleSpecification string

// Common Definitions https://developer.oanda.com/rest-live-v20/pricing-common-df/

// PriceValue is a string representation of a decimal number that represents a price value.
type PriceValue string

// PriceBucket represents a price available for a specified liquidity amount.
type PriceBucket struct {
	// Price is the price offered.
	Price PriceValue `json:"price"`
	// Liquidity is the amount of liquidity offered.
	Liquidity int32 `json:"liquidity"`
}

// ---------------------------------------------------------------
// Endpoints https://developer.oanda.com/rest-live-v20/pricing-ep/
// ---------------------------------------------------------------

type priceService struct {
	client *Client
}

func newPriceService(client *Client) *priceService {
	return &priceService{client}
}

type PriceLatestCandlesticksRequest struct {
	specifications    []CandleSpecification
	units             *DecimalNumber
	smooth            bool
	dailyAlignment    *int
	alignmentTimezone *string
	weeklyAlignment   *WeeklyAlignment
}

func NewPriceLatestCandlesticksRequest() *PriceLatestCandlesticksRequest {
	return &PriceLatestCandlesticksRequest{
		specifications: make([]CandleSpecification, 0),
		smooth:         false,
	}
}

func (r *PriceLatestCandlesticksRequest) Specification(specs ...CandleSpecification) *PriceLatestCandlesticksRequest {
	r.specifications = append(r.specifications, specs...)
	return r
}

func (r *PriceLatestCandlesticksRequest) Units(units DecimalNumber) *PriceLatestCandlesticksRequest {
	r.units = &units
	return r
}

func (r *PriceLatestCandlesticksRequest) Smooth() *PriceLatestCandlesticksRequest {
	r.smooth = true
	return r
}

func (r *PriceLatestCandlesticksRequest) DailyAlignment(dailyAlignment int) *PriceLatestCandlesticksRequest {
	r.dailyAlignment = &dailyAlignment
	return r
}

func (r *PriceLatestCandlesticksRequest) AlignmentTimezone(alignmentTimezone string) *PriceLatestCandlesticksRequest {
	r.alignmentTimezone = &alignmentTimezone
	return r
}

func (r *PriceLatestCandlesticksRequest) WeeklyAlignment(weeklyAlignment WeeklyAlignment) *PriceLatestCandlesticksRequest {
	r.weeklyAlignment = &weeklyAlignment
	return r
}

func (r *PriceLatestCandlesticksRequest) validate() error {
	if len(r.specifications) == 0 {
		return errors.New("missing specifications")
	}
	if r.dailyAlignment != nil {
		if *r.dailyAlignment < 0 || *r.dailyAlignment > 23 {
			return fmt.Errorf("daily alignment must be between 0 and 23")
		}
	}
	if r.alignmentTimezone != nil {
		if _, err := time.LoadLocation(*r.alignmentTimezone); err != nil {
			return err
		}
	}
	return nil
}

func (r *PriceLatestCandlesticksRequest) values() (url.Values, error) {
	values := url.Values{}
	if err := r.validate(); err != nil {
		return nil, err
	}
	var specs []string
	for _, spec := range r.specifications {
		specs = append(specs, string(spec))
	}
	values.Set("candleSpecifications", strings.Join(specs, ","))
	if r.units != nil {
		values.Set("units", string(*r.units))
	}
	if r.smooth {
		values.Set("smooth", "True")
	}
	if r.dailyAlignment != nil {
		values.Set("dailyAlignment", strconv.Itoa(*r.dailyAlignment))
	}
	if r.alignmentTimezone != nil {
		values.Set("alignmentTimezone", *r.alignmentTimezone)
	}
	if r.weeklyAlignment != nil {
		values.Set("weeklyAlignment", string(*r.weeklyAlignment))
	}
	return values, nil
}

type PriceLatestCandlesticksResponse struct {
	LatestCandles []CandlestickResponse `json:"latestCandles"`
}

func (s *priceService) LatestCandlesticks(ctx context.Context, req *PriceLatestCandlesticksRequest) ([]CandlestickResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%s/candles/latest", s.client.accountID)
	values, err := req.values()
	if err != nil {
		return nil, err
	}
	resp, err := doGet[PriceLatestCandlesticksResponse](s.client, ctx, path, values)
	if err != nil {
		return nil, err
	}
	return resp.LatestCandles, nil
}

type PriceInformationRequest struct {
	instruments            []InstrumentName
	since                  *DateTime
	includeHomeConversions bool
}

func NewPriceInformationRequest() *PriceInformationRequest {
	return &PriceInformationRequest{
		instruments:            make([]InstrumentName, 0),
		includeHomeConversions: false,
	}
}

func (r *PriceInformationRequest) Instruments(instruments ...InstrumentName) *PriceInformationRequest {
	r.instruments = append(r.instruments, instruments...)
	return r
}

func (r *PriceInformationRequest) Since(since DateTime) *PriceInformationRequest {
	r.since = &since
	return r
}

func (r *PriceInformationRequest) IncludeHomeConversions() *PriceInformationRequest {
	r.includeHomeConversions = true
	return r
}

func (r *PriceInformationRequest) validate() error {
	if len(r.instruments) == 0 {
		return errors.New("missing instruments")
	}
	return nil
}

func (r *PriceInformationRequest) values() (url.Values, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("instruments", strings.Join(r.instruments, ","))
	if r.since != nil {
		values.Set("since", r.since.String())
	}
	if r.includeHomeConversions {
		values.Set("includeHomeConversions", "true")
	}
	return values, nil
}

type PriceInformationResponse struct {
	Prices          []ClientPrice     `json:"prices"`
	HomeConversions []HomeConversions `json:"homeConversions"`
	Time            DateTime          `json:"time"`
}

func (s *priceService) Information(ctx context.Context, req *PriceInformationRequest) (*PriceInformationResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%s/pricing", s.client.accountID)
	values, err := req.values()
	if err != nil {
		return nil, err
	}
	return doGet[PriceInformationResponse](s.client, ctx, path, values)
}

type PriceCandlesticksRequest struct {
	CandlesticksRequest
	units *int
}

func NewPriceCandlesticksRequest(instrument InstrumentName, granularity CandlestickGranularity) *PriceCandlesticksRequest {
	return &PriceCandlesticksRequest{
		CandlesticksRequest: *NewCandlesticksRequest(instrument, granularity),
	}
}

func (r *PriceCandlesticksRequest) Units(units int) *PriceCandlesticksRequest {
	r.units = &units
	return r
}

func (r *PriceCandlesticksRequest) values() (url.Values, error) {
	values, err := r.CandlesticksRequest.values()
	if err != nil {
		return nil, err
	}
	if r.units != nil {
		if *r.units <= 0 {
			return nil, errors.New("units must be greater than 0")
		}
		values.Set("units", strconv.Itoa(*r.units))
	}
	return values, nil
}

func (s *priceService) Candlesticks(ctx context.Context, req *PriceCandlesticksRequest) (*CandlestickResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%s/instruments/%s/candles", s.client.accountID, req.Instrument)
	values, err := req.values()
	if err != nil {
		return nil, err
	}
	return doGet[CandlestickResponse](s.client, ctx, path, values)
}

type priceStreamService struct {
	client *StreamClient
}

func newPriceStreamService(client *StreamClient) *priceStreamService {
	return &priceStreamService{client: client}
}

type PriceStreamRequest struct {
	instruments           []InstrumentName
	snapShot              bool
	includeHomeConversion bool
}

func NewPriceStreamRequest(instruments ...InstrumentName) *PriceStreamRequest {
	return &PriceStreamRequest{
		instruments:           instruments,
		snapShot:              true,
		includeHomeConversion: false,
	}
}

func (r *PriceStreamRequest) DisableSnapShot() *PriceStreamRequest {
	r.snapShot = false
	return r
}

func (r *PriceStreamRequest) IncludeHomeConversion() *PriceStreamRequest {
	r.includeHomeConversion = true
	return r
}

func (r *PriceStreamRequest) values() (url.Values, error) {
	values := url.Values{}
	if len(r.instruments) == 0 {
		return nil, errors.New("missing instruments")
	}
	values.Set("instruments", strings.Join(r.instruments, ","))
	if !r.snapShot {
		values.Set("snapShot", "False")
	}
	if r.includeHomeConversion {
		values.Set("includeHomeConversions", "True")
	}
	return values, nil
}

type PriceStreamItem interface {
	GetType() string
	GetTime() DateTime
}

func (s *priceStreamService) Stream(ctx context.Context, req *PriceStreamRequest, ch chan<- PriceStreamItem, done <-chan struct{}) error {
	path := fmt.Sprintf("/v3/account/%s/pricing/stream", s.client.accountID)
	values, err := req.values()
	if err != nil {
		return err
	}
	u, err := joinURL(s.client.baseURL, path, values)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return err
	}
	s.client.setHeaders(httpReq)
	httpResp, err := s.client.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %w", err)
	}
	defer closeBody(httpResp)
	dec := json.NewDecoder(httpResp.Body)
	for {
		select {
		case <-done:
			slog.Info("price stream closed")
			return nil
		case <-ctx.Done():
			return ctx.Err()
		default:
			var typeOnly struct {
				Type string `json:"type"`
			}
			if err := dec.Decode(&typeOnly); err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("failed to decode JSON response: %w", err)
			}
			switch typeOnly.Type {
			case "PRICE":
				var price ClientPrice
				if err := dec.Decode(&price); err != nil {
					return err
				}
				ch <- price
			case "HEARTBEAT":
				var heartbeat PricingHeartbeat
				if err := dec.Decode(&heartbeat); err != nil {
					return err
				}
				ch <- heartbeat
			}
		}
	}
}
