package comm

import (
	"io"
)

// ReadStream reads data from io.Reader. Function onRead will be invoked when some data have been read into buffer.
func ReadStream(r io.Reader, bufferSize int, onRead func(buf []byte) error) error {
	buf := make([]byte, bufferSize)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n > 0 {
			err := onRead(buf[:n])
			if err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}
