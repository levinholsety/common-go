package comm

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
func Transform(dstFileName string, srcFileName string, transform func(w io.Writer, r io.Reader) (int64, error)) (n int64, err error) {
	err = OpenRead(srcFileName, func(srcFile *os.File) error {
		return OpenWrite(dstFileName, func(dstFile *os.File) (err error) {
			n, err = transform(dstFile, srcFile)
			return
		})
	})
	return
}

// CopyFile copies a file to another place.
func CopyFile(dstFileName, srcFileName string) (n int64, err error) {
	err = OpenRead(srcFileName, func(srcFile *os.File) error {
		return OpenWrite(dstFileName, func(dstFile *os.File) (err error) {
			n, err = io.Copy(dstFile, srcFile)
			return
		})
	})
	return
}