package aes

import (
	"hash"

	"github.com/levinholsety/common-go/util"
)

//NewKey creates a 256 bits AES key.
func NewKey() []byte {
	return util.RandomBytes(32)
}

//NewIV creates a IV.
func NewIV() []byte {
	return util.RandomBytes(16)
}

//GenerateKey generates key.
func GenerateKey(password, salt []byte, alg hash.Hash, key []byte) {
	for i := 0; i < len(key); i += alg.Size() {
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
