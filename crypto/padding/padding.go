package padding

import "github.com/levinholsety/common-go/crypto"

// PKCS7Padding implements PKCS#7 padding algorithm.
var (
	PKCS7Padding crypto.Padding = &pkcs7Padding{}
)
