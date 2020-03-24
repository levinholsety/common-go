package rsa

import (
	"crypto/rand"
	"crypto/rsa"
)

var (
	kfXML   KeyFormat        = &xmlFormatKey{}
	kfPKCS1 PrivateKeyFormat = &pemPKCS1FormatPrivateKey{}
	kfPKCS8 PrivateKeyFormat = &pemPKCS8FormatPrivateKey{}
	kfPEM   PublicKeyFormat  = &pemFormatPublicKey{}
)

// XMLKeyFormat returns an instance of XML key format.
func XMLKeyFormat() KeyFormat {
	return kfXML
}

// PKCS1PrivateKeyFormat returns an instance of PKCS1 private key format.
func PKCS1PrivateKeyFormat() PrivateKeyFormat {
	return kfPKCS1
}

// PKCS8PrivateKeyFormat returns an instance of PKCS8 private key format.
func PKCS8PrivateKeyFormat() PrivateKeyFormat {
	return kfPKCS8
}

// PEMPublicKeyFormat returns an instance of PEM public key format.
func PEMPublicKeyFormat() PublicKeyFormat {
	return kfPEM
}

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
