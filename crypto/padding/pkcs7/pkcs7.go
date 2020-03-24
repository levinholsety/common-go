// Package pkcs7 implements PKCS #7 padding algorithm.
package pkcs7

import (
	"errors"

	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/crypto"
)

type padding struct{}

// AddPadding adds padding to last block and returns it.
func (p *padding) AddPadding(data []byte, blockSize int) (result []byte) {
	dataLen := len(data)
	resultLen := (dataLen + blockSize) / blockSize * blockSize
	result = make([]byte, resultLen)
	copy(result, data)
	paddingByte := byte(resultLen - dataLen)
	comm.FillByteArray(result[dataLen:], paddingByte)
	return
}

// RemovePadding removes padding from data.
func (p *padding) RemovePadding(data []byte) (result []byte, err error) {
	dataLen := len(data)
	paddingByte := data[dataLen-1]
	resultLen := dataLen - int(paddingByte)
	if resultLen < 0 {
		err = errors.New("bad padding")
		return
	}
	for _, b := range data[resultLen:] {
		if b != paddingByte {
			err = errors.New("bad padding")
			return
		}
	}
	result = data[:resultLen]
	return
}

// NewPadding creates and returns an instance of PKCS #7 padding.
func NewPadding() crypto.Padding {
	return &padding{}
}
