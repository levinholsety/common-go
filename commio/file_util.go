package commio

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

// Transform copies a file to another place and performs specified transformation.
func Transform(transform func(w io.Writer, r io.Reader) (int64, error), dst string, src string) (n int64, err error) {
	err = OpenRead(src, func(srcFile *os.File) error {
		return OpenWrite(dst, func(dstFile *os.File) (err error) {
			n, err = transform(dstFile, srcFile)
			return
		})
	})
	return
}

// CopyFile copies a file to another place.
func CopyFile(dst, src string) (n int64, err error) {
	err = OpenRead(src, func(srcFile *os.File) error {
		return OpenWrite(dst, func(dstFile *os.File) (err error) {
			n, err = io.Copy(dstFile, srcFile)
			return
		})
	})
	return
}
