package Ftx

import (
	"encoding/json"
)

type CancelAllOrders struct {
	cl *HttpClient
}

type CancelAllOrdersPayload struct {
	Market *string `json:"market"`
	Side   *string `json:"side"`
}

func NewCancelAllOrders() *CancelAllOrders {
	return &CancelAllOrders{
		cl: NewHttpClient(),
	}
}

func(self *CancelAllOrders) Do(payload *CancelAllOrdersPayload) (bool, error) {
	payloadJson, _ := json.Marshal(payload)
	resp, err := self.cl.Delete("/api/orders", nil, string(payloadJson))
	if err != nil {
		return false, err
	}
	err = CheckApiError(resp.String())
	if err != nil {
		return false, err
	}
	return true, nil
}