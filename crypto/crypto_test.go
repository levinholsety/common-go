package crypto_test

import (
	"bytes"
	"crypto/aes"
	"io"
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/crypto"
)

func Test_EncryptAndDecrypt_With_AES_ECB(tb *testing.T) {
	key, _ := newKeyAndIV(tb)
	b, err := aes.NewCipher(key)
	assert.NoError(tb, err)
	padding := new(crypto.PKCS7Padding)
	test(tb, func(data []byte) {
		encrypted := &bytes.Buffer{}
		err = crypto.Encrypt(b, padding, encrypted, bytes.NewReader(data))
		assert.NoError(tb, err)
		decrypted := &bytes.Buffer{}
		err = crypto.Decrypt(b, padding, decrypted, bytes.NewReader(encrypted.Bytes()))
		assert.NoError(tb, err)
		assert.ByteArrayEqual(tb, data, decrypted.Bytes())
	})
}

func Test_EncryptAndDecryptByteArray_With_AES_CBC(tb *testing.T) {
	key, iv := newKeyAndIV(tb)
	b, err := aes.NewCipher(key)
	assert.NoError(tb, err)
	padding := new(crypto.PKCS7Padding)
	test(tb, func(data []byte) {
		encrypted, err := crypto.EncryptByteArray(crypto.NewCBC(b, iv), padding, data)
		assert.NoError(tb, err)
		decrypted, err := crypto.DecryptByteArray(crypto.NewCBC(b, iv), padding, encrypted)
		assert.NoError(tb, err)
		assert.ByteArrayEqual(tb, data, decrypted)
	})
}

func Test_EncryptionWriterAndDecryptionReader_With_AES_CBC(tb *testing.T) {
	key, iv := newKeyAndIV(tb)
	b, err := aes.NewCipher(key)
	assert.NoError(tb, err)
	padding := new(crypto.PKCS7Padding)
	test(tb, func(data []byte) {
		encrypted := &bytes.Buffer{}
		w := crypto.NewEncryptionWriter(encrypted, crypto.NewCBC(b, iv), padding)
		_, err = io.Copy(w, bytes.NewReader(data))
		assert.NoError(tb, err)
		err = w.Close()
		assert.NoError(tb, err)
		decrypted := &bytes.Buffer{}
		r, err := crypto.NewDecryptionReader(bytes.NewReader(encrypted.Bytes()), crypto.NewCBC(b, iv), padding)
		assert.NoError(tb, err)
		_, err = io.Copy(decrypted, r)
		assert.NoError(tb, err)
		assert.ByteArrayEqual(tb, data, decrypted.Bytes())
	})
}

func newKeyAndIV(tb testing.TB) (key, iv []byte) {
	key, err := comm.RandomBytes(32)
	assert.NoError(tb, err)
	iv, err = comm.RandomBytes(16)
	assert.NoError(tb, err)
	return
}
func test(tb testing.TB, f func(data []byte)) {
	// 测试固定数据
	f([]byte("Hello 世界！"))
	// 测试随机数据
	lengthArray := []int{0, 1, 15, 16, 17, 255, 256, 257}
	for _, length := range lengthArray {
		data, err := comm.RandomBytes(length)
		assert.NoError(tb, err)
		f(data)
	}
}
