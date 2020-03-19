package comm

import (
	"crypto/rand"
	"io"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	sf = &Snowflake{
		Epoch:          TimeMillis(time.Date(2016, 11, 1, 0, 0, 0, 0, time.UTC)),
		NodeIDBits:     10,
		SequenceIDBits: 12,
	}
)

// ExecutablePath returns executable path.
func ExecutablePath() (executablePath string, err error) {
	executablePath, err = exec.LookPath(os.Args[0])
	if err != nil {
		return
	}
	executablePath, err = filepath.Abs(executablePath)
	return
}

// StartupPath returns startup path.
func StartupPath() (string, error) {
	return filepath.Abs(".")
}

// UniqueIDForNode generates an unique ID for specified node.
func UniqueIDForNode(nodeID int64) int64 {
	return sf.NewIDForNode(nodeID)
}

// UniqueID generates an unique ID.
func UniqueID() int64 {
	return sf.NewID()
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

// ExecutePiplineCommands executes commands in pipline mode.
func ExecutePiplineCommands(output io.Writer, commands ...*exec.Cmd) (err error) {
	if len(commands) > 0 {
		src := commands[0]
		for _, dst := range commands[1:] {
			dst.Stdin, src.Stdout = io.Pipe()
			src = dst
		}
		src.Stdout = output
		for _, cmd := range commands {
			if err = cmd.Start(); err != nil {
				return
			}
		}
		for _, cmd := range commands {
			if err = cmd.Wait(); err != nil {
				return
			}
			if obj, ok := cmd.Stdin.(io.Closer); ok {
				obj.Close()
			}
			if obj, ok := cmd.Stdout.(io.Closer); ok {
				obj.Close()
			}
		}
	}
	return nil
}
