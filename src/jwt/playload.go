package jwt

import (
	"encoding/base64"
	"encoding/json"
)

type JwtPlayload struct {
	Iss string `json:"iss"`
	Iat int64  `json:"iat"`
}

func (p JwtPlayload) Base64Content() (string, error) {

	playloadByte, err := json.Marshal(p)

	if nil != err {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(playloadByte), nil
}
