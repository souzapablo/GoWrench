package keys_load

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/youmark/pkcs8"
)

var privateKeys map[string]*rsa.PrivateKey

func LoadEncryptedPrivateKey(keyId, privateRsakeyBase64, passphrase string) (*rsa.PrivateKey, error) {

	if privateKeys == nil {
		privateKeys = make(map[string]*rsa.PrivateKey)
	}

	pemBytes, err := base64.StdEncoding.DecodeString(privateRsakeyBase64)
	if err != nil {
		return nil, fmt.Errorf("read key: %w", err)
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("no PEM block")
	}

	// Parse & decrypt PKCS#8 with passphrase
	key, err := pkcs8.ParsePKCS8PrivateKey(block.Bytes, []byte(passphrase))
	if err != nil {
		return nil, fmt.Errorf("parse pkcs8: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not RSA")
	}

	privateKeys[keyId] = rsaKey

	return rsaKey, nil
}

func GetPrivateKey(keyId string) (*rsa.PrivateKey, error) {
	key, ok := privateKeys[keyId]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", keyId)
	}
	return key, nil
}
