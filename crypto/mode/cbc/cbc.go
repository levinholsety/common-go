// Package cbc implements CBC mode.
package cbc

import (
	"crypto/cipher"
	"runtime"
	"unsafe"
)

// cbcBlock (Cipher Block Chaining) is a cipher mode.
type cbcBlock struct {
	block cipher.Block
	iv    []byte
	buf   []byte
}

var _ cipher.Block = (*cbcBlock)(nil)

// NewCipher creates a cipher with CBC mode.
func NewCipher(b cipher.Block, iv []byte) cipher.Block {
	blockSize := b.BlockSize()
	cbc := &cbcBlock{
		block: b,
		iv:    make([]byte, blockSize),
		buf:   make([]byte, blockSize),
	}
	copy(cbc.iv, iv)
	return cbc
}

// BlockSize returns block size.
func (p *cbcBlock) BlockSize() int {
	return p.block.BlockSize()
}

// Encrypt encrypts a block.
func (p *cbcBlock) Encrypt(dst, src []byte) {
	xorBytes(p.buf, src, p.iv)
	p.block.Encrypt(dst, p.buf)
	copy(p.iv, dst)
}

// Decrypt decrypts a block.
func (p *cbcBlock) Decrypt(dst, src []byte) {
	p.block.Decrypt(p.buf, src)
	xorBytes(dst, p.buf, p.iv)
	copy(p.iv, src)
}

const wordSize = int(unsafe.Sizeof(uintptr(0)))
const supportsUnaligned = runtime.GOARCH == "386" ||
	runtime.GOARCH == "amd64" ||
	runtime.GOARCH == "ppc64" ||
	runtime.GOARCH == "ppc64le" ||
	runtime.GOARCH == "s390x"

func xorBytes(dst, a, b []byte) int {
	if supportsUnaligned {
		return fastXORBytes(dst, a, b)
	}
	return safeXORBytes(dst, a, b)
}

func fastXORBytes(dst, a, b []byte) int {
	n := len(a)
	w := n / wordSize
	if w > 0 {
		dw := *(*[]uintptr)(unsafe.Pointer(&dst))
		aw := *(*[]uintptr)(unsafe.Pointer(&a))
		bw := *(*[]uintptr)(unsafe.Pointer(&b))
		for i := 0; i < w; i++ {
			dw[i] = aw[i] ^ bw[i]
		}
	}
	for i := (n - n%wordSize); i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
}

func safeXORBytes(dst, a, b []byte) int {
	n := len(a)
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
}
