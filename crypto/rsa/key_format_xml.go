package rsa

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/xml"
	"math/big"
)

type xmlFormatKey struct{}

type rsaKeyValue struct {
	XMLName  string `xml:"RSAKeyValue"`
	Modulus  string `xml:"Modulus,omitempty"`
	Exponent string `xml:"Exponent,omitempty"`
	P        string `xml:"P,omitempty"`
	Q        string `xml:"Q,omitempty"`
	DP       string `xml:"DP,omitempty"`
	DQ       string `xml:"DQ,omitempty"`
	InverseQ string `xml:"InverseQ,omitempty"`
	D        string `xml:"D,omitempty"`
}

func (f *xmlFormatKey) EncodePrivateKey(key *rsa.PrivateKey) []byte {
	data, _ := xml.Marshal(rsaKeyValue{
		Modulus:  bigIntToBase64String(key.N),
		Exponent: bigIntToBase64String(big.NewInt(int64(key.E))),
		P:        bigIntToBase64String(key.Primes[0]),
		Q:        bigIntToBase64String(key.Primes[1]),
		DP:       bigIntToBase64String(key.Precomputed.Dp),
		DQ:       bigIntToBase64String(key.Precomputed.Dq),
		InverseQ: bigIntToBase64String(key.Precomputed.Qinv),
		D:        bigIntToBase64String(key.D),
	})
	return data
}
func (f *xmlFormatKey) DecodePrivateKey(xmlString []byte) (key *rsa.PrivateKey, err error) {
	kv := new(rsaKeyValue)
	if err = xml.Unmarshal(xmlString, kv); err != nil {
		return
	}
	bigIntArray, err := batchBase64StringToBigInt(
		kv.Modulus,
		kv.Exponent,
		kv.P,
		kv.Q,
		kv.DP,
		kv.DQ,
		kv.InverseQ,
		kv.D,
	)
	if err != nil {
		return
	}
	key = &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: bigIntArray[0],
			E: int(bigIntArray[1].Int64()),
		},
		Primes: []*big.Int{
			bigIntArray[2],
			bigIntArray[3],
		},
		Precomputed: rsa.PrecomputedValues{
			Dp:        bigIntArray[4],
			Dq:        bigIntArray[5],
			Qinv:      bigIntArray[6],
			CRTValues: []rsa.CRTValue{},
		},
		D: bigIntArray[7],
	}
	return
}
func (f *xmlFormatKey) EncodePublicKey(key *rsa.PublicKey) []byte {
	data, _ := xml.Marshal(rsaKeyValue{
		Modulus:  bigIntToBase64String(key.N),
		Exponent: bigIntToBase64String(big.NewInt(int64(key.E))),
	})
	return data
}
func (f *xmlFormatKey) DecodePublicKey(xmlString []byte) (pub *rsa.PublicKey, err error) {
	kv := new(rsaKeyValue)
	if err = xml.Unmarshal(xmlString, kv); err != nil {
		return
	}
	bigIntArray, err := batchBase64StringToBigInt(kv.Modulus, kv.Exponent)
	if err != nil {
		return
	}
	pub = &rsa.PublicKey{
		N: bigIntArray[0],
		E: int(bigIntArray[1].Int64()),
	}
	return
}

func bigIntToBase64String(val *big.Int) string {
	return base64.StdEncoding.EncodeToString(val.Bytes())
}

func base64StringToBigInt(base64String string) (v *big.Int, err error) {
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return
	}
	v = new(big.Int).SetBytes(data)
	return
}

func batchBase64StringToBigInt(base64StringArray ...string) (values []*big.Int, err error) {
	values = make([]*big.Int, len(base64StringArray))
	for i, base64String := range base64StringArray {
		if values[i], err = base64StringToBigInt(base64String); err != nil {
			return
		}
	}
	return
}
