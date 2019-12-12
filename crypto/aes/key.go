package aes

import (
	"hash"

	"github.com/levinholsety/common-go/comm"
)

// GenerateKey generates a 256 bits AES key.
func GenerateKey() ([]byte, error) {
	return comm.RandomBytes(32)
}

// GenerateIV generates a IV.
func GenerateIV() ([]byte, error) {
	return comm.RandomBytes(16)
}

// NewKey creates key from password.
func NewKey(password, salt []byte, alg hash.Hash, key []byte) {
	keySize := len(key)
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
