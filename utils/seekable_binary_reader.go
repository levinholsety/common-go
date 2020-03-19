package utils

import (
	"encoding/binary"
	"io"
)

// NewSeekableBinaryReader creates an instance SeekableBinaryReader and returns it.
func NewSeekableBinaryReader(r io.ReadSeeker, o binary.ByteOrder) *SeekableBinaryReader {
	return &SeekableBinaryReader{
		BinaryReader: NewBinaryReader(r, o),
		seeker:       r,
	}
}

// SeekableBinaryReader adds seeking capability.
type SeekableBinaryReader struct {
	*BinaryReader
	seeker io.Seeker
}

// Seek seeks in internal object.
func (p *SeekableBinaryReader) Seek(offset int64, whence int) (int64, error) {
	return p.seeker.Seek(offset, whence)
}
