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

// InstrumentName is a string containing the base currency and quote currency delimited by a "_"
// (e.g., EUR_USD, GBP_JPY).
type InstrumentName string

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
	tm, err := time.Parse("2006-01-02T15:04:05.999999999Z", s)
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

// Tag represents a categorization tag for an entity.
type Tag struct {
	// Type is the type of the tag.
	Type string `json:"type"`
	// Name is the name of the tag.
	Name string `json:"name"`
}

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
