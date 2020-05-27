// Package crypto provides cryptography methods.
package crypto

import (
	"errors"
)

// PaddingAlgorithm is the interface that wraps the padding methods.
type PaddingAlgorithm interface {
	AddPadding(block []byte, dataLen int) (err error)
	RemovePadding(data []byte) (result []byte, err error)
}

// BlockSizeInfo is the interface that wraps the methods to get size info of block.
type BlockSizeInfo interface {
	BlockSize() int
	CipherBlockSize() int
}

// Encryptor is the interface that wraps the methods to encrypt data.
type Encryptor interface {
	BlockSizeInfo
	Encrypt(dst, src []byte) (err error)
}

// Decryptor is the interface that wraps the methods to decrypt data.
type Decryptor interface {
	BlockSizeInfo
	Decrypt(dst, src []byte) (n int, err error)
}

// Errors
var (
	ErrBadPadding       = errors.New("bad padding")
	ErrIllegalBlockSize = errors.New("illegal block size")
)

func readBlocks(data []byte, blockSize int, paddingAlg PaddingAlgorithm, onRead func(block []byte) error) (err error) {
	for len(data) >= blockSize {
		err = onRead(data[:blockSize])
		if err != nil {
			return
		}
		data = data[blockSize:]
	}
	if paddingAlg == nil {
		err = onRead(data)
	} else {
		block := make([]byte, blockSize)
		err = paddingAlg.AddPadding(block, copy(block, data))
		if err != nil {
			return
		}
		err = onRead(block)
	}
	return
}

// Encrypt encrypts data.
func Encrypt(data []byte, encryptor Encryptor, paddingAlg PaddingAlgorithm) (result []byte, err error) {
	blockSize := encryptor.BlockSize()
	cipherBlockSize := encryptor.CipherBlockSize()
	result = make([]byte, (len(data)+blockSize)/blockSize*cipherBlockSize)
	offset := 0
	err = readBlocks(data, blockSize, paddingAlg, func(block []byte) (err error) {
		err = encryptor.Encrypt(result[offset:], block)
		if err != nil {
			return
		}
		offset += cipherBlockSize
		return
	})
	return
}

// Decrypt decrypts data.
func Decrypt(data []byte, decryptor Decryptor, paddingAlg PaddingAlgorithm) (result []byte, err error) {
	cipherBlockSize := decryptor.CipherBlockSize()
	dataLen := len(data)
	if dataLen%cipherBlockSize != 0 {
		err = ErrBadPadding
		return
	}
	blockSize := decryptor.BlockSize()
	result = make([]byte, dataLen/cipherBlockSize*blockSize)
	offset := 0
	for len(data) > cipherBlockSize {
		_, err = decryptor.Decrypt(result[offset:], data[:cipherBlockSize])
		if err != nil {
			return
		}
		data = data[cipherBlockSize:]
		offset += blockSize
	}
	n, err := decryptor.Decrypt(result[offset:], data)
	if err != nil {
		return
	}
	if paddingAlg == nil {
		result = result[:len(result)-blockSize+n]
	} else {
		result, err = paddingAlg.RemovePadding(result)
	}
	return
}
