package oanda

import (
	"encoding/json"
	"time"
)

// Definitions https://developer.oanda.com/rest-live-v20/primitives-df/

// DecimalNumber is a decimal number encoded as a string. The amount of precision provided depends
// on what the number represents.
type DecimalNumber string

// AccountUnits is a quantity of an Account's home currency. This is a DecimalNumber encoded as a
// string. The amount of precision provided depends on the Account's home currency.
type AccountUnits string

// Currency represents a currency name identifier. This is an ISO 4217 currency code (e.g., USD, EUR, JPY).
type Currency string

// Tag represents a categorization tag for an entity.
type Tag struct {
	// Type is the type of the tag.
	Type string `json:"type"`
	// Name is the name of the tag.
	Name string `json:"name"`
}

// InstrumentName is a string containing the base currency and quote currency delimited by a "_"
// (e.g., EUR_USD, GBP_JPY).
type InstrumentName = string

// InstrumentType represents the type of an Instrument.
type InstrumentType string

const (
	// InstrumentTypeCurrency represents a currency pair.
	InstrumentTypeCurrency InstrumentType = "CURRENCY"
	// InstrumentTypeCFD represents a Contract For Difference.
	InstrumentTypeCFD InstrumentType = "CFD"
	// InstrumentTypeMetal represents a metal (e.g., XAU, XAG).
	InstrumentTypeMetal InstrumentType = "METAL"
)

// DayOfWeek represents a day of the week.
type DayOfWeek string

const (
	// DayOfWeekSunday represents Sunday.
	DayOfWeekSunday DayOfWeek = "SUNDAY"
	// DayOfWeekMonday represents Monday.
	DayOfWeekMonday DayOfWeek = "MONDAY"
	// DayOfWeekTuesday represents Tuesday.
	DayOfWeekTuesday DayOfWeek = "TUESDAY"
	// DayOfWeekWednesday represents Wednesday.
	DayOfWeekWednesday DayOfWeek = "WEDNESDAY"
	// DayOfWeekThursday represents Thursday.
	DayOfWeekThursday DayOfWeek = "THURSDAY"
	// DayOfWeekFriday represents Friday.
	DayOfWeekFriday DayOfWeek = "FRIDAY"
	// DayOfWeekSaturday represents Saturday.
	DayOfWeekSaturday DayOfWeek = "SATURDAY"
)

// FinancingDayOfWeek represents a day of the week when financing charges are debited or credited.
type FinancingDayOfWeek struct {
	// DayOfWeek is the day of the week to charge the financing.
	DayOfWeek DayOfWeek `json:"dayOfWeek"`
	// DaysCharged is the number of days worth of financing to be charged on dayOfWeek.
	DaysCharged int `json:"daysCharged"`
}

// InstrumentFinancing contains financing data for an Instrument.
type InstrumentFinancing struct {
	// LongRate is the financing rate to be used for a long position for the instrument. The value
	// is in decimal rather than percentage points (e.g. 5% is represented as 0.05).
	LongRate DecimalNumber `json:"longRate"`
	// ShortRate is the financing rate to be used for a short position for the instrument. The value
	// is in decimal rather than percentage points (e.g. 5% is represented as 0.05).
	ShortRate DecimalNumber `json:"shortRate"`
	// FinancingDaysOfWeek contains the days of the week to debit or credit financing charges; the
	// exact time of day at which to apply the financing is set in the DivisionTradingGroup for the
	// client's account.
	FinancingDaysOfWeek []FinancingDayOfWeek `json:"financingDaysOfWeek"`
}

