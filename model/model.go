package model

import "time"

type Candle struct {
	Time  time.Time
	Price int
}

type Flag struct {
	Flag string `json:"flag"`
}

type CandleMap map[string][]Candle
