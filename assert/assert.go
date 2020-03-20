// Package assert provides assertion methods for testing.
package assert

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// NoError asserts err is nil.
func NoError(tb testing.TB, err error) {
	if err != nil {
		tb.Errorf("Assert failed: %w", err)
	}
}

// StringEqual asserts expected string equals actural string.
func StringEqual(tb testing.TB, expected, actrual string) {
	if expected != actrual {
		tb.Errorf("Assert failed:\nexpected: %s\nactrual: %s\n", expected, actrual)
	}
}

// ByteArrayEqual asserts expected byte array equals actual byte array.
func ByteArrayEqual(tb testing.TB, expected, actrual []byte) {
	if !bytes.Equal(expected, actrual) {
		tb.Errorf("Assert failed:\nexpected: %s\nactrual: %s\n", hex.EncodeToString(expected), hex.EncodeToString(actrual))
	}
}
