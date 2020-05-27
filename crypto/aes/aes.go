// Package aes implements AES encryption and decryption algorithm.
package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"io"

	"github.com/levinholsety/common-go/crypto"
	"github.com/levinholsety/common-go/crypto/mode/cbc"
	"github.com/levinholsety/common-go/crypto/padding/pkcs7"
)

func prepareCipher(key, iv []byte, f func(b cipher.Block) error) error {
	b, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	b = cbc.NewCipher(b, iv)
	return f(b)
}

// Encrypt encrypts data with AES CBC algorithm.
func Encrypt(data, key, iv []byte) (result []byte, err error) {
	err = prepareCipher(key, iv, func(b cipher.Block) (err error) {
		result, err = crypto.Encrypt(data, &block{b}, pkcs7.NewPaddingAlgorithm())
		return
	})
	return
}

// Decrypt decrypts data with AES CBC algorithm.
func Decrypt(data, key, iv []byte) (result []byte, err error) {
	err = prepareCipher(key, iv, func(b cipher.Block) (err error) {
		result, err = crypto.Decrypt(data, &block{b}, pkcs7.NewPaddingAlgorithm())
		return
	})
	return
}

// NewEncryptionWriter creates and returns an encryption writer.
// The writer wraps w which holds the data to be encrypted.
// When write data into it, the data will be encrypted with AES/CBC/PKCS7Padding algorithm.
func NewEncryptionWriter(w io.Writer, key, iv []byte) (ew io.WriteCloser, err error) {
	err = prepareCipher(key, iv, func(b cipher.Block) (err error) {
		ew = crypto.NewEncryptionWriter(w, &block{b}, pkcs7.NewPaddingAlgorithm())
		return
	})
	return
}

// NewDecryptionReader creates and returns a decryption reader.
// The reader wraps r which holds the data to be decrypted.
// When read data from it, the data will be decrypted with AES/CBC/PKCS7Padding algorithm.
func NewDecryptionReader(r io.Reader, key, iv []byte) (dr io.Reader, err error) {
	err = prepareCipher(key, iv, func(b cipher.Block) (err error) {
		dr, err = crypto.NewDecryptionReader(r, &block{b}, pkcs7.NewPaddingAlgorithm())
		return
	})
	return
}
