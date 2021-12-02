package apple

import (
	"bytes"
	"golang.org/x/net/http2"
	"net/http"
	"reflect"
	"time"
)

const (
	Dev  = "https://api.sandbox.push.apple.com"
	Gold = "https://api.push.apple.com"
)

// @See https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/sending_notification_requests_to_apns
type Header struct {
	Method         string `json:"method"`
	Path           string `json:"path"`
	Authorization  string `json:"authorization"`
	ApnsPushType   string `json:"apns-push-type"` // Has 6 valid values
	ApnsId         string `json:"apns-id"`
	ApnsExpiration string `json:"apns-expiration"`
	ApnsPriority   string `json:"apns-priority"`
	ApnsTopic      string `json:"apns-topic"`
	ApnsCollapseId string `json:"apns-collapse-id"`
}

type Target struct {
	Env         string
	HttpHeader  http.Header
	HttpPayload []byte
	Token       string
}

func (h Header) Build() (http.Header, error) {
	header := http.Header{}
	t := reflect.TypeOf(h)
	v := reflect.ValueOf(h)
	for i := 0; i < t.NumField(); i++ {
		if "" != v.Field(i).String() {
			header.Set(t.Field(i).Tag.Get("json"), v.Field(i).String())
		}
	}

	return header, nil
}

func (t *Target) Notify() (*http.Response, error) {
	r, err := http.NewRequest("POST", t.Env+"/3/device/"+t.Token, bytes.NewReader(t.HttpPayload))
	if nil != err {
		return nil, err
	}

	r.Header = t.HttpHeader

	client := &http.Client{
		Timeout:   time.Second * 3,
		Transport: &http2.Transport{},
	}

	result, err := client.Do(r)

	if nil != err {
		return nil, err
	}

	defer result.Body.Close()

	return result, nil
}
