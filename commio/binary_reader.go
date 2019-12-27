package commio

import (
	"encoding/binary"
	"io"
)

// Delimiters.
const (
	DelimNull byte = 0
)

// NewBinaryReader creates an instance of BinaryReader.
func NewBinaryReader(r io.Reader, o binary.ByteOrder) *BinaryReader {
	return &BinaryReader{
		reader:    r,
		ByteOrder: o,
	}
}

// BinaryReader provides methods to read data.
type BinaryReader struct {
	reader    io.Reader
	ByteOrder binary.ByteOrder
}

// Read reads content into data.
func (p *BinaryReader) Read(data interface{}) (err error) {
	err = binary.Read(p.reader, p.ByteOrder, data)
	return
}

// ReadUInt64 reads a uint64 value.
func (p *BinaryReader) ReadUInt64() (result uint64, err error) {
	err = p.Read(&result)
	return
}

// ReadInt64 reads a int64 value.
func (p *BinaryReader) ReadInt64() (result int64, err error) {
	err = p.Read(&result)
	return
}

// ReadUInt32 reads a uint32 value.
func (p *BinaryReader) ReadUInt32() (result uint32, err error) {
	err = p.Read(&result)
	return
}

// ReadInt32 reads a int32 value.
func (p *BinaryReader) ReadInt32() (result int32, err error) {
	err = p.Read(&result)
	return
}

// ReadUInt16 reads a uint16 value.
func (p *BinaryReader) ReadUInt16() (result uint16, err error) {
	err = p.Read(&result)
	return
}

// ReadInt16 reads a int16 value.
func (p *BinaryReader) ReadInt16() (result int16, err error) {
	err = p.Read(&result)
	return
}

// ReadUInt8 reads a uint8 value.
func (p *BinaryReader) ReadUInt8() (result uint8, err error) {
	err = p.Read(&result)
	return
}

// ReadInt8 reads a int8 value.
func (p *BinaryReader) ReadInt8() (result int8, err error) {
	err = p.Read(&result)
	return
}

// ReadByte reads a byte.
func (p *BinaryReader) ReadByte() (result byte, err error) {
	err = p.Read(&result)
	return
}

// ReadByteArray reads byte array in specified size.
func (p *BinaryReader) ReadByteArray(size int) (result []byte, err error) {
	result = make([]byte, size)
	_, err = io.ReadFull(p.reader, result)
	return
}

// ReadByteArrayUntil reads byte array until delim occurs.
func (p *BinaryReader) ReadByteArrayUntil(delim byte) (result []byte, err error) {
	var b byte
	for err = p.Read(&b); err == nil && b != delim; err = p.Read(&b) {
		result = append(result, b)
	}
	return
}

// ReadStringFixed reads string in specified size.
func (p *BinaryReader) ReadStringFixed(size int) (result string, err error) {
	buf, err := p.ReadByteArray(size)
	if err != nil {
		return
	}
	result = string(buf)
	return
}

// ReadString reads string until null character occurs.
func (p *BinaryReader) ReadString() (result string, err error) {
	buf, err := p.ReadByteArrayUntil(DelimNull)
	if err != nil {
		return
	}
	result = string(buf)
	return
}
