package ioutil

import (
	"io"
)

// ReadBlocks reads from io.Reader in blocks.
func ReadBlocks(r io.Reader, blockSize int, onReadBlock func(block []byte) error) error {
	buf := make([]byte, blockSize)
	for true {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 && err == io.EOF {
			break
		}
		err = onReadBlock(buf[:n])
		if err != nil {
			return err
		}
	}
	return nil
}
