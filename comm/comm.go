package comm

import (
	"crypto/rand"
	"errors"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
)

var defaultIDGenerator = NewIDGenerator(0)

var (
	//ExecutablePath returns the executable file path.
	ExecutablePath = func() (executablePath string) {
		executablePath, err := exec.LookPath(os.Args[0])
		if err != nil {
			panic(err)
		}
		executablePath, err = filepath.Abs(executablePath)
		if err != nil {
			panic(err)
		}
		return
	}()
	//ExecutableDir returns the directory of executable file.
	ExecutableDir = filepath.Dir(ExecutablePath)
	//StartupPath returns the current working directory.
	StartupPath = func() (startupPath string) {
		startupPath, err := filepath.Abs(filepath.Dir("."))
		if err != nil {
			panic(err)
		}
		return
	}()
)

// GenerateID generates a new ID with default IDGenerator.
func GenerateID() int64 {
	return defaultIDGenerator.GenerateID()
}

//Random fill random bytes in buffer.
func Random(buf []byte) {
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
}

//RandomBytes returns random bytes in specified length.
func RandomBytes(length int) (buf []byte) {
	buf = make([]byte, length)
	Random(buf)
	return
}

//RandomInt returns a random int value from 0 to max.
func RandomInt(max int) int {
	value, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	return int(value.Int64())
}

func Try(f func()) (err error) {
	defer func() {
		if obj := recover(); obj != nil {
			err = obj.(error)
		}
	}()
	f()
	return
}

func Throw(errMsg string) {
	panic(errors.New(errMsg))
}
