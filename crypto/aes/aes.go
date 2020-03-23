package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"io"

	"github.com/levinholsety/common-go/crypto"
)

func newBlock(key, iv []byte) (b cipher.Block, err error) {
	if b, err = aes.NewCipher(key); err != nil {
		return
	}
	b = crypto.NewCBC(b, iv)
	return
}

// Encrypt encrypts data from io.Reader to io.Writer with AES CBC algorithm.
func Encrypt(key, iv []byte, w io.Writer, r io.Reader) (err error) {
	b, err := newBlock(key, iv)
	if err != nil {
		return
	}
	return crypto.Encrypt(b, new(crypto.PKCS7Padding), w, r)
}

// EncryptByteArray encrypts data with AES CBC algorithm.
func EncryptByteArray(key, iv, data []byte) (result []byte, err error) {
	b, err := newBlock(key, iv)
	if err != nil {
		return
	}
	return crypto.EncryptByteArray(b, new(crypto.PKCS7Padding), data)
}

// Decrypt decrypts data from io.Reader to io.Writer with AES CBC algorithm.
func Decrypt(key, iv []byte, w io.Writer, r io.Reader) (err error) {
	b, err := newBlock(key, iv)
	if err != nil {
		return
	}
	return crypto.Decrypt(b, new(crypto.PKCS7Padding), w, r)
}

// DecryptByteArray decrypts data with AES CBC algorithm.
func DecryptByteArray(key, iv, data []byte) (result []byte, err error) {
	b, err := newBlock(key, iv)
	if err != nil {
		return
	}
	return crypto.DecryptByteArray(b, new(crypto.PKCS7Padding), data)
}
