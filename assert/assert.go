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

// StringEqual asserts that expected string and actrual string are equal.
func StringEqual(tb testing.TB, expected, actrual string) {
	if expected != actrual {
		tb.Fatalf("Assert failed:\nexpected: %s\n actrual: %s\n%s\n",
			expected, actrual, debug.Stack())
	}
}

// IntEqual asserts that expected int and actrual int are equal.
func IntEqual(tb testing.TB, expected, actrual int) {
	if expected != actrual {
		tb.Fatalf("Assert failed:\nexpected: %d\n actrual: %d\n%s\n",
			expected, actrual, debug.Stack())
	}
}

// ByteArrayEqual asserts that expected byte array and actrual byte array are equal.
func ByteArrayEqual(tb testing.TB, expected, actrual []byte) {
	if !bytes.Equal(expected, actrual) {
		tb.Fatalf("Assert failed:\nexpected: (%d)%s\n actrual: (%d)%s\n%s\n",
			len(expected), hex.EncodeToString(expected),
			len(actrual), hex.EncodeToString(actrual),
			debug.Stack())
	}
}
