package stream

import (
	"io"
)

// Read reads data from io.Reader.
func Read(r io.Reader, onRead func(buf []byte) error) (err error) {
	return read(r, 8192, onRead)
}

func read(r io.Reader, bufferSize int, onRead func(buf []byte) error) (err error) {
	buf := make([]byte, bufferSize)
	var n int
	for true {
		n, err = r.Read(buf)
		if err != nil && err != io.EOF {
			return
		}
		if n > 0 {
			err = onRead(buf[:n])
			if err != nil {
				return
			}
		}
		if err == io.EOF {
			break
		}
	}
	return
}
