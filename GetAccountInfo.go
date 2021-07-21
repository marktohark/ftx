package Ftx

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type GetAccountInfo struct {}

type GetAccountInfoResponse struct {
	BackstopProvider             bool        `json:"backstopProvider"`
	Collateral                   float64     `json:"collateral"`
	FreeCollateral               float64     `json:"freeCollateral"`
	InitialMarginRequirement     float64     `json:"initialMarginRequirement"`
	Leverage                     float64         `json:"leverage"`
	Liquidating                  bool        `json:"liquidating"`
	MaintenanceMarginRequirement float64     `json:"maintenanceMarginRequirement"`
	MakerFee                     float64     `json:"makerFee"`
	MarginFraction               float64     `json:"marginFraction"`
	OpenMarginFraction           float64     `json:"openMarginFraction"`
	TakerFee                     float64     `json:"takerFee"`
	TotalAccountValue            float64     `json:"totalAccountValue"`
	TotalPositionSize            float64     `json:"totalPositionSize"`
	Username                     string      `json:"username"`
	Positions                    []GetAccountInfoPositions `json:"positions"`
}

type GetAccountInfoPositions struct {
	Cost                         float64 `json:"cost"`
	EntryPrice                   float64 `json:"entryPrice"`
	Future                       string  `json:"future"`
	InitialMarginRequirement     float64 `json:"initialMarginRequirement"`
	LongOrderSize                float64 `json:"longOrderSize"`
	MaintenanceMarginRequirement float64 `json:"maintenanceMarginRequirement"`
	NetSize                      float64 `json:"netSize"`
	OpenSize                     float64 `json:"openSize"`
	RealizedPnl                  float64 `json:"realizedPnl"`
	ShortOrderSize               float64 `json:"shortOrderSize"`
	Side                         string  `json:"side"`
	Size                         float64 `json:"size"`
	UnrealizedPnl                float64     `json:"unrealizedPnl"`
}


func NewGetAccountInfo() *GetAccountInfo {
	return &GetAccountInfo{}
}

func(self *GetAccountInfo) Do() (*GetAccountInfoResponse, error) {
	cl := NewHttpClient()
	resp, err := cl.Get("/api/account", nil)
	if err != nil {
		return nil, err
	}
	err = CheckApiError(resp.String())
	if err != nil {
		return nil, err
	}
	var result GetAccountInfoResponse
	err = json.Unmarshal([]byte(gjson.Get(resp.String(), "result").Raw), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}