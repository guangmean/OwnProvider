package main

import (
	"OwnProvider/apple"
	"OwnProvider/jwt"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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
	env := r.FormValue("env")
	pushType := r.FormValue("voip")
	deviceToken := r.FormValue("token")
	payload := r.FormValue("payload")
	apnsTopic := r.FormValue("topic")
	if len(apnsTopic) < 10 {
		apnsTopic = "com.angularcorp.iCupid"
	}

	if "voip" != pushType {
		pushType = "alert"
	} else {
		apnsTopic = apnsTopic + ".voip"
	}

	jwtHeader := jwt.Header{
		Alg: "ES256",
		Kid: "P5BN65T68Y", // Your Kid
	}
	jwtPayload := jwt.Payload{
		Iss: "T5RB2UCSAR", // Your Iss - Team Id
		Iat: time.Now().Unix(),
	}

	jwToken, err := jwt.Token(jwtHeader, jwtPayload, "")
	if nil != err {
		log.Println("Build JWT token failure before push")
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
	}
	t := apple.Target{server, httpHeader, []byte(payload), deviceToken}
	resp, err := t.Notify()

	if nil != err {
		log.Printf("Network erro: %v", err)
	}
	resp.Body.Close()

	respHttpBody, err := ioutil.ReadAll(resp.Body)

	log.Println("推送环境【" + server + "】|" + kind + "|" + resp.Status + ":" + string(respHttpBody[:]) + "|" + deviceToken)

	w.Write([]byte("Push Success"))
}
