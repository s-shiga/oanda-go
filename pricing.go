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

// GetType returns the price type string.
func (p ClientPrice) GetType() string {
	return p.Type
}

// GetTime returns the time of the price.
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

// GetType returns the heartbeat type string ("HEARTBEAT").
func (p PricingHeartbeat) GetType() string {
	return p.Type
}

// GetTime returns the time of the heartbeat.
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
	Liquidity string `json:"liquidity"`
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

// PriceLatestCandlesticksRequest represents a request for the latest candlestick data.
// Use [NewPriceLatestCandlesticksRequest] to create one, then chain setters.
type PriceLatestCandlesticksRequest struct {
	specifications    []CandleSpecification
	units             *DecimalNumber
	smooth            bool
	dailyAlignment    *int
	alignmentTimezone *string
	weeklyAlignment   *WeeklyAlignment
}

// NewPriceLatestCandlesticksRequest creates a new empty [PriceLatestCandlesticksRequest].
func NewPriceLatestCandlesticksRequest() *PriceLatestCandlesticksRequest {
	return &PriceLatestCandlesticksRequest{
		specifications: make([]CandleSpecification, 0),
		smooth:         false,
	}
}

// Specification adds one or more candle specifications (e.g. "EUR_USD:S10:BM") to the request.
func (r *PriceLatestCandlesticksRequest) Specification(specs ...CandleSpecification) *PriceLatestCandlesticksRequest {
	r.specifications = append(r.specifications, specs...)
	return r
}

// Units sets the number of units used to calculate the volume-weighted average bid and ask prices.
func (r *PriceLatestCandlesticksRequest) Units(units DecimalNumber) *PriceLatestCandlesticksRequest {
	r.units = &units
	return r
}

// Smooth enables smoothing, which uses the previous candle's close as the open price.
func (r *PriceLatestCandlesticksRequest) Smooth() *PriceLatestCandlesticksRequest {
	r.smooth = true
	return r
}

// DailyAlignment sets the hour of the day (0-23) used for granularities that have daily alignment.
func (r *PriceLatestCandlesticksRequest) DailyAlignment(dailyAlignment int) *PriceLatestCandlesticksRequest {
	r.dailyAlignment = &dailyAlignment
	return r
}

// AlignmentTimezone sets the timezone to use for the daily alignment parameter.
func (r *PriceLatestCandlesticksRequest) AlignmentTimezone(alignmentTimezone string) *PriceLatestCandlesticksRequest {
	r.alignmentTimezone = &alignmentTimezone
	return r
}

// WeeklyAlignment sets the day of the week used for granularities that have weekly alignment.
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

// PriceLatestCandlesticksResponse is the response returned by [priceService.LatestCandlesticks].
type PriceLatestCandlesticksResponse struct {
	LatestCandles []CandlestickResponse `json:"latestCandles"`
}

// LatestCandlesticks retrieves the latest candlestick data for the Account configured via [WithAccountID].
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/candles/latest
//
// Reference: https://developer.oanda.com/rest-live-v20/pricing-ep/#collapse_endpoint_1
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

// PriceInformationRequest represents a request for pricing information for one or more instruments.
// Use [NewPriceInformationRequest] to create one, then chain setters.
type PriceInformationRequest struct {
	Instruments            []InstrumentName
	Since                  *DateTime
	IncludeHomeConversions bool
}

// NewPriceInformationRequest creates a new empty [PriceInformationRequest].
func NewPriceInformationRequest() *PriceInformationRequest {
	return &PriceInformationRequest{
		Instruments:            make([]InstrumentName, 0),
		IncludeHomeConversions: false,
	}
}

// AddInstruments adds one or more instruments to retrieve pricing for.
func (r *PriceInformationRequest) AddInstruments(instruments ...InstrumentName) *PriceInformationRequest {
	r.Instruments = append(r.Instruments, instruments...)
	return r
}

// SetSince filters to return only prices that have changed since the given time.
func (r *PriceInformationRequest) SetSince(since DateTime) *PriceInformationRequest {
	r.Since = &since
	return r
}

// SetIncludeHomeConversions enables inclusion of home conversion factors in the response.
func (r *PriceInformationRequest) SetIncludeHomeConversions() *PriceInformationRequest {
	r.IncludeHomeConversions = true
	return r
}

