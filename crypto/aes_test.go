package crypto_test

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/commio"
	"github.com/levinholsety/common-go/crypto/aes"
	"github.com/levinholsety/common-go/crypto/cipher"
	"github.com/levinholsety/common-go/crypto/mode"
	"github.com/levinholsety/common-go/crypto/paddings"
)

var (
	key, _                 = aes.GenerateKey()
	iv, _                  = aes.GenerateIV()
	data, _                = comm.RandomBytes(1024*1024*5 - 1)
	encryptedAESData, _    = aes.Encrypt(data, key)
	encryptedAESCBCData, _ = aes.EncryptCBC(data, key, iv)
	dataArray, _           = func(lengths ...int) (result [][]byte, err error) {
		for _, length := range lengths {
			var array []byte
			if array, err = comm.RandomBytes(length); err != nil {
				return
			}
			result = append(result, array)
		}
		return
	}(0xfe, 0xff, 0x100, 0xfffe, 0xffff, 0x10000)
)

func Test_AES_ECB_EncryptAndDecrypt(t *testing.T) {
	for _, data := range dataArray {
		encrypted, err := aes.Encrypt(data, key)
		assert.NoError(t, err)
		decrypted, err := aes.Decrypt(encrypted, key)
		assert.NoError(t, err)
		assert.ByteArrayEqual(t, data, decrypted)
	}
}

func Test_AES_ECB_EncryptAndDecryptStream(t *testing.T) {
	path1 := filepath.Join(os.TempDir(), "aes_ecb_test_path1")
	defer os.Remove(path1)
	path2 := filepath.Join(os.TempDir(), "aes_ecb_test_path2")
	defer os.Remove(path2)
	path3 := filepath.Join(os.TempDir(), "aes_ecb_test_path3")
	defer os.Remove(path3)
	for _, data := range dataArray {
		err := ioutil.WriteFile(path1, data, 0644)
		assert.NoError(t, err)
		_, err = commio.Transform(newEncryptor(key, nil), path2, path1)
		assert.NoError(t, err)
		_, err = commio.Transform(newDecryptor(key, nil), path3, path2)
		assert.NoError(t, err)
		decrypted, err := ioutil.ReadFile(path3)
		assert.NoError(t, err)
		assert.ByteArrayEqual(t, data, decrypted)
	}
}

func Test_AES_CBC_EncryptAndDecrypt(t *testing.T) {
	for _, data := range dataArray {
		encrypted, err := aes.EncryptCBC(data, key, iv)
		assert.NoError(t, err)
		decrypted, err := aes.DecryptCBC(encrypted, key, iv)
		assert.NoError(t, err)
		assert.ByteArrayEqual(t, data, decrypted)
	}
}

func Test_AES_CBC_EncryptAndDecryptStream(t *testing.T) {
	path1 := filepath.Join(os.TempDir(), "aes_cbc_test_path1")
	defer os.Remove(path1)
	path2 := filepath.Join(os.TempDir(), "aes_cbc_test_path2")
	defer os.Remove(path2)
	path3 := filepath.Join(os.TempDir(), "aes_cbc_test_path3")
	defer os.Remove(path3)
	for _, data := range dataArray {
		err := ioutil.WriteFile(path1, data, 0644)
		assert.NoError(t, err)
		commio.Transform(newEncryptor(key, iv), path2, path1)
		assert.NoError(t, err)
		commio.Transform(newDecryptor(key, iv), path3, path2)
		assert.NoError(t, err)
		decrypted, err := ioutil.ReadFile(path3)
		assert.NoError(t, err)
		assert.ByteArrayEqual(t, data, decrypted)
	}
}

func Benchmark_AES_ECB_Encrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		aes.Encrypt(data, key)
	}
}

func Benchmark_AES_ECB_Decrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		aes.Decrypt(encryptedAESData, key)
	}
}

func Benchmark_AES_CBC_Encrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		aes.EncryptCBC(data, key, iv)
	}
}

func Benchmark_AES_CBC_Decrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		aes.DecryptCBC(encryptedAESCBCData, key, iv)
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
		commio.Transform(newEncryptor(key, nil), dstPath, srcPath)
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
		commio.Transform(newDecryptor(key, nil), dstPath, srcPath)
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
		commio.Transform(newEncryptor(key, iv), dstPath, srcPath)
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
		commio.Transform(newDecryptor(key, iv), dstPath, srcPath)
	}
}

type encryptor struct {
	key []byte
	iv  []byte
}

func newEncryptor(key, iv []byte) func(w io.Writer, r io.Reader) (n int64, err error) {
	p := &encryptor{
		key: key,
		iv:  iv,
	}
	return p.transform
}

func (p *encryptor) transform(w io.Writer, r io.Reader) (n int64, err error) {
	b, err := aes.NewECB(key)
	if iv != nil {
		b = mode.NewCBC(b, iv)
	}
	cw := cipher.NewCipherWriter(w, b, paddings.PKCS7)
	if n, err = io.Copy(cw, r); err != nil {
		return
	}
	err = cw.Close()
	return
}

type decryptor struct {
	key []byte
	iv  []byte
}

func newDecryptor(key, iv []byte) func(w io.Writer, r io.Reader) (n int64, err error) {
	p := &decryptor{
		key: key,
		iv:  iv,
	}
	return p.transform
}

func (p *decryptor) transform(w io.Writer, r io.Reader) (n int64, err error) {
	b, err := aes.NewECB(key)
	if iv != nil {
		b = mode.NewCBC(b, iv)
	}
	cr, err := cipher.NewCipherReader(r, b, paddings.PKCS7)
	if err != nil {
		return
	}
	n, err = io.Copy(w, cr)
	return
}
