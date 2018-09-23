package jwt

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"math/big"
)

type JwtApple struct {
	TeamId     string `json:"teamid"`
	PrivateKey string `json:"privatekey"`
}

type EcdsaSignature struct {
	R, S *big.Int
}

func Signature(header string, playload string, privatekey string) (string, error) {
	h := sha256.New()
	h.Write([]byte(header + "." + playload))
	hash := h.Sum(nil)

	data, err := ioutil.ReadFile(privatekey)
	if nil != err {
		return "", err
	}
	block, _ := pem.Decode(data)
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes) // Parse an unencrypted, PKCS#8 private key
	if nil != err {
		return "", err
	}

	r, s, err := ecdsa.Sign(rand.Reader, priv.(*ecdsa.PrivateKey), hash)
	if nil != err {
		return "", err
	}
	sig, err := asn1.Marshal(EcdsaSignature{r, s})
	if nil != err {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sig), nil
}
