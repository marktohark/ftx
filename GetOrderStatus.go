package Ftx

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"strconv"
)

type GetOrderStatus struct {}


type GetOrderStatusResponse struct {
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

func NewGetOrderStatus() *GetOrderStatus {
	return &GetOrderStatus{}
}

func(self *GetOrderStatus) Do(orderId int64) (*GetOrderStatusResponse, error) {
	cl := NewHttpClient()
	resp, err := cl.Get("/api/orders/" + strconv.FormatInt(orderId, 10), nil)
	if err != nil {
		return nil, err
	}
	err = CheckApiError(resp.String())
	if err != nil {
		return nil, err
	}
	var result GetOrderStatusResponse
	err = json.Unmarshal([]byte(gjson.Get(resp.String(), "result").Raw), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}