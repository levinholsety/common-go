package fileutil

import (
	"io"
	"os"
)

// ReadBlocks reads file in blocks.
func ReadBlocks(filename string, blockSize int, onReadBlock func(block []byte) error) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	buf := make([]byte, blockSize)
	for true {
		n, err := file.Read(buf)
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
