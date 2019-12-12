package cipher

import (
	"crypto/cipher"
	"io"

	"github.com/levinholsety/common-go/crypto"
)

func NewCipherWriter(w io.Writer, b cipher.Block, padding crypto.Padding) io.WriteCloser {
	blockSize := b.BlockSize()
	return &cipherWriter{
		writer:       w,
		block:        b,
		padding:      padding,
		blockSize:    blockSize,
		buffer:       make([]byte, blockSize),
		encData:      make([]byte, blockSize),
		bufferLength: 0,
	}
}

type cipherWriter struct {
	writer       io.Writer
	block        cipher.Block
	padding      crypto.Padding
	blockSize    int
	buffer       []byte
	bufferLength int
	encData      []byte
}

var _ io.WriteCloser = new(cipherWriter)

func (p *cipherWriter) Write(data []byte) (n int, err error) {
	n = len(data)
	for len(data) > 0 {
		if p.bufferLength < p.blockSize {
			copied := copy(p.buffer[p.bufferLength:], data)
			p.bufferLength += copied
			data = data[copied:]
			if p.bufferLength == p.blockSize {
				p.block.Encrypt(p.encData, p.buffer)
				if _, err = p.writer.Write(p.encData); err != nil {
					return
				}
				p.bufferLength = 0
			}
		}
	}
	return
}

func (p *cipherWriter) Close() (err error) {
	if p.buffer == nil {
		return
	}
	p.buffer = p.buffer[:p.bufferLength]
	p.buffer = p.padding.AddPadding(p.buffer, p.blockSize)
	p.block.Encrypt(p.encData, p.buffer)
	if _, err = p.writer.Write(p.encData); err != nil {
		return
	}
	p.buffer = nil
	return
}
