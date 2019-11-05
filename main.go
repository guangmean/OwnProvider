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

	olog, err := os.OpenFile("/tmp/log_ownprovider.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if nil != err {
		fmt.Println("Can not open log file")
	}
	log.SetOutput(olog)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("------------------------------------------")

	// Home Page
	http.HandleFunc("/v2/push", Push)

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
		Iat: time.Now().Unix(),
	}

	jwToken, err := jwt.Token(jwtHeader, jwtPayload)
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

	server := apple.Gold
	if "dev" == env {
		server = apple.Dev
	}
	t := apple.Target{server, httpHeader, []byte(customPayload), deviceToken}
	resp, err := t.Notify()

	if nil != err {
		fmt.Println("Remote Notification Push Failure")
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

}
