package Ftx

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type GetBalances struct {}

type GetBalancesData struct {
	Coin                   string  `json:"coin"`
	Free                   float64 `json:"free"`
	SpotBorrow             float64 `json:"spotBorrow"`
	Total                  float64 `json:"total"`
	UsdValue               float64 `json:"usdValue"`
	AvailableWithoutBorrow float64 `json:"availableWithoutBorrow"`
}

func NewGetBalances() *GetBalances {
	return &GetBalances{}
}

func(self *GetBalances) Do() ([]GetBalancesData, error) {
	cl := NewHttpClient()
	resp, err := cl.Get("/api/wallet/balances", nil)
	if err != nil {
		return nil, err
	}
	err = CheckApiError(resp.String())
	if err != nil {
		return nil, err
	}
	var result []GetBalancesData
	err = json.Unmarshal([]byte(gjson.Get(resp.String(), "result").Raw), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
