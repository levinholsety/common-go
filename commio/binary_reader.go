package commio

import (
	"encoding/binary"
	"errors"
	"io"
)

// Delimiters.
const (
	DelimNull byte = 0
)

var (
	errInvalidSize = errors.New("invalid size")
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
}

// Read reads content into data.
func (p *BinaryReader) Read(data interface{}) (err error) {
	err = binary.Read(p.reader, p.ByteOrder, data)
	return
}

// MustRead reads content into data.
func (p *BinaryReader) MustRead(data interface{}) {
	if err := binary.Read(p.reader, p.ByteOrder, data); err != nil {
		panic(err)
	}
	return
}

// MustReadUInt64 reads a uint64 value.
func (p *BinaryReader) MustReadUInt64() (result uint64) {
	p.MustRead(&result)
	return
}

// MustReadInt64 reads a int64 value.
func (p *BinaryReader) MustReadInt64() (result int64) {
	p.MustRead(&result)
	return
}

// MustReadUInt32 reads a uint32 value.
func (p *BinaryReader) MustReadUInt32() (result uint32) {
	p.MustRead(&result)
	return
}

// MustReadInt32 reads a int32 value.
func (p *BinaryReader) MustReadInt32() (result int32) {
	p.MustRead(&result)
	return
}

// MustReadUInt16 reads a uint16 value.
func (p *BinaryReader) MustReadUInt16() (result uint16) {
	p.MustRead(&result)
	return
}

// MustReadInt16 reads a int16 value.
func (p *BinaryReader) MustReadInt16() (result int16) {
	p.MustRead(&result)
	return
}

// MustReadUInt8 reads a uint8 value.
func (p *BinaryReader) MustReadUInt8() (result uint8) {
	p.MustRead(&result)
	return
}

// MustReadInt8 reads a int8 value.
func (p *BinaryReader) MustReadInt8() (result int8) {
	p.MustRead(&result)
	return
}

// MustReadByte reads a byte.
func (p *BinaryReader) MustReadByte() (result byte) {
	p.MustRead(&result)
	return
}

// ReadUInt reads byte array in specified size and returns uint.
func (p *BinaryReader) ReadUInt(size int) (result uint, err error) {
	switch size {
	case 0:
		result = 0
	case 1:
		var v uint8
		if err = p.Read(&v); err != nil {
			return
		}
		result = uint(v)
	case 2:
		var v uint16
		if err = p.Read(&v); err != nil {
			return
		}
		result = uint(v)
	case 4:
		var v uint32
		if err = p.Read(&v); err != nil {
			return
		}
		result = uint(v)
	case 8:
		var v uint64
		if err = p.Read(&v); err != nil {
			return
		}
		result = uint(v)
	default:
		err = errInvalidSize
	}
	return
}

// MustReadUInt reads byte array in specified size and returns uint.
func (p *BinaryReader) MustReadUInt(size int) uint {
	result, err := p.ReadUInt(size)
	if err != nil {
		panic(err)
	}
	return result
}

// ReadByteArray reads byte array in specified size.
func (p *BinaryReader) ReadByteArray(size int) (result []byte, err error) {
	result = make([]byte, size)
	_, err = io.ReadFull(p.reader, result)
	return
}

// MustReadByteArray reads byte array in specified size.
func (p *BinaryReader) MustReadByteArray(size int) []byte {
	result, err := p.ReadByteArray(size)
	if err != nil {
		panic(err)
	}
	return result
}

// ReadByteArrayUntil reads byte array until delim occurs.
func (p *BinaryReader) ReadByteArrayUntil(delim byte) (result []byte, err error) {
	var b byte
	for err = p.Read(&b); err == nil && b != delim; err = p.Read(&b) {
		result = append(result, b)
	}
	return
}

// MustReadByteArrayUntil reads byte array until delim occurs.
func (p *BinaryReader) MustReadByteArrayUntil(delim byte) []byte {
	result, err := p.ReadByteArrayUntil(delim)
	if err != nil {
		panic(err)
	}
	return result
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

// MustReadStringFixed reads string in specified size.
func (p *BinaryReader) MustReadStringFixed(size int) string {
	result, err := p.ReadStringFixed(size)
	if err != nil {
		panic(err)
	}
	return result
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

// MustReadString reads string until null character occurs.
func (p *BinaryReader) MustReadString() string {
	result, err := p.ReadString()
	if err != nil {
		panic(err)
	}
	return result
}
