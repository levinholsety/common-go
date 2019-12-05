package aes

import (
	"crypto/aes"

	"github.com/levinholsety/common-go/crypto"
	"github.com/levinholsety/common-go/crypto/padding"
)

// NewAES creates a cipher with ECB mode to encrypt or decrypt data.
func NewAES(key []byte) *crypto.Cipher {
	cb, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return &crypto.Cipher{
		Block:   cb,
		Padding: padding.PKCS7Padding,
	}
}

// NewAESCBC creates a cipher with CBC mode to encrypt or decrypt data.
// This cipher cannot be reused.
func NewAESCBC(key, iv []byte) *crypto.Cipher {
	cb, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return &crypto.Cipher{
		Block:   crypto.NewCBC(cb, iv),
		Padding: padding.PKCS7Padding,
	}
}

// // NewWriter creates a AES writer.
// func NewWriter(w io.Writer, key, iv []byte) (aesWriter io.WriteCloser, err error) {
// 	c := NewAES(key)
// 	if iv != nil {
// 		c = c.CBC(iv)
// 	}
// 	aesWriter = crypto.NewCipherWriter(w, c.Block, c.Padding)
// 	return
// }

// // NewReader creates a AES reader.
// func NewReader(r io.Reader, key, iv []byte) (aesReader io.Reader, err error) {
// 	c := NewAES(key)
// 	if iv != nil {
// 		c = c.CBC(iv)
// 	}
// 	aesReader = crypto.NewCipherReader(r, c.Block, c.Padding)
// 	return
// }
