package fileutil

import (
	"io"
	"os"
)

// OpenRead opens a file for reading.
func OpenRead(filename string, onOpen func(file *os.File) error) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return onOpen(file)
}

// OpenWrite opens a file for writing.
func OpenWrite(filename string, onOpen func(file *os.File) error) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return onOpen(file)
}

// ReadBlocks reads file in blocks.
func ReadBlocks(filename string, blockSize int, onReadBlock func(block []byte) error) error {
	return OpenRead(filename, func(file *os.File) error {
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
	})
}
