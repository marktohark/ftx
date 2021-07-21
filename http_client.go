package Ftx

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"strings"
	"time"
)

type HttpClient struct {
	cl *resty.Client
}

type httpQuery map[string]string

func NewHttpClient() *HttpClient {
	return &HttpClient{
		resty.New(),
	}
}

func(self *HttpClient) getDomain() string {
	return "https://ftx.com"
}

func(self *HttpClient) defaultHeader(method string, path string, query httpQuery, payload string) *resty.Request {
	r := self.cl.R()
	r.SetHeader("content-type", "application/json")
	r.SetBody(payload)
	pathAndQuery := ""
	if query != nil {
		r.SetQueryParams(query)
		params := r.QueryParam.Encode()
		pathAndQuery = fmt.Sprintf("%s?%s", path, params)
	} else {
		pathAndQuery = path
	}
	method = strings.ToUpper(method)
	ts := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	signaturePayload := fmt.Sprintf("%s%s%s%s", ts, method, pathAndQuery, payload)
	_hmac := ComputeHmacSha256(signaturePayload, ApiSecret)
	request := r.SetHeaders(map[string]string{
		"FTX-KEY": ApiKey,
		"FTX-SIGN": _hmac,
		"FTX-TS": ts,
	})
	return request
}

func(self *HttpClient) Get(path string, query httpQuery) (*resty.Response, error) {
	req := self.defaultHeader("GET", path, query, "")
	resp, err := req.Get(self.getDomain() + path)
	return  resp, err
}

func(self *HttpClient) Post(path string, query httpQuery, payload string) (*resty.Response, error) {
	req := self.defaultHeader("POST", path, query, payload)
	resp, err := req.Post(self.getDomain() + path)
	return  resp, err
}

func(self *HttpClient) Delete(path string, query httpQuery, payload string) (*resty.Response, error) {
	req := self.defaultHeader("DELETE", path, query, payload)
	resp, err := req.Delete(self.getDomain() + path)
	return  resp, err
}