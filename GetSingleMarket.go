package Ftx

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
)

type GetSingleMarket struct {}

type GetSingleMarketResponse struct {
	Name           string      `json:"name"`
	BaseCurrency   *string `json:"baseCurrency"`
	QuoteCurrency  *string `json:"quoteCurrency"`
	Type           string      `json:"type"`
	Underlying     string      `json:"underlying"`
	Enabled        bool        `json:"enabled"`
	Ask            float64     `json:"ask"`
	Bid            float64         `json:"bid"`
	Last           float64     `json:"last"`
	PostOnly       bool        `json:"postOnly"`
	PriceIncrement float64     `json:"priceIncrement"`
	SizeIncrement  float64     `json:"sizeIncrement"`
	Restricted     bool        `json:"restricted"`
}

func NewGetSingleMarket() *GetSingleMarket {
	return &GetSingleMarket{}
}

func(self *GetSingleMarket) Do(market string) (*GetSingleMarketResponse, error) {
	cl := NewHttpClient()
	resp, err := cl.Get(fmt.Sprintf("/api/markets/%s", market), nil)
	if err != nil {
		return nil, err
	}
	err = CheckApiError(resp.String())
	if err != nil {
		return nil, err
	}
	var result GetSingleMarketResponse
	err = json.Unmarshal([]byte(gjson.Get(resp.String(), "result").Raw), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}