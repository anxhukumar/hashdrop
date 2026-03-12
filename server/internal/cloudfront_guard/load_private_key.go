package cloudfrontguard

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func loadPrivateKey(key string) (*rsa.PrivateKey, error) {

	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	privateKey := parsedKey.(*rsa.PrivateKey)

	return privateKey, nil
}
