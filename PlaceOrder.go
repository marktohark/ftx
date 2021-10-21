package Ftx

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type PlaceOrder struct {
	cl *HttpClient
}

type PlaceOrderPayload struct {
	Market     string `json:"market"`
	Side       string      `json:"side"`
	Price      *float64     `json:"price"`
	Type       string      `json:"type"`
	Size       float64     `json:"size"`
	ReduceOnly *bool        `json:"reduceOnly,omitempty"`
	Ioc        *bool        `json:"ioc,omitempty"`
	PostOnly   *bool        `json:"postOnly,omitempty"`
	ClientID   *string `json:"clientId,omitempty"`
	RejectOnPriceBand *bool `json:"rejectOnPriceBand,omitempty"`
}

type PlaceOrderResponse struct {
	CreatedAt     string      `json:"createdAt"`
	FilledSize    float64     `json:"filledSize"`
	Future        string      `json:"future"`
	ID            int64       `json:"id"`
	Market        string      `json:"market"`
	Price         float64     `json:"price"`
	RemainingSize float64     `json:"remainingSize"`
	Side          string      `json:"side"`
	Size          float64     `json:"size"`
	Status        string      `json:"status"`
	Type          string      `json:"type"`
	ReduceOnly    bool        `json:"reduceOnly"`
	Ioc           bool        `json:"ioc"`
	PostOnly      bool        `json:"postOnly"`
	ClientID      *string      `json:"clientId"`
}

func NewPlaceOrder() *PlaceOrder {
	return &PlaceOrder{
		cl: NewHttpClient(),
	}
}

func(self *PlaceOrder) Do(payload *PlaceOrderPayload) (*PlaceOrderResponse, error) {
	payloadJson, _ := json.Marshal(payload)
	resp, err := self.cl.Post("/api/orders", nil, string(payloadJson))
	if err != nil {
		return nil, err
	}
	err = CheckApiError(resp.String())
	if err != nil {
		return nil, err
	}
	var result PlaceOrderResponse
	err = json.Unmarshal([]byte(gjson.Get(resp.String(), "result").Raw), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}