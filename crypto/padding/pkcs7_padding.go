package padding

import (
	"bytes"
)

type pkcs7Padding struct{}

// AddPadding adds padding to last block and returns it.
func (p *pkcs7Padding) AddPadding(data []byte, blockSize int) []byte {
	dataLength := len(data)
	modulus := dataLength % blockSize
	paddingLength := blockSize - modulus
	blockWithPadding := bytes.Repeat([]byte{byte(paddingLength)}, blockSize)
	if paddingLength < blockSize {
		copy(blockWithPadding, data[dataLength-modulus:])
	}
	return blockWithPadding
}

// RemovePadding removes padding from data.
func (p *pkcs7Padding) RemovePadding(data []byte) []byte {
	dataLength := len(data)
	paddingLength := int(data[dataLength-1])
	return data[:dataLength-paddingLength]
}
