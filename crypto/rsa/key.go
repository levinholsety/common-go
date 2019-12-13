package rsa

import (
	"crypto/rand"
	"crypto/rsa"
)

//Key format conversion tools.
var (
	XML       KeyFormat        = &xmlFormatKey{}
	PEM_PKCS1 PrivateKeyFormat = &pemPKCS1FormatPrivateKey{}
	PEM_PKCS8 PrivateKeyFormat = &pemPKCS8FormatPrivateKey{}
	PEM       PublicKeyFormat  = &pemFormatPublicKey{}
)

//PrivateKeyFormat provides methods for private key format conversion.
type PrivateKeyFormat interface {
	EncodePrivateKey(key *rsa.PrivateKey) []byte
	DecodePrivateKey(data []byte) (*rsa.PrivateKey, error)
}

//PublicKeyFormat provides methods for public key format conversion.
type PublicKeyFormat interface {
	EncodePublicKey(key *rsa.PublicKey) []byte
	DecodePublicKey(data []byte) (*rsa.PublicKey, error)
}

//KeyFormat wraps PrivateKeyFormat interface and PublicKeyFormat interface.
type KeyFormat interface {
	PrivateKeyFormat
	PublicKeyFormat
}

//NewPrivateKey generates a RSA 2048 bits private key.
func NewPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}
