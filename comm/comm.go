package comm

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
)

var defaultIDGenerator = NewIDGenerator(0)

var (
	//ExecutablePath returns the executable file path.
	ExecutablePath = func() (executablePath string) {
		executablePath, _ = exec.LookPath(os.Args[0])
		executablePath, _ = filepath.Abs(executablePath)
		return
	}()
	//ExecutableDir returns the directory of executable file.
	ExecutableDir = filepath.Dir(ExecutablePath)
	//StartupPath returns the current working directory.
	StartupPath, _ = filepath.Abs(".")
)

// GenerateID generates a new ID with default IDGenerator.
func GenerateID() int64 {
	return defaultIDGenerator.GenerateID()
}

//Random fill random bytes in buffer.
func Random(buf []byte) (err error) {
	_, err = rand.Read(buf)
	return
}

//RandomBytes returns random bytes in specified length.
func RandomBytes(length int) (result []byte, err error) {
	result = make([]byte, length)
	err = Random(result)
	return
}

//RandomInt returns a random int value from 0 to max.
func RandomInt(max int) (result int, err error) {
	value, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	result = int(value.Int64())
	return
}

// Try catches panic in f and returns error.
func Try(f func()) (err error) {
	defer func() {
		if msg := recover(); msg != nil {
			var ok bool
			if err, ok = msg.(error); !ok {
				err = fmt.Errorf("%v", msg)
			}
		}
	}()
	f()
	return
}

// Throw panics with specified message.
func Throw(msg interface{}) {
	if err, ok := msg.(error); ok {
		panic(err)
	}
	panic(fmt.Errorf("%v", msg))
}
