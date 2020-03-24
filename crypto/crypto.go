// Package crypto provides cryptography methods.
package crypto

// Padding is the interface that wraps the padding methods.
type Padding interface {
	AddPadding(data []byte, blockSize int) (result []byte)
	RemovePadding(data []byte) (result []byte, err error)
}

// BlockSizeInfo is the interface that wraps the methods to get size info of block.
type BlockSizeInfo interface {
	DataBlockSize() int
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

// Encrypt encrypts data.
func Encrypt(data []byte, encryptor Encryptor, padding Padding) (result []byte, err error) {
	dataBlockSize := encryptor.DataBlockSize()
	cipherBlockSize := encryptor.CipherBlockSize()
	dataLen := len(data)
	dataBeginOff, dataEndOff := 0, dataBlockSize
	var resultLen int
	if padding == nil {
		resultLen = (dataLen + dataBlockSize - 1) / dataBlockSize * cipherBlockSize
	} else {
		resultLen = (dataLen + dataBlockSize) / dataBlockSize * cipherBlockSize
	}
	result = make([]byte, resultLen)
	resultBeginOff := 0
	for dataEndOff <= dataLen {
		if err = encryptor.Encrypt(result[resultBeginOff:], data[dataBeginOff:dataEndOff]); err != nil {
			return
		}
		resultBeginOff += cipherBlockSize
		dataBeginOff = dataEndOff
		dataEndOff += dataBlockSize
	}
	if padding == nil {
		if dataBeginOff < dataLen {
			err = encryptor.Encrypt(result[resultBeginOff:], data[dataBeginOff:])
		}
	} else {
		block := padding.AddPadding(data[dataBeginOff:], dataBlockSize)
		err = encryptor.Encrypt(result[resultBeginOff:], block)
	}
	return
}

// Decrypt decrypts data.
func Decrypt(data []byte, decryptor Decryptor, padding Padding) (result []byte, err error) {
	dataBlockSize := decryptor.DataBlockSize()
	cipherBlockSize := decryptor.CipherBlockSize()
	dataLen := len(data)
	dataBeginOff := 0
	result = make([]byte, dataLen/cipherBlockSize*dataBlockSize)
	resultBeginOff := 0
	for dataBeginOff < dataLen {
		var n int
		if n, err = decryptor.Decrypt(result[resultBeginOff:], data[dataBeginOff:dataBeginOff+cipherBlockSize]); err != nil {
			return
		}
		resultBeginOff += n
		dataBeginOff += cipherBlockSize
	}
	result = result[:resultBeginOff]
	if padding != nil {
		result, err = padding.RemovePadding(result)
	}
	return
}
