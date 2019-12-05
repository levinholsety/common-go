package crypto

import (
	"bufio"
	"crypto/cipher"
	"io"

	"github.com/levinholsety/common-go/ioutil"
)

// Padding provides methods for padding.
type Padding interface {
	AddPadding(data []byte, blockSize int) []byte
	RemovePadding(data []byte) []byte
}

// Cipher provides methods to encrypt or decrypt data.
type Cipher struct {
	Block   cipher.Block
	Padding Padding
}

// Encrypt encrypts bytes.
func (p *Cipher) Encrypt(src []byte) []byte {
	blockSize := p.Block.BlockSize()
	buf := make([]byte, (len(src)/blockSize+1)*blockSize)
	dst := buf
	for len(src) >= blockSize {
		p.Block.Encrypt(dst, src[:blockSize])
		dst = dst[blockSize:]
		src = src[blockSize:]
	}
	paddedBlock := p.Padding.AddPadding(src, blockSize)
	p.Block.Encrypt(dst, paddedBlock)
	return buf
}

// Decrypt decrypts bytes.
func (p *Cipher) Decrypt(src []byte) []byte {
	blockSize := p.Block.BlockSize()
	buf := make([]byte, len(src))
	dst := buf
	for len(src) > 0 {
		p.Block.Decrypt(dst, src[:blockSize])
		dst = dst[blockSize:]
		src = src[blockSize:]
	}
	return p.Padding.RemovePadding(buf)
}

// EncryptStream encrypts data.
func (p *Cipher) EncryptStream(w io.Writer, r io.Reader) (n int64, err error) {
	if _, ok := r.(interface{}).(bufio.Reader); !ok {
		r = bufio.NewReader(r)
	}
	blockSize := p.Block.BlockSize()
	var lastBlock []byte
	buf := make([]byte, blockSize)
	err = ioutil.ReadBlocks(r, blockSize, func(block []byte) error {
		if len(block) == blockSize {
			p.Block.Encrypt(buf, block)
			_, err = w.Write(buf)
			if err != nil {
				return err
			}
		} else {
			lastBlock = block
		}
		return nil
	})
	if err != nil {
		return
	}
	lastBlock = p.Padding.AddPadding(lastBlock, blockSize)
	p.Block.Encrypt(buf, lastBlock)
	_, err = w.Write(buf)
	return
}

// DecryptStream decrypts data.
func (p *Cipher) DecryptStream(w io.Writer, r io.Reader) (n int64, err error) {
	if _, ok := r.(interface{}).(bufio.Reader); !ok {
		r = bufio.NewReader(r)
	}
	blockSize := p.Block.BlockSize()
	var buf []byte
	err = ioutil.ReadBlocks(r, blockSize, func(block []byte) error {
		if buf != nil {
			_, err = w.Write(buf)
			if err != nil {
				return err
			}
		} else {
			buf = make([]byte, blockSize)
		}
		p.Block.Decrypt(buf, block)
		return nil
	})
	if buf != nil {
		buf = p.Padding.RemovePadding(buf)
		w.Write(buf)
	}
	if err != nil {
		return
	}
	return
}
