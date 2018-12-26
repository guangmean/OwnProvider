package api

import (
	"bytes"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
)

type ApplePush struct {
	HttpHeader  http.Header
	HttpPayload []byte
	DeviceToken string
	Env         string
}

func (ap ApplePush) Notify() (int, string, []byte, error) {
	req, err := http.NewRequest("POST", ap.Env+"/3/device/"+ap.DeviceToken, bytes.NewReader(ap.HttpPayload))
	if nil != err {
		return 406, "", []byte("Init Request Failure"), err
	}

	req.Header = ap.HttpHeader

	client := &http.Client{
		Transport: &http2.Transport{},
	}

	resp, err := client.Do(req)
	if nil != err {
		return http.StatusOK, "", []byte("Do Request Failure"), err
	}

	code := resp.StatusCode
	header := resp.Header.Get("apns-id")
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if nil != err {
		return code, "", []byte("Read Response Error"), err
	}

	return code, header, body, nil
}
