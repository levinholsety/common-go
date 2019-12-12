package paddings

import "github.com/levinholsety/common-go/crypto"

// PKCS7Padding implements PKCS#7 padding algorithm.
var (
	PKCS7 crypto.Padding = &pkcs7Padding{}
)