// Instrument represents the full specification of a tradeable instrument.
type Instrument struct {
	// Name is the name of the Instrument.
	Name InstrumentName `json:"name"`
	// Type is the type of the Instrument.
	Type InstrumentType `json:"type"`
	// DisplayName is the display name of the Instrument.
	DisplayName string `json:"displayName"`
	// PipLocation is the location of the "pip" for this instrument. The decimal position of the pip
	// in this Instrument's price can be found at 10^pipLocation (e.g. -4 pipLocation results in a
	// pip size of 10^-4 = 0.0001).
	PipLocation int `json:"pipLocation"`
	// DisplayPrecision is the number of decimal places that should be used to display prices for
	// this instrument.
	DisplayPrecision int `json:"displayPrecision"`
	// TradeUnitsPrecision is the amount of decimal places that may be provided when specifying the
	// number of units traded for this instrument.
	TradeUnitsPrecision int `json:"tradeUnitsPrecision"`
	// MinimumTradeSize is the smallest number of units allowed to be traded for this instrument.
	MinimumTradeSize DecimalNumber `json:"minimumTradeSize"`
	// MaximumTrailingStopDistance is the maximum trailing stop distance allowed for a trailing
	// stop loss created for this instrument. Specified in price units.
	MaximumTrailingStopDistance DecimalNumber `json:"maximumTrailingStopDistance"`
	// MinimumGuaranteedStopLossDistance is the minimum distance allowed between the Trade's fill
	// price and the configured price for guaranteed Stop Loss Orders created for this instrument.
	// Specified in price units.
	MinimumGuaranteedStopLossDistance DecimalNumber `json:"minimumGuaranteedStopLossDistance"`
	// MinimumTrailingStopDistance is the minimum trailing stop distance allowed for a trailing
	// stop loss created for this instrument. Specified in price units.
	MinimumTrailingStopDistance DecimalNumber `json:"minimumTrailingStopDistance"`
	// MaximumPositionSize is the maximum position size allowed for this instrument. Specified in units.
	MaximumPositionSize DecimalNumber `json:"maximumPositionSize"`
	// MaximumOrderUnits is the maximum units allowed for an Order placed for this instrument.
	// Specified in units.
	MaximumOrderUnits DecimalNumber `json:"maximumOrderUnits"`
	// MarginRate is the margin rate for this instrument.
	MarginRate DecimalNumber `json:"marginRate"`
	// Commission is the commission structure for this instrument.
	Commission InstrumentCommission `json:"commission"`
	// GuaranteedStopLossOrderMode describes the current guaranteed Stop Loss Order mode of the Account
	// for this Instrument.
	GuaranteedStopLossOrderMode GuaranteedStopLossOrderModeForInstrument `json:"guaranteedStopLossOrderMode"`
	// GuaranteedStopLossOrderExecutionPremium is the amount that is charged to the account if a
	// guaranteed Stop Loss Order is triggered and filled. The value is in price units and is charged
	// for each unit of the Trade.
	GuaranteedStopLossOrderExecutionPremium DecimalNumber `json:"guaranteedStopLossOrderExecutionPremium"`
	// GuaranteedStopLossOrderLevelRestriction contains the guaranteed Stop Loss Order level restriction
	// for this instrument.
	GuaranteedStopLossOrderLevelRestriction GuaranteedStopLossOrderLevelRestriction `json:"guaranteedStopLossOrderLevelRestriction"`
	// Financing contains the financing data for this instrument.
	Financing InstrumentFinancing `json:"financing"`
	// Tags are the tags associated with this instrument.
	Tags []Tag `json:"tags"`
}

// DateTime represents a date and time value in RFC 3339 format. The DateTime format is used for
// fields representing specific points in time.
type DateTime time.Time

// UnmarshalJSON implements custom JSON unmarshaling for DateTime to handle both RFC3339 format
// and the special "0" value which represents an unset/zero time.
func (t *DateTime) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "0" {
		return nil
	}
	tm, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}
	*t = DateTime(tm)
	return nil
}

// AcceptDatetimeFormat specifies how DateTime fields should be represented in HTTP responses.
type AcceptDatetimeFormat string

