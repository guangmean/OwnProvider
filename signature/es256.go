package signature

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

type EcdsaSignature struct {
	R, S *big.Int
}

func Sign(header string, playload string, keyfile string) (string, error) {
	// Sign with SHA256 hash algorithms
	h := sha256.New()
	h.Write([]byte(header + "." + playload))
	hash := h.Sum(nil)
	// Read private key content from file
	p8 := os.Getenv("OWNPROVIDERP8")
	if "" != keyfile {
		p8 = keyfile
	}
	data, err := ioutil.ReadFile(p8)
	if nil != err {
		return "", err
	}
	// The .p8 is a PKCS#8 format's pem(The base64 format of DER(Encoded ASN.1 structure)) file
	// ASN.1	- Abstract Syntax Notation One
	// BER		- Basic Encoding Rules
	// DER		- Distinguished Encoding Rules
	block, _ := pem.Decode(data)
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes) // Parse an unencrypted, PKCS#8 private key
	if nil != err {
		return "", err
	}

	// Do Sign
	ecdsaKey, _ := priv.(*ecdsa.PrivateKey)
	if r, s, err := ecdsa.Sign(rand.Reader, ecdsaKey, hash); err == nil {

		curveBits := ecdsaKey.Curve.Params().BitSize

		keyBytes := curveBits / 8
		if curveBits%8 > 0 {
			keyBytes += 1
		}

		// We serialize the outpus (r and s) into big-endian byte arrays and pad
		// them with zeros on the left to make sure the sizes work out. Both arrays
		// must be keyBytes long, and the output must be 2*keyBytes long.
		rBytes := r.Bytes()
		rBytesPadded := make([]byte, keyBytes)
		copy(rBytesPadded[keyBytes-len(rBytes):], rBytes)

		sBytes := s.Bytes()
		sBytesPadded := make([]byte, keyBytes)
		copy(sBytesPadded[keyBytes-len(sBytes):], sBytes)

		out := append(rBytesPadded, sBytesPadded...)

		return strings.TrimRight(base64.URLEncoding.EncodeToString(out), "="), nil

	} else {

		return "", err

	}

}
