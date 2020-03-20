package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"

	"github.com/levinholsety/common-go/crypto"
)

func encrypt(b cipher.Block, src []byte) (dst []byte, err error) {
	w := bytes.NewBuffer(make([]byte, 0, len(src)))
	cw := crypto.NewCipherWriter(w, b, new(crypto.PKCS7Padding))
	r := bytes.NewReader(src)
	if _, err = io.Copy(cw, r); err != nil {
		return
	}
	if err = cw.Close(); err != nil {
		return
	}
	dst = w.Bytes()
	return
}

func decrypt(b cipher.Block, src []byte) (dst []byte, err error) {
	w := bytes.NewBuffer(make([]byte, 0, len(src)))
	r, err := crypto.NewCipherReader(bytes.NewReader(src), b, new(crypto.PKCS7Padding))
	if err != nil {
		return
	}
	if _, err = io.Copy(w, r); err != nil {
		return
	}
	dst = w.Bytes()
	return
}

func Encrypt(src, key []byte) ([]byte, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return encrypt(b, src)
}

func Decrypt(src, key []byte) ([]byte, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return decrypt(b, src)
}

func EncryptCBC(src, key, iv []byte) (dst []byte, err error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return encrypt(crypto.NewCBC(b, iv), src)
}

func DecryptCBC(src, key, iv []byte) ([]byte, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return decrypt(crypto.NewCBC(b, iv), src)
}

func NewECB(key []byte) (cipher.Block, error) {
	return aes.NewCipher(key)
}
