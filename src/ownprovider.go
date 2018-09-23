package main

import (
	"api"
	"apns"
	"encoding/json"
	"flag"
	"fmt"
	"jwt"
	"net/http"
)

var (
	port     = flag.String("port", "8888", "The service port.")
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

	// API Service
	http.HandleFunc("/api/notify", push)
	http.ListenAndServe(":9527", nil)
}

func push(w http.ResponseWriter, r *http.Request) {

	type AppleResult struct {
		Code int
		Body string
	}
	var appleResult = AppleResult{200, "Success"}

	topic := r.FormValue("topic")
	deviceToken := r.FormValue("token")
	playload := r.FormValue("playload")

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

	if "" == playload {
		appleResult.Code = 406
		appleResult.Body = "Please specify a playload"
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
		[]byte(playload),
		deviceToken,
	}

	status, body, _ := push.Notify()
	appleResult.Code = status
	appleResult.Body = string(body[:])

	jsonBytes, _ := json.Marshal(appleResult)
	w.Write([]byte(jsonBytes))
}
