package model

type TradeQ struct {
	ID      int64
	Account string
	Volume  float64
	Open    float64
	Close   float64
	Side    string
}
