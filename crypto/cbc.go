package crypto

import (
	"crypto/cipher"
	"runtime"
	"unsafe"
)

// CBC represents AES CBC mode.
type CBC struct {
	cb  cipher.Block
	iv  []byte
	buf []byte
}

var _ cipher.Block = new(CBC)

// NewCBC creates a cipher with CBC mode.
func NewCBC(cipherBlock cipher.Block, iv []byte) cipher.Block {
	blockSize := cipherBlock.BlockSize()
	cbc := &CBC{
		cb:  cipherBlock,
		iv:  make([]byte, blockSize),
		buf: make([]byte, blockSize),
	}
	copy(cbc.iv, iv)
	return cbc
}

// BlockSize returns block size.
func (p *CBC) BlockSize() int {
	return p.cb.BlockSize()
}

// Encrypt encrypts a block.
func (p *CBC) Encrypt(dst, src []byte) {
	xorBytes(p.buf, src, p.iv)
	p.cb.Encrypt(dst, p.buf)
	copy(p.iv, dst)
}

// Decrypt decrypts a block.
func (p *CBC) Decrypt(dst, src []byte) {
	p.cb.Decrypt(p.buf, src)
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
