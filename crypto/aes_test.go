package crypto_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/crypto/aes"
	"github.com/levinholsety/common-go/fileutil"
	"github.com/levinholsety/common-go/util"
)

var (
	key                 = aes.NewKey()
	iv                  = aes.NewIV()
	data                = util.RandomBytes(1024*1024*5 - 1)
	encryptedAESData    = aes.NewAES(key).Encrypt(data)
	encryptedAESCBCData = aes.NewAESCBC(key, iv).Encrypt(data)
	dataArray           = [][]byte{
		util.RandomBytes(0xfe),
		util.RandomBytes(0xff),
		util.RandomBytes(0x100),
		util.RandomBytes(0xfffe),
		util.RandomBytes(0xffff),
		util.RandomBytes(0x10000),
	}
)

func Test_AES_ECB_EncryptAndDecrypt(t *testing.T) {
	for _, data := range dataArray {
		encrypted := aes.NewAES(key).Encrypt(data)
		decrypted := aes.NewAES(key).Decrypt(encrypted)
		assert.ByteArrayEqual(t, data, decrypted)
	}
}

func Test_AES_ECB_EncryptAndDecryptStream(t *testing.T) {
	path1 := filepath.Join(os.TempDir(), "aes_ecb_test_path1")
	path2 := filepath.Join(os.TempDir(), "aes_ecb_test_path2")
	path3 := filepath.Join(os.TempDir(), "aes_ecb_test_path3")
	for _, data := range dataArray {
		err := ioutil.WriteFile(path1, data, 0644)
		assert.NoError(t, err)
		_, err = fileutil.Transform(aes.NewAES(key).EncryptStream, path2, path1)
		assert.NoError(t, err)
		_, err = fileutil.Transform(aes.NewAES(key).DecryptStream, path3, path2)
		assert.NoError(t, err)
		decrypted, err := ioutil.ReadFile(path3)
		assert.NoError(t, err)
		assert.ByteArrayEqual(t, data, decrypted)
	}
	os.Remove(path1)
	os.Remove(path2)
	os.Remove(path3)
}

func Test_AES_CBC_EncryptAndDecrypt(t *testing.T) {
	for _, data := range dataArray {
		encrypted := aes.NewAESCBC(key, iv).Encrypt(data)
		decrypted := aes.NewAESCBC(key, iv).Decrypt(encrypted)
		assert.ByteArrayEqual(t, data, decrypted)
	}
}

func Test_AES_CBC_EncryptAndDecryptStream(t *testing.T) {
	path1 := filepath.Join(os.TempDir(), "aes_cbc_test_path1")
	path2 := filepath.Join(os.TempDir(), "aes_cbc_test_path2")
	path3 := filepath.Join(os.TempDir(), "aes_cbc_test_path3")
	for _, data := range dataArray {
		err := ioutil.WriteFile(path1, data, 0644)
		assert.NoError(t, err)
		fileutil.Transform(aes.NewAESCBC(key, iv).EncryptStream, path2, path1)
		assert.NoError(t, err)
		fileutil.Transform(aes.NewAESCBC(key, iv).DecryptStream, path3, path2)
		assert.NoError(t, err)
		decrypted, err := ioutil.ReadFile(path3)
		assert.NoError(t, err)
		assert.ByteArrayEqual(t, data, decrypted)
	}
	os.Remove(path1)
	os.Remove(path2)
	os.Remove(path3)
}

func Benchmark_AES_ECB_Encrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		aes.NewAES(key).Encrypt(data)
	}
}

func Benchmark_AES_ECB_Decrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		aes.NewAES(key).Decrypt(encryptedAESData)
	}
}

func Benchmark_AES_CBC_Encrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		aes.NewAESCBC(key, iv).Encrypt(data)
	}
}

func Benchmark_AES_CBC_Decrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		aes.NewAESCBC(key, iv).Decrypt(encryptedAESCBCData)
	}
}

func Benchmark_AES_ECB_EncryptStream(b *testing.B) {
	b.StopTimer()
	srcPath := filepath.Join(os.TempDir(), "aes_cbc_encrypt_src")
	err := ioutil.WriteFile(srcPath, data, 0644)
	assert.NoError(b, err)
	dstPath := filepath.Join(os.TempDir(), "aes_cbc_encrypt_dst")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fileutil.Transform(aes.NewAES(key).EncryptStream, dstPath, srcPath)
	}
}

func Benchmark_AES_ECB_DecryptStream(b *testing.B) {
	b.StopTimer()
	srcPath := filepath.Join(os.TempDir(), "aes_cbc_decrypt_src")
	err := ioutil.WriteFile(srcPath, encryptedAESData, 0644)
	assert.NoError(b, err)
	dstPath := filepath.Join(os.TempDir(), "aes_cbc_decrypt_dst")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fileutil.Transform(aes.NewAES(key).DecryptStream, dstPath, srcPath)
	}
}

func Benchmark_AES_CBC_EncryptStream(b *testing.B) {
	b.StopTimer()
	srcPath := filepath.Join(os.TempDir(), "aes_cbc_encrypt_src")
	err := ioutil.WriteFile(srcPath, data, 0644)
	assert.NoError(b, err)
	dstPath := filepath.Join(os.TempDir(), "aes_cbc_encrypt_dst")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fileutil.Transform(aes.NewAESCBC(key, iv).EncryptStream, dstPath, srcPath)
	}
}

func Benchmark_AES_CBC_DecryptStream(b *testing.B) {
	b.StopTimer()
	srcPath := filepath.Join(os.TempDir(), "aes_cbc_decrypt_src")
	err := ioutil.WriteFile(srcPath, encryptedAESCBCData, 0644)
	assert.NoError(b, err)
	dstPath := filepath.Join(os.TempDir(), "aes_cbc_decrypt_dst")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fileutil.Transform(aes.NewAESCBC(key, iv).DecryptStream, dstPath, srcPath)
	}
}
