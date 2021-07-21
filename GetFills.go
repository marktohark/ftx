package Ftx

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type GetFills struct {}

type GetFillsQuery struct {
	Market *string `json:"market,omitempty"`
	StartTime *int64 `json:"start_time,omitempty"`
	EndTime *int64 `json:"end_time,omitempty"`
	Order *string `json:"order,omitempty"`
	OrderId int64 `json:"orderId"`
}

type GetFillsData struct {
	Fee           float64     `json:"fee"`
	FeeCurrency   string      `json:"feeCurrency"`
	FeeRate       float64     `json:"feeRate"`
	Future        string      `json:"future"`
	ID            int64         `json:"id"`
	Liquidity     string      `json:"liquidity"`
	Market        string      `json:"market"`
	BaseCurrency  *string `json:"baseCurrency"`
	QuoteCurrency *string `json:"quoteCurrency"`
	OrderID       int64         `json:"orderId"`
	TradeID       int64         `json:"tradeId"`
	Price         float64     `json:"price"`
	Side          string      `json:"side"`
	Size          float64         `json:"size"`
	Time          string   `json:"time"`
	Type          string      `json:"type"`
}

func NewGetFills() *GetFills {
	return &GetFills{}
}

func(self *GetFills) Do(query *GetFillsQuery) ([]GetFillsData, error) {
	cl := NewHttpClient()
	resp, err := cl.Get("/api/fills", Struct2MapString(query))
	if err != nil {
		return nil, err
	}
	err = CheckApiError(resp.String())
	if err != nil {
		return nil, err
	}
	var result []GetFillsData
	err = json.Unmarshal([]byte(gjson.Get(resp.String(), "result").Raw), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}