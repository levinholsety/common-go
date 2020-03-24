package aes

import (
	"hash"

	"github.com/levinholsety/common-go/comm"
)

// NewKey creates a 256 bits AES key.
func NewKey() (key []byte, err error) {
	return comm.RandomBytes(32)
}

// NewIV creates a IV.
func NewIV() (iv []byte, err error) {
	return comm.RandomBytes(16)
}

// GenerateKey generates key from password.
func GenerateKey(password, salt []byte, alg hash.Hash, keySize int) (key []byte) {
	key = make([]byte, keySize)
	for i := 0; i < keySize; i += alg.Size() {
		alg.Reset()
		if i > 0 {
			alg.Write(key[i-alg.Size() : i])
		}
		if len(password) > 0 {
			alg.Write(password)
		}
		if len(salt) > 0 {
			alg.Write(salt)
		}
		copy(key[i:], alg.Sum(nil))
	}
	return
}
