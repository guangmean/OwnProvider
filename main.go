package main

import (
	"OwnProvider/apple"
	"OwnProvider/jwt"
	"OwnProvider/loger"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {

	p8 := os.Getenv("OWNPROVIDERP8")
	if p8 == "" {
		fmt.Println("Error: ENV - OWNPROVIDERP8 is empty")
		return
	}

	// Home Page
	http.HandleFunc("/ownprovider/inner/push", Push)

	// Server
	http.ListenAndServe(":27953", nil)

	loger.WriteLog(loger.LOG_LEVEL_INFO, "OwnProvider service is launching...")
}

func Push(w http.ResponseWriter, r *http.Request) {
	env := r.FormValue("env")
	pushType := r.FormValue("voip")
	deviceToken := r.FormValue("token")
	payload := r.FormValue("payload")
	topic := r.FormValue("topic")

	apnsTopic := "com.example.app"

	if len(topic) > 10 {
		apnsTopic = topic
	}

	if "voip" != pushType {
		pushType = "alert"
	} else {
		apnsTopic += ".voip"
	}

	jwtHeader := jwt.Header{
		Alg: "ES256",
		Kid: "***", // Your Kid
	}
	jwtPayload := jwt.Payload{
		Iss: "***", // Your Iss - Team Id
		Iat: time.Now().Unix(),
	}

	jwToken, err := jwt.Token(jwtHeader, jwtPayload, "")
	if nil != err {
		loger.WriteLog(loger.LOG_LEVEL_ERROR, err)
		w.Write([]byte("Error before push."))
		return
	}

	apnsId := uuid.New().String()
	var result map[string]interface{}
	json.Unmarshal([]byte(payload), &result)
	aps, ok := result["aps"].(map[string]interface{})
	if ok {
		if nil != aps["apnsid"] {
			id := aps["apnsid"].(string)
			if 36 == len(id) {
				apnsId = id
			}
		}
	}

	h := apple.Header{
		Method:         "POST",
		Path:           "/3/device/" + deviceToken,
		Authorization:  "bearer " + jwToken,
		ApnsPushType:   pushType, // alert | background | voip | complication | fileprovider | mdm
		ApnsId:         apnsId,
		ApnsExpiration: "0",
		ApnsPriority:   "10",
		ApnsTopic:      apnsTopic,
		//ApnsCollapseId: "",
	}
	httpHeader, err := h.Build()

	//log.Printf("HTTP HEADER : %v", httpHeader)

	server := apple.Gold
	if "sandbox" == env {
		server = apple.Dev
	}
	t := apple.Target{server, httpHeader, []byte(payload), deviceToken}
	resp, err := t.Notify()

	if nil != err {
		loger.WriteLog(loger.LOG_LEVEL_ERROR, err)
		w.Write([]byte("Push Failured"))
		return
	}

	appleReplyApnsId := resp.Header.Get("Apns-Id")

	respHttpBody, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		loger.WriteLog(loger.LOG_LEVEL_ERROR, err)
		w.Write([]byte("Push Requested"))
		return
	}
	resp.Body.Close()

	if "" == string(respHttpBody[:]) {
		loger.WriteLog(loger.LOG_LEVEL_INFO, deviceToken+", "+apnsId+"【"+appleReplyApnsId+"】"+" push success")
		w.Write([]byte("Push Success - " + apnsId + "【" + appleReplyApnsId + "】"))
	} else {
		loger.WriteLog(loger.LOG_LEVEL_INFO, deviceToken+", "+apnsId+"【"+appleReplyApnsId+"】"+" error, "+string(respHttpBody[:]))
		w.Write([]byte("Push Error - " + apnsId + "【" + appleReplyApnsId + "】"))
	}
}
