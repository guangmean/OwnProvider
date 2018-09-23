package jwt

import (
	"encoding/base64"
	"encoding/json"
)

type JwtHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
}

func (h JwtHeader) Base64Content() (string, error) {

	headerByte, err := json.Marshal(h)

	if nil != err {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(headerByte), nil
}
