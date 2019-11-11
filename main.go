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
	fmt.Println("OwnProvider starting...")

	logger, err := os.OpenFile("/tmp/log_ownprovider_v2.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if nil != err {
		fmt.Println("Can not open log file")
	}
	log.SetOutput(logger)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("------------------------------------------")

	// Home Page
	http.HandleFunc("/ownprovider/v2/innerpush", Push)
	http.HandleFunc("/ownprovider/v2/remote", Remote)

	// Server
	http.ListenAndServe(":27952", nil)

	fmt.Println("OwnProvider server ready!!!")
}

func Push(w http.ResponseWriter, r *http.Request) {
	iss := r.FormValue("iss")
	kid := r.FormValue("kid")
	env := r.FormValue("env")
	apnsTopic := r.FormValue("topic")
	deviceToken := r.FormValue("token")
	customPayload := r.FormValue("payload")

	jwtHeader := jwt.Header{
		Alg: "ES256",
		Kid: kid,
	}
	jwtPayload := jwt.Payload{
		Iss: iss,
		Iat: time.Now().Unix() + 600,
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
		ApnsTopic:      apnsTopic,
		//ApnsCollapseId: "",
	}
	httpHeader, err := h.Build()

	log.Printf("HTTP HEADER : %v", httpHeader)

	server := apple.Gold
	if "dev" == env {
		server = apple.Dev
	}
	t := apple.Target{server, httpHeader, []byte(customPayload), deviceToken}
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

func Remote(w http.ResponseWriter, r *http.Request) {
	iss := r.FormValue("iss")
	kid := r.FormValue("kid")
	p8 := r.FormValue("p8")
	env := r.FormValue("env")
	apnsType := r.FormValue("type")
	apnsTopic := r.FormValue("topic")
	deviceToken := r.FormValue("token")
	customPayload := r.FormValue("payload")

	if "" == apnsType {
		apnsType = "alert"
	}

	if "" == iss || "" == kid || "" == p8 || "" == apnsTopic || "" == deviceToken || "" == customPayload {
		w.Write([]byte("Please fill all args"))
		return
	}

	jwtHeader := jwt.Header{
		Alg: "ES256",
		Kid: kid,
	}
	jwtPayload := jwt.Payload{
		Iss: iss,
		Iat: time.Now().Unix() + 600,
	}

	p8file := "/tmp/ownprovider_" + iss + ".p8"
	p8Bytes := []byte(p8)
	err := ioutil.WriteFile(p8file, p8Bytes, 0644)
	if nil != err {
		fmt.Println("Build p8file failure")
	}

	jwToken, err := jwt.Token(jwtHeader, jwtPayload, p8file)
	if nil != err {
		fmt.Println("Build JWT token failure")
	}

	server := apple.Gold
	if "dev" == env {
		server = apple.Dev
	}

	h := apple.Header{
		Method:         "POST",
		Path:           "/3/device/" + deviceToken,
		Authorization:  "bearer " + jwToken,
		ApnsPushType:   apnsType, // alert | background | voip | complication | fileprovider | mdm        //ApnsId:         "",
		ApnsExpiration: "0",
		ApnsPriority:   "10",
		ApnsTopic:      apnsTopic,
		//ApnsCollapseId: "",
	}
	httpHeader, err := h.Build()
	t := apple.Target{server, httpHeader, []byte(customPayload), deviceToken}
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

	log.Println(status)
	log.Println(header)
	log.Println(body)

	w.Write([]byte(status + "\n" + header + "\n" + body))
}
