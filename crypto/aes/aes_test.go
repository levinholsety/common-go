package aes_test

import (
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/crypto/aes"
)

func TestAES(tb *testing.T) {
	key, err := aes.NewKey()
	assert.NoError(tb, err)
	iv, err := aes.NewIV()
	assert.NoError(tb, err)
	// 固定数据测试
	testAES(tb, key, iv, []byte("Hello 世界！"))
	// 随机数据测试
	lengthArray := []int{0, 1, 15, 16, 17, 255, 256, 257}
	for _, length := range lengthArray {
		data, err := comm.RandomBytes(length)
		assert.NoError(tb, err)
		testAES(tb, key, iv, data)
	}
}

func testAES(tb *testing.T, key, iv, data []byte) {
	encrypted, err := aes.EncryptByteArray(key, iv, data)
	assert.NoError(tb, err)
	decrypted, err := aes.DecryptByteArray(key, iv, encrypted)
	assert.NoError(tb, err)
	assert.ByteArrayEqual(tb, data, decrypted)
}
