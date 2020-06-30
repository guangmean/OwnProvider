package main

import (
	"OwnProvider/apple"
	"OwnProvider/jwt"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("OwnProvider inner only starting...")

	p8 := os.Getenv("OWNPROVIDERP8")
	if p8 == "" {
		fmt.Println("Error: ENV - OWNPROVIDERP8 is empty")
		return
	}

	logger, err := os.OpenFile("/tmp/log_ownprovider_inner.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if nil != err {
		fmt.Println("Can not open log file")
	}
	log.SetOutput(logger)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("------------------------------------------")

	// Home Page
	http.HandleFunc("/ownprovider/inner/push", Push)

	// Server
	http.ListenAndServe(":27953", nil)

	fmt.Println("OwnProvider inner only server ready!!!")
}

func Push(w http.ResponseWriter, r *http.Request) {
	deviceToken := r.FormValue("token")
	payload := r.FormValue("payload")

	jwtHeader := jwt.Header{
		Alg: "ES256",
		Kid: "***",
	}
	jwtPayload := jwt.Payload{
		Iss: "***",
		Iat: time.Now().Unix(),
	}

	jwToken, err := jwt.Token(jwtHeader, jwtPayload, "")
	if nil != err {
		fmt.Println("Build JWT token failure")
	}

	h := apple.Header{
		Method:        "POST",
		Path:          "/3/device/" + deviceToken,
		Authorization: "bearer " + jwToken,
		ApnsPushType:  "alert", // alert | background | voip | complication | fileprovider | mdm
		//ApnsId:         "",
		ApnsExpiration: "0",
		ApnsPriority:   "10",
		ApnsTopic:      "***",
		//ApnsCollapseId: "",
	}
	httpHeader, err := h.Build()

	log.Printf("HTTP HEADER : %v", httpHeader)

	server := apple.Gold
	t := apple.Target{server, httpHeader, []byte(payload), deviceToken}
	resp, err := t.Notify()

	if nil != err {
		log.Printf("Remote Notification Push Failure: %v", err)
	}

	respHttpBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	status := resp.Proto + " " + resp.Status
	header := ""
	for k, v := range resp.Header {
		header += k + ": " + strings.Join(v, " ")
	}
	body := string(respHttpBody[:])

	log.Println("inner ---> " + status)
	log.Println("inner ---> " + header)
	log.Println("inner ---> " + body)

	w.Write([]byte("Push Success"))
}
