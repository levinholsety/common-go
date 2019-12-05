package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type pemFormatPublicKey struct{}

func (f *pemFormatPublicKey) EncodePublicKey(key *rsa.PublicKey) []byte {
	data, _ := x509.MarshalPKIXPublicKey(key)
	return pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: data,
	})
}
func (f *pemFormatPublicKey) DecodePublicKey(pemData []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemData)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}
