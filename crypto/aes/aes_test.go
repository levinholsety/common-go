package aes_test

import (
	"bytes"
	aes2 "crypto/aes"
	"crypto/cipher"
	"io"
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/crypto/aes"
	"github.com/levinholsety/common-go/crypto/padding/pkcs7"
)

func prepareData(tb testing.TB, f func(data []byte)) {
	f([]byte("Hello 世界！"))
	for _, length := range []int{0, 1, 15, 16, 17, 255, 256, 257} {
		data, err := comm.RandomBytes(length)
		assert.NoError(tb, err)
		f(data)
	}
}

func prepareKey(tb testing.TB, f func(key, iv []byte)) {
	key, err := aes.NewKey()
	assert.NoError(tb, err)
	iv, err := aes.NewIV()
	assert.NoError(tb, err)
	f(key, iv)
}

func TestEncryptAndDecrypt(tb *testing.T) {
	prepareKey(tb, func(key, iv []byte) {
		prepareData(tb, func(data []byte) {
			encrypted, err := aes.Encrypt(data, key, iv)
			assert.NoError(tb, err)
			decrypted, err := aes.Decrypt(encrypted, key, iv)
			assert.NoError(tb, err)
			assert.ByteArrayEqual(tb, data, decrypted)
		})
	})
}

func TestEncryptionWriterAndDecryptionReader(tb *testing.T) {
	prepareKey(tb, func(key, iv []byte) {
		prepareData(tb, func(data []byte) {
			r := bytes.NewReader(data)
			w := &bytes.Buffer{}
			ew, err := aes.NewEncryptionWriter(w, key, iv)
			assert.NoError(tb, err)
			_, err = io.Copy(ew, r)
			assert.NoError(tb, err)
			err = ew.Close()
			assert.NoError(tb, err)
			r = bytes.NewReader(w.Bytes())
			dr, err := aes.NewDecryptionReader(r, key, iv)
			assert.NoError(tb, err)
			w = &bytes.Buffer{}
			_, err = io.Copy(w, dr)
			assert.NoError(tb, err)
			assert.ByteArrayEqual(tb, data, w.Bytes())
		})
	})
}

func BenchmarkSystemAESEncrypt(tb *testing.B) {
	tb.StopTimer()
	key, err := aes.NewKey()
	assert.NoError(tb, err)
	iv, err := aes.NewIV()
	assert.NoError(tb, err)
	data, err := comm.RandomBytes(999)
	assert.NoError(tb, err)
	tb.StartTimer()
	for i := 0; i < tb.N; i++ {
		var b cipher.Block
		b, err = aes2.NewCipher(key)
		assert.NoError(tb, err)
		d := pkcs7.NewPadding().AddPadding(data, b.BlockSize())
		result := make([]byte, len(d))
		cipher.NewCBCEncrypter(b, iv).CryptBlocks(result, d)
	}
}

func BenchmarkAESEncrypt(tb *testing.B) {
	tb.StopTimer()
	key, err := aes.NewKey()
	assert.NoError(tb, err)
	iv, err := aes.NewIV()
	assert.NoError(tb, err)
	data, err := comm.RandomBytes(999)
	assert.NoError(tb, err)
	tb.StartTimer()
	for i := 0; i < tb.N; i++ {
		_, err = aes.Encrypt(data, key, iv)
		assert.NoError(tb, err)
	}
}
