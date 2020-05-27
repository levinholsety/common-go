// Package pkcs7 implements PKCS #7 padding algorithm.
package pkcs7

import (
	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/crypto"
)

type paddingAlgorithm struct{}

// AddPadding adds padding to last block and returns it.
func (p *paddingAlgorithm) AddPadding(data []byte, blockSize int) (result []byte, err error) {
	if blockSize < 0x01 || blockSize > 0xff {
		err = crypto.ErrIllegalBlockSize
		return
	}
	data = data[len(data)/blockSize*blockSize:]
	result = make([]byte, blockSize)
	n := copy(result, data)
	paddingByte := byte(blockSize - len(data))
	comm.FillByteArray(result[n:], paddingByte)
	return
}

// RemovePadding removes padding from data.
func (p *paddingAlgorithm) RemovePadding(data []byte) (result []byte, err error) {
	dataLen := len(data)
	paddingByte := data[dataLen-1]
	resultLen := dataLen - int(paddingByte)
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
