package rsa_test

import (
	"crypto/rsa"
	"encoding/json"
	"testing"

	"github.com/levinholsety/common-go/assert"
	alg "github.com/levinholsety/common-go/crypto/rsa"
)

func Test_RSA_KeyFormatConversion(t *testing.T) {
	key, err := alg.NewPrivateKey()
	assert.NoError(t, err)
	testPrivateKey(t, "PEM PKCS#1", key, alg.PKCS1PrivateKeyFormat())
	testPrivateKey(t, "PEM PKCS#8", key, alg.PKCS8PrivateKeyFormat())
	testPrivateKey(t, "XML", key, alg.XMLKeyFormat())
	testPublicKey(t, "PEM", &key.PublicKey, alg.PEMPublicKeyFormat())
	testPublicKey(t, "XML", &key.PublicKey, alg.XMLKeyFormat())
}

func testPrivateKey(t *testing.T, name string, key *rsa.PrivateKey, f alg.PrivateKeyFormat) {
	original := marshal(t, key)
	encodedKey := f.EncodePrivateKey(key)
	decodedKey, err := f.DecodePrivateKey(encodedKey)
	assert.NoError(t, err)
	current := marshal(t, decodedKey)
	assert.ByteArrayEqual(t, original, current)
}

func testPublicKey(t *testing.T, name string, key *rsa.PublicKey, f alg.PublicKeyFormat) {
	original := marshal(t, key)
	encodedKey := f.EncodePublicKey(key)
	decodedKey, err := f.DecodePublicKey(encodedKey)
	assert.NoError(t, err)
	current := marshal(t, decodedKey)
	assert.ByteArrayEqual(t, original, current)
}

func marshal(t *testing.T, key interface{}) []byte {
	encodedKey, err := json.Marshal(key)
	assert.NoError(t, err)
	return encodedKey
}
