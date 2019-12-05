package assert

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func NoError(tb testing.TB, err error) {
	if err != nil {
		tb.Errorf("Assert failed:\nerror: %s\n", err.Error())
	}
}

func StringEqual(tb testing.TB, expected, actrual string) {
	if expected != actrual {
		tb.Errorf("Assert failed:\nexpected: %s\nactrual: %s\n", expected, actrual)
	}
}

func ByteArrayEqual(tb testing.TB, expected, actrual []byte) {
	if !bytes.Equal(expected, actrual) {
		tb.Errorf("Assert failed:\nexpected: %s\nactrual: %s\n", hex.EncodeToString(expected), hex.EncodeToString(actrual))
	}
}
