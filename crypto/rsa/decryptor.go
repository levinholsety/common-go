package rsa

import (
	"crypto/rand"
	"crypto/rsa"
)

type decryptor struct {
	privateKey *rsa.PrivateKey
}

func (p *decryptor) DataBlockSize() int {
	return p.privateKey.Size() - 11
}

func (p *decryptor) CipherBlockSize() int {
	return p.privateKey.Size()
}

func (p *decryptor) Encrypt(dst, src []byte) (err error) {
	return
}

func (p *decryptor) Decrypt(dst, src []byte) (n int, err error) {
	buf, err := rsa.DecryptPKCS1v15(rand.Reader, p.privateKey, src)
	if err != nil {
		return
	}
	n = copy(dst, buf)
	return
}
