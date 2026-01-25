package oanda

type TransactionID string

type ClientID string

type ClientTag string

type ClientComment string

type ClientExtensions struct {
	ID      ClientID      `json:"id"`
	Tag     ClientTag     `json:"tag"`
	Comment ClientComment `json:"comment"`
}
