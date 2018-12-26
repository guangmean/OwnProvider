package main

import (
	"api"
	"apns"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"jwt"
	"net/http"
	"time"
)

var (
	port     = flag.String("port", "9696", "The service port.")
	tomlConf = flag.String("config", "/tmp/ownprovider.toml", "TOML format config file.")
	logFile  = flag.String("log", "/tmp/log_ownprovider.log", "The log file.")
)

func main() {

	flag.Parse()

	if *port == "" {
		fmt.Println("Please specify the service port")
	}
	if *tomlConf == "" {
		fmt.Println("The appid and business secrect configure file must be specified")
		return
	}

	fmt.Println("Own Provider Server Start...")

	// API Service
	http.HandleFunc("/api/notify", push)
	http.ListenAndServe(":9696", nil)
}

func push(w http.ResponseWriter, r *http.Request) {

	type AppleResult struct {
		Code   int
		Header string
		Body   string
	}
	var appleResult = AppleResult{200, "", ""}

	topic := r.FormValue("topic")
	deviceToken := r.FormValue("token")
	payload := r.FormValue("payload")

	if "" == topic {
		appleResult.Code = 406
		appleResult.Body = "Please specify a topic"
		jsonBytes, _ := json.Marshal(appleResult)
		w.Write([]byte(jsonBytes))
		return
	}

	if "" == deviceToken {
		appleResult.Code = 406
		appleResult.Body = "Please specify a token"
		jsonBytes, _ := json.Marshal(appleResult)
		w.Write([]byte(jsonBytes))
		return
	}

	if "" == payload {
		appleResult.Code = 406
		appleResult.Body = "Please specify a payload"
		jsonBytes, _ := json.Marshal(appleResult)
		w.Write([]byte(jsonBytes))
		return
	}

	jwtToken, err := jwt.Token(*tomlConf)
	if nil != err {
		appleResult.Code = 406
		appleResult.Body = "Create Jwt Token Failure"
	}

	h := apns.ApnsHeader{"POST", "/3/device/" + deviceToken, "bearer " + jwtToken, topic}
	httpHeader, err := h.Build4HttpRequest()
	if nil != err {
		appleResult.Code = 406
		appleResult.Body = "Create Request HTTP Header Failure"
	}

	var push = api.ApplePush{
		httpHeader,
		[]byte(payload),
		deviceToken,
		apns.ServerGold,
	}

	status, header, body, _ := push.Notify()
	appleResult.Code = status
	appleResult.Header = header
	appleResult.Body = string(body[:])

	jsonBytes, _ := json.Marshal(appleResult)
	w.Write([]byte(jsonBytes))
}
