package crypto

import (
	"bufio"
	"crypto/cipher"
	"io"
)

// NewDecryptionReader creates a decryption reader to read data from encrypted data.
func NewDecryptionReader(r io.Reader, block cipher.Block, padding Padding) (dr io.Reader, err error) {
	blockSize := block.BlockSize()
	buffer := make([]byte, blockSize)
	eof := false
	if _, err = io.ReadFull(r, buffer); err != nil {
		if err == io.EOF {
			eof = true
		} else {
			return
		}
	}
	dr = bufio.NewReaderSize(&decryptionReader{
		reader:    r,
		block:     block,
		padding:   padding,
		blockSize: blockSize,
		buffer:    buffer,
		eof:       eof,
	}, blockSize)
	return
}

type decryptionReader struct {
	reader    io.Reader
	block     cipher.Block
	padding   Padding
	blockSize int
	buffer    []byte
	eof       bool
}

func (p *decryptionReader) Read(buffer []byte) (n int, err error) {
	if p.eof {
		err = io.EOF
		return
	}
	p.block.Decrypt(buffer, p.buffer)
	if _, err = io.ReadFull(p.reader, p.buffer); err != nil {
		if err == io.EOF {
			p.eof = true
			if buffer, err = p.padding.RemovePadding(buffer[:p.blockSize]); err != nil {
				return
			}
			n = len(buffer)
			err = io.EOF
		}
		return
	}
	n = p.blockSize
	return
}