const (
	// AcceptDatetimeFormatUnix means DateTime fields will be specified as UNIX timestamps (fractional
	// seconds since January 1, 1970 UTC).
	AcceptDatetimeFormatUnix AcceptDatetimeFormat = "UNIX"
	// AcceptDatetimeFormatRFC3339 means DateTime fields will be specified in RFC 3339 format
	// (e.g., "2006-01-02T15:04:05.999999999Z").
	AcceptDatetimeFormatRFC3339 AcceptDatetimeFormat = "RFC3339"
)

// InstrumentCommission represents the commission structure for an Instrument.
type InstrumentCommission struct {
	// Commission is the commission amount (in the Account's home currency) charged per unitsTraded
	// of the instrument.
	Commission DecimalNumber `json:"commission"`
	// UnitsTraded is the number of units traded that the commission amount is based on.
	UnitsTraded DecimalNumber `json:"unitsTraded"`
	// MinimumCommission is the minimum commission amount (in the Account's home currency) that is
	// charged when an Order is filled for this instrument.
	MinimumCommission DecimalNumber `json:"minimumCommission"`
}

// GuaranteedStopLossOrderModeForInstrument describes the guaranteed Stop Loss Order mode for an Instrument.
type GuaranteedStopLossOrderModeForInstrument string

const (
	// GuaranteedStopLossOrderModeForInstrumentDisabled means the Account is not permitted to create
	// guaranteed Stop Loss Orders for this Instrument.
	GuaranteedStopLossOrderModeForInstrumentDisabled GuaranteedStopLossOrderModeForInstrument = "DISABLED"
	// GuaranteedStopLossOrderModeForInstrumentAllowed means the Account is able, but not required to
	// have guaranteed Stop Loss Orders for open Trades for this Instrument.
	GuaranteedStopLossOrderModeForInstrumentAllowed GuaranteedStopLossOrderModeForInstrument = "ALLOWED"
	// GuaranteedStopLossOrderModeForInstrumentRequired means the Account is required to have guaranteed
	// Stop Loss Orders for all open Trades for this Instrument.
	GuaranteedStopLossOrderModeForInstrumentRequired GuaranteedStopLossOrderModeForInstrument = "REQUIRED"
)

// GuaranteedStopLossOrderLevelRestriction represents the guaranteed Stop Loss Order level restriction
// for an Instrument.
type GuaranteedStopLossOrderLevelRestriction struct {
	// Volume is the total allowed trade volume for guaranteed Stop Loss Orders within the specified
	// price range.
	Volume DecimalNumber `json:"volume"`
	// PriceRange is the price range the volume applies to. This value is in price units.
	PriceRange DecimalNumber `json:"priceRange"`
}

type Direction string

const (
	DirectionLong  Direction = "LONG"
	DirectionShort Direction = "SHORT"
)

// PricingComponent to get Candlestick data for
// Can contain any combination of the characters “M” (midpoint candles) “B” (bid candles) and “A” (ask candles).
type PricingComponent = string

// ConversionFactor represents a conversion factor used in currency conversions.
type ConversionFactor struct {
	// Factor is the conversion factor value.
	Factor DecimalNumber `json:"factor"`
}

// HomeConversionFactors represents the factors to use to convert quantities of a given currency
// into the Account's home currency. The conversion factor depends on the scenario the conversion
// is required for.
type HomeConversionFactors struct {
	// GainQuoteHome is the ConversionFactor in effect for the Account for converting any gains
	// realized in Instrument quote units into units of the Account's home currency.
	GainQuoteHome ConversionFactor `json:"gainQuoteHome"`
	// LossQuoteHome is the ConversionFactor in effect for the Account for converting any losses
	// realized in Instrument quote units into units of the Account's home currency.
	LossQuoteHome ConversionFactor `json:"lossQuoteHome"`
	// GainBaseHome is the ConversionFactor in effect for the Account for converting any gains
	// realized in Instrument base units into units of the Account's home currency.
	GainBaseHome ConversionFactor `json:"gainBaseHome"`
	// LossBaseHome is the ConversionFactor in effect for the Account for converting any losses
	// realized in Instrument base units into units of the Account's home currency.
	LossBaseHome ConversionFactor `json:"lossBaseHome"`
}
