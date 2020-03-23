// Package assert provides assertion methods for testing.
package assert

import (
	"bytes"
	"encoding/hex"
	"runtime/debug"
	"testing"
)

// NoError asserts err is nil.
func NoError(tb testing.TB, err error) {
	if err != nil {
		tb.Fatalf("Assert failed: %v\n%s\n", err, debug.Stack())
	}
}

// StringEqual asserts expected string equals actural string.
func StringEqual(tb testing.TB, expected, actrual string) {
	if expected != actrual {
		tb.Fatalf("Assert failed:\nexpected: %s\n actrual: %s\n%s\n", expected, actrual, debug.Stack())
	}
}

// ByteArrayEqual asserts expected byte array equals actual byte array.
func ByteArrayEqual(tb testing.TB, expected, actrual []byte) {
	if !bytes.Equal(expected, actrual) {
		tb.Fatalf("Assert failed:\nexpected: (%d)%s\n actrual: (%d)%s\n%s\n",
			len(expected), hex.EncodeToString(expected),
			len(actrual), hex.EncodeToString(actrual),
			debug.Stack())
	}
}
