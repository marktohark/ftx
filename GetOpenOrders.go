package Ftx

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type GetOpenOrders struct {}

type GetOpenOrdersData struct {
	CreatedAt     string   `json:"createdAt"`
	FilledSize    float64         `json:"filledSize"`
	Future        string      `json:"future"`
	ID            int64         `json:"id"`
	Market        string      `json:"market"`
	Price         float64     `json:"price"`
	AvgFillPrice  float64     `json:"avgFillPrice"`
	RemainingSize float64         `json:"remainingSize"`
	Side          string      `json:"side"`
	Size          float64         `json:"size"`
	Status        string      `json:"status"`
	Type          string      `json:"type"`
	ReduceOnly    bool        `json:"reduceOnly"`
	Ioc           bool        `json:"ioc"`
	PostOnly      bool        `json:"postOnly"`
	ClientID      *string `json:"clientId"`
}

func NewGetOpenOrders() *GetOpenOrders {
	return &GetOpenOrders{}
}

func(self *GetOpenOrders) Do(market string) ([]GetOpenOrdersData, error) {
	cl := NewHttpClient()
	resp, err := cl.Get("/api/orders?market=" + market, nil)
	if err != nil {
		return nil, err
	}
	err = CheckApiError(resp.String())
	if err != nil {
		return nil, err
	}
	var result []GetOpenOrdersData
	err = json.Unmarshal([]byte(gjson.Get(resp.String(), "result").Raw), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}