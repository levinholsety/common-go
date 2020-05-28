// Package pkcs7 implements PKCS #7 padding algorithm.
package pkcs7

import (
	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/crypto"
)

type paddingAlgorithm struct{}

// AddPadding adds padding to data and returns it.
func (p *paddingAlgorithm) AddPadding(block []byte, dataLen int) (err error) {
	blockSize := len(block)
	if blockSize < 0x01 || blockSize > 0xff {
		err = crypto.ErrIllegalBlockSize
		return
	}
	comm.FillBytes(block[dataLen:], byte(blockSize-dataLen))
	return
}

// RemovePadding removes padding from data.
func (p *paddingAlgorithm) RemovePadding(data []byte) (result []byte, err error) {
	length := len(data)
	if length < 1 {
		err = crypto.ErrBadPadding
		return
	}
	paddingByte := data[length-1]
	resultLen := length - int(paddingByte)
	if resultLen < 0 {
		err = crypto.ErrBadPadding
		return
	}
	for _, b := range data[resultLen:] {
		if b != paddingByte {
			err = crypto.ErrBadPadding
			return
		}
	}
	result = data[:resultLen]
	return
}

// NewPaddingAlgorithm creates and returns an instance of PKCS #7 padding.
func NewPaddingAlgorithm() crypto.PaddingAlgorithm {
	return &paddingAlgorithm{}
}
