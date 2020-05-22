// Package assert provides assertion methods for testing.
package assert

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/num"
)

// NoError asserts err is nil.
func NoError(tb testing.TB, err error) {
	if err != nil {
		print(tb, nil, err)
	}
}

// True asserts that value is true.
func True(tb testing.TB, value bool) {
	Equal(tb, assertBool(true), assertBool(value))
}

// StringEqual asserts that expected string and actrual string are equal.
func StringEqual(tb testing.TB, expected, actrual string) {
	Equal(tb, assertString(expected), assertString(actrual))
}

// IntEqual asserts that expected int and actrual int are equal.
func IntEqual(tb testing.TB, expected, actrual int) {
	Equal(tb, assertInt(expected), assertInt(actrual))
}

// Float64Equal asserts that expected value and actrual value are equal.
func Float64Equal(tb testing.TB, expected, actrual float64) {
	Equal(tb, assertFloat64(expected), assertFloat64(actrual))
}

// ByteArrayEqual asserts that expected byte array and actrual byte array are equal.
func ByteArrayEqual(tb testing.TB, expected, actrual []byte) {
	Equal(tb, assertByteArray(expected), assertByteArray(actrual))
}

// Equal asserts that expected value and actrual value are equal.
func Equal(tb testing.TB, expected, actrual comm.Equalizer) {
	if !expected.Equal(actrual) {
		print(tb, expected, actrual)
	}
}

func print(tb testing.TB, expected, actrual interface{}) {
	fmt.Println("    Assert failed:")
	fmt.Printf("    expected: %v\n", expected)
	fmt.Printf("    actrual: %v\n", actrual)
	stack(os.Stdout)
	tb.FailNow()
}

func stack(w io.StringWriter) {
	skip := 0
	pc, file, line, ok := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()
	currentPackage := funcName[:strings.LastIndexByte(funcName, '.')]
	for {
		skip++
		pc, file, line, ok = runtime.Caller(skip)
		if !ok {
			break
		}
		funcName = runtime.FuncForPC(pc).Name()
		pkg := funcName[:strings.LastIndexByte(funcName, '.')]
		if pkg == currentPackage {
			continue
		}
		w.WriteString("    ")
		w.WriteString(funcName)
		w.WriteString("\n        ")
		w.WriteString(file)
		w.WriteString(fmt.Sprintf(":%d\n", line))
	}
}

type assertByteArray []byte

func (p assertByteArray) Equal(argA comm.Equalizer) bool {
	a := argA.(assertByteArray)
	if len(p) != len(a) {
		return false
	}
	for i, v := range p {
		if v != a[i] {
			return false
		}
	}
	return true
}

func (p assertByteArray) String() string {
	return "0x" + hex.EncodeToString([]byte(p))
}

type assertBool bool

func (v assertBool) Equal(argA comm.Equalizer) bool {
	return v == argA
}

type assertString string

func (v assertString) Equal(argA comm.Equalizer) bool {
	return v == argA
}

type assertInt int

func (v assertInt) Equal(argA comm.Equalizer) bool {
	return v == argA
}

type assertFloat64 float64

func (v assertFloat64) Equal(argA comm.Equalizer) bool {
	return num.Float64Equal(float64(v), float64(argA.(assertFloat64)))
}
