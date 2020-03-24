package crypto

import (
	"errors"
	"io"
)

type encryptionWriter struct {
	writer          io.Writer
	encryptor       Encryptor
	padding         Padding
	dataBlockSize   int
	dataBlock       []byte
	dataBlockOffset int
	cipherBlockSize int
	cipherBlock     []byte
	closed          bool
}

func (p *encryptionWriter) Write(data []byte) (n int, err error) {
	if p.closed {
		err = errors.New("Writer has been closed")
		return
	}
	length := len(data)
	for n < length {
		count := copy(p.dataBlock[p.dataBlockOffset:], data[n:])
		p.dataBlockOffset += count
		n += count
		if p.dataBlockOffset == p.dataBlockSize {
			if err = p.encryptor.Encrypt(p.cipherBlock, p.dataBlock); err != nil {
				return
			}
			if _, err = p.writer.Write(p.cipherBlock); err != nil {
				return
			}
			p.dataBlockOffset = 0
		}
	}
	return
}

func (p *encryptionWriter) Close() (err error) {
	if p.closed {
		return
	}
	if p.padding == nil {
		if p.dataBlockOffset > 0 {
			if err = p.encryptor.Encrypt(p.cipherBlock, p.dataBlock[:p.dataBlockOffset]); err != nil {
				return
			}
			if _, err = p.writer.Write(p.cipherBlock); err != nil {
				return
			}
		}
	} else {
		if err = p.encryptor.Encrypt(p.cipherBlock, p.padding.AddPadding(p.dataBlock[:p.dataBlockOffset], p.dataBlockSize)); err != nil {
			return
		}
		if _, err = p.writer.Write(p.cipherBlock); err != nil {
			return
		}
	}
	p.closed = true
	return
}

// NewEncryptionWriter wraps w and returns an encryption writer to encrypt data.
// w holds the data to be encrypted.
// The data will be encrypted after it has been written to the encryption writer.
// Remember to close the encryption writer at the end.
func NewEncryptionWriter(w io.Writer, encryptor Encryptor, padding Padding) io.WriteCloser {
	dataBlockSize := encryptor.DataBlockSize()
	cipherBlockSize := encryptor.CipherBlockSize()
	return &encryptionWriter{
		writer:          w,
		encryptor:       encryptor,
		padding:         padding,
		dataBlockSize:   dataBlockSize,
		dataBlock:       make([]byte, dataBlockSize),
		dataBlockOffset: 0,
		cipherBlockSize: cipherBlockSize,
		cipherBlock:     make([]byte, cipherBlockSize),
		closed:          false,
	}
}
