package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type pemPKCS8FormatPrivateKey struct{}

func (f *pemPKCS8FormatPrivateKey) EncodePrivateKey(key *rsa.PrivateKey) []byte {
	data, _ := x509.MarshalPKCS8PrivateKey(key)
	return pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: data,
	})
}
func (f *pemPKCS8FormatPrivateKey) DecodePrivateKey(pemData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key.(*rsa.PrivateKey), nil
}
