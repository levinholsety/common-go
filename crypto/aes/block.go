package aes

import "crypto/cipher"

type block struct {
	block cipher.Block
}

func (p *block) DataBlockSize() int {
	return p.block.BlockSize()
}

func (p *block) CipherBlockSize() int {
	return p.block.BlockSize()
}

func (p *block) Encrypt(dst, src []byte) (err error) {
	p.block.Encrypt(dst, src)
	return
}

func (p *block) Decrypt(dst, src []byte) (n int, err error) {
	p.block.Decrypt(dst, src)
	n = len(src)
	return
}
