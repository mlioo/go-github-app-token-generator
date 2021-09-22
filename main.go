package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

var (
	appID string
)

func main() {

	b, ok := os.LookupEnv("APP_PRIVATE_KEY")
	if !ok {
		fmt.Printf("::error title=Private key env not set::Private key env variable has not been set")
		return
	}

	appID, ok = os.LookupEnv("APP_ID")
	if !ok {
		fmt.Printf("::error title=App ID not set::App ID envrionment var is not set")
		return
	}

	pemBytes, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		fmt.Printf("::error title=Base64 decode failed::PEM secret should be base64 encoded")
		return
	}

	key, err := LoadPEMFromBytes(pemBytes)
	if err != nil {
		fmt.Printf("::error title=PEM invalid::Unable to load PEM %s", err)
		return
	}

	jwt := IssueJWTFromPEM(key)

	token, err := GetInstallationToken(jwt)
	if err != nil {
		fmt.Printf("::error title=Installation token error::Unable to get intallation token")
		//todo handle error
	}

	fmt.Printf("::set-output name=token::%s", *token)
}

func LoadPEMFromBytes(key []byte) (*rsa.PrivateKey, error) {
	b, _ := pem.Decode(key)
	if b != nil {
		key = b.Bytes
	}

	parsedKey, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("private key should be a PKCS1 key; parse error: %v", err)
	}

	return parsedKey, nil
}
