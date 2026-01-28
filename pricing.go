package oanda

// Definitions https://developer.oanda.com/rest-live-v20/pricing-df/

// ClientPrice represents the price available for an Account at a given time.
type ClientPrice struct {
	// Bids are the bid prices available.
	Bids []PriceBucket `json:"bids"`
	// Asks are the ask prices available.
	Asks []PriceBucket `json:"asks"`
	// CloseoutBid is the closeout bid price.
	CloseoutBid PriceValue `json:"closeoutBid"`
	// CloseoutAsk is the closeout ask price.
	CloseoutAsk PriceValue `json:"closeoutAsk"`
	// Timestamp is the date/time when the price was generated.
	Timestamp DateTime `json:"timestamp"`
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
	Liquidity int `json:"liquidity"`
}
