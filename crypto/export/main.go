package main

import "C"

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"reflect"
	"unsafe"

	"github.com/levinholsety/common-go/crypto/aes"
	"github.com/levinholsety/common-go/crypto/rsa"
	"github.com/levinholsety/common-go/fileutil"
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
	err := fileutil.ReadBlocks(filename, 0x10000, func(block []byte) (err error) {
		_, err = h.Write(block)
		return
	})
	if err != nil {
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
	return newCString(hex.EncodeToString(aes.NewKey()))
}

//export aesNewIV
func aesNewIV() *cchar {
	return newCString(hex.EncodeToString(aes.NewIV()))
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
	return C.int(copy(buffer, aes.NewAES(key).Encrypt(data)))
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
	return C.int(copy(buffer, aes.NewAES(key).Decrypt(data)))
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
	return C.int(copy(buffer, aes.NewAESCBC(key, iv).Encrypt(data)))
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
	return C.int(copy(buffer, aes.NewAESCBC(key, iv).Decrypt(data)))
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
