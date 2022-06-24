package main

import (
	"OwnProvider/apple"
	"OwnProvider/jwt"
	"OwnProvider/loger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func init() {
	http.DefaultClient.Timeout = time.Second * 3
}

func main() {
	fmt.Println("OwnProvider inner only starting...")

	p8 := os.Getenv("OWNPROVIDERP8")
	if p8 == "" {
		fmt.Println("Error: ENV - OWNPROVIDERP8 is empty")
		return
	}

	// Home Page
	http.HandleFunc("/ownprovider/inner/push", Push)

	// Server
	http.ListenAndServe(":27953", nil)

	fmt.Println("OwnProvider inner only server ready!!!")
}

func Push(w http.ResponseWriter, r *http.Request) {
	env := r.FormValue("env")
	pushType := r.FormValue("voip")
	deviceToken := r.FormValue("token")
	payload := r.FormValue("payload")
	apnsTopic := r.FormValue("topic")

	teamid := "***"
	kid := "***"
	kfile := ""
	switch apnsTopic {
	case "com.topic.yours":
		teamid = "***"
		kid = "***"
		kfile = "***"
	default:
		apnsTopic = "com.example.app"
		teamid = "***"
		kid = "***"
		kfile = "***"
	}

	if "voip" != pushType {
		pushType = "alert"
	} else {
		apnsTopic = apnsTopic + ".voip"
	}

	jwtHeader := jwt.Header{
		Alg: "ES256",
		Kid: kid, // Your Kid
	}
	jwtPayload := jwt.Payload{
		Iss: teamid, // Your Iss - Team Id
		Iat: time.Now().Unix(),
	}

	jwToken, err := jwt.Token(jwtHeader, jwtPayload, kfile)
	if nil != err {
		loger.WriteLog(loger.LOG_LEVEL_ERROR, err)
		w.Write([]byte("Error before push."))
		return
	}

	kind := "UNKNOWN"
	apnsId := uuid.New().String()
	var result map[string]interface{}
	json.Unmarshal([]byte(payload), &result)
	aps, ok := result["aps"].(map[string]interface{})
	if ok {
		if nil != aps["type"] {
			tp := aps["type"].(float64)
			if 2 == tp {
				kind = "CHAT"
			} else if 3 == tp {
				kind = "ADMIRER"
			} else if 4 == tp {
				kind = "GIFTSENT"
			} else if 403 == tp {
				kind = "DECLINE"
			} else if 21 == tp {
				kind = "METION"
			} else {
				kind = strconv.Itoa(int(tp))
			}
		}

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
	if "sandbox" == env || "develop" == env {
		server = apple.Dev
	} else {
		env = "gold"
	}
	t := apple.Target{
		Env:         server,
		HttpHeader:  httpHeader,
		HttpPayload: []byte(payload),
		Token:       deviceToken,
	}
	resp, err := t.Notify()

	if nil != err {
		loger.WriteLog(loger.LOG_LEVEL_ERROR, err)
	}
	resp.Body.Close()

	respHttpBody, err := ioutil.ReadAll(resp.Body)
	respHttpBodyStr := string(respHttpBody[:])
	if 0 == len(respHttpBodyStr) {
		respHttpBodyStr = "NOBODY"
	}
	loger.WriteLog(loger.LOG_LEVEL_INFO, "推送【"+env+"->"+pushType+"】|"+kind+"|"+resp.Status+"|"+respHttpBodyStr+"|"+deviceToken)

	w.Write([]byte("Push Success"))
}
