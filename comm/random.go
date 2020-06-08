package comm

import (
	"crypto/rand"
	"io"
	"math/big"
)

//Random fill random bytes in buffer.
func Random(buf []byte) (err error) {
	_, err = io.ReadFull(rand.Reader, buf)
	return
}

//RandomBytes returns random bytes in specified length.
func RandomBytes(length uint) (result []byte, err error) {
	result = make([]byte, length)
	err = Random(result)
	return
}

//RandomInt returns a random int value from 0 to max.
func RandomInt(min, max int) (result int, err error) {
	value, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	result = int(value.Int64()) + min
	return
}
