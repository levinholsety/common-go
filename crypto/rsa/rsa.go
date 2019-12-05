package rsa

import (
	"crypto/rand"
	"crypto/rsa"
)

//Encrypt encrypts data with RSA algorithm.
func Encrypt(data []byte, publicKey *rsa.PublicKey) (encryptedData []byte, err error) {
	keySize := publicKey.N.BitLen() / 8
	blockSize := keySize - 11
	encryptedData = make([]byte, computeEncryptedDataSize(len(data), blockSize, keySize))
	buf := encryptedData
	var encryptedBlock []byte
	for len(data) >= blockSize {
		encryptedBlock, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, data[:blockSize])
		if err != nil {
			return
		}
		copy(buf, encryptedBlock)
		buf = buf[keySize:]
		data = data[blockSize:]
	}
	if len(data) > 0 {
		encryptedBlock, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
		if err != nil {
			return
		}
		copy(buf, encryptedBlock)
	}
	return
}

//Decrypt decrypts data with RSA algorithm.
func Decrypt(data []byte, privateKey *rsa.PrivateKey) (decryptedData []byte, err error) {
	keySize := privateKey.N.BitLen() / 8
	blockSize := keySize - 11
	decryptedData = make([]byte, len(data)/keySize*blockSize)
	buf := decryptedData
	n := 0
	var decryptedBlock []byte
	for len(data) > 0 {
		decryptedBlock, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, data[:keySize])
		if err != nil {
			return
		}
		n += copy(buf, decryptedBlock)
		buf = buf[blockSize:]
		data = data[keySize:]
	}
	decryptedData = decryptedData[:n]
	return
}

func computeEncryptedDataSize(dataLength, inSize, outSize int) int {
	n := dataLength / inSize
	if dataLength%inSize == 0 {
		return outSize * n
	}
	return outSize * (n + 1)
}
