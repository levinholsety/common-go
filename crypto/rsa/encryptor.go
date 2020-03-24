package rsa

import (
	"crypto/rand"
	"crypto/rsa"
)

type encryptor struct {
	publicKey *rsa.PublicKey
}

func (p *encryptor) DataBlockSize() int {
	return p.publicKey.Size() - 11
}

func (p *encryptor) CipherBlockSize() int {
	return p.publicKey.Size()
}

func (p *encryptor) Encrypt(dst, src []byte) (err error) {
	buf, err := rsa.EncryptPKCS1v15(rand.Reader, p.publicKey, src)
	if err != nil {
		return
	}
	copy(dst, buf)
	return
}

func (p *encryptor) Decrypt(dst, src []byte) (n int, err error) {
	return
}
