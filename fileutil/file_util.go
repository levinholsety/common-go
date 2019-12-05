package fileutil

import (
	"io"
	"os"

	"github.com/levinholsety/common-go/ioutil"
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
		return ioutil.ReadBlocks(file, blockSize, onReadBlock)
	})
}

func Transform(transform func(w io.Writer, r io.Reader) (int64, error), dst string, src string) (n int64, err error) {
	err = OpenRead(src, func(srcFile *os.File) error {
		return OpenWrite(dst, func(dstFile *os.File) (err error) {
			n, err = transform(dstFile, srcFile)
			return
		})
	})
	return
}
