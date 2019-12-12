package paddings

import (
	"bytes"
	"errors"

	"github.com/levinholsety/common-go/crypto"
)

var (
	errIllegalBlockSize = errors.New("illegal block size")
	errBadPadding       = errors.New("bad padding")
)

type pkcs7Padding struct{}

var _ crypto.Padding = new(pkcs7Padding)

// AddPadding adds padding to last block and returns it.
func (p *pkcs7Padding) AddPadding(src []byte, blockSize int) (dst []byte) {
	dataLength := len(src)
	paddingLength := blockSize - dataLength%blockSize
	padding := bytes.Repeat([]byte{byte(paddingLength)}, paddingLength)
	dst = make([]byte, dataLength+paddingLength)
	copy(dst, src)
	copy(dst[dataLength:], padding)
	return
}

// RemovePadding removes padding from data.
func (p *pkcs7Padding) RemovePadding(src []byte, blockSize int) (dataLength int, err error) {
	encDataLength := len(src)
	if encDataLength < blockSize || encDataLength%blockSize != 0 {
		err = errIllegalBlockSize
		return
	}
	paddingLength := int(src[encDataLength-1])
	if paddingLength > blockSize {
		err = errBadPadding
		return
	}
	dataLength = encDataLength - paddingLength
	padding := bytes.Repeat([]byte{byte(paddingLength)}, paddingLength)
	if !bytes.Equal(padding, src[dataLength:encDataLength]) {
		err = errBadPadding
		return
	}
	return
}
