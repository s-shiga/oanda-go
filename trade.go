package oanda

type TradeID string

type TradeSummary struct {
	ID         TradeID        `json:"id"`
	Instrument InstrumentName `json:"instrument_name"`
}
