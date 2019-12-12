package rsa

import (
	"bufio"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/levinholsety/common-go/crypto/aes"
)

type pemPKCS1FormatPrivateKey struct{}

func (f *pemPKCS1FormatPrivateKey) EncodePrivateKey(key *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
}
func (f *pemPKCS1FormatPrivateKey) DecodePrivateKey(pemData []byte) (*rsa.PrivateKey, error) {
	return f.DecodeEncryptedPrivateKey(pemData, nil)
}
func (f *pemPKCS1FormatPrivateKey) DecodeEncryptedPrivateKey(pemData, password []byte) (privateKey *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(pemData)
	if block.Type != "RSA PRIVATE KEY" {
		err = fmt.Errorf("invalid pem type: %s", block.Type)
		return
	}
	var data []byte
	v, ok := block.Headers["Proc-Type"]
	if ok && v == "4,ENCRYPTED" {
		v, ok = block.Headers["DEK-Info"]
		if ok {
			values := strings.SplitN(v, ",", 2)
			if len(values) == 2 && values[0] == "AES-256-CBC" {
				var iv []byte
				iv, err = hex.DecodeString(values[1])
				if err != nil {
					return
				}
				salt := iv[:8]
				key := make([]byte, 32)
				if len(password) == 0 {
					fmt.Print("Password: ")
					scanner := bufio.NewScanner(os.Stdin)
					if scanner.Scan() {
						password = []byte(scanner.Text())
					}
				}
				aes.NewKey(password, salt, md5.New(), key)
				data, err = aes.DecryptCBC(block.Bytes, key, iv)
			} else {
				err = fmt.Errorf("unsupported DEK-Info: %s", v)
				return
			}
		}
	} else {
		data = block.Bytes
	}
	privateKey, err = x509.ParsePKCS1PrivateKey(data)
	return
}
