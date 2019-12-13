package commio

import (
	"encoding/binary"
	"io"
)

// NewBinaryReader creates an instance of BinaryReader.
func NewBinaryReader(r io.Reader, order binary.ByteOrder) *BinaryReader {
	return &BinaryReader{
		reader:    r,
		ByteOrder: order,
	}
}

// BinaryReader provides methods to read data.
type BinaryReader struct {
	reader    io.Reader
	ByteOrder binary.ByteOrder
	n         int
}

// Reader returns internal reader.
func (p *BinaryReader) Reader() io.Reader {
	return p.reader
}

// Number returns number of byte has been read.
func (p *BinaryReader) Number() int {
	return p.n
}

// Read reads content into data.
func (p *BinaryReader) Read(data interface{}) (err error) {
	err = binary.Read(p.reader, p.ByteOrder, data)
	p.n += binary.Size(data)
	return
}

// MustRead reads content into data.
func (p *BinaryReader) MustRead(data interface{}) {
	if err := binary.Read(p.reader, p.ByteOrder, data); err != nil {
		panic(err)
	}
	p.n += binary.Size(data)
	return
}

// MustReadUInt64 reads a uint64 value.
func (p *BinaryReader) MustReadUInt64() (v uint64) {
	p.MustRead(&v)
	return
}

// MustReadInt64 reads a int64 value.
func (p *BinaryReader) MustReadInt64() (v int64) {
	p.MustRead(&v)
	return
}

// MustReadUInt32 reads a uint32 value.
func (p *BinaryReader) MustReadUInt32() (v uint32) {
	p.MustRead(&v)
	return
}

// MustReadInt32 reads a int32 value.
func (p *BinaryReader) MustReadInt32() (v int32) {
	p.MustRead(&v)
	return
}

// MustReadUInt16 reads a uint16 value.
func (p *BinaryReader) MustReadUInt16() (v uint16) {
	p.MustRead(&v)
	return
}

// MustReadInt16 reads a int16 value.
func (p *BinaryReader) MustReadInt16() (v int16) {
	p.MustRead(&v)
	return
}

// MustReadUInt8 reads a uint8 value.
func (p *BinaryReader) MustReadUInt8() (v uint8) {
	p.MustRead(&v)
	return
}

// MustReadInt8 reads a int8 value.
func (p *BinaryReader) MustReadInt8() (v int8) {
	p.MustRead(&v)
	return
}

// MustReadByteArray reads byte array in specified size.
func (p *BinaryReader) MustReadByteArray(size int) (v []byte) {
	v = make([]byte, size)
	p.MustRead(&v)
	return
}

// MustReadUInt reads byte array in specified size and returns uint.
func (p *BinaryReader) MustReadUInt(size int) (v uint) {
	buf := p.MustReadByteArray(size)
	for i, b := range buf {
		v |= (uint(b) << ((size - i - 1) * 8))
	}
	return
}

// MustReadString reads string in specified size.
func (p *BinaryReader) MustReadString(size int) (v string) {
	return string(p.MustReadByteArray(size))
}
