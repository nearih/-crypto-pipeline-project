package model

import "time"

type Ticker struct {
	Symbol    string
	Bid       float64
	Ask       float64
	BidSize   float64
	AskSize   float64
	Last      float64
	Timestamp time.Time
}

type Request struct {
	Symbol string `json:"symbol"`
}
