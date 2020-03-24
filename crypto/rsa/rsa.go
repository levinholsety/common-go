// Package rsa implements RSA encryption and decryption algorithm.
package rsa

import (
	"crypto/rsa"
	"io"

	"github.com/levinholsety/common-go/crypto"
)

// Encrypt encrypts data with RSA algorithm.
func Encrypt(data []byte, publicKey *rsa.PublicKey) (result []byte, err error) {
	return crypto.Encrypt(data, &encryptor{publicKey}, nil)
}

// Decrypt decrypts data with RSA algorithm.
func Decrypt(data []byte, privateKey *rsa.PrivateKey) (result []byte, err error) {
	return crypto.Decrypt(data, &decryptor{privateKey}, nil)
}

// NewEncryptionWriter creates and returns an encryption writer to encrypt data with RSA.
func NewEncryptionWriter(w io.Writer, publicKey *rsa.PublicKey) io.WriteCloser {
	return crypto.NewEncryptionWriter(w, &encryptor{publicKey}, nil)
}

// NewDecryptionReader creates and returns a decryption reader to decrypt data with RSA.
func NewDecryptionReader(r io.Reader, privateKey *rsa.PrivateKey) (io.Reader, error) {
	return crypto.NewDecryptionReader(r, &decryptor{privateKey}, nil)
}