func (r *PriceInformationRequest) validate() error {
	if len(r.Instruments) == 0 {
		return errors.New("missing instruments")
	}
	return nil
}

func (r *PriceInformationRequest) values() (url.Values, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("instruments", strings.Join(r.Instruments, ","))
	if r.Since != nil {
		values.Set("since", r.Since.String())
	}
	if r.IncludeHomeConversions {
		values.Set("includeHomeConversions", "true")
	}
	return values, nil
}

// PriceInformationResponse is the response returned by [priceService.Information].
type PriceInformationResponse struct {
	Prices          []ClientPrice     `json:"prices"`
	HomeConversions []HomeConversions `json:"homeConversions"`
	Time            DateTime          `json:"time"`
}

// Information retrieves pricing information for the specified instruments.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/pricing
//
// Reference: https://developer.oanda.com/rest-live-v20/pricing-ep/#collapse_endpoint_2
func (s *priceService) Information(ctx context.Context, req *PriceInformationRequest) (*PriceInformationResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%s/pricing", s.client.accountID)
	values, err := req.values()
	if err != nil {
		return nil, err
	}
	return doGet[PriceInformationResponse](s.client, ctx, path, values)
}

// PriceCandlesticksRequest represents a request for account-specific candlestick data.
// It extends [CandlesticksRequest] with an optional units parameter.
type PriceCandlesticksRequest struct {
	CandlesticksRequest
	units *int
}

// NewPriceCandlesticksRequest creates a new [PriceCandlesticksRequest] with the given instrument and granularity.
func NewPriceCandlesticksRequest(instrument InstrumentName, granularity CandlestickGranularity) *PriceCandlesticksRequest {
	return &PriceCandlesticksRequest{
		CandlesticksRequest: *NewCandlesticksRequest(instrument, granularity),
	}
}

// Units sets the number of units used to calculate the volume-weighted average bid and ask prices.
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

// Candlesticks retrieves candlestick data for an instrument using account-specific pricing.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/instruments/{instrument}/candles
//
// Reference: https://developer.oanda.com/rest-live-v20/pricing-ep/#collapse_endpoint_4
func (s *priceService) Candlesticks(ctx context.Context, req *PriceCandlesticksRequest) (*CandlestickResponse, error) {
	path := fmt.Sprintf("/v3/accounts/%s/instruments/%s/candles", s.client.accountID, req.Instrument)
	values, err := req.values()
	if err != nil {
		return nil, err
	}
	return doGet[CandlestickResponse](s.client, ctx, path, values)
}

// PriceStreamRequest represents a request to open a pricing stream.
// Use [NewPriceStreamRequest] to create one.
type PriceStreamRequest struct {
	instruments           []InstrumentName
	snapShot              bool
	includeHomeConversion bool
}

// NewPriceStreamRequest creates a new [PriceStreamRequest] for the given instruments.
// Snapshot is enabled by default.
func NewPriceStreamRequest(instruments ...InstrumentName) *PriceStreamRequest {
	return &PriceStreamRequest{
		instruments:           instruments,
		snapShot:              true,
		includeHomeConversion: false,
	}
}

// DisableSnapShot disables the initial snapshot of prices when the stream starts.
func (r *PriceStreamRequest) DisableSnapShot() *PriceStreamRequest {
	r.snapShot = false
	return r
}

// IncludeHomeConversion enables inclusion of home conversion factors in the stream.
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

// PriceStreamItem is the interface implemented by items received from the pricing stream
// ([ClientPrice] and [PricingHeartbeat]).
type PriceStreamItem interface {
	GetType() string
	GetTime() DateTime
}

// Price opens a streaming connection for pricing data. Items are sent to ch until
// done is closed or the context is cancelled.
//
// This corresponds to the OANDA API endpoint: GET /v3/accounts/{accountID}/pricing/stream
//
// Reference: https://developer.oanda.com/rest-live-v20/pricing-ep/#collapse_endpoint_3
func (c *StreamClient) Price(ctx context.Context, req *PriceStreamRequest, ch chan<- PriceStreamItem, done <-chan struct{}) error {
	path := fmt.Sprintf("/v3/accounts/%s/pricing/stream", c.accountID)
	values, err := req.values()
	if err != nil {
		return err
	}
	u, err := joinURL(c.baseURL, path, values)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return err
	}
	c.setHeaders(httpReq)
	httpResp, err := c.httpClient.Do(httpReq)
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
