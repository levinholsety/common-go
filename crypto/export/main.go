package main

import "C"

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"os"
	"reflect"
	"unsafe"

	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/crypto/aes"
	"github.com/levinholsety/common-go/crypto/rsa"
)

func main() {}

type cchar C.char

func (p *cchar) goBytes(len C.int) []byte {
	var v []byte
	h := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	h.Data = uintptr(unsafe.Pointer(p))
	h.Len = int(len)
	h.Cap = h.Len
	return v
}

func (p *cchar) goString() string {
	return C.GoString((*C.char)(p))
}

func newCString(value string) *cchar {
	return (*cchar)(C.CString(value))
}

func digest(filename string, h hash.Hash) string {
	if err := comm.OpenRead(filename, func(file *os.File) error {
		return comm.ReadBlocks(file, 0x10000, func(block []byte) (err error) {
			_, err = h.Write(block)
			return
		})
	}); err != nil {
		panic(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

//export md5Digest
func md5Digest(cFilename *cchar) *cchar {
	return newCString(digest(cFilename.goString(), md5.New()))
}

//export sha1Digest
func sha1Digest(cFilename *cchar) *cchar {
	return newCString(digest(cFilename.goString(), sha1.New()))
}

//export sha256Digest
func sha256Digest(cFilename *cchar) *cchar {
	return newCString(digest(cFilename.goString(), sha256.New()))
}

//export aesNewKey
func aesNewKey() *cchar {
	key, err := aes.GenerateKey()
	if err != nil {
		panic(err)
	}
	return newCString(hex.EncodeToString(key))
}

//export aesNewIV
func aesNewIV() *cchar {
	iv, err := aes.GenerateIV()
	if err != nil {
		panic(err)
	}
	return newCString(hex.EncodeToString(iv))
}

//export aesEncrypt
func aesEncrypt(cData *cchar, cDataLength C.int, cKey *cchar, cBuffer *cchar, cBufferLength C.int) C.int {
	keyHexString := cKey.goString()
	key, err := hex.DecodeString(keyHexString)
	if err != nil {
		panic(err)
	}
	data := cData.goBytes(cDataLength)
	buffer := cBuffer.goBytes(cBufferLength)
	encData, err := aes.Encrypt(data, key)
	if err != nil {
		panic(err)
	}
	return C.int(copy(buffer, encData))
}

//export aesDecrypt
func aesDecrypt(cData *cchar, cDataLength C.int, cKey *cchar, cBuffer *cchar, cBufferLength C.int) C.int {
	keyHexString := cKey.goString()
	key, err := hex.DecodeString(keyHexString)
	if err != nil {
		panic(err)
	}
	data := cData.goBytes(cDataLength)
	buffer := cBuffer.goBytes(cBufferLength)
	decData, err := aes.Decrypt(data, key)
	if err != nil {
		panic(err)
	}
	return C.int(copy(buffer, decData))
}

//export aesCBCEncrypt
func aesCBCEncrypt(cData *cchar, cDataLength C.int, cKey *cchar, cIV *cchar, cBuffer *cchar, cBufferLength C.int) C.int {
	keyHexString := cKey.goString()
	key, err := hex.DecodeString(keyHexString)
	if err != nil {
		panic(err)
	}
	ivHexString := cIV.goString()
	iv, err := hex.DecodeString(ivHexString)
	if err != nil {
		panic(err)
	}
	data := cData.goBytes(cDataLength)
	buffer := cBuffer.goBytes(cBufferLength)
	encData, err := aes.EncryptCBC(data, key, iv)
	if err != nil {
		panic(err)
	}
	return C.int(copy(buffer, encData))
}

//export aesCBCDecrypt
func aesCBCDecrypt(cData *cchar, cDataLen C.int, cKey *cchar, cIV *cchar, cBuffer *cchar, cBufferLength C.int) C.int {
	keyHexString := cKey.goString()
	key, err := hex.DecodeString(keyHexString)
	if err != nil {
		panic(err)
	}
	ivHexString := cIV.goString()
	iv, err := hex.DecodeString(ivHexString)
	if err != nil {
		panic(err)
	}
	data := cData.goBytes(cDataLen)
	buffer := cBuffer.goBytes(cBufferLength)
	decData, err := aes.DecryptCBC(data, key, iv)
	if err != nil {
		panic(err)
	}
	return C.int(copy(buffer, decData))
}

//export rsaNewPrivateKey
func rsaNewPrivateKey() *cchar {
	privateKey := rsa.NewPrivateKey()
	return newCString(string(rsa.PEM_PKCS1.EncodePrivateKey(privateKey)))
}

//export rsaExportPublicKey
func rsaExportPublicKey(privateKey *cchar) *cchar {
	key, err := rsa.PEM_PKCS1.DecodePrivateKey([]byte(privateKey.goString()))
	if err != nil {
		panic(err)
	}
	return newCString(string(rsa.PEM.EncodePublicKey(&key.PublicKey)))
}

//export rsaEncrypt
func rsaEncrypt(cData *cchar, cDataLength C.int, cKey *cchar, cBuffer *cchar, cBufferLength C.int) C.int {
	key, err := rsa.PEM.DecodePublicKey([]byte(cKey.goString()))
	if err != nil {
		panic(err)
	}
	data := cData.goBytes(cDataLength)
	buffer := cBuffer.goBytes(cBufferLength)
	encryptedData, err := rsa.Encrypt(data, key)
	if err != nil {
		panic(err)
	}
	return C.int(copy(buffer, encryptedData))
}

//export rsaDecrypt
func rsaDecrypt(cData *cchar, cDataLength C.int, cKey *cchar, cBuffer *cchar, cBufferLength C.int) C.int {
	key, err := rsa.PEM_PKCS1.DecodePrivateKey([]byte(cKey.goString()))
	if err != nil {
		panic(err)
	}
	data := cData.goBytes(cDataLength)
	buffer := cBuffer.goBytes(cBufferLength)
	decryptedData, err := rsa.Decrypt(data, key)
	if err != nil {
		panic(err)
	}
	return C.int(copy(buffer, decryptedData))
}
