package api

import (
	"apns"
	"bytes"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
)

type ApplePush struct {
	HttpHeader   http.Header
	HttpPlayload []byte
	DeviceToken  string
}

func (ap ApplePush) Notify() (int, []byte, error) {
	req, err := http.NewRequest("POST", apns.ServerDev+"/3/device/"+ap.DeviceToken, bytes.NewReader(ap.HttpPlayload))
	if nil != err {
		return 406, []byte("Init Request Failure"), err
	}

	req.Header = ap.HttpHeader

	client := &http.Client{
		Transport: &http2.Transport{},
	}

	resp, err := client.Do(req)
	if nil != err {
		return http.StatusOK, []byte("Do Request Failure"), err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if nil != err {
		return http.StatusOK, []byte("Read Response Error"), err
	}

	return http.StatusOK, body, nil
}
