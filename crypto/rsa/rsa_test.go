package rsa_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/crypto/rsa"
)

func prepareData(tb testing.TB, f func(data []byte)) {
	f([]byte("Hello 世界！"))
	for _, length := range []int{0, 1, 15, 16, 17, 255, 256, 257} {
		data, err := comm.RandomBytes(length)
		assert.NoError(tb, err)
		f(data)
	}
}

func TestEncryptAndDecrypt(tb *testing.T) {
	privateKey, err := rsa.NewPrivateKey()
	assert.NoError(tb, err)
	prepareData(tb, func(data []byte) {
		encrypted, err := rsa.Encrypt(data, &privateKey.PublicKey)
		assert.NoError(tb, err)
		decrypted, err := rsa.Decrypt(encrypted, privateKey)
		assert.NoError(tb, err)
		assert.ByteArrayEqual(tb, data, decrypted)
	})
}

func TestEncryptionWriterAndDecryptionReader(tb *testing.T) {
	privateKey, err := rsa.NewPrivateKey()
	assert.NoError(tb, err)
	prepareData(tb, func(data []byte) {
		r := bytes.NewReader(data)
		w := &bytes.Buffer{}
		ew := rsa.NewEncryptionWriter(w, &privateKey.PublicKey)
		_, err = io.Copy(ew, r)
		assert.NoError(tb, err)
		err = ew.Close()
		assert.NoError(tb, err)
		r = bytes.NewReader(w.Bytes())
		dr, err := rsa.NewDecryptionReader(r, privateKey)
		assert.NoError(tb, err)
		w = &bytes.Buffer{}
		_, err = io.Copy(w, dr)
		assert.NoError(tb, err)
		assert.ByteArrayEqual(tb, data, w.Bytes())
	})
}
