package crypto_test

import (
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/crypto/rsa"
)

func Test_RSA_EncryptAndDecrypt(t *testing.T) {
	data, _ := comm.RandomBytes(10000)
	privateKey := rsa.NewPrivateKey()
	encrypted, err := rsa.Encrypt(data, &privateKey.PublicKey)
	assert.NoError(t, err)
	decrypted, err := rsa.Decrypt(encrypted, privateKey)
	assert.NoError(t, err)
	assert.ByteArrayEqual(t, decrypted, data)
}
