package Ftx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"reflect"
	"strconv"
	"strings"
)

var (
	ApiKey = ""
	ApiSecret = ""
)

type APIError struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("<APIError> code=%d, msg=%s", e.Code, e.Message)
}

func IsAPIError(e error) bool {
	_, ok := e.(*APIError)
	return ok
}

func NewApiError(code int64, msg string) *APIError {
	return &APIError{
		Code: code,
		Message: msg,
	}
}

func ComputeHmacSha256(str string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(str))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func SetApiKey(key, secret string) {
	ApiKey, ApiSecret = key, secret
}

func CheckApiError(respJson string) error {
	success := gjson.Get(respJson, "success").Bool()
	if success {
		return nil
	}
	return NewApiError(0, gjson.Get(respJson, "error").String())
}

func Struct2MapString(data interface{}) map[string]string {
	result := make(map[string]string)
	values := reflect.ValueOf(data)
	types := reflect.TypeOf(data)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
		types = types.Elem()
	}
	numOfField := values.NumField()
	for i := 0; i < numOfField; i++ {
		v := values.Field(i)
		tag := types.Field(i).Tag.Get("json")
		if tag == "" {
			continue
		}
		fieldName := strings.Split(tag, ",")[0]
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				continue
			}
			result[fieldName] = TypeToString(v.Elem())
		} else {
			result[fieldName] = TypeToString(v)
		}
	}
	return result
}

func TypeToString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Int64:
		return strconv.FormatInt(value.Interface().(int64), 10)
	case reflect.Float64:
		return decimal.NewFromFloat(value.Interface().(float64)).String()
	case reflect.String:
		return value.Interface().(string)
	default:
		return ""
	}
}

func StrPtr(v string) *string {
	return &v
}

func Float64Ptr(v float64) *float64 {
	return &v
}

func Int64Ptr(v int64) *int64 {
	return &v
}