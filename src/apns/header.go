package apns

import (
	"net/http"
	"reflect"
)

const (
	ServerDev  = "https://api.sandbox.push.apple.com"
	ServerGold = "https://api.push.apple.com"
)

type ApnsHeader struct {
	Method        string `json:"method"`
	Path          string `json:"path"`
	Authorization string `json:"authorization"`
	Topic         string `json:"apns-topic"`
}

// Build4HttpRequest build the ApnsHeader struct into a HTTP request header map format
func (h ApnsHeader) Build4HttpRequest() (http.Header, error) {
	httpHeader := http.Header{}
	t := reflect.TypeOf(h)
	v := reflect.ValueOf(h)
	for i := 0; i < t.NumField(); i++ {
		httpHeader.Set(t.Field(i).Tag.Get("json"), v.Field(i).String())
	}

	return httpHeader, nil
}
