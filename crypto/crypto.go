// Package crypto provides cryptography methods.
package crypto

import (
	"bytes"
	"crypto/cipher"
	"io"
)

// Encrypt reads data from reader and encrypts to writer.
func Encrypt(b cipher.Block, padding Padding, w io.Writer, r io.Reader) (err error) {
	src := make([]byte, b.BlockSize())
	dst := make([]byte, b.BlockSize())
	var n int
	for n, err = io.ReadFull(r, src); err == nil; n, err = io.ReadFull(r, src) {
		b.Encrypt(dst, src)
		_, err = w.Write(dst)
	}
	if err == io.ErrUnexpectedEOF || err == io.EOF {
		padding.AddPadding(src, n)
		b.Encrypt(dst, src)
		_, err = w.Write(dst)
	}
	return
}

// EncryptByteArray encrypts byte array.
func EncryptByteArray(b cipher.Block, padding Padding, data []byte) (encrypted []byte, err error) {
	bs := b.BlockSize()
	s := len(data)
	as := s % bs
	if as > 0 {
		s += bs - as
	}
	w := bytes.NewBuffer(make([]byte, 0, s))
	if err = Encrypt(b, padding, w, bytes.NewReader(data)); err != nil {
		return
	}
	encrypted = w.Bytes()
	return
}

// Decrypt reads data from reader and decrypts to writer.
func Decrypt(b cipher.Block, padding Padding, w io.Writer, r io.Reader) (err error) {
	src := make([]byte, b.BlockSize())
	dst := make([]byte, b.BlockSize())
	if _, err = io.ReadFull(r, src); err != nil {
		return
	}
	b.Decrypt(dst, src)
	for true {
		if _, err = io.ReadFull(r, src); err != nil {
			break
		}
		if _, err = w.Write(dst); err != nil {
			return
		}
		b.Decrypt(dst, src)
	}
	if err == io.EOF {
		if dst, err = padding.RemovePadding(dst); err != nil {
			return
		}
		_, err = w.Write(dst)
	}
	return
}

// DecryptByteArray decrypts byte array.
func DecryptByteArray(b cipher.Block, padding Padding, data []byte) (result []byte, err error) {
	w := bytes.NewBuffer(make([]byte, 0, len(data)))
	if err = Decrypt(b, padding, w, bytes.NewReader(data)); err != nil {
		return
	}
	result = w.Bytes()
	return
}
