package crypto

import (
	"bufio"
	"crypto/cipher"
	"io"
)

func NewCipherReader(r io.Reader, b cipher.Block, padding Padding) (p io.Reader, err error) {
	br, err := newBlockReader(r, b, padding)
	if err != nil {
		return
	}
	p = bufio.NewReaderSize(br, br.blockSize)
	return
}

type blockReader struct {
	reader    io.Reader
	block     cipher.Block
	padding   Padding
	blockSize int
	buffer    []byte
}

var _ io.Reader = new(blockReader)

func newBlockReader(r io.Reader, b cipher.Block, padding Padding) (p *blockReader, err error) {
	blockSize := b.BlockSize()
	p = &blockReader{
		reader:    r,
		block:     b,
		padding:   padding,
		blockSize: blockSize,
		buffer:    make([]byte, blockSize),
	}
	_, err = io.ReadFull(r, p.buffer)
	if err != nil {
		if err == io.EOF {
			p.buffer = nil
		}
		return
	}
	return
}

func (p *blockReader) Read(buf []byte) (n int, err error) {
	if p.buffer == nil {
		err = io.EOF
		return
	}
	p.block.Decrypt(buf, p.buffer)
	n = p.blockSize
	_, err = io.ReadFull(p.reader, p.buffer)
	if err != nil {
		if err == io.EOF {
			if n, err = p.padding.RemovePadding(buf[:p.blockSize], p.blockSize); err != nil {
				return
			}
			p.buffer = nil
		}
		return
	}
	return
}