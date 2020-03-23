package crypto

import (
	"crypto/cipher"
	"errors"
	"io"
)

// NewEncryptionWriter creates an encryption writer to encrypt data.
func NewEncryptionWriter(w io.Writer, b cipher.Block, padding Padding) io.WriteCloser {
	blockSize := b.BlockSize()
	return &encryptionWriter{
		writer:    w,
		block:     b,
		padding:   padding,
		blockSize: blockSize,
		buffer:    make([]byte, blockSize),
		offset:    0,
		encBuffer: make([]byte, blockSize),
		closed:    false,
	}
}

type encryptionWriter struct {
	writer    io.Writer
	block     cipher.Block
	padding   Padding
	blockSize int
	buffer    []byte
	offset    int
	encBuffer []byte
	closed    bool
}

func (p *encryptionWriter) Write(data []byte) (n int, err error) {
	if p.closed {
		err = errors.New("Writer has been closed")
		return
	}
	length := len(data)
	for n < length {
		count := copy(p.buffer[p.offset:], data[n:])
		p.offset += count
		n += count
		if p.offset == p.blockSize {
			p.block.Encrypt(p.encBuffer, p.buffer)
			if _, err = p.writer.Write(p.encBuffer); err != nil {
				return
			}
			p.offset = 0
		}
	}
	return
}

func (p *encryptionWriter) Close() (err error) {
	if p.closed {
		return
	}
	p.padding.AddPadding(p.buffer, p.offset)
	p.block.Encrypt(p.encBuffer, p.buffer)
	if _, err = p.writer.Write(p.encBuffer); err != nil {
		return
	}
	p.closed = true
	return
}
