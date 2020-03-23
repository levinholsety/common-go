package crypto

import (
	"errors"

	"github.com/levinholsety/common-go/comm"
)

// var (
// 	errIllegalBlockSize = errors.New("illegal block size")
// 	errBadPadding       = errors.New("bad padding")
// )

// PKCS7Padding represents PKCS #7 padding.
type PKCS7Padding struct{}

var _ Padding = new(PKCS7Padding)

// AddPadding adds padding to last block and returns it.
func (p *PKCS7Padding) AddPadding(block []byte, size int) {
	b := byte(len(block) - size)
	comm.FillByte(block[size:], b)
}

// func (p *PKCS7Padding) AddPadding(src []byte, blockSize int) (dst []byte) {
// 	dataLength := len(src)
// 	paddingLength := blockSize - dataLength%blockSize
// 	padding := bytes.Repeat([]byte{byte(paddingLength)}, paddingLength)
// 	dst = make([]byte, dataLength+paddingLength)
// 	copy(dst, src)
// 	copy(dst[dataLength:], padding)
// 	return
// }

// RemovePadding removes padding from data.
func (p *PKCS7Padding) RemovePadding(block []byte) ([]byte, error) {
	blockSize := len(block)
	padding := block[blockSize-1]
	dataSize := blockSize - int(padding)
	if dataSize < 0 {
		return nil, errors.New("Bad padding")
	}
	for _, b := range block[dataSize:] {
		if b != padding {
			return nil, errors.New("Bad padding")
		}
	}
	return block[:dataSize], nil
}

// func (p *PKCS7Padding) RemovePadding(src []byte, blockSize int) (dataLength int, err error) {
// 	encDataLength := len(src)
// 	if encDataLength < blockSize || encDataLength%blockSize != 0 {
// 		err = errIllegalBlockSize
// 		return
// 	}
// 	paddingLength := int(src[encDataLength-1])
// 	if paddingLength > blockSize {
// 		err = errBadPadding
// 		return
// 	}
// 	dataLength = encDataLength - paddingLength
// 	padding := bytes.Repeat([]byte{byte(paddingLength)}, paddingLength)
// 	if !bytes.Equal(padding, src[dataLength:encDataLength]) {
// 		err = errBadPadding
// 		return
// 	}
// 	return
// }
