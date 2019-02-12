package main

import (
	"api"
	"apns"
	"encoding/json"
	"flag"
	"fmt"
	"jwt"
	"log"
	"net/http"
	"os"
)

var (
	port     = flag.String("port", "9696", "The service port.")
	tomlConf = flag.String("config", "/tmp/ownprovider.toml", "TOML format config file.")
	logFile  = flag.String("log", "/tmp/log_ownprovider.log", "The log file.")
)

func main() {

	flag.Parse()
	logFilePtr, logErr := os.OpenFile(*logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if nil != logErr {
		fmt.Println("Failed to open ", *logFile, " for log")
	}
	log.SetOutput(logFilePtr)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("--------------------OwnProvider Log--------------------")

	if *port == "" {
		fmt.Println("Please specify the service port")
	}
	if *tomlConf == "" {
		fmt.Println("The appid and business secrect configure file must be specified")
		return
	}

	log.Println("OwnProvider Server Start...")

	// API Service
	http.HandleFunc("/api/notify", push)
	http.HandleFunc("/api/v2/notify", ownprovider)
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

	if "" == topic || "" == deviceToken || "" == payload {
		log.Println("Missing params topic, token or payload which are required")
		appleResult.Code = 406
		appleResult.Body = "Topic, Device Token and Payload are required"
		jsonBytes, _ := json.Marshal(appleResult)
		w.Write([]byte(jsonBytes))
		return
	}

	jwtToken, err := jwt.Token(*tomlConf)
	if nil != err {
		log.Println(fmt.Sprintf("Create Jwt Token failure: %s", err))
		appleResult.Code = 406
		appleResult.Body = "Create Jwt Token Failure"
	}

	h := apns.ApnsHeader{"POST", "/3/device/" + deviceToken, "bearer " + jwtToken, topic}
	httpHeader, err := h.Build4HttpRequest()
	if nil != err {
		log.Println(fmt.Sprintf("Create HTTP Request Header failure: %s", err))
		appleResult.Code = 406
		appleResult.Body = "Create HTTP Request Header Failure"
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

func ownprovider(w http.ResponseWriter, r *http.Request) {

	type AppleResult struct {
		Code   int
		Header string
		Body   string
	}
	var appleResult = AppleResult{200, "", ""}

	jsonBytes, _ := json.Marshal(appleResult)
	w.Write([]byte(jsonBytes))
}
