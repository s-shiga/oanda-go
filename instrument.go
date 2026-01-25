package oanda

type CandlestickGranularity string

const (
	S5  CandlestickGranularity = "S5"
	S10 CandlestickGranularity = "S10"
	S15 CandlestickGranularity = "S15"
	S30 CandlestickGranularity = "S30"
	M1  CandlestickGranularity = "M1"
	M2  CandlestickGranularity = "M2"
	M4  CandlestickGranularity = "M4"
	M5  CandlestickGranularity = "M5"
	M10 CandlestickGranularity = "M10"
	M15 CandlestickGranularity = "M15"
	M30 CandlestickGranularity = "M30"
	H1  CandlestickGranularity = "H1"
	H2  CandlestickGranularity = "H2"
	H3  CandlestickGranularity = "H3"
	H4  CandlestickGranularity = "H4"
	H6  CandlestickGranularity = "H6"
	H8  CandlestickGranularity = "H8"
	H12 CandlestickGranularity = "H12"
	D   CandlestickGranularity = "D"
	W   CandlestickGranularity = "W"
	M   CandlestickGranularity = "M"
)

type WeeklyAlignment string

const (
	WeeklyAlignmentMonday    WeeklyAlignment = "Monday"
	WeeklyAlignmentTuesday   WeeklyAlignment = "Tuesday"
	WeeklyAlignmentWednesday WeeklyAlignment = "Wednesday"
	WeeklyAlignmentThursday  WeeklyAlignment = "Thursday"
	WeeklyAlignmentFriday    WeeklyAlignment = "Friday"
	WeeklyAlignmentSaturday  WeeklyAlignment = "Saturday"
	WeeklyAlignmentSunday    WeeklyAlignment = "Sunday"
)

type Candlestick struct {
	Time DateTime        `json:"time"`
	Bid  CandlestickData `json:"bid"`
}

type CandlestickData struct {
	O PriceValue `json:"o"`
	H PriceValue `json:"h"`
	L PriceValue `json:"l"`
	C PriceValue `json:"c"`
}

type CandlestickResponse struct {
	Instrument  InstrumentName         `json:"instrument"`
	Granularity CandlestickGranularity `json:"granularity"`
	Candles     []Candlestick          `json:"candles"`
}
